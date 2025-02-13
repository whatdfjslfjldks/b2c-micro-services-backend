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

//type ProductDetail struct {
//	ProductId    int32
//	ProductName  string
//	ProductImg   []*pb.ProductImg
//	ProductPrice float64
//	ProductType  []*pb.ProductType
//	Sold         int32
//}

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
	SecImg           []*pb.SecImg
	Time             int32 // 场次
}

// Products 顶层products表
type Products struct {
	ID          int32
	Name        string
	Kind        int32
	Category    int32
	Description string
	CreateTime  string
}

type NormalProducts struct {
	ID          int32
	Name        string
	Description string
	Kind        int32
	Category    int32
	ImageList   []*pb.PImg
	Price       float64
	OriginPrice float64
	Stock       int32
	Sold        int32
	CreateTime  string
}

type Test struct {
	PPP []NormalProducts
}

type SecProducts struct {
	ImageList   []*pb.PImg
	Price       float64
	OriginPrice float64
	Stock       int32
	Sold        int32
	StartTime   string
	Duration    string
	SessionId   int32
}

//type GetProductById struct {
//	ID          int32
//	Name        string
//	Price       float64
//	OriginPrice float64
//	CategoryID  int32
//	KindID      int32
//	Description string
//	Sold        int32
//	Stock       int32
//	StartTime   string
//	Duration    string
//	SessionID   int32
//	ImageList   []*pb.ProductImg
//	TypeList    []*pb.ProductType
//}
