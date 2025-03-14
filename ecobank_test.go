package ecobank

import (
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockHTTPClient is a custom RoundTripper that intercepts HTTP requests and returns mocked responses.
type mockHTTPClient struct {
	mockResponse   string
	statusCode     int
	requestHandler func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.requestHandler != nil {
		return m.requestHandler(req)
	}

	resp := httptest.NewRecorder()
	resp.WriteHeader(m.statusCode)
	_, err := resp.WriteString(m.mockResponse)
	if err != nil {
		return nil, err
	}

	return resp.Result(), nil
}

// newMockClient creates a Client with a mock HTTP client.
func newMockClient(t *testing.T, response string, statusCode int) *Client {
	t.Helper()
	mockTransport := &mockHTTPClient{
		mockResponse: response,
		statusCode:   statusCode,
	}

	// Create a client with NewClient
	mockClient, err := NewClient("mock-client-id", "mock-secret", "mock-lab-key")
	require.NoError(t, err)
	mockClient.token = "mock-token"
	mockClient.tokenExpiresAt = time.Now().Add(time.Hour)

	// Inject our mock HTTP client
	mockClient.client.HTTPClient = &http.Client{Transport: mockTransport}

	return mockClient
}

func TestGenerateSecureHashFrom(t *testing.T) {
	type PaymentHeader struct {
		Amount      string `json:"amount"`
		Currency    string `json:"currency"`
		Beneficiary string `json:"beneficiary"`
		Reference   string `json:"reference"`
		IgnoreField string `json:"ignore" securehash:"ignore"`
		DashField   string `json:"-"`
		SecureHash  string `json:"secureHash"`
	}

	type TestStructWithPaymentHeader struct {
		RequestID     string        `json:"requestId"`
		AffiliateCode string        `json:"affiliateCode"`
		PaymentHeader PaymentHeader `json:"paymentHeader"`
		OtherField    string        `json:"otherField"`
	}

	type TestStructWithoutPaymentHeader struct {
		RequestID     string `json:"requestId"`
		AffiliateCode string `json:"affiliateCode"`
		OtherField    string `json:"otherField"`
	}

	testCases := []struct {
		name     string
		input    any
		key      string
		expected string
	}{
		{
			name: "PaymentHeaderOnly",
			input: TestStructWithPaymentHeader{
				RequestID:     "REQ123",
				AffiliateCode: "AFF",
				PaymentHeader: PaymentHeader{
					Amount:      "100.00",
					Currency:    "USD",
					Beneficiary: "Test Beneficiary",
					Reference:   "REF456",
					IgnoreField: "ignore-me",
					DashField:   "dash-me",
					SecureHash:  "secure-me",
				},
				OtherField: "other-value",
			},
			key:      "testKey",
			expected: generateSecureHash("100.00USDTest BeneficiaryREF456", "testKey"),
		},
		{
			name: "NoPaymentHeader",
			input: TestStructWithoutPaymentHeader{
				RequestID:     "REQ123",
				AffiliateCode: "AFF",
				OtherField:    "other-value",
			},
			key:      "testKey",
			expected: generateSecureHash("REQ123AFFother-value", "testKey"),
		},
		{
			name:     "EmptyStruct",
			input:    TestStructWithoutPaymentHeader{},
			key:      "testKey",
			expected: generateSecureHash("", "testKey"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := generateSecureHashFrom(tc.input, tc.key)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestGenerateSecureHash(t *testing.T) {
	data := "testData"
	key := "testKey"
	expected := func() string {
		hash := sha512.Sum512([]byte(data + key))
		return hex.EncodeToString(hash[:])
	}()

	actual := generateSecureHash(data, key)
	assert.Equal(t, expected, actual)
}
