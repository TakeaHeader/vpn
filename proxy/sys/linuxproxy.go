//go:build linux

package sys

func SetGlobalProxy(proxyServer string, bypasses ...string) error {
	return nil
}

func Off() error {
	return nil

}

func Flush() error {
	return nil
}
