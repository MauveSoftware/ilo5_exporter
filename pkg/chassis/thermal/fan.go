package thermal

import (
	"github.com/MauveSoftware/ilo5_exporter/pkg/common"
)

type Fan struct {
	Name           string        `json:"Name"`
	CurrentReading float64       `json:"Reading"`
	Status         common.Status `json:"Status"`
}
