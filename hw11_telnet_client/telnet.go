package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	conn    net.Conn
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *telnetClient) Connect() error {
	var err error

	c.conn, err = net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}

	return nil
}

func (c *telnetClient) Close() error {
	c.in.Close()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *telnetClient) Send() error {
	buf := make([]byte, 1024)

	n, err := c.in.Read(buf)
	if err != nil {
		return err
	}

	if _, err = c.conn.Write(buf[:n]); err != nil {
		return err
	}

	return nil
}

func (c *telnetClient) Receive() error {
	buf := make([]byte, 1024)

	n, err := c.conn.Read(buf)
	if err != nil {
		return err
	}

	if _, err = c.out.Write(buf[:n]); err != nil {
		return err
	}

	return nil
}
