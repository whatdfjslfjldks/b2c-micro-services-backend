package handler

import (
	pb "micro-services/pkg/proto/risk-server"
)

type Server struct {
	pb.UnimplementedRiskServiceServer
}
