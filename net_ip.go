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

func NetworkList() (netIfs []NetInterfaces, err error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, it := range ifs {
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
			netIf.Address = append(netIf.Address, AddressInfo{
				IP:   ipAddr.IP.String(),
				Mask: ipAddr.Mask.String(),
			})
		}
		netIfs = append(netIfs, netIf)
	}
	return
}
