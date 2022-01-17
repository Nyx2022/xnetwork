package xnetwork

import (
	"context"
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"github.com/xiecat/xnetwork/input"
	"github.com/xiecat/xnetwork/testutils/tcp"
	"log"
	"net"
	"strconv"
	"strings"
	"testing"
)

func TestClient_SendBytes(t *testing.T) {

	ts := tcp.NewTCPServer(func(conn net.Conn) {
		defer conn.Close()

		data := make([]byte, 4)
		if _, err := conn.Read(data); err != nil {
			return
		}
		if string(data) == "PING" {
			_, _ = conn.Write([]byte("PONG"))
		}
	})
	defer ts.Close()
	tmp := strings.SplitN(ts.URL, ":", 2)

	client := NewClient()

	ctx := context.Background()
	port, _ := strconv.ParseInt(tmp[1], 10, 64)
	addr := &input.ServiceAsset{
		Host:    tmp[0],
		Port:    int(port),
		Network: "tcp",
	}
	err := client.Dial(ctx, addr)
	require.Nil(t, err, "could not  dial network")

	req1 := &Request{
		Raw: []byte("PING"),
	}
	resp1, err := client.Do(req1)
	require.Nil(t, err, "could not do network request")
	require.Equal(t, "PONG", string(resp1.RawResponse), "could not get correct network response")
}

func TestClient_SendHex(t *testing.T) {

	ts := tcp.NewTCPServer(func(conn net.Conn) {
		defer conn.Close()
		data := make([]byte, 20)
		if _, err := conn.Read(data); err != nil {
			return
		}
		_, _ = conn.Write([]byte("pong"))
	})
	defer ts.Close()

	tmp := strings.SplitN(ts.URL, ":", 2)

	options := DefaultClientOptions()
	NetOptions = options

	client := NewClient()

	ctx := context.Background()
	port, _ := strconv.ParseInt(tmp[1], 10, 64)
	addr := &input.ServiceAsset{
		Host:    tmp[0],
		Port:    int(port),
		Network: "tcp",
	}

	err := client.Dial(ctx, addr)
	require.Nil(t, err, "could not  dial network")

	h, _ := hex.DecodeString("11AACC")
	req := &Request{
		Raw: h,
	}

	resp, err := client.Do(req)
	require.Nil(t, err, "could not do network request")
	require.Equal(t, "pong", string(resp.RawResponse), "could not get correct network response")
}

func TestClient_SendMultiSteps(t *testing.T) {
	var routerErr error
	ts := tcp.NewTCPServer(func(conn net.Conn) {
		defer conn.Close()

		data := make([]byte, 5)
		if _, err := conn.Read(data); err != nil {
			routerErr = err
			return
		}
		if string(data) == "FIRST" {
			_, _ = conn.Write([]byte("PING"))
		}

		data = make([]byte, 6)
		if _, err := conn.Read(data); err != nil {
			routerErr = err
			return
		}
		if string(data) == "SECOND" {
			_, _ = conn.Write([]byte("PONG"))
		}
	})
	defer ts.Close()

	if routerErr != nil {
		log.Println(routerErr)
		return
	}

	tmp := strings.SplitN(ts.URL, ":", 2)

	client := NewClient()

	ctx := context.Background()
	port, _ := strconv.ParseInt(tmp[1], 10, 64)
	addr := &input.ServiceAsset{
		Host:    tmp[0],
		Port:    int(port),
		Network: "tcp",
	}

	err := client.Dial(ctx, addr)
	require.Nil(t, err, "could not  dial network")

	req1 := &Request{
		Raw: []byte("FIRST"),
	}
	resp1, err := client.Do(req1)
	require.Nil(t, err, "could not do network request")
	require.Equal(t, "PING", string(resp1.RawResponse), "could not get correct network response")
	req2 := &Request{
		Raw: []byte("SECOND"),
	}
	resp2, err := client.Do(req2)
	require.Nil(t, err, "could not do network request")
	require.Equal(t, "PONG", string(resp2.RawResponse), "could not get correct network response")

}

