package panicrecovery

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
)

// RecoverInterceptor is a gRPC interceptor that recovers from panics in gRPC methods.
func RecoverInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic error with stack trace.
				log.Printf("Recovered from panic in gRPC method %s: %v\nStack Trace: %s", info.FullMethod, err, string(debug.Stack()))
			}
		}()
		// Call the handler to execute the RPC method
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
			if err := recover(); err != nil {
				// Log the panic error with stack trace for stream-based methods.
				log.Printf("Recovered from panic in gRPC streaming method %s: %v\nStack Trace: %s", info.FullMethod, err, string(debug.Stack()))
			}
		}()
		// Call the handler to execute the streaming RPC method
		return handler(srv, stream)
	}
}
