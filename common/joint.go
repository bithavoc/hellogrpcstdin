package common

import (
	"io"
	"net"
	"time"
)

type StdinAddr struct {
	s string
}

func NewStdinAddr(s string) *StdinAddr {
	return &StdinAddr{s}
}
func (a *StdinAddr) Network() string {
	return "stdio"
}

func (a *StdinAddr) String() string {
	return a.s
}

type StdStreamJoint struct {
	in     io.Reader
	out    io.Writer
	closed bool
	local  *StdinAddr
	remote *StdinAddr
}

func NewStdStreamJoint(in io.Reader, out io.Writer) *StdStreamJoint {
	return &StdStreamJoint{
		local:  NewStdinAddr("local"),
		remote: NewStdinAddr("remote"),
		in:     in,
		out:    out,
	}
}

func (s *StdStreamJoint) LocalAddr() net.Addr {
	return s.local
}

func (s *StdStreamJoint) RemoteAddr() net.Addr {
	return s.remote
}

func (s *StdStreamJoint) Read(b []byte) (n int, err error) {
	return s.in.Read(b)
}

func (s *StdStreamJoint) Write(b []byte) (n int, err error) {
	return s.out.Write(b)
}

func (s *StdStreamJoint) Close() error {
	s.closed = true
	return nil
}

func (s *StdStreamJoint) SetDeadline(t time.Time) error {
	return nil
}

func (s *StdStreamJoint) SetReadDeadline(t time.Time) error {
	return nil
}

func (s *StdStreamJoint) SetWriteDeadline(t time.Time) error {
	return nil
}
