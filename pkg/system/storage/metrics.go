// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2022. Licensed under [MIT](LICENSE) license.
//
// SPDX-License-Identifier: MIT

package storage

import (
	"sync"

	"github.com/MauveSoftware/ilo5_exporter/pkg/client"
	"github.com/MauveSoftware/ilo5_exporter/pkg/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	prefix = "ilo5_storage_"
)

var (
	diskDriveCapacityDesc *prometheus.Desc
	diskDriveHealthyDesc  *prometheus.Desc
)

func init() {
	dl := []string{"host", "location", "model", "media_type"}
	diskDriveCapacityDesc = prometheus.NewDesc(prefix+"disk_capacity_byte", "Capacity of the disk in bytes", dl, nil)
	diskDriveHealthyDesc = prometheus.NewDesc(prefix+"disk_healthy", "Health status of the diks", dl, nil)
}

// Describe describes all metrics for the storage package
func Describe(ch chan<- *prometheus.Desc) {
	ch <- diskDriveCapacityDesc
	ch <- diskDriveHealthyDesc
}

// Collect collects metrics for storage controllers
func Collect(parentPath string, cl client.Client, ch chan<- prometheus.Metric, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()

	collectStorage(parentPath, cl, ch, errCh)
}

func collectStorage(parentPath string, cl client.Client, ch chan<- prometheus.Metric, errCh chan<- error) {
	p := parentPath + "/Storage"
	crtls := common.MemberList{}

	err := cl.Get(p, &crtls)
	if err != nil {
		errCh <- errors.Wrap(err, "could not get storage controller summary")
		return
	}

	for _, l := range crtls.Members {
		collectStorageController(l.Path, cl, ch, errCh)
	}
}

func collectStorageController(path string, cl client.Client, ch chan<- prometheus.Metric, errCh chan<- error) {
	strg := StorageInfo{}

	err := cl.Get(path, &strg)
	if err != nil {
		errCh <- errors.Wrap(err, "could not get storage controller summary")
		return
	}

	for _, drv := range strg.Drives {
		collectDiskDrive(drv.Path, cl, ch, errCh)
	}
}

func collectDiskDrive(path string, cl client.Client, ch chan<- prometheus.Metric, errCh chan<- error) {
	d := DiskDrive{}
	err := cl.Get(path, &d)
	if err != nil {
		errCh <- errors.Wrapf(err, "could not get drive information from %s", path)
		return
	}

	l := []string{cl.HostName(), d.LocationString(), d.Model, d.MediaType}
	ch <- prometheus.MustNewConstMetric(diskDriveCapacityDesc, prometheus.GaugeValue, float64(d.CapacityBytes), l...)
	ch <- prometheus.MustNewConstMetric(diskDriveHealthyDesc, prometheus.GaugeValue, d.Status.HealthValue(), l...)
}
