package telnet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Client struct {
	config    *Config
	conn      net.Conn
	closeOnce sync.Once
}

var ErrNoConnected = errors.New("not connected")

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) Connect() error {
	url := net.JoinHostPort(c.config.Host, c.config.Port)
	conn, err := net.DialTimeout("tcp", url, c.config.Timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", url, err)
	}

	c.conn = conn
	return nil
}

func (c *Client) Run(ctx context.Context, in io.Reader, out io.Writer) error {
	if c.conn == nil {
		return ErrNoConnected
	}

	g, _ := errgroup.WithContext(ctx)

	serverClosed := make(chan struct{})

	g.Go(func() error {
		defer close(serverClosed)

		if _, err := io.Copy(out, c.conn); err != nil {
			return fmt.Errorf("reading from connection: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		stdinDone := make(chan error, 1)

		go func() {
			_, err := io.Copy(c.conn, in)
			stdinDone <- err
		}()

		select {
		case <-ctx.Done():
			c.Close()
			return ctx.Err()
		case <-serverClosed:
			c.Close()
			return nil
		case err := <-stdinDone:
			if tcpConn, ok := c.conn.(*net.TCPConn); ok {
				tcpConn.CloseWrite()
			}

			if err != nil {
				return fmt.Errorf("writing to connection: %w", err)
			}

			<-serverClosed

			return nil
		}
	})

	return g.Wait()
}

func (c *Client) Close() error {
	var err error
	c.closeOnce.Do(func() {
		if c.conn != nil {
			err = c.conn.Close()
		}
	})
	return err
}
