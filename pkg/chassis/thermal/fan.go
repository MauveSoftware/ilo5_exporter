// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2022. Licensed under [MIT](LICENSE) license.
//
// SPDX-License-Identifier: MIT

package thermal

import (
	"github.com/MauveSoftware/ilo5_exporter/pkg/common"
)

type Fan struct {
	Name           string        `json:"Name"`
	CurrentReading float64       `json:"Reading"`
	Status         common.Status `json:"Status"`
}
