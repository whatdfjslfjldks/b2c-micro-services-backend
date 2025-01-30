package dto

import pb "micro-services/pkg/proto/product-server"

// Product 一键导入商品
type Product struct {
	Name        string
	Category    int32
	Description string
	Price       float64
	Stock       int32
}

// ProductById 商品详情
type ProductById struct {
	Name        string
	Cover       string
	Category    int32
	Description string
	Price       float64
}

type ProductDetail struct {
	ProductId    int32
	ProductName  string
	ProductImg   []*pb.ProductImg
	ProductPrice float64
	ProductType  []*pb.ProductType
	Sold         int32
}

type img struct {
	ImgId  int32
	ImgUrl string
}
type productType struct {
	TypeId   int32
	TypeName string
}

// SecKillProduct 秒杀商品
type SecKillProduct struct {
	SecName          string
	SecDescription   string
	SecPrice         float64
	SecOriginalPrice float64
	Stock            int32
	StartTime        string
	EndTime          string
	SecType          []*pb.SecType
	SecImg           []*pb.SecImg
	Time             int32 // 场次
}
