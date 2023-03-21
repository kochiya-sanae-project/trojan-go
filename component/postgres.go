//go:build postgres || full || mini
// +build postgres full mini

package build

import (
	_ "github.com/p4gefau1t/trojan-go/statistic/postgres"
)
