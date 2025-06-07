package interceptors

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// NewUnaryServerLoggingInterceptor создает серверный интерсептор для логирования unary RPC.
func NewUnaryServerLoggingInterceptor(logger *log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now().UTC()
		ts := start.Format(server.LogTimestampFormat)

		var ip string
		if p, ok := peer.FromContext(ctx); ok {
			ip = p.Addr.String()
			if idx := strings.LastIndex(ip, ":"); idx != -1 {
				ip = ip[:idx]
			}
		} else {
			ip = "unknown"
		}

		var userAgent string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if ua := md["user-agent"]; len(ua) > 0 {
				userAgent = ua[0]
			}
		}

		resp, err := handler(ctx, req)

		status := "OK"
		if err != nil {
			status = fmt.Sprintf("ERROR: %v", err)
		}

		latency := time.Since(start).Milliseconds()

		logLine := fmt.Sprintf("%s [%s] %s %q %d %q",
			ip,
			ts,
			info.FullMethod,
			status,
			latency,
			userAgent,
		)

		logger.Println(logLine)

		return resp, err
	}
}
