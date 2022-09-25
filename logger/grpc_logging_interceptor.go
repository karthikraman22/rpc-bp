package logger

import (
	"context"
	"net"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func durationToMilliseconds(duration time.Duration) float32 {
	return float32(duration.Nanoseconds()/1000) / 1000
}

type contextKey string

const (
	CORRELATION_ID contextKey = "x-correlation-id"
)

func (c contextKey) String() string {
	return string(c)
}

// Unary server interceptor
func UnaryServerInterceptor(log Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		md, _ := metadata.FromIncomingContext(ctx)
		peer_ip, peer_port, scheme := getRemoteAddressFromMetaData(md, ctx)
		service := path.Dir(info.FullMethod)[1:]
		method := path.Base(info.FullMethod)

		correlationId := md[CORRELATION_ID.String()]

		newCtx := context.WithValue(ctx, CORRELATION_ID, correlationId)

		// Calls the handler
		resp, err := handler(newCtx, req)

		duration := durationToMilliseconds(time.Since(start))
		code := status.Code(err)

		if err != nil {
			log.Info("failed", CORRELATION_ID.String(), correlationId, "meta", md, "duration_ms", duration, "code", code, "service", service, "method", method, "raddr", (peer_ip + ":" + peer_port), "scheme", scheme, "error", err)
		} else {
			log.Info("success", CORRELATION_ID.String(), correlationId, "meta", md, "duration_ms", duration, "code", code, "service", service, "method", method, "raddr", (peer_ip + ":" + peer_port), "scheme", scheme)
		}

		return resp, err
	}
}

func getRemoteAddressFromMetaData(md metadata.MD, ctx context.Context) (ip, port, netType string) {
	ip, port = getRemoteAddressSetFromMeta(md)
	// no ip and port were passed through gateway

	if len(ip) < 1 {
		ip, port, netType = "0.0.0.0", "0", ""
		if peer, ok := peer.FromContext(ctx); ok {
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

		forwardedRemoteIPList := md["x-forwarded-for"]

		// Deal with forwarded remote ip
		if len(forwardedRemoteIPList) > 0 {
			forwardedRemoteIP := forwardedRemoteIPList[0]
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

func getRemoteAddressSetFromMeta(md metadata.MD) (ip, port string) {
	if v := md.Get("x-forwarded-remote-addr"); len(v) > 0 {
		ip, port, _ = net.SplitHostPort(v[0])
	}

	if ip == "::1" {
		ip = "localhost"
	}

	return ip, port
}

// GetGwInfo Extract gateway related information from metadata.
func getGwInfo(md metadata.MD) (gwMethod, gwPath, gwScheme, gwUserAgent string) {
	gwMethod, gwPath, gwScheme, gwUserAgent = "", "", "", ""

	if tokens := md["x-forwarded-method"]; len(tokens) > 0 {
		gwMethod = tokens[0]
	}

	if tokens := md["x-forwarded-path"]; len(tokens) > 0 {
		gwPath = tokens[0]
	}

	if tokens := md["x-forwarded-scheme"]; len(tokens) > 0 {
		gwScheme = tokens[0]
	}

	if tokens := md["x-forwarded-user-agent"]; len(tokens) > 0 {
		gwUserAgent = tokens[0]
	}

	return gwMethod, gwPath, gwScheme, gwUserAgent
}
