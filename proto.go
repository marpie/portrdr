package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
)

func redirectHandlerTCP(l net.Listener, rdr *Redirect) {
	for {
		lConn, err := l.Accept()
		if err != nil {
			continue
		}

		rConn, err := net.Dial(rdr.Protocol, rdr.RemoteAddr)
		if err != nil {
			lConn.Close()
			continue
		}

		go handleCopy(lConn, rConn)
		go handleCopy(rConn, lConn)
	}
}

func redirectHandlerUDP(l *net.UDPConn, rAddr *net.UDPAddr) {
	for {
		buf := make([]byte, 4096)

		rlen, _, err := l.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		go l.WriteToUDP(buf[:rlen], rAddr)
	}
}

func redirectHandlerTLS(l net.Listener, rdr *Redirect, config *tls.Config) {
	for {
		lConn, err := l.Accept()
		if err != nil {
			ErrorOut(err)
			continue
		}
		conn := tls.Server(lConn, config)
		if err := conn.Handshake(); err != nil {
			ErrorOut(err)
			conn.Close()
			continue
		}
		state := conn.ConnectionState()

		var addr string
		for _, v := range rdr.RemoteAddrs {
			if v.CertId == state.ServerName {
				addr = v.RemoteAddr
				break
			}
		}
		if len(addr) < 1 {
			addr = rdr.RemoteAddrs[0].RemoteAddr
		}

		var rConn net.Conn
		if rdr.Protocol == "tls2tls" {
			cltConfig := &tls.Config{}
			cltConfig.InsecureSkipVerify = rdr.SslSkipVerify
			cltConfig.PreferServerCipherSuites = true
			rConn, err = tls.Dial("tcp", addr, cltConfig)
		} else if rdr.Protocol == "tls2tcp" {
			rConn, err = net.Dial("tcp", addr)
		}

		if err != nil {
			ErrorOut(err)
			conn.Close()
			continue
		}

		go handleCopy(conn, rConn)
		go handleCopy(rConn, conn)
	}
}

func SetupRedirect(rdr Redirect) error {
	fmt.Printf("[*] Setting up [%s] %s <-> %s ...\n", rdr.Protocol, rdr.LocalAddr, rdr.RemoteAddr)

	if (rdr.Protocol == "tcp") || (rdr.Protocol == "tcp6") {
		l, err := net.Listen(rdr.Protocol, rdr.LocalAddr)
		if err != nil {
			return err
		}

		c, err := net.Dial(rdr.Protocol, rdr.RemoteAddr)
		if err != nil {
			return err
		}
		c.Close()

		go redirectHandlerTCP(l, &rdr)
	} else if (rdr.Protocol == "udp") || (rdr.Protocol == "udp6") {
		lAddr, err := net.ResolveUDPAddr(rdr.Protocol, rdr.LocalAddr)
		if err != nil {
			return err
		}

		rAddr, err := net.ResolveUDPAddr(rdr.Protocol, rdr.RemoteAddr)
		if err != nil {
			return err
		}

		l, err := net.ListenUDP(rdr.Protocol, lAddr)
		if err != nil {
			return err
		}

		go redirectHandlerUDP(l, rAddr)
	} else if strings.HasPrefix(rdr.Protocol, "tls") {
		config := &tls.Config{}
		if len(rdr.ApplicationProtocols) > 0 {
			config.NextProtos = rdr.ApplicationProtocols[:]
		}
		config.Rand = rand.Reader

		config.Certificates = make([]tls.Certificate, len(rdr.Certs))
		i := 0
		var err error
		for _, v := range rdr.Certs {
			config.Certificates[i], err = tls.LoadX509KeyPair(v.CertFile, v.KeyFile)
			if err != nil {
				return err
			}
			i++
		}
		config.BuildNameToCertificate()

		conn, err := tls.Listen("tcp", rdr.LocalAddr, config)
		if err != nil {
			return err
		}

		go redirectHandlerTLS(conn, &rdr, config)
	} else {
		return NewError("Unknown protocol: %v", rdr.Protocol)
	}

	return nil
}
