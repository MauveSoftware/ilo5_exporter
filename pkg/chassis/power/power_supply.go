package power

import (
	"github.com/MauveSoftware/ilo5_exporter/pkg/common"
)

type PowerSupply struct {
	SerialNumber string        `json:"SerialNumber"`
	Status       common.Status `json:"Status"`
}
