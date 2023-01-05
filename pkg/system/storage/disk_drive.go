package storage

import (
	"github.com/MauveSoftware/ilo5_exporter/pkg/common"
)

type DiskDrive struct {
	MediaType string `json:"MediaType"`
	Model     string `json:"Model"`
	Location  []struct {
		Info string `json:"Info"`
	} `json:"Location"`
	CapacityBytes uint64        `json:"CapacityBytes"`
	Status        common.Status `json:"Status"`
}

func (drv *DiskDrive) LocationString() string {
	if len(drv.Location) == 0 {
		return ""
	}

	return drv.Location[0].Info
}
