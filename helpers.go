package ecobank

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
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
//
// When marshaling, the time is formatted using the provided layout. If no layout is provided,
// the time.DateTime layout is used since that is the format mostly expected by the API.
type Time struct {
	time time.Time

	layout string
}

// NewTime returns a new Time.
func NewTime(time time.Time) Time {
	return NewTimeWithLayout(time, "")
}

// NewTimeWithLayout returns a new Time with the given layout.
func NewTimeWithLayout(time time.Time, layout string) Time {
	return Time{time: time, layout: layout}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	if t.layout != "" {
		// try first with the provided layout
		t.time, err = time.Parse(t.layout, s)
		if err == nil {
			return nil
		}
	}

	for _, format := range formats {
		nt, formatErr := time.Parse(format, s)
		if formatErr == nil {
			t.time = nt
			err = nil // found a format that works, reset the error
			break
		}
		err = errors.Join(err, formatErr)
	}

	return err
}

// MarshalJSON implements the json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", t.String())), nil
}

// GetTime returns the time.Time value.
func (t Time) GetTime() time.Time {
	return t.time
}

// String returns the strings representation of the time.
func (t Time) String() string {
	layout := time.DateTime
	if t.layout != "" {
		layout = t.layout
	}
	return t.time.Format(layout)
}

// Date is a wrapper around time.Time.
type Date struct {
	Time
}

// NewDate returns a new Date with the default layout YYYYMMDD used by the API
func NewDate(time time.Time) Date {
	return Date{Time: NewTimeWithLayout(time, dateFormat)}
}

// MarshalJSON implements the json.Marshaler interface.
func (date Date) MarshalJSON() ([]byte, error) {
	return date.time.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (date *Date) UnmarshalJSON(b []byte) (err error) {
	return date.time.UnmarshalJSON(b)
}

func checkErr1[A any](_ A, err error) error {
	return err
}

func formatToStr(v any) string {
	switch s := v.(type) {
	case string:
		return s
	case decimal.Decimal:
		return s.String()
	case bool:
		return strconv.FormatBool(s)
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32)
	case int:
		return strconv.Itoa(s)
	case int64:
		return strconv.FormatInt(s, 10)
	case int32:
		return strconv.Itoa(int(s))
	case int16:
		return strconv.FormatInt(int64(s), 10)
	case int8:
		return strconv.FormatInt(int64(s), 10)
	case uint:
		return strconv.FormatUint(uint64(s), 10)
	case uint64:
		return strconv.FormatUint(s, 10)
	case uint32:
		return strconv.FormatUint(uint64(s), 10)
	case uint16:
		return strconv.FormatUint(uint64(s), 10)
	case uint8:
		return strconv.FormatUint(uint64(s), 10)
	case json.Number:
		return s.String()
	case []byte:
		return string(s)
	case fmt.Stringer:
		return s.String()
	default:
		return ""
	}
}
