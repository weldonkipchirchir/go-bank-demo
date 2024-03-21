package grpc_api

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	user_agent    = "grpcgateway-user-agent"
	grpcUserAgent = "user-agent"
	client_ip     = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Fatal("failed to get metadata")
	}
	if userAgent := md.Get(grpcUserAgent); len(userAgent) > 0 {
		mtdt.UserAgent = userAgent[0]
	}
	if clientIP := md.Get(client_ip); len(clientIP) > 0 {
		mtdt.ClientIp = clientIP[0]
	}

	if peer, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIp = peer.Addr.String()
	}
	return mtdt
}
