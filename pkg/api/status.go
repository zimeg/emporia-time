package api

import (
	"net/http"

	"github.com/zimeg/emporia-time/internal/errors"
)

// StatusURL points to a file that appears during maintenance
const StatusURL string = "https://s3.amazonaws.com/com.emporiaenergy.manual.ota/maintenance/maintenance.json"

// Status returns if the Emporia API is available
//
// https://github.com/magico13/PyEmVue/blob/master/api_docs.md#detection-of-maintenance
func (emp *Emporia) Status() (bool, error) {
	resp, err := http.Get(StatusURL)
	if err != nil {
		return false, errors.Wrap(errors.ErrEmporiaStatus, err)
	}
	defer resp.Body.Close()
	status := resp.StatusCode == 403
	return status, nil
}
