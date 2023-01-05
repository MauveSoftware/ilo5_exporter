package memory

import "github.com/MauveSoftware/ilo5_exporter/pkg/common"

type MemoryDIMM struct {
	Name   string        `json:"Name"`
	Status common.Status `json:"Status"`
	SizeMB uint64        `json:"CapacityMiB"`
}
