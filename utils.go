package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func ErrorOut(err error) {
	fmt.Fprintf(os.Stderr, "[E] %v\n", err)
}

func ErrorExit(err error) {
	ErrorOut(err)
	os.Exit(1)
}

func handleCopy(src, dst net.Conn) {
	io.Copy(dst, src)
}
