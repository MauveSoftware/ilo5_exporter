// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2022. Licensed under [MIT](LICENSE) license.
//
// SPDX-License-Identifier: MIT

package memory

import (
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/MauveSoftware/ilo5_exporter/pkg/client"
	"github.com/MauveSoftware/ilo5_exporter/pkg/common"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	prefix = "ilo5_memory_"
)

var (
	healthyDesc     *prometheus.Desc
	totalMemory     *prometheus.Desc
	dimmHealthyDesc *prometheus.Desc
	dimmSizeDesc    *prometheus.Desc
)

func init() {
	l := []string{"host"}
	healthyDesc = prometheus.NewDesc(prefix+"healthy", "Health status of the memory", l, nil)
	totalMemory = prometheus.NewDesc(prefix+"total_byte", "Total memory installed in bytes", l, nil)

	l = append(l, "name")
	dimmHealthyDesc = prometheus.NewDesc(prefix+"dimm_healthy", "Health status of processor", l, nil)
	dimmSizeDesc = prometheus.NewDesc(prefix+"dimm_byte", "DIMM size in bytes", l, nil)
}

// Describe describes all metrics for the memory package
func Describe(ch chan<- *prometheus.Desc) {
	ch <- healthyDesc
	ch <- totalMemory
	ch <- dimmHealthyDesc
	ch <- dimmSizeDesc
}

// Collect collects metrics for memory modules
func Collect(parentPath string, cl client.Client, ch chan<- prometheus.Metric, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()

	p := parentPath + "/Memory"
	mem := common.MemberList{}

	err := cl.Get(p, &mem)
	if err != nil {
		errCh <- errors.Wrap(err, "could not get memory summary")
		return
	}

	wg.Add(len(mem.Members))

	for _, m := range mem.Members {
		go collectForDIMM(m.Path, cl, ch, wg, errCh)
	}
}

func collectForDIMM(link string, cl client.Client, ch chan<- prometheus.Metric, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()

	i := strings.Index(link, "Systems/")
	p := link[i:]

	d := MemoryDIMM{}

	err := cl.Get(p, &d)
	if err != nil {
		errCh <- errors.Wrapf(err, "could not get memory information from %s", link)
		return
	}

	l := []string{cl.HostName(), d.Name}

	if d.Status.State == "" {
		return
	}

	ch <- prometheus.MustNewConstMetric(dimmHealthyDesc, prometheus.GaugeValue, d.Status.HealthValue(), l...)
	ch <- prometheus.MustNewConstMetric(dimmSizeDesc, prometheus.GaugeValue, float64(d.SizeMB<<20), l...)
}
