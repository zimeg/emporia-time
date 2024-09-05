package api

func (em *EmporiaMock) GetCustomerDevices() (devices []Device, err error) {
	args := em.Called()
	return args.Get(0).([]Device), args.Error(1)
}
