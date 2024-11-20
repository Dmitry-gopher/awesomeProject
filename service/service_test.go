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

	// тест самого сервиса
	err := srv.Run()
	require.NoError(t, err)
	mockProd.AssertCalled(t, "Produce")
	mockPres.AssertCalled(t, "Present", []string{"http://***********", "normal text"})
}

func TestService_Run_ProducerError(t *testing.T) {
	mockProd := new(MockProducer)
	mockPres := new(MockPresenter)

	// ошибка в Produce
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

	// успешный Produce, но ошибка в Present
	mockProd.On("Produce").Return([]string{"http://example.com", "normal text"}, nil)
	mockPres.On("Present", mock.Anything).Return(errors.New("write error"))
	srv := NewService(mockProd, mockPres)
	err := srv.Run()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Presenter error")
	mockProd.AssertCalled(t, "Produce")
	mockPres.AssertCalled(t, "Present", mock.Anything)
}

// тест функции маскирования отдельно
func TestService_ReplaceLinks(t *testing.T) {
	srv := &Service{}
	input := []string{
		"http://example.com",
		"normal text",
		"http://anotherlink.com/path",
		"no link here",
	}
	expected := []string{
		"http://***********",
		"normal text",
		"http://********************",
		"no link here",
	}

	output := srv.replaceLinks(input)
	require.Equal(t, expected, output)
}
