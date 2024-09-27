package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type IPInfo struct {
	IP        string `json:"ip"`
	Subnet    string `json:"subnet"`
	Network   string `json:"network"`
	Broadcast string `json:"broadcast"`
	Class     string `json:"class"`
	IsPrivate bool   `json:"is_private"`
}

func main() {
	http.HandleFunc("/api/ipinfo", handleIPInfo)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleIPInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ipAddress := r.URL.Query().Get("ip")
	if ipAddress == "" {
		http.Error(w, "IP address is required", http.StatusBadRequest)
		return
	}

	info, err := getIPInfo(ipAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(info)
}

func getIPInfo(ipAddress string) (IPInfo, error) {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return IPInfo{}, fmt.Errorf("Invalid IP address")
	}

	info := IPInfo{
		IP:        ip.String(),
		IsPrivate: ip.IsPrivate(),
		Class:     getIPClass(ip),
	}

	if ipv4 := ip.To4(); ipv4 != nil {
		mask := net.CIDRMask(32, 32) // Start with a /32 mask
		for i := 32; i > 0; i-- {
			testMask := net.CIDRMask(i, 32)
			if ipv4.Mask(testMask).Equal(ipv4.Mask(mask)) {
				mask = testMask
			} else {
				break
			}
		}
		network := ipv4.Mask(mask)
		broadcast := getIPv4Broadcast(network, mask)

		info.Subnet = fmt.Sprintf("%s/%d", network.String(), maskBits(mask))
		info.Network = network.String()
		info.Broadcast = broadcast.String()
	} else {
		// For IPv6, we'll use a /64 prefix as it's commonly used for subnets
		mask := net.CIDRMask(64, 128)
		network := ip.Mask(mask)
		info.Subnet = fmt.Sprintf("%s/64", network.String())
		info.Network = network.String()
		info.Broadcast = "N/A for IPv6"
	}

	return info, nil
}

func getIPClass(ip net.IP) string {
	if ip.To4() == nil {
		return "IPv6"
	}
	firstOctet := ip[0]
	switch {
	case firstOctet >= 1 && firstOctet <= 126:
		return "A"
	case firstOctet >= 128 && firstOctet <= 191:
		return "B"
	case firstOctet >= 192 && firstOctet <= 223:
		return "C"
	case firstOctet >= 224 && firstOctet <= 239:
		return "D (Multicast)"
	case firstOctet >= 240 && firstOctet <= 255:
		return "E (Reserved)"
	default:
		return "Unknown"
	}
}

func maskBits(mask net.IPMask) int {
	bits, _ := mask.Size()
	return bits
}

func getIPv4Broadcast(network net.IP, mask net.IPMask) net.IP {
	broadcast := make(net.IP, 4)
	for i := 0; i < 4; i++ {
		broadcast[i] = network[i] | ^mask[i]
	}
	return broadcast
}
