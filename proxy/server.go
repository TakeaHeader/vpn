package proxy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

func NewHttpsProxyServer(context context.Context, port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	con := make(chan net.Conn)
	go func() {
		for {
			client, err := ln.Accept()
			if err != nil {
				log.Printf(fmt.Sprintf("%s,%v", "accept  failed", err))
				continue
			}
			con <- client
		}
	}()
	for {
		select {
		case <-context.Done():
			return nil
		case client := <-con:
			go handleConn(client)
		}
	}
}

func handleConn(client net.Conn) {
	remote, err := onHandshake(client)
	if err != nil {
		log.Printf(fmt.Sprintf("%s,%v", "handshake failed", err))
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
