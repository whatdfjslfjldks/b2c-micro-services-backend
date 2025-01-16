package repository

import (
	"database/sql"
	"errors"
	"log"
	"micro-services/pkg/utils"
	"micro-services/product-server/pkg/config"
)

type Product struct {
	Name        string
	Category    int32
	Description string
	Price       float64
	Stock       int32
}

type ProductById struct {
	Name        string
	Cover       string
	Category    int32
	Description string
	Price       float64
}

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
