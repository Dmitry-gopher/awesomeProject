package service

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Produce() ([]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

type MockPresenter struct {
	mock.Mock
}

func (m *MockPresenter) Present(data []string) error {
	args := m.Called(data)
	return args.Error(0)
}

func TestService_Run_Success(t *testing.T) {
	mockProd := new(MockProducer)
	mockPres := new(MockPresenter)
	mockProd.On("Produce").Return([]string{"http://example.com", "normal text"}, nil)
	mockPres.On("Present", []string{"http://***********", "normal text"}).Return(nil)
	srv := NewService(mockProd, mockPres)

	// Тест самого сервиса
	err := srv.Run()
	require.NoError(t, err)
	mockProd.AssertCalled(t, "Produce")
	mockPres.AssertCalled(t, "Present", []string{"http://***********", "normal text"})
}

func TestService_Run_ProducerError(t *testing.T) {
	mockProd := new(MockProducer)
	mockPres := new(MockPresenter)

	// Ошибка в Produce
	mockProd.On("Produce").Return(nil, errors.New("file not found"))
	srv := NewService(mockProd, mockPres)
	err := srv.Run()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Producer error")
	mockProd.AssertCalled(t, "Produce")
}

func TestService_Run_PresenterError(t *testing.T) {
	mockProd := new(MockProducer)
	mockPres := new(MockPresenter)

	// Успешный Produce, но ошибка в Present
	mockProd.On("Produce").Return([]string{"http://example.com", "normal text"}, nil)
	mockPres.On("Present", mock.Anything).Return(errors.New("write error"))
	srv := NewService(mockProd, mockPres)
	err := srv.Run()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Presenter error")
	mockProd.AssertCalled(t, "Produce")
	mockPres.AssertCalled(t, "Present", mock.Anything)
}

// Тест функции маскирования отдельно (табличным тестированием)
func TestService_ReplaceLinks(t *testing.T) {
	srv := &Service{}

	// Таблица тестов
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Link at start",
			input:    "http://example.com",
			expected: "http://***********",
		},
		{
			name:     "Link in middle",
			input:    "Visit http://example.com for details",
			expected: "Visit http://*********** for details",
		},
		{
			name:     "No link",
			input:    "No links here",
			expected: "No links here",
		},
		{
			name:     "Multiple links",
			input:    "http://example.com and http://test.com",
			expected: "http://*********** and http://********",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	// Запуск тестов
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := srv.replaceLinks(tt.input)
			require.Equal(t, tt.expected, output, "replaceLinks should return correct masked string")
		})
	}
}
