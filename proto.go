package main

import (
	"net"
)

type Redirect interface {
	SetupRedirect() error
}

type tcp2tcpRedirect portRedirect

func (rdr *tcp2tcpRedirect) SetupRedirect() error {
	// Start Listener
	l, err := net.Listen("tcp", rdr.LocalAddr)
	if err != nil {
		return err
	}

	// Test Connection
	c, err := net.Dial("tcp", rdr.RemoteAddr)
	if err != nil {
		return err
	}
	c.Close()

	go func() {
		for {
			lConn, err := l.Accept()
			if err != nil {
				continue
			}

			rConn, err := net.Dial("tcp", rdr.RemoteAddr)
			if err != nil {
				lConn.Close()
				continue
			}

			go handleCopy(lConn, rConn)
			go handleCopy(rConn, lConn)
		}
	}()

	return nil
}

type tcp2udpRedirect portRedirect

func (rdr *tcp2udpRedirect) SetupRedirect() error {
	lAddr, err := net.ResolveUDPAddr("udp", rdr.LocalAddr)
	if err != nil {
		return err
	}

	rAddr, err := net.ResolveUDPAddr("udp", rdr.RemoteAddr)
	if err != nil {
		return err
	}

	l, err := net.ListenUDP("udp", lAddr)
	if err != nil {
		return err
	}

	go func() {
		for {
			buf := make([]byte, 4096)

			rlen, _, err := l.ReadFromUDP(buf)
			if err != nil {
				continue
			}

			go l.WriteToUDP(buf[:rlen], rAddr)
		}
	}()

	return nil
}

type udp2udpRedirect portRedirect

func (rdr *udp2udpRedirect) SetupRedirect() error {
	return ERR_NOT_IMPLEMENTED
}

type udp2tcpRedirect portRedirect

func (rdr *udp2tcpRedirect) SetupRedirect() error {
	return ERR_NOT_IMPLEMENTED
}
