package handler

import (
	"context"
	"fmt"
	"log"
	productServerProto "micro-services/pkg/proto/product-server"
	pb "micro-services/pkg/proto/recommend-server"
	"micro-services/recommend-server/internal/service/recommend"
	"micro-services/recommend-server/pkg/instance"
	"strconv"
	"strings"
)

func (s *Server) GetRecommendProductList(ctx context.Context, req *pb.GetRecommendProductListRequest) (
	*pb.GetRecommendProductListResponse, error) {
	resp := &pb.GetRecommendProductListResponse{}

	// 获取相似用户表
	simUserId, err := recommend.GetSimUserId(req.UserId)
	if err != nil {
		resp.Code = 500
		resp.StatusCode = "GLB-003"
		resp.Msg = err.Error()
		return resp, nil
	}

	var totalId []string
	for _, val := range simUserId {
		// 转换为整数类型
		v, _ := strconv.Atoi(val)
		//fmt.Println("v: ", v)

		// 获取相似的产品ID
		simProductId, err := recommend.GetSimProductId(int64(v))
		if err != nil {
			return nil, err
		}

		// 将返回的产品ID添加到totalId，并去重
		for _, pid := range simProductId {
			if !contains(totalId, pid) {
				totalId = append(totalId, pid)
			}
		}
	}

	var products []*pb.ProductListItem2
	for _, item := range totalId {
		productId := strings.TrimPrefix(item, "product_")
		//fmt.Println("Extracted product ID:", productId)
		pId, err := strconv.Atoi(productId)
		if err != nil {
			log.Printf("Error converting product ID to integer: %v", err)
			continue
		}
		//fmt.Println("productId: ", pId)
		product, err := instance.GrpcClient.GetProductById(&productServerProto.GetProductByIdRequest{
			ProductId: int32(pId),
		})
		if err != nil {
			log.Printf("Error getting product by ID: %v", err)
			continue
		}
		fmt.Println("product: ", product)
		products = append(products, &pb.ProductListItem2{
			ProductId:         int32(pId),
			ProductName:       product.ProductName,
			ProductCover:      product.ProductCover,
			ProductPrice:      product.ProductPrice,
			ProductCategoryId: product.ProductCategoryId,
			Description:       product.Description,
		})
	}
	//fmt.Println("products: ", products)
	resp.Code = 200
	resp.Msg = "获取推荐商品成功！"
	resp.StatusCode = "GLB-000"
	resp.ProductList = products
	return resp, nil
}

// 判断切片中是否已经存在某个元素
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
