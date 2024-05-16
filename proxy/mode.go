package proxy

import (
	"context"
	"game/proxy/sys"
	"log"
	"os/signal"
	"syscall"
)

func NewModeServer(mode string, addr string, port string) {
	SignalBackground()
	SetProxySetting()
	if "client" == mode {
		NewHttpsProxyClient(addr, port)
	}
	if "server" == mode {
		NewHttpsProxyServer(port)
	}
}

func SetProxySetting() error {
	if err := sys.SetGlobalProxy("127.0.0.1:9999", "<local>"); err != nil {
		log.Printf("set Proxy Setting failed ")
		return err
	}
	if err := sys.Flush(); err != nil {
		log.Printf("flush Proxy Setting failed ")
		return err
	}
	return nil
}

func ClearProxySetting() error {
	if err := sys.Off(); err != nil {
		log.Printf("off Proxy Setting failed ")
		return err
	}
	if err := sys.Flush(); err != nil {
		log.Printf("flush Proxy Setting failed ")
		return err
	}
	return nil
}

func SignalBackground() (ctx context.Context, stop context.CancelFunc) {
	baseCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
	go func() {
		<-baseCtx.Done()
		log.Printf("close Proxy Server ")
		ClearProxySetting()
	}()
	return ctx, stop
}
