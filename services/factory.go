package services

import "url-shortener-api/models"

// ServiceFactory handles the creation and dependency injection of services
type ServiceFactory struct{}

// NewServiceFactory creates a new instance of ServiceFactory
func NewServiceFactory() *ServiceFactory {
	return &ServiceFactory{}
}

// CreateURLService creates a new URLService with all its dependencies
func (f *ServiceFactory) CreateURLService() models.URLService {
	storage := NewURLStorage()
	generator := NewShortCodeGenerator()
	validator := NewURLValidator()

	return &URLServiceImpl{
		storage:   storage,
		generator: generator,
		validator: validator,
	}
}
