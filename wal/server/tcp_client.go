package server

import (
	"fmt"
	"net"
	"os"
	"time"
)

func hostname() (string, error) {
	return os.Hostname()
}

func WriteHostName() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s", ":9999"))
	if err != nil {
		return err
	}
	h, err := hostname()
	if err != nil {
		return err
	}
	fmt.Printf("hostname:%s\n", h)
	connDeadLine, err := time.Parse(time.RFC3339, "20ms")
	if err != nil {
		return err
	}
	if err := conn.SetDeadline(connDeadLine); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "trying to register hostname-[%s]", h)
	return write(conn, h, 10)
}

func write(conn net.Conn, hName string, retryCount int) error {
	// start a counter
	i := 0
	// for first run set it to zero
	if i < 1 {
		i = 0
	}
	if n, err := conn.Write([]byte(hName)); n < 1 || err != nil {
		if err == nil {
			i += 1
			// when counter matches retryCount exit
			if retryCount > i {
				write(conn, hName, retryCount)
			}
			// break recursive tries
			return fmt.Errorf("error writing hostname")
		}
		return err
	}
	return nil
}
