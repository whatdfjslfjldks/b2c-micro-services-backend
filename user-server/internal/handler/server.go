package handler

import pb "micro-services/pkg/proto/user-server"

type Server struct {
	pb.UnimplementedUserServiceServer
}
