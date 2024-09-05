package api

import (
	"github.com/stretchr/testify/mock"
)

type EmporiaMock struct {
	mock.Mock
}

func (em *EmporiaMock) SetDevice(deviceID string) {
	em.Called(deviceID)
}

func (em *EmporiaMock) SetToken(token string) {
	em.Called(token)
}
