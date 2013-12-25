package main

import (
	"fmt"
	"net"
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
	} else {
		return NewError("Unknown protocol: %v", rdr.Protocol)
	}

	return nil
}
