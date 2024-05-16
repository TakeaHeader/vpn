package proxy

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

func NewHttpsProxyServer(port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	for {
		client, err := ln.Accept()
		if err != nil {
			log.Printf(fmt.Sprintf("%s,%v", "accept  failed", err))
			continue
		}
		go handleConn(client)
	}
}

func handleConn(client net.Conn) {
	remote, err := onHandshake(client)
	if err != nil {
		log.Printf(fmt.Sprintf("%s,%v", "server handshake failed", err))
		return
	}
	clientConn := BatConn{client}
	go PipeThenClose(clientConn, remote)
	PipeThenClose(remote, clientConn)
}

// 握手协议
func onHandshake(conn net.Conn) (remote net.Conn, err error) {
	var bat *Bat
	if bat, err = DecryptBat(conn); err != nil {
		return nil, err
	}
	if bat.Cmd != CMD_CLIENT_HOST {
		return nil, errors.New("exchange host failed")
	}
	host := string(bat.Packet)
	bat = ZeroBat(CMD_OK, nil)
	if _, err := bat.WriteEncrypt(conn); err != nil {
		return nil, err
	}
	remote, err = net.DialTimeout("tcp", host, time.Second*5)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s,%v", host, err))
	}
	log.Printf("accept : %s", host)
	return remote, nil
}
