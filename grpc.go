package panicrecovery

import (
	"context"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Context key
const PanicDetailsKey = "panicDetails"

// RecoverInterceptor is a gRPC interceptor that recovers from panics in gRPC methods.
func RecoverInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		defer func() {
			ctx = handlegRPCPanic(ctx, info.FullMethod)
		}()
		// Call the handler to execute the RPC method and pass ctx to downstream
		return handler(ctx, req)
	}
}

// RecoverStreamInterceptor is a gRPC interceptor that recovers from panics in gRPC streaming methods.
func RecoverStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		defer func() {
			_ = handlegRPCPanic(stream.Context(), info.FullMethod)
		}()
		// Call the handler to execute the streaming RPC method
		return handler(srv, stream)
	}
}

func handlegRPCPanic(ctx context.Context, method string) context.Context {
	if err := recover(); err != nil {
		Logger.Error("Recovered from panic in gRPC",
			zap.Any("error", err),
			zap.String("method", method),
			zap.String("stack_trace", string(debug.Stack())),
			zap.Time("timestamp", time.Now()))

		// Atach recovery info to ctx
		return context.WithValue(ctx, PanicDetailsKey, err)
	}
	return ctx
}
