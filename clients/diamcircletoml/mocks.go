package diamcircletoml

import "github.com/stretchr/testify/mock"

// MockClient is a mockable diamcircletoml client.
type MockClient struct {
	mock.Mock
}

// GetDiamcircleToml is a mocking a method
func (m *MockClient) GetDiamcircleToml(domain string) (*Response, error) {
	a := m.Called(domain)
	return a.Get(0).(*Response), a.Error(1)
}

// GetDiamcircleTomlByAddress is a mocking a method
func (m *MockClient) GetDiamcircleTomlByAddress(address string) (*Response, error) {
	a := m.Called(address)
	return a.Get(0).(*Response), a.Error(1)
}
