package service

import (
	"database/sql"
	"log"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/product-server/pkg/config"
)

// GetProductList 获取商品列表，支持按分类和价格区间过滤
func GetProductList(currentPage int32, pageSize int32, categoryId int32, sort int32) (
	[]*pb.ProductListItem, int32, error) {
	// 构建 SQL 查询语句，根据 categoryId 和 priceRange 判断是否过滤
	log.Printf("sort is %d", sort)
	var query string
	var countQuery string
	var args []interface{}

	// 基本的查询结构
	query = "SELECT product_id, name, description, cover_url, category_id, prices " +
		"FROM b2c_product.products WHERE 1=1"
	// TODO 1=1 避免判断是否为第一个查询，加上1=1可以等加and
	countQuery = "SELECT COUNT(*) FROM b2c_product.products where 1=1"

	// 根据 categoryId 判断是否过滤分类
	if categoryId != 0 {
		query += " AND category_id = ?"
		countQuery += " AND category_id = ?"
		args = append(args, categoryId)
	}

	// 根据 sort 参数构建排序条件
	switch sort {
	case 1:
		// 按价格升序排序
		query += " ORDER BY prices ASC"
	case 2:
		query += " ORDER BY create_at DESC"
	default:
		// 默认不排序（或者按默认排序，比如按 ID 排序）
		query += " ORDER BY product_id ASC"
	}

	// 加入分页查询
	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, (currentPage-1)*pageSize)

	// 执行商品查询
	row, err := config.MySqlClient.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer row.Close()

	// 获取商品总数
	var totalItems int32
	if categoryId == 0 {
		err = config.MySqlClient.QueryRow(countQuery).Scan(&totalItems)
	} else {
		err = config.MySqlClient.QueryRow(countQuery, categoryId).Scan(&totalItems)
	}

	if err != nil {
		return nil, 0, err
	}

	// 处理查询结果
	var productList []*pb.ProductListItem
	for row.Next() {
		var productId int32
		var productName string
		var productDescription string
		var productCover sql.NullString // 使用 NullString 类型来处理可能的 NULL 值
		var productCategoryId int32
		var price float64
		err = row.Scan(&productId, &productName, &productDescription, &productCover, &productCategoryId, &price)
		if err != nil {
			return nil, 0, err
		}

		// 如果 cover_url 是 NULL，则设置为默认值
		var coverURL string
		if productCover.Valid {
			coverURL = productCover.String
		} else {
			coverURL = "" // 使用空字符串作为默认值
		}

		productList = append(productList, &pb.ProductListItem{
			ProductId:         productId,
			ProductName:       productName,
			Description:       productDescription,
			ProductCover:      coverURL,
			ProductCategoryId: productCategoryId,
			ProductPrice:      price,
		})
	}

	//log.Printf("productList is %v", productList)

	return productList, totalItems, nil
}
