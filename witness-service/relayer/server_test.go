package relayer

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/iotexproject/ioTube/witness-service/grpc/services"
)

func TestPublicRelayHandler(t *testing.T) {
	tests := []struct {
		path        string
		wantStatus  int
		wantHandled bool
	}{
		{path: "/list", wantStatus: http.StatusNoContent, wantHandled: true},
		{path: "/lookup", wantStatus: http.StatusNoContent, wantHandled: true},
		{path: "/submit", wantStatus: http.StatusNotFound},
		{path: "/reset", wantStatus: http.StatusNotFound},
		{path: "/stale_heights", wantStatus: http.StatusNotFound},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			handled := false
			handler := publicRelayHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				handled = true
				w.WriteHeader(http.StatusNoContent)
			}))
			request := httptest.NewRequest(http.MethodPost, test.path, nil)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			require.Equal(t, test.wantStatus, response.Code)
			require.Equal(t, test.wantHandled, handled)
		})
	}
}

func TestAdminUnaryInterceptor(t *testing.T) {
	const token = "test-admin-token"
	interceptor := adminUnaryInterceptor(token)
	info := &grpc.UnaryServerInfo{FullMethod: services.RelayService_Reset_FullMethodName}
	handlerCalled := false
	handler := func(context.Context, interface{}) (interface{}, error) {
		handlerCalled = true
		return "ok", nil
	}

	_, err := interceptor(context.Background(), nil, info, handler)
	require.Equal(t, codes.PermissionDenied, status.Code(err))
	require.False(t, handlerCalled)

	remoteContext := peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{IP: net.ParseIP("192.0.2.10"), Port: 1234},
	})
	remoteContext = metadata.NewIncomingContext(remoteContext, metadata.Pairs("authorization", "Bearer "+token))
	_, err = interceptor(remoteContext, nil, info, handler)
	require.Equal(t, codes.PermissionDenied, status.Code(err))
	require.False(t, handlerCalled)

	localContext := peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1234},
	})
	_, err = interceptor(localContext, nil, info, handler)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
	require.False(t, handlerCalled)

	localContext = metadata.NewIncomingContext(localContext, metadata.Pairs("authorization", "Bearer "+token))
	response, err := interceptor(localContext, nil, info, handler)
	require.NoError(t, err)
	require.Equal(t, "ok", response)
	require.True(t, handlerCalled)
}

func TestAdminUnaryInterceptorAllowsNonAdminRPC(t *testing.T) {
	interceptor := adminUnaryInterceptor("")
	info := &grpc.UnaryServerInfo{FullMethod: services.RelayService_Submit_FullMethodName}
	response, err := interceptor(context.Background(), nil, info, func(context.Context, interface{}) (interface{}, error) {
		return "ok", nil
	})
	require.NoError(t, err)
	require.Equal(t, "ok", response)
}

func TestAdminUnaryInterceptorDisablesAdminRPCWithoutToken(t *testing.T) {
	interceptor := adminUnaryInterceptor("")
	info := &grpc.UnaryServerInfo{FullMethod: services.RelayService_Reset_FullMethodName}
	localContext := peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1234},
	})
	_, err := interceptor(localContext, nil, info, func(context.Context, interface{}) (interface{}, error) {
		return "unexpected", nil
	})
	require.Equal(t, codes.PermissionDenied, status.Code(err))
}
