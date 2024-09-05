package api

import (
	"github.com/zimeg/emporia-time/pkg/energy"
	"github.com/zimeg/emporia-time/pkg/times"
)

func (em *EmporiaMock) GetChartUsage(times times.TimeMeasurement) (results energy.EnergyResult, err error) {
	args := em.Called(times)
	return args.Get(0).(energy.EnergyResult), args.Error(1)
}
