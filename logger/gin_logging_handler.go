package logger

import (
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/peer"
)

// Gin Logging handler
func GinLoggingHandler(log Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process the next
		c.Next()

		service := c.Request.RequestURI
		method := c.Request.Method
		duration := durationToMilliseconds(time.Since(start))
		code := c.Writer.Status()
		peer_ip, peer_port, scheme := getRemoteAddressFromGinContext(c)

		correlationId := c.Writer.Header().Get("x-correlation-id")

		if c.Writer.Status() >= 500 {
			log.Info("failed", "x-correlation-id", correlationId, "duration_ms", duration, "code", code, "service", service, "method", method, "raddr", (peer_ip + ":" + peer_port), "scheme", scheme, "error", c.Errors.String())
		} else {
			log.Info("success", "x-correlation-id", correlationId, "duration_ms", duration, "code", code, "service", service, "method", method, "raddr", (peer_ip + ":" + peer_port), "scheme", scheme)
		}

	}
}

func getRemoteAddressFromGinContext(c *gin.Context) (ip, port, netType string) {
	ip, port = getRemoteAddressSet(c)
	// no ip and port were passed through gateway

	if len(ip) < 1 {
		ip, port, netType = "0.0.0.0", "0", ""
		if peer, ok := peer.FromContext(c.Request.Context()); ok {
			netType = peer.Addr.Network()
			// Here is the tricky part
			// We only try to parse IPV4 style Address
			// Rest of peer.Addr implementations are not well formatted string
			// and in this case, we leave port as zero and IP as the returned
			// String from Addr.String() function
			//
			// BTW, just skip the error since it would not impact anything
			// Operators could observe this error from monitor dashboards by
			// validating existence of IP & PORT fields
			ip, port, _ = net.SplitHostPort(peer.Addr.String())
		}

		forwardedRemoteIP := c.Request.Header.Get("x-forwarded-for")

		// Deal with forwarded remote ip
		if len(forwardedRemoteIP) > 0 {
			if forwardedRemoteIP == "::1" {
				forwardedRemoteIP = "localhost"
			}

			ip = forwardedRemoteIP
		}

		if ip == "::1" {
			ip = "localhost"
		}
	}

	return ip, port, netType
}

func getRemoteAddressSet(c *gin.Context) (ip, port string) {
	if v := c.Request.Header.Get("x-forwarded-remote-addr"); len(v) > 0 {
		ip, port, _ = net.SplitHostPort(v)
	}

	if ip == "::1" {
		ip = "localhost"
	}

	return ip, port
}
