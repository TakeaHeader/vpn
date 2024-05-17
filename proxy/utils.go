package proxy

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func PipeThenClose(src, dst net.Conn) {
	defer dst.Close()
	for {
		src.SetReadDeadline(time.Now().Add(time.Second * 5))
		buf := make([]byte, 1024*32)
		n, err := src.Read(buf)
		if n > 0 {
			dst.SetWriteDeadline(time.Now().Add(time.Second * 3))
			if _, err := dst.Write(buf[0:n]); err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
	return
}

func HijackConn(w http.ResponseWriter) (net.Conn, error) {
	hijab, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("not http.Hijacker")
	}
	client, _, e := hijab.Hijack()
	if e != nil {
		return nil, errors.New(fmt.Sprintf("Hijack err: %v ", e))

	}
	return client, e
}

// targetHost 目标主机,服务器
func makeHTTPStunnelConn(targetHost, server string) (net.Conn, error) {
	proxyServer, err := net.DialTimeout("tcp", server, time.Second*10)
	if err != nil {
		log.Printf("proxy Dial Err: %v", err)
		return nil, err
	}
	req := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Host: targetHost},
		Header: make(http.Header),
		Host:   targetHost,
		Body:   nil,
	}
	if err := req.Write(proxyServer); err != nil {
		proxyServer.Close()
		log.Printf("proxy proxyServer Err: %v", err)
		return nil, err
	}
	resp, err := http.ReadResponse(bufio.NewReader(proxyServer), req)
	if err != nil || resp.StatusCode != http.StatusOK {
		proxyServer.Close()
		return nil, errors.New(fmt.Sprintf("Connect build failed : %d ", resp.StatusCode))
	}
	return proxyServer, nil
}
