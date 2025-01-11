package repository

import (
	"errors"
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
