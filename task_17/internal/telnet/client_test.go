package telnet

import (
	"bytes"
	"context"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func startEchoServer(t *testing.T) (string, func()) {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "failed to start echo server")

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				io.Copy(c, c)
			}(conn)
		}
	}()

	return ln.Addr().String(), func() { ln.Close() }
}

func connectToEcho(t *testing.T) (*Client, func()) {
	t.Helper()

	addr, stop := startEchoServer(t)
	host, port, _ := net.SplitHostPort(addr)
	cfg := NewConfig(host, port, 5*time.Second)
	client := NewClient(cfg)

	require.NoError(t, client.Connect())

	return client, func() {
		client.Close()
		stop()
	}
}

func TestConnectTimeout(t *testing.T) {
	cfg := NewConfig("192.0.2.1", "12345", 100*time.Millisecond)
	client := NewClient(cfg)

	start := time.Now()
	err := client.Connect()
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.Less(t, elapsed, 2*time.Second, "timeout not respected")
}

func TestConnectRefused(t *testing.T) {
	cfg := NewConfig("127.0.0.1", "1", 2*time.Second)
	client := NewClient(cfg)

	err := client.Connect()
	require.Error(t, err)
}

func TestRunNotConnected(t *testing.T) {
	cfg := NewConfig("127.0.0.1", "1234", time.Second)
	client := NewClient(cfg)

	err := client.Run(context.Background(), strings.NewReader(""), &bytes.Buffer{})
	require.ErrorIs(t, err, ErrNoConnected)
}

func TestRunEcho(t *testing.T) {
	client, cleanup := connectToEcho(t)
	defer cleanup()

	input := "hello\nworld\n"
	var out bytes.Buffer

	err := client.Run(context.Background(), strings.NewReader(input), &out)
	require.NoError(t, err)
	assert.Equal(t, input, out.String())
}

func TestRunContextCancel(t *testing.T) {
	client, cleanup := connectToEcho(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())

	pr, pw := io.Pipe()
	defer pw.Close()
	defer pr.Close()

	done := make(chan error, 1)
	go func() {
		done <- client.Run(ctx, pr, &bytes.Buffer{})
	}()

	cancel()

	select {
	case err := <-done:
		require.Error(t, err, "Run should return error on context cancel")
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not finish after context cancel")
	}
}

func TestRunServerCloses(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()

	serverMsg := "goodbye\n"

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		conn.Write([]byte(serverMsg))
		conn.Close()
	}()

	host, port, _ := net.SplitHostPort(ln.Addr().String())
	cfg := NewConfig(host, port, 5*time.Second)
	client := NewClient(cfg)

	require.NoError(t, client.Connect())
	defer client.Close()

	pr, pw := io.Pipe()
	defer pw.Close()
	defer pr.Close()

	var out bytes.Buffer

	done := make(chan error, 1)
	go func() {
		done <- client.Run(context.Background(), pr, &out)
	}()

	select {
	case err := <-done:
		require.NoError(t, err)
		assert.Equal(t, serverMsg, out.String())
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not finish after server closed connection")
	}
}
