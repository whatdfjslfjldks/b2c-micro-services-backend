package handler

import (
	pb "micro-services/pkg/proto/pay-server"
)

type Server struct {
	pb.UnimplementedPayServiceServer
}
