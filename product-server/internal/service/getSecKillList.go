package service

import (
	pb "micro-services/pkg/proto/product-server"
	"micro-services/product-server/internal/repository"
)

// GetSecListAndTotalItems 获取秒杀商品列表和总数
func GetSecListAndTotalItems(currentPage int32, pageSize int32, time int32) (
	[]*pb.SecListItem, int32, error) {

	// 获取秒杀商品列表
	products, err := repository.GetSecList(currentPage, pageSize, time)
	if err != nil {
		return nil, 0, err
	}

	return products, 0, nil
}
