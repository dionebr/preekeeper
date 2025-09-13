// Package fasthttpproxy - local copy do original
package internal

import (
	"net"
	"net/url"
)

// FasthttpHTTPDialer returns a dialer func for fasthttp that supports HTTP proxies
func FasthttpHTTPDialer(proxyAddr string) func(addr string) (net.Conn, error) {
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		return nil
	}
	return func(addr string) (net.Conn, error) {
		conn, err := net.Dial("tcp", proxyURL.Host)
		if err != nil {
			return nil, err
		}
		// HTTP CONNECT
		_, err = conn.Write([]byte("CONNECT " + addr + " HTTP/1.1\r\nHost: " + addr + "\r\n\r\n"))
		if err != nil {
			conn.Close()
			return nil, err
		}
		// Read response (simplificado)
		buf := make([]byte, 4096)
		_, err = conn.Read(buf)
		if err != nil {
			conn.Close()
			return nil, err
		}
		return conn, nil
	}
}
