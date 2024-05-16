package proxy

import (
	"net"
)

type BatConn struct {
	net.Conn
}

func (b BatConn) Read(p []byte) (n int, err error) {
	var bat *Bat
	if bat, err = DecryptBat(b.Conn); err != nil {
		return 0, err
	}
	copy(p, bat.Packet)
	return len(bat.Packet), nil
}

func (b BatConn) Write(p []byte) (n int, err error) {
	var bat *Bat
	bat = ZeroBat(CMD_CLIENT_EXCHANGE, p)
	if n, err = bat.WriteEncrypt(b.Conn); err != nil {
		return 0, err
	}
	return n, nil
}
