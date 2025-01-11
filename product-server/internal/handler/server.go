package handler

import pb "micro-services/pkg/proto/product-server"

type Server struct {
	pb.UnimplementedProductServiceServer
}
