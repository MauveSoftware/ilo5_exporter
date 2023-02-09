// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2022. Licensed under [MIT](LICENSE) license.
//
// SPDX-License-Identifier: MIT

package power

import (
	"github.com/MauveSoftware/ilo5_exporter/pkg/client"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	prefix = "ilo5_power_"
)

var (
	powerCurrentDesc       *prometheus.Desc
	powerAvgDesc           *prometheus.Desc
	powerMinDesc           *prometheus.Desc
	powerMaxDesc           *prometheus.Desc
	powerSupplyHealthyDesc *prometheus.Desc
	powerSupplyEnabledDesc *prometheus.Desc
	powerCapacityDesc      *prometheus.Desc
)

func init() {
	l := []string{"host"}
	lpm := append(l, "id")
	powerCurrentDesc = prometheus.NewDesc(prefix+"current_watt", "Current power consumption in watt", lpm, nil)
	powerAvgDesc = prometheus.NewDesc(prefix+"average_watt", "Average power consumption in watt", lpm, nil)
	powerMinDesc = prometheus.NewDesc(prefix+"min_watt", "Minimum power consumption in watt", lpm, nil)
	powerMaxDesc = prometheus.NewDesc(prefix+"max_watt", "Maximum power consumption in watt", lpm, nil)
	powerCapacityDesc = prometheus.NewDesc(prefix+"capacity_watt", "Power capacity in watt", lpm, nil)

	l = append(l, "serial")
	powerSupplyHealthyDesc = prometheus.NewDesc(prefix+"supply_healthy", "Health status of the power supply", l, nil)
	powerSupplyEnabledDesc = prometheus.NewDesc(prefix+"supply_enabled", "Status of the power supply", l, nil)
}

func Describe(ch chan<- *prometheus.Desc) {
	ch <- powerCurrentDesc
	ch <- powerAvgDesc
	ch <- powerMinDesc
	ch <- powerMaxDesc
	ch <- powerCapacityDesc
	ch <- powerSupplyHealthyDesc
	ch <- powerSupplyEnabledDesc
}

func Collect(parentPath string, cl client.Client, ch chan<- prometheus.Metric) error {
	pwr := Power{}

	err := cl.Get(parentPath+"/Power", &pwr)
	if err != nil {
		return errors.Wrap(err, "could not get power data")
	}

	l := []string{cl.HostName()}

	for _, pwc := range pwr.PowerControl {
		la := append(l, pwc.ID)
		ch <- prometheus.MustNewConstMetric(powerCurrentDesc, prometheus.GaugeValue, pwc.PowerConsumedWatts, la...)
		ch <- prometheus.MustNewConstMetric(powerAvgDesc, prometheus.GaugeValue, pwc.Metrics.AverageConsumedWatts, la...)
		ch <- prometheus.MustNewConstMetric(powerMinDesc, prometheus.GaugeValue, pwc.Metrics.MinConsumedWatts, la...)
		ch <- prometheus.MustNewConstMetric(powerMaxDesc, prometheus.GaugeValue, pwc.Metrics.MaxConsumedWatts, la...)
		ch <- prometheus.MustNewConstMetric(powerCapacityDesc, prometheus.GaugeValue, pwc.PowerCapacityWatts, la...)
	}

	for _, sup := range pwr.PowerSupplies {
		if sup.Status.State == "Absent" {
			continue
		}

		la := append(l, sup.SerialNumber)
		ch <- prometheus.MustNewConstMetric(powerSupplyEnabledDesc, prometheus.GaugeValue, sup.Status.EnabledValue(), la...)
		ch <- prometheus.MustNewConstMetric(powerSupplyHealthyDesc, prometheus.GaugeValue, sup.Status.HealthValue(), la...)
	}

	return nil
}
