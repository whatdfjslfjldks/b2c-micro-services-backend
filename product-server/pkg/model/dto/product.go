package dto

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
