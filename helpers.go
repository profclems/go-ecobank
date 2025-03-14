package ecobank

import (
	"errors"
	"time"
)

// secureHasher is an interface for types that can set a secure hash.
type secureHasher interface {
	SetHash(string)
	GetHash() string
}

type secureHashOption struct {
	SecureHash string `json:"secureHash" securehash:"ignore"` // optional: generated if not provided
}

// SetHash sets the secure hash.
func (opt *secureHashOption) SetHash(hash string) {
	opt.SecureHash = hash
}

// GetHash returns the secure hash.
func (opt *secureHashOption) GetHash() string {
	return opt.SecureHash
}

var _ secureHasher = (*secureHashOption)(nil)

const (
	timeFormat = "2006-01-02T15:04:05.999"
	dateFormat = "20060102"
)

// This is the order in which we try to parse the timestamps.
// The API is quite inconsistent with the timestamp formats it returns, so we have to try multiple formats.
var formats = []string{
	timeFormat,
	time.DateTime,
	time.RFC3339, // this is enough for both time.RFC3339 and time.RFC3339Nano layouts
	time.DateOnly,
}

// AddTimeFormat adds the new time layouts to the list of formats to try when parsing timestamps.
//
// This is not thread-safe and should be called before making any requests.
func AddTimeFormat(layouts ...string) {
	formats = append(formats, layouts...)
}

// Time is a wrapper around time.Time.
//
// The API returns timestamps in the format "2006-01-02T15:04:05.999" which is
// not the default format used by Go. This type allows us to unmarshal the
// timestamps correctly.
type Time struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	for _, format := range formats {
		nt, formatErr := time.Parse(format, s)
		if formatErr == nil {
			t.Time = nt
			err = nil // found a format that works, reset the error
			break
		}
		err = errors.Join(err, formatErr)
	}

	return err
}

// Date is a wrapper around time.Time.
type Date struct {
	time.Time
}

// NewDate returns a new Date.
func NewDate(time time.Time) Date {
	return Date{Time: time}
}

// MarshalJSON implements the json.Marshaler interface.
func (date Date) MarshalJSON() ([]byte, error) {
	return []byte(date.Time.Format(dateFormat)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (date *Date) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	date.Time, err = time.Parse(dateFormat, s)
	return err
}

// String returns the string representation of the date in the default format used by the API
// which is YYYYMMDD.
//
// This is the same as calling date.Format("20060102").
func (date Date) String() string {
	return date.Format(dateFormat)
}

// Format returns the string representation of the date using the given format.
func (date Date) Format(format string) string {
	return date.Time.Format(format)
}

func checkErr1[A any](_ A, err error) error {
	return err
}
