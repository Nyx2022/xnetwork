package xnetwork

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/kataras/golog"
	"github.com/xiecat/xhttp/xtls"
	"github.com/xiecat/xnetwork/input"
	"net"
	"time"
)

type (
	// ResponseMiddleware run after receive response
	ResponseMiddleware func(*Response, *Client) error
)

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Client struct
//_______________________________________________________________________

type Client struct {
	ClientOptions *ClientOptions
	Debug         bool

	afterResponse []ResponseMiddleware
	addr          *input.ServiceAsset
	ctx           context.Context
	dialer        *net.Dialer
	conn          net.Conn
	attempt       int
}

func NewClient() *Client {
	options := GetNetOptions()
	dialer := &net.Dialer{
		Timeout:   time.Duration(options.DialTimeout) * time.Second,
		KeepAlive: time.Duration(options.KeepAlive) * time.Second,
	}
	c := &Client{
		ClientOptions: options,
		dialer:        dialer,
		Debug:         options.Debug,
		conn:          nil,
	}
	c.afterResponse = []ResponseMiddleware{
		responseLogger,
	}
	return c
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Client methods get
//_______________________________________________________________________

func (c *Client) Dial(ctx context.Context, addr *input.ServiceAsset) error {
	return c.dial(ctx, false, addr)
}

func (c *Client) DialTLS(ctx context.Context, addr *input.ServiceAsset) error {
	return c.dial(ctx, true, addr)
}

func (c *Client) GetContext() context.Context {
	return c.ctx
}

func (c *Client) Do(request *Request) (*Response, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("connection not established, try dial before do request")
	}
	golog.Debugf("%s %s", c.addr.Network, c.addr.String())
	var err error
	request.sendAt = time.Now()
	_, err = c.conn.Write(request.Raw)
	if err != nil {
		return nil, err
	}

	readBuf := make([]byte, c.ClientOptions.ReadSize)
	n, err := c.conn.Read(readBuf)
	if err != nil {
		return nil, err
	}
	response := &Response{
		Conn:        c.conn,
		Request:     request,
		RawResponse: readBuf[:n],
		Length:      n,
	}
	response.readAt = time.Now()
	//ResponseMiddleware
	for _, f := range c.afterResponse {
		if err = f(response, c); err != nil {
			return nil, err
		}
	}
	return response, nil
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Client methods set
//_______________________________________________________________________

func (c *Client) AfterResponse(fn ResponseMiddleware) {
	c.afterResponse = append(c.afterResponse, fn)
}

func (c *Client) setAddr(addr *input.ServiceAsset) *Client {
	c.addr = addr
	return c
}

func (c *Client) setContext(ctx context.Context) *Client {
	c.ctx = ctx
	return c
}

func (c *Client) dial(ctx context.Context, isTLS bool, addr *input.ServiceAsset) error {

	var (
		conn    net.Conn
		errDial error
	)

	c.setAddr(addr).setContext(ctx)
	c.attempt = 0

	err := c.ClientOptions.Limiter.Wait(c.GetContext())
	if err != nil {
		return err
	}

	network := addr.GetNetWork()
	if network == "" {
		return fmt.Errorf("target network(tcp/udp) not set")
	}

	addrStr := addr.String()
	// do request with retry
	for i := 0; ; i++ {
		if isTLS {
			tlsConfig, err := xtls.NewTLSConfig(c.ClientOptions.TlsOptions)
			if err != nil {
				return err
			}
			conn, errDial = tls.DialWithDialer(c.dialer, network, addrStr, tlsConfig)
		} else {
			conn, errDial = c.dialer.DialContext(ctx, network, addrStr)
		}
		if errDial == nil {
			break
		}
		remain := c.ClientOptions.FailRetries - i
		if remain <= 0 {
			break
		}
		c.attempt++
		// waitTime todo
		select {
		case <-time.After(time.Duration(100)):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if errDial != nil {
		return fmt.Errorf("giving up connect to %s %s after %d attempt(s): %v",
			network, addrStr, c.attempt, err)
	}

	c.conn = conn

	deadLineTime := time.Now().Add(time.Duration(c.ClientOptions.ReadTimeout) * time.Second)
	err = conn.SetReadDeadline(deadLineTime)

	if err != nil {
		return err
	}

	return nil
}
