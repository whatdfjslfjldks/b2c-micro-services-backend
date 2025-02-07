package handler

import (
	pb "micro-services/pkg/proto/order-server"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
}
