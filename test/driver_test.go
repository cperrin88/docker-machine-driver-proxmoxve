package test

import (
	"github.com/cperrin88/docker-machine-driver-proxmoxve/pkg/driver"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetConfigFromFlags(t *testing.T) {
	d := driver.NewDriver("default", "path")

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			"proxmoxve-api-url": "https://pve:8006/api2/json",
			"proxmoxve-user":    "root",
			"proxmoxve-pass":    "pass",
		},
		CreateFlags: d.GetCreateFlags(),
	}

	err := d.SetConfigFromFlags(checkFlags)

	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)
}
