package remotedialer

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// ConnectAuthorizer custom for authorization
type ConnectAuthorizer func(proto, address string) bool

// ClientConnect connect to WS and wait 5 seconds when error
func ClientConnect(ctx context.Context, wsURL string, headers http.Header, dialer *websocket.Dialer,
	auth ConnectAuthorizer, onConnect func(context.Context) error) error {
	if err := ConnectToProxy(ctx, wsURL, headers, auth, dialer, onConnect); err != nil {
		GetLogger().Errorf("Remotedialer proxy error, %s", err)
		time.Sleep(time.Duration(5) * time.Second)
		return err
	}
	return nil
}

// ConnectToProxy connect to websocket server
func ConnectToProxy(rootCtx context.Context, proxyURL string, headers http.Header, auth ConnectAuthorizer, dialer *websocket.Dialer, onConnect func(context.Context) error) error {
	GetLogger().Infof("Connecting to proxy, url: %s", proxyURL)

	if dialer == nil {
		dialer = &websocket.Dialer{Proxy: http.ProxyFromEnvironment, HandshakeTimeout: HandshakeTimeOut}
	}
	ws, resp, err := dialer.Dial(proxyURL, headers)
	if err != nil {
		if resp == nil {
			GetLogger().Errorf("Failed to connect to proxy. Empty dialer response, err %s", err)
		} else {
			rb, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				GetLogger().Errorf("Failed to connect to proxy, %s. Response status: %v - %v. Couldn't read response body (err: %v)", err, resp.StatusCode, resp.Status, err2)
			} else {
				GetLogger().Errorf("Failed to connect to proxy, %s. Response status: %v - %v. Response body: %s", err, resp.StatusCode, resp.Status, rb)
			}
		}
		return err
	}
	defer ws.Close()

	result := make(chan error, 2)

	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	if onConnect != nil {
		go func() {
			if err := onConnect(ctx); err != nil {
				result <- err
			}
		}()
	}

	session := NewClientSession(auth, ws)
	defer session.Close()

	go func() {
		_, err = session.Serve(ctx)
		result <- err
	}()

	select {
	case <-ctx.Done():
		GetLogger().Infof("Proxy done, url: %s, err %s", proxyURL, ctx.Err())
		return nil
	case err := <-result:
		return err
	}
}
