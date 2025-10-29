package server

import (
	"net"
	"testing"
)

func TestListen(t *testing.T) {
	tests := []struct {
		name   string
		conn   *net.UDPConn
		buf    []byte
		expErr bool
	}{
		{
			name:   "empty conn",
			buf:    make([]byte, 100),
			expErr: true,
		},
		{
			name:   "valid conn no ip",
			buf:    make([]byte, 100),
			conn:   &net.UDPConn{},
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Listen(test.conn, test.buf)
			if test.expErr {
				t.Logf("err:%v\n", err)
				if err == nil {
					t.Fatalf("expected err\tgot:%v\n", err)
				}
				return
			}

		})

	}

}
