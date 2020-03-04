package peers

import (
	"fmt"
	"net"
)

//ListenUDP sets up a UDP server
//on the given port.
func ListenUDP(port string) {
	serverAddr, e := net.ResolveUDPAddr("udp", port)
	if e != nil {
		panic(e)
	}
	conn, e := net.ListenUDP("udp", serverAddr)
	defer conn.Close()
	buf := make([]byte, 4096)
	for {
		n, addr, e := conn.ReadFromUDP(buf)
		fmt.Println("Got packet from: ", addr.String())
		fmt.Println("Message: ", string(buf[:n]))
		if e != nil {
			fmt.Println(e.Error())
		}
	}
}
