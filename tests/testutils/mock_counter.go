package testutils

// MockCounterService is a mock implementation of CounterService
type MockCounterService struct {
	counter int64
}

// NewMockCounterService creates a new mock counter service
func NewMockCounterService() *MockCounterService {
	return &MockCounterService{
		counter: 0,
	}
}

// GetNextCounter increments and returns the next counter value
func (m *MockCounterService) GetNextCounter() (int64, error) {
	m.counter++
	return m.counter, nil
}

// GetCurrentCounter returns the current counter value
func (m *MockCounterService) GetCurrentCounter() (int64, error) {
	return m.counter, nil
}

// InitializeCounter initializes the counter
func (m *MockCounterService) InitializeCounter() error {
	m.counter = 0
	return nil
}
