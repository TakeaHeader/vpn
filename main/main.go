package main

import (
	"flag"
	"game/proxy"
	"log"
	"strings"
)

/*
func main() {
	token := flag.String("token", "", "set xddq login token")
	flag.Parse()
	if len(strings.TrimSpace(*token)) == 0 {
		log.Fatalf("token 参数必须")
	}
	packet, err := protocol.PacketFromBase64(*token)
	if err != nil {
		log.Fatalf("解析token参数失败<%v>", err)
	}
	init := protobuf.InitInfo{Token: *token, PlayerID: packet.PlayerID}
	if err := protobuf.NewSocketClient(init); err != nil {
		log.Println(err)
	}
	fmt.Println()
}
*/

func main() {
	mode := flag.String("mode", "", "set run mode")
	addr := flag.String("addr", "", "set server addr")
	port := flag.String("port", "", "set lisener port ")
	flag.Parse()
	if len(strings.TrimSpace(*mode)) == 0 {
		log.Fatalf("mode 参数必须")
	}
	if !(*mode == "server" || *mode == "client") {
		log.Fatalf("mode 参数必须是('server','client')")
	}
	if *mode == "client" && len(strings.TrimSpace(*addr)) == 0 {
		log.Fatalf("mode 'client' addr 参数必须")
	}
	if *port == "" {
		log.Fatalf("mode port 参数必须")
	}
	proxy.NewModeServer(*mode, *addr, *port)
}
