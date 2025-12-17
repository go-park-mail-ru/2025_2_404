// Package readyookassaip provides functionality to check if an IP address belongs to YooKassa's known IP ranges.
package readyookassaip

import (
	"net"
	"strings"
)

var yookassaIPRanges = []string{
	"185.71.76.0/27",
	"185.71.77.0/27",
	"77.75.153.0/25",
	"77.75.156.11/32",
	"77.75.156.35/32",
	"77.75.154.128/25",
	"2a02:5180::/32",
}

var yookassaNetworks []*net.IPNet

func init() {
	for _, cidr := range yookassaIPRanges {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic("invalid CIDR in yookassaIPRanges: " + cidr)
		}
		yookassaNetworks = append(yookassaNetworks, ipNet)
	}
}

func IsYooKassaIP(ipStr string) bool {
	if strings.Contains(ipStr, ":") && !strings.Contains(ipStr, "]") {
		ipStr = strings.Split(ipStr, ":")[0]
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, net := range yookassaNetworks {
		if net.Contains(ip) {
			return true
		}
	}
	return false
}
