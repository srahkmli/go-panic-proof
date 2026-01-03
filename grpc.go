package panicrecovery

import (
	"context"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Context key
const PanicDetailsKey = "panicDetails"

// RecoverInterceptor is a gRPC interceptor that recovers from panics in gRPC methods.
func RecoverInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				ctx, err = handlegRPCPanic(ctx, r, info.FullMethod)
			}
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
	) (err error) {
		defer func() {
			if r := recover(); r != nil {
				_, err = handlegRPCPanic(stream.Context(), r, info.FullMethod)
			}
		}()
		// Call the handler to execute the streaming RPC method
		return handler(srv, stream)
	}
}

func handlegRPCPanic(ctx context.Context, r any, method string) (context.Context, error) {
	Logger.Error("Recovered from panic in gRPC",
		zap.Any("error", r),
		zap.String("method", method),
		zap.String("stack_trace", string(debug.Stack())),
		zap.Time("timestamp", time.Now()))

	// Atach recovery info to ctx
	ctx = context.WithValue(ctx, PanicDetailsKey, r)

	return ctx, status.Errorf(codes.Internal, "panic occurred: %v", r)
}
