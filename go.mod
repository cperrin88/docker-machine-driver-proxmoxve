module github.com/cperrin88/docker-machine-driver-proxmoxve

go 1.13

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190717161051-705d9623b7c1

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Telmate/proxmox-api-go v0.0.0-20200116224409-320525bf3340
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/machine v0.16.2
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/stretchr/testify v1.2.2
	github.com/xorcare/pointer v1.0.0
	golang.org/x/crypto v0.0.0-20200115085410-6d4e4cb37c7d // indirect
)
