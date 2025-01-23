package panicrecovery

import (
	"context"
	"testing"

	"google.golang.org/grpc"
)

type mockHandler struct {
	shouldPanic bool
}

type mockServerStream struct {
	ctx context.Context
	grpc.ServerStream
}

func (m *mockServerStream) Context() context.Context {
	return m.ctx
}

func (m *mockHandler) Handle(ctx context.Context, req any) (any, error) {
	if m.shouldPanic {
		panic("test panic")
	}

	return "Hale Haji", nil
}

func TestRecoverInterceptor(t *testing.T) {
	tests := []struct {
		name        string
		shouldPanic bool
	}{
		{
			name:        "no panic",
			shouldPanic: false,
		},
		{
			name:        "with panic",
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := RecoverInterceptor()
			handler := &mockHandler{shouldPanic: tt.shouldPanic}

			ctx := context.Background()
			info := &grpc.UnaryServerInfo{FullMethod: "TestMethod"}

			_, err := interceptor(ctx, "test request", info, handler.Handle)

			if tt.shouldPanic && err == nil {
				t.Error("Expected error from panic recovery, got nil")
			}

			if !tt.shouldPanic && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestRecoverStreamInterceptor(t *testing.T) {
	tests := []struct {
		name        string
		shouldPanic bool
	}{
		{
			name:        "no panic",
			shouldPanic: false,
		},
		{
			name:        "with panic",
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := RecoverStreamInterceptor()
			stream := &mockServerStream{ctx: context.Background()}
			info := &grpc.StreamServerInfo{FullMethod: "TestStreamMethod"}

			handler := func(srv any, stream grpc.ServerStream) error {
				if tt.shouldPanic {
					panic("test panic")
				}
				return nil
			}

			err := interceptor(nil, stream, info, handler)

			if tt.shouldPanic && err == nil {
				t.Error("Expected error from panic recovery, got nil")
			}

			if !tt.shouldPanic && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}
