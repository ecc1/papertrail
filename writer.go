package papertrail

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type pt struct {
	dest string
	conn net.Conn
}

// Writer returns an io.Writer that writes to papertrailapp.com
// using the destination specified in the PAPERTRAIL environment variable.
func Writer() (io.Writer, error) {
	const ptEnvVar = "PAPERTRAIL"
	dest := os.Getenv(ptEnvVar)
	if len(dest) == 0 {
		return nil, fmt.Errorf("%s environment variable is not set", ptEnvVar)
	}
	return &pt{dest: dest}, nil
}

func (p *pt) Write(data []byte) (int, error) {
	if p.conn == nil {
		conn, err := net.Dial("udp", p.dest)
		if err != nil {
			return 0, err
		}
		p.conn = conn
	}
	n, err := p.conn.Write(data)
	if err != nil {
		// Better luck next time?
		_ = p.conn.Close()
		p.conn = nil
	}
	return n, err
}

// StartLogging is a convenience function that calls log.SetOutput
// to start logging on both os.Stderr and papertrail.
func StartLogging() {
	local := os.Stderr
	remote, err := Writer()
	if err != nil {
		log.Print(err)
		return
	}
	log.SetOutput(io.MultiWriter(local, remote))
}
