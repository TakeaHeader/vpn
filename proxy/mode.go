package proxy

import (
	"context"
	"fmt"
	"game/proxy/sys"
	"log"
	"os/signal"
	"syscall"
)

func NewModeServer(mode string, addr string, port int) {
	baseCtx, stop := SignalBackground()
	defer stop()
	if "client" == mode {
		SetProxySetting(fmt.Sprintf("%s:%d", "127.0.0.1", port))
		NewHttpsProxyClient(baseCtx, addr, port)
	}
	if "server" == mode {
		NewHttpsProxyServer(baseCtx, port)
	}
}

func SetProxySetting(server string) error {
	if err := sys.SetGlobalProxy(server, "<local>"); err != nil {
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

func SignalBackground() (context.Context, context.CancelFunc) {
	baseCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
	go func() {
		<-baseCtx.Done()
		log.Printf("close Proxy Server ")
		ClearProxySetting()
	}()
	return baseCtx, stop
}
