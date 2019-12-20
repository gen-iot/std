package std

import (
	"fmt"
	"net"
)

type AddressInfo struct {
	IP   string `json:"ip"`
	Mask string `json:"mask"`
}

func (this *AddressInfo) String() string {
	return fmt.Sprintf("IP=%s,Mask=%s", this.IP, this.Mask)
}

type NetInterfaces struct {
	Name    string        `json:"name"`
	Mac     string        `json:"mac"`
	Address []AddressInfo `json:"address"`
}

func (this *NetInterfaces) String() string {
	return fmt.Sprintf("Name=%s,Mac=%s", this.Name, this.Mac)
}

type ListNetFlag = int

const (
	SkipLoopbackIP ListNetFlag = 1
	SkipNoMac      ListNetFlag = 2
	SkipDefault                = SkipLoopbackIP | SkipNoMac
)

func NetworkList(flag ListNetFlag) (netIfs []NetInterfaces, err error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, it := range ifs {
		macStr := it.HardwareAddr.String()
		if flag&SkipNoMac != 0 && len(macStr) == 0 {
			continue
		}
		netIf := NetInterfaces{
			Name: it.Name,
			Mac:  it.HardwareAddr.String(),
		}
		address, err := it.Addrs()
		if err != nil {
			return nil, err
		}
		for _, it := range address {
			ipAddr := it.(*net.IPNet)
			if flag&SkipLoopbackIP != 0 && ipAddr.IP.IsLoopback() {
				continue
			}
			netIf.Address = append(netIf.Address, AddressInfo{
				IP:   ipAddr.IP.String(),
				Mask: ipAddr.Mask.String(),
			})
		}
		netIfs = append(netIfs, netIf)
	}
	return
}
