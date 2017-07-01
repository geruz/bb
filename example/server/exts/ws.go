package exts

import (
	"bufio"
	"errors"
	"github.com/valyala/fasthttp"
	"net"
	"strings"
	"time"
)

type HandshakeError struct {
	message string
}

func (e HandshakeError) Error() string { return e.message }

// Upgrader specifies parameters for upgrading an HTTP connection to a
// WebSocket connection.
type Upgrader struct {
	HandshakeTimeout                time.Duration
	ReadBufferSize, WriteBufferSize int
	Subprotocols                    []string
}

func printErr(ctx fasthttp.RequestCtx, status int, reason string) {
	ctx.SetBody([]byte(reason))
	ctx.SetStatusCode(status)
}

func (u *Upgrader) selectSubprotocol(ctx fasthttp.RequestCtx) string {
	if u.Subprotocols != nil {
		clientProtocols := Subprotocols(r)
		for _, serverProtocol := range u.Subprotocols {
			for _, clientProtocol := range clientProtocols {
				if clientProtocol == serverProtocol {
					return clientProtocol
				}
			}
		}
	} else if responseHeader != nil {
		return responseHeader.Get("Sec-Websocket-Protocol")
	}
	return ""
}

func (u *Upgrader) Upgrade(ctx fasthttp.RequestCtx) *Conn {
	if !ctx.IsGet() {
		printErr(ctx, fasthttp.StatusMethodNotAllowed, "websocket: not a websocket handshake: request method is not GET")
		return nil
	}
	headers := ctx.Request.Header
	ext := headers.Peek("Sec-Websocket-Extensions")
	if ext != nil {
		printErr(ctx, fasthttp.StatusInternalServerError, "websocket: application specific 'Sec-Websocket-Extensions' headers are unsupported")
		return nil
	}

	if !tokenListContainsValue(headers, "Connection", "upgrade") {
		printErr(ctx, fasthttp.StatusBadRequest, "websocket: not a websocket handshake: 'upgrade' token not found in 'Connection' header")
		return nil
	}

	if !tokenListContainsValue(headers, "Upgrade", "websocket") {
		printErr(ctx, fasthttp.StatusBadRequest, "websocket: not a websocket handshake: 'websocket' token not found in 'Upgrade' header")
		return nil
	}

	if !tokenListContainsValue(headers, "Sec-Websocket-Version", "13") {
		printErr(ctx, fasthttp.StatusBadRequest, "websocket: unsupported version: 13 not found in 'Sec-Websocket-Version' header")
		return nil
	}

	challengeKey := headers.Peek("Sec-Websocket-Key")
	if challengeKey == nil {
		printErr(ctx, fasthttp.StatusBadRequest, "websocket: not a websocket handshake: `Sec-Websocket-Key' header is missing or blank")
		return nil
	}

	subprotocol := u.selectSubprotocol(r, responseHeader)

	// Negotiate PMCE
	var compress bool
	if u.EnableCompression {
		for _, ext := range parseExtensions(r.Header) {
			if ext[""] != "permessage-deflate" {
				continue
			}
			compress = true
			break
		}
	}

	var (
		netConn net.Conn
		err     error
	)

	h, ok := w.(http.Hijacker)
	if !ok {
		return u.returnError(w, r, http.StatusInternalServerError, "websocket: response does not implement http.Hijacker")
	}
	var brw *bufio.ReadWriter
	netConn, brw, err = h.Hijack()
	if err != nil {
		return u.returnError(w, r, http.StatusInternalServerError, err.Error())
	}

	if brw.Reader.Buffered() > 0 {
		netConn.Close()
		return nil, errors.New("websocket: client sent data before handshake is complete")
	}

	c := newConnBRW(netConn, true, u.ReadBufferSize, u.WriteBufferSize, brw)
	c.subprotocol = subprotocol

	if compress {
		c.newCompressionWriter = compressNoContextTakeover
		c.newDecompressionReader = decompressNoContextTakeover
	}

	p := c.writeBuf[:0]
	p = append(p, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
	p = append(p, computeAcceptKey(challengeKey)...)
	p = append(p, "\r\n"...)
	if c.subprotocol != "" {
		p = append(p, "Sec-Websocket-Protocol: "...)
		p = append(p, c.subprotocol...)
		p = append(p, "\r\n"...)
	}
	if compress {
		p = append(p, "Sec-Websocket-Extensions: permessage-deflate; server_no_context_takeover; client_no_context_takeover\r\n"...)
	}
	for k, vs := range responseHeader {
		if k == "Sec-Websocket-Protocol" {
			continue
		}
		for _, v := range vs {
			p = append(p, k...)
			p = append(p, ": "...)
			for i := 0; i < len(v); i++ {
				b := v[i]
				if b <= 31 {
					// prevent response splitting.
					b = ' '
				}
				p = append(p, b)
			}
			p = append(p, "\r\n"...)
		}
	}
	p = append(p, "\r\n"...)

	// Clear deadlines set by HTTP server.
	netConn.SetDeadline(time.Time{})

	if u.HandshakeTimeout > 0 {
		netConn.SetWriteDeadline(time.Now().Add(u.HandshakeTimeout))
	}
	if _, err = netConn.Write(p); err != nil {
		netConn.Close()
		return nil, err
	}
	if u.HandshakeTimeout > 0 {
		netConn.SetWriteDeadline(time.Time{})
	}

	return c, nil
}

// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
//
// This function is deprecated, use websocket.Upgrader instead.
//
// The application is responsible for checking the request origin before
// calling Upgrade. An example implementation of the same origin policy is:
//
//	if req.Header.Get("Origin") != "http://"+req.Host {
//		http.Error(w, "Origin not allowed", 403)
//		return
//	}
//
// If the endpoint supports subprotocols, then the application is responsible
// for negotiating the protocol used on the connection. Use the Subprotocols()
// function to get the subprotocols requested by the client. Use the
// Sec-Websocket-Protocol response header to specify the subprotocol selected
// by the application.
//
// The responseHeader is included in the response to the client's upgrade
// request. Use the responseHeader to specify cookies (Set-Cookie) and the
// negotiated subprotocol (Sec-Websocket-Protocol).
//
// The connection buffers IO to the underlying network connection. The
// readBufSize and writeBufSize parameters specify the size of the buffers to
// use. Messages can be larger than the buffers.
//
// If the request is not a valid WebSocket handshake, then Upgrade returns an
// error of type HandshakeError. Applications should handle this error by
// replying to the client with an HTTP error response.
func Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header, readBufSize, writeBufSize int) (*Conn, error) {
	u := Upgrader{ReadBufferSize: readBufSize, WriteBufferSize: writeBufSize}
	return u.Upgrade(w, r, responseHeader)
}

// Subprotocols returns the subprotocols requested by the client in the
// Sec-Websocket-Protocol header.
func Subprotocols(r *http.Request) []string {
	h := strings.TrimSpace(r.Header.Get("Sec-Websocket-Protocol"))
	if h == "" {
		return nil
	}
	protocols := strings.Split(h, ",")
	for i := range protocols {
		protocols[i] = strings.TrimSpace(protocols[i])
	}
	return protocols
}

func tokenListContainsValue(header fasthttp.RequestHeader, name string, value string) bool {
	val := header.Peek(name)
	if val == nil {
		return false
	}
	return string(val) == value
}
