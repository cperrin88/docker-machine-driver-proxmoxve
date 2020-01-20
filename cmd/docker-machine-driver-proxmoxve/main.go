package main

import (
	"github.com/cperrin88/docker-machine-driver-proxmoxve/pkg/driver"
	"github.com/docker/machine/libmachine/drivers/plugin"
)

func main() {
	plugin.RegisterDriver(driver.NewDriver("default", ""))
}
