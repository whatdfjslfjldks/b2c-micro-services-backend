package repository

import (
	"database/sql"
	"errors"
	"log"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/pkg/utils"
	"micro-services/product-server/pkg/config"
)

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

//type img struct {
//	ImgId  int32
//	ImgUrl string
//}
//type productType struct {
//	TypeId   int32
//	TypeName string
//}

func SaveProductsBatch(products []Product) error {
	// 生成批量插入的SQL语句
	sql := "INSERT INTO b2c_product.products (name, category_id, description, prices, stock, create_at, update_at) VALUES "
	var args []interface{}

	// 构建SQL语句和参数
	for _, p := range products {
		sql += "(?, ?, ?, ?, ?, ?, ?),"
		args = append(args, p.Name, p.Category, p.Description, p.Price, p.Stock, utils.GetTime(), utils.GetTime())
	}

	// 去掉最后一个逗号
	sql = sql[:len(sql)-1]

	// 执行批量插入
	_, err := config.MySqlClient.Exec(sql, args...)
	if err != nil {
		//fmt.Printf("could not insert to products: %v", err)
		return errors.New("GLB-003")
	}
	return nil
}

func GetProductById(id int32) (ProductById, error) {
	var product ProductById
	var cover sql.NullString // 使用 sql.NullString 来处理 cover_url 字段

	// 查询数据库
	err := config.MySqlClient.QueryRow("SELECT name, category_id, description, prices, cover_url FROM b2c_product.products WHERE product_id = ?", id).Scan(&product.Name, &product.Category, &product.Description, &product.Price, &cover)
	if err != nil {
		log.Printf("Failed to get product by ID-%d: %v", id, err)
		return product, errors.New("GLB-003")
	}

	// 如果 cover 是 NULL，将其设为空字符串
	if cover.Valid {
		product.Cover = cover.String
	} else {
		product.Cover = "" // NULL 处理为空字符串
	}

	return product, nil
}

// GetProductDetailById 获取详情页商品数据
func GetProductDetailById(id int32) (ProductDetail, error) {
	var product ProductDetail
	// 获取商品基本信息
	err := config.MySqlClient.QueryRow(
		"SELECT name, prices FROM b2c_product.products WHERE product_id = ?", id).Scan(&product.ProductName, &product.ProductPrice)
	if err != nil {
		log.Printf("Failed to get product by ID-%d: %v", id, err)
		return product, errors.New("GLB-003")
	}

	// 获取商品类型
	row, err := config.MySqlClient.Query("SELECT type_description, id FROM b2c_product.product_type WHERE product_id=?", id)
	if err != nil {
		log.Printf("Failed to get product type by ID-%d: %v", id, err)
		return product, errors.New("GLB-003")
	}
	defer row.Close()

	// 使用普通结构体，而不是指针
	for row.Next() {
		var pType pb.ProductType // 直接使用结构体，而不是指针
		err = row.Scan(&pType.ProductTypeName, &pType.ProductTypeId)
		if err != nil {
			log.Printf("Failed to get product type by ID-%d: %v", id, err)
			return product, errors.New("GLB-003")
		}
		product.ProductType = append(product.ProductType, &pType)
	}

	// 获取商品图片
	rowImg, err := config.MySqlClient.Query("SELECT img, id FROM b2c_product.product_image WHERE product_id=?", id)
	if err != nil {
		log.Printf("Failed to get product images by ID-%d: %v", id, err)
		return product, errors.New("GLB-003")
	}
	defer rowImg.Close()

	// 使用普通结构体，而不是指针
	for rowImg.Next() {
		var img pb.ProductImg // 直接使用结构体，而不是指针
		err = rowImg.Scan(&img.ImgUrl, &img.ImgId)
		if err != nil {
			log.Printf("Failed to get product image by ID-%d: %v", id, err)
			return product, errors.New("GLB-003")
		}
		product.ProductImg = append(product.ProductImg, &img)
	}

	return product, nil
}
