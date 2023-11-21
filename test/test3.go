package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	fmt.Println(ParseIP("192.168.10.31"))

	fmt.Println(strings.TrimRight("hello world", "wd"))
}

func ParseIP(ip string) string {
	if strings.Contains(ip, "/") == true {
		if strings.Contains(ip, "/32") == true {
			nip := strings.Replace(ip, "/32", "", -1)
			address := net.ParseIP(nip)
			if address == nil {
				log.Fatal("illegal ip address")
			}
			return address.String()
		}
	}
	address := net.ParseIP(ip)
	if address == nil {
		log.Fatal("illegal ip address")
	}
	return address.String()
}
