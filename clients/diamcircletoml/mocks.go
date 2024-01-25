package diamcircletoml

import "github.com/stretchr/testify/mock"

// MockClient is a mockable diamcircletoml client.
type MockClient struct {
	mock.Mock
}

// GetdiamcircleToml is a mocking a method
func (m *MockClient) GetdiamcircleToml(domain string) (*Response, error) {
	a := m.Called(domain)
	return a.Get(0).(*Response), a.Error(1)
}

// GetdiamcircleTomlByAddress is a mocking a method
func (m *MockClient) GetdiamcircleTomlByAddress(address string) (*Response, error) {
	a := m.Called(address)
	return a.Get(0).(*Response), a.Error(1)
}
