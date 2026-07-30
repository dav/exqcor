// Package netinfo detects the LAN IP addresses phones can use to reach this
// laptop's server.
package netinfo

import (
	"net"
)

type Candidate struct {
	Interface string `json:"interface"`
	IP        string `json:"ip"`
	Default   bool   `json:"default"`
}

// Candidates returns the up, non-loopback, private IPv4 addresses of this
// machine. The address on the default outbound route (if any) is flagged.
func Candidates() []Candidate {
	def := defaultRouteIP()
	var out []Candidate
	ifaces, err := net.Interfaces()
	if err != nil {
		return out
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipnet.IP.To4()
			if ip == nil || !ip.IsPrivate() {
				continue
			}
			out = append(out, Candidate{
				Interface: iface.Name,
				IP:        ip.String(),
				Default:   def != "" && ip.String() == def,
			})
		}
	}
	return out
}

// Best returns the most likely LAN IP: the default-route address if present,
// otherwise the first candidate, otherwise "".
func Best() string {
	cs := Candidates()
	for _, c := range cs {
		if c.Default {
			return c.IP
		}
	}
	if len(cs) > 0 {
		return cs[0].IP
	}
	return ""
}

// defaultRouteIP finds the local address the OS would use for outbound
// traffic. No packets are sent; the UDP dial only resolves a route.
func defaultRouteIP() string {
	conn, err := net.Dial("udp4", "192.0.2.1:9")
	if err != nil {
		return ""
	}
	defer conn.Close()
	if la, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		if ip := la.IP.To4(); ip != nil {
			return ip.String()
		}
	}
	return ""
}
