package api

func (em *EmporiaMock) Status() (available bool, err error) {
	args := em.Called()
	return args.Bool(0), args.Error(1)
}
