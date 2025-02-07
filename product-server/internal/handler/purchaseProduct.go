package handler

import (
	"context"
	"fmt"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/product-server/internal/repository"
)

func (s *Server) PurchaseSecKill(ctx context.Context, req *pb.PurchaseSecKillRequest) (
	*pb.PurchaseSecKillResponse, error) {

	for {
		a := repository.Test()
		fmt.Println("111: ", a)
	}

	fmt.Println("purchaseSecKill: ", req.Id, req.AccessToken)
	return nil, nil
}
