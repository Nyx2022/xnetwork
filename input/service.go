package input

import (
	"fmt"
	"github.com/kataras/golog"
	"net"
	"strconv"
	"strings"
)

// ServiceAsset 服务资产
type ServiceAsset struct {
	AssetID        string          `json:"asset_id" #:"资产编号 (必填)"`
	Host           string          `json:"host" #:"ip or domain (必填)"`
	Port           int             `json:"port" #:"port (必填)"`
	IsIPv4         bool            `json:"is_ipv4" #:"是否为ipv4 (必填)"`
	Network        string          `json:"network" #:"tcp or udp (必填)"`
	ServiceFingers *ServiceFingers `json:"service_finger" #:"服务指纹"`
	ApplyFingers   *ApplyFingers   `json:"apply_finger" #:"端口上的应用指纹"`
}

func (a *ServiceAsset) SetNetwork(network string) *ServiceAsset {
	a.Network = network
	return a
}

func (a *ServiceAsset) String() string {
	return net.JoinHostPort(a.Host, strconv.Itoa(a.Port))
}

func (a *ServiceAsset) SetFinger(name string) *ServiceAsset {
	sf := ServiceFingers{}
	sf = append(sf, &ServiceFinger{
		ProtocolName: name,
	})
	a.ServiceFingers = &sf
	return a
}

func (a *ServiceAsset) GetNetWork() string {
	return a.Network
}

func (a *ServiceAsset) GetProtocol() string {
	sfs := *a.ServiceFingers
	if len(sfs) > 1 {
		golog.Errorf("asset %s multi service fingers", a.AssetID)
		return ""
	}
	return sfs[0].ProtocolName
}

func NewServiceAsset(addr string, network string) (*ServiceAsset, error) {
	tmp := strings.SplitN(addr, ":", 2)
	host := tmp[0]
	port, err := strconv.Atoi(tmp[1])
	if err != nil || host == "" {
		return nil, fmt.Errorf("invalidation service addr %s %w", addr, err)
	}
	if err != nil {
		return nil, err
	}
	return &ServiceAsset{
		Host:           host,
		Port:           port,
		IsIPv4:         true,
		Network:        network,
		ServiceFingers: nil,
	}, nil
}
