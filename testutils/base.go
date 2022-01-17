package testutils

import (
	"totem/v1/pkg/global"
	"totem/v1/pkg/input"
	"totem/v1/pkg/utils"
)

func init() {
	global.TotemOptions = *global.DefaultOptions()
}

var (
	VulnHost = "192.168.123.30"
	Options  = global.DefaultOptions()
)

func NewTestServiceAsset(port int) *input.ServiceAsset {
	addr := &input.ServiceAsset{
		Host:    VulnHost,
		Port:    port,
		Network: "tcp",
	}
	addr.SetMeta(&input.TaskMeta{TaskID: utils.RandStr(8)})
	return addr
}
