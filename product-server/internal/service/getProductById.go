package service

import (
	"errors"
	"fmt"
	"log"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/product-server/internal/repository"
)

func GetProductById(id int32) (
	*pb.ProductById, error) {
	// 查顶层products表，并判断商品kind
	products, err := repository.GetInfoFromProducts(id)
	//fmt.Println("products: ", products)
	if err != nil {
		log.Printf("GetInfoFromProducts err: %v", err)
		return nil, err
	}
	switch products.Kind {
	case 1:
		// 普通商品 四表联查
		productDetail, err := repository.GetProductDetail(id)
		if err != nil {
			//log.Println("GetProductDetail error:", err)
			return nil, err
		}
		return &pb.ProductById{
			Id:            products.ID,
			Name:          products.Name,
			Price:         productDetail.Price,
			OriginalPrice: productDetail.OriginPrice,
			CategoryId:    products.Category,
			KindId:        products.Kind,
			Description:   products.Description,
			Sold:          productDetail.Sold,
			Stock:         productDetail.Stock,
			PImg:          productDetail.ImageList,
			PType:         productDetail.TypeList,
		}, nil
	case 2:
		productDetail, err := repository.GetSecProductDetail(id)
		if err != nil {
			log.Println("GetSecProductDetail error:", err)
			return nil, err
		}
		return &pb.ProductById{
			Id:            products.ID,
			Name:          products.Name,
			Price:         productDetail.Price,
			OriginalPrice: productDetail.OriginPrice,
			CategoryId:    products.Category,
			KindId:        products.Kind,
			Description:   products.Description,
			Sold:          productDetail.Sold,
			Stock:         productDetail.Stock,
			PImg:          productDetail.ImageList,
			PType:         productDetail.TypeList,
			StartTime:     productDetail.StartTime,
			Duration:      productDetail.Duration,
			SessionId:     productDetail.SessionId,
		}, nil
	case 3:
		fmt.Println("kind is: ", id)
	default:
		return nil, errors.New("GLB-001")
	}

	return nil, nil
}
