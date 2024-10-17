package appstore

import (
	"net"
)

func getAppStoreWhitelist() []net.IPNet {
	whitelist := "17.0.0.0/8"
	_, ipNet, _ := net.ParseCIDR(whitelist)
	return []net.IPNet{*ipNet}
}

func IsAppStoreWhitelistedIp(ip string) bool {
	whitelist := getAppStoreWhitelist()

	userIP := net.ParseIP(ip)
	for _, network := range whitelist {
		if network.Contains(userIP) {
			return true
		}
	}

	return false
}
