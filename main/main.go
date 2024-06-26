package main

import (
	"flag"
	"game/proxy"
	"log"
	"strings"
)

func main() {
	mode := flag.String("mode", "", "set run mode")
	addr := flag.String("addr", "", "set server addr")
	port := flag.Int("port", 0, "set listener port ")
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
	if *port <= 0 {
		log.Fatalf("mode port 参数必须或参数不合法")
	}
	proxy.NewModeServer(*mode, *addr, *port)
}
