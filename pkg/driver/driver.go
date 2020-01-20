package driver

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/state"
	"github.com/xorcare/pointer"
	"time"
)

type Driver struct {
	*drivers.BaseDriver
	ApiUrl  string
	User    string
	Pass    string
	OTP     string
	Node    string
	IsoPath string
	client  *proxmox.Client
	vmRef   *proxmox.VmRef
}

const (
	defaultUser = "root@pam"
	timeout     = 240
)

func (d *Driver) DriverName() string {
	return "proxmoxve"
}

func (d *Driver) Create() error {
	client, err := d.getClient()

	if err != nil {
		return err
	}
	vmId, err := client.GetNextID(0)
	if err != nil {
		return err
	}
	vmr := proxmox.NewVmRef(vmId)
	vmr.SetNode(d.Node)

	config := &proxmox.ConfigQemu{
		Name:        d.MachineName,
		Description: "Docker Machine VM",
		Memory:      2048,
		Boot:        "d",
		Onboot:      true,
		QemuCpu:     "host",
		QemuOs:      "l26",
		QemuCores:   1,
		QemuSockets: 1,
		QemuVcpus:   1,
		QemuIso:     "NFS:iso/rancheros-proxmoxve-autoformat.iso",
		Agent:       1,
		QemuDisks: map[int]map[string]interface{}{
			0: {
				"type":         "virtio",
				"storage":      "local",
				"storage_type": "dir",
				"size":         "30G",
				"backup":       true,
				"cache":        "none",
				"format":       "qcow2",
			},
		},
		QemuNetworks: map[int]map[string]interface{}{
			0: {
				"model":  "virtio",
				"bridge": "vmbr0",
			},
		},
	}

	err = config.CreateVm(vmr, client)
	if err != nil {
		return err
	}

	d.vmRef = vmr

	err = d.Start()
	if err != nil {
		return err
	}

	timeoutEnd := time.Now().Add(time.Second * timeout)

	for time.Now().Before(timeoutEnd) {
		st, err := d.GetState()
		if err != nil {
			return err
		}

		ip, err := d.GetIP()
		if err != nil {
			return err
		}

		if st == state.Running && ip != "" {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return errors.New(fmt.Sprintf("Timeout reached. VM didnt start within %d seconds", timeout))
}

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			Name:   "proxmoxve-api-url",
			Usage:  "Proxmox VE API URL (example: https://hostname:8006/api2/json)",
			EnvVar: "PROXMOXVE_API_URL",
			Value:  "",
		},
		mcnflag.StringFlag{
			Name:   "proxmoxve-user",
			Usage:  "Proxmox VE connection user (default: root)",
			EnvVar: "PROXMOXVE_USER",
			Value:  defaultUser,
		},
		mcnflag.StringFlag{
			Name:   "proxmoxve-pass",
			Usage:  "Proxmox VE connection password",
			EnvVar: "PROXMOXVE_PASS",
			Value:  "",
		},
		mcnflag.StringFlag{
			Name:   "proxmoxve-otp",
			Usage:  "Proxmox VE OTP Token (optional)",
			EnvVar: "PROXMOXVE_OTP",
			Value:  "",
		},
		mcnflag.StringFlag{
			Name:   "proxmoxve-node",
			Usage:  "Proxmox VE node (default: pve)",
			EnvVar: "PROXMOXVE_NODE",
			Value:  "pve",
		},
	}
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", nil
	}
	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

func (d *Driver) GetState() (state.State, error) {
	client, err := d.getClient()
	if err != nil {
		return state.Error, err
	}

	vmr, err := d.getVmRef()
	if err != nil {
		return state.Error, err
	}

	st, err := client.GetVmState(vmr)

	if err != nil {
		return state.Error, err
	}

	switch st["status"] {
	case "stopped":
		return state.Stopped, nil
	case "running":
		return state.Running, nil
	case "starting":
		return state.Starting, nil
	default:
		return state.None, nil
	}
}

func (d *Driver) Kill() error {
	client, err := d.getClient()
	if err != nil {
		return err
	}

	vmr, err := d.getVmRef()
	if err != nil {
		return err
	}

	_, err = client.StopVm(vmr)

	return err
}

func (d *Driver) Remove() error {
	client, err := d.getClient()
	if err != nil {
		return err
	}

	vmr, err := d.getVmRef()
	if err != nil {
		return err
	}

	_, err = client.DeleteVm(vmr)
	return err
}

func (d *Driver) Restart() error {
	client, err := d.getClient()
	if err != nil {
		return err
	}

	vmr, err := d.getVmRef()
	if err != nil {
		return err
	}

	_, err = client.ShutdownVm(vmr)
	if err != nil {
		return err
	}

	_, err = client.StartVm(vmr)
	return err
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.ApiUrl = flags.String("proxmoxve-api-url")
	d.User = flags.String("proxmoxve-user")
	d.Pass = flags.String("proxmoxve-pass")
	d.OTP = flags.String("proxmoxve-otp")
	d.Node = flags.String("proxmoxve-node")

	d.SetSwarmConfigFromFlags(flags)

	return nil
}

func (d *Driver) Start() error {
	client, err := d.getClient()
	if err != nil {
		return err
	}

	vmr, err := d.getVmRef()
	if err != nil {
		return err
	}

	_, err = client.StartVm(vmr)
	return err
}

func (d *Driver) Stop() error {
	client, err := d.getClient()
	if err != nil {
		return err
	}

	vmr, err := d.getVmRef()
	if err != nil {
		return err
	}

	_, err = client.ShutdownVm(vmr)
	return err
}

func (d *Driver) GetIP() (string, error) {
	if d.IPAddress != "" {
		return d.IPAddress, nil
	}
	client, err := d.getClient()
	if err != nil {
		return "", err
	}

	vmr, err := d.getVmRef()
	if err != nil {
		return "", err
	}

	info, err := client.GetVmAgentNetworkInterfaces(vmr)
	if err != nil {
		return "", nil
	}

	for _, v := range info {
		if v.Name == "eth0" {
			d.IPAddress = v.IPAddresses[0].String()
			return d.IPAddress, nil
		}
	}
	return "", nil
}

func (d *Driver) getClient() (*proxmox.Client, error) {
	if d.client != nil {
		return d.client, nil
	}
	tlsconf := &tls.Config{InsecureSkipVerify: true}
	client, err := proxmox.NewClient(d.ApiUrl, nil, tlsconf)
	if err != nil {
		return nil, err
	}

	err = client.Login(d.User, d.Pass, d.OTP)
	if err != nil {
		return nil, err
	}

	d.client = client

	return client, nil
}

func (d *Driver) getVmRef() (*proxmox.VmRef, error) {
	if d.vmRef != nil {
		return d.vmRef, nil
	}
	client, err := d.getClient()
	if err != nil {
		return nil, err
	}

	vmr, err := client.GetVmRefByName(d.MachineName)
	if err != nil {
		return nil, err
	}

	d.vmRef = vmr

	return vmr, nil
}

func NewDriver(hostName, storePath string) drivers.Driver {
	proxmox.Debug = pointer.Bool(true)
	return NewDerivedDriver(hostName, storePath)
}

func NewDerivedDriver(hostName, storePath string) *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
}
