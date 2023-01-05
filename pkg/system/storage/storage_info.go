package storage

import "github.com/MauveSoftware/ilo5_exporter/pkg/common"

type StorageInfo struct {
	Drives []common.Member `json:"Drives"`
}
