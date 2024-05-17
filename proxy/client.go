package proxy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func NewHttpsProxyClient(context context.Context, serAddr string, port int) error {
	client := Client{add: serAddr}
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: client}
	go func() {
		<-context.Done()
		server.Shutdown(context)
	}()
	return server.ListenAndServe()
}

type Client struct {
	add string
}

func (c Client) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	method, host, _ := request.Method, request.RequestURI, request.Proto
	if !("CONNECT" == method) {
		host = request.Host
		if strings.Index(host, ":") == -1 {
			host = fmt.Sprintf("%s:%d", host, 80)
		}
	}
	proxyServer, err := handshake(host, c.add)
	if err != nil {
		log.Printf("handshake failed: %v ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("proxy: %s ", host)
	remote := BatConn{proxyServer}
	if "CONNECT" == method {
		writer.WriteHeader(http.StatusOK)
	} else {
		//---写请求
		err = request.WriteProxy(remote)
		if err != nil {
			log.Printf("proxy err: %v ", err)
			return
		}
	}

	client, err := HijackConn(writer)
	if err != nil {
		log.Printf("HijackConn failed: %v ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	go PipeThenClose(remote, client)
	PipeThenClose(client, remote)
}

// 握手协议
func handshake(target, server string) (proxyServer net.Conn, err error) {
	proxyServer, err = net.DialTimeout("tcp", server, time.Second*3)
	if err != nil {
		log.Printf("proxy Dial Err: %v", err)
		return nil, err
	}
	var bat *Bat
	bat = ZeroBat(CMD_CLIENT_HOST, []byte(target))
	if _, err = bat.WriteEncrypt(proxyServer); err != nil {
		return nil, err
	}
	if bat, err = DecryptBat(proxyServer); err != nil {
		proxyServer.Close()
		return nil, err
	}
	if bat.Cmd != CMD_OK {
		proxyServer.Close()
		return nil, errors.New("client handshake failed.")
	}
	return
}
