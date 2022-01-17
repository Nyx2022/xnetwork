package input

/**
 * @Author: jweny
 * @Author: https://github.com/jweny
 * @Date: 2021/11/4 17:00
 * @Desc:
 */

type Finger interface {
	MatchInfo() []string
}

type ApplyFingers []*ApplyFinger

func (s *ApplyFingers) MatchInfo() []string {
	if s == nil {
		return nil
	}
	var matchList []string
	for _, finger := range *s {
		matchList = append(matchList, finger.Name)
	}
	return matchList
}

func (s *ApplyFingers) Clone() *ApplyFingers {
	newFinger := *s
	return &newFinger
}

// ApplyFinger 端口上的应用指纹
type ApplyFinger struct {
	Name     string `json:"name" yaml:"name" #:"名称 component.app_name"`
	Uuid     string `json:"uuid" yaml:"uuid" #:"uuid"`
	Category string `json:"category" yaml:"-" #:"分类 component.app_category"`
	Version  string `json:"version" yaml:"-" #:"版本 component.app_version"`
}

type ServiceFingers []*ServiceFinger

func (s *ServiceFingers) MatchInfo() []string {
	var matchList []string
	for _, finger := range *s {
		matchList = append(matchList, finger.ProtocolName)
	}
	return matchList
}

func (s *ServiceFingers) Clone() *ServiceFingers {
	newFinger := *s
	return &newFinger
}

// ServiceFinger 服务指纹
type ServiceFinger struct {
	Tunnel        string `json:"tunnel" #:"隧道, 例如ssl tunnel"`
	ProtocolName  string `json:"protocol_name" #:"协议名称, 例如ftp(关键字段, 弱口令检测依赖)"`
	ProtocolType  string `json:"protocol_type" #:"协议类型, tcp / udp"`
	DeviceName    string `json:"device_name" #:"设备名称"`
	DeviceFactory string `json:"device_factory" #:"设备厂商"`
	OSName        string `json:"os_name" #:"操作系统名称"`
	OSFamily      string `json:"os_family" #:"操作系统类型"`
	OSVendor      string `json:"os_vendor" #:"操作系统厂商"`
}

func (s *ServiceFinger) Clone() *ServiceFinger {
	newFinger := *s
	return &newFinger
}
