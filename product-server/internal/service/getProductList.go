package service

import (
	pb "micro-services/pkg/proto/product-server"
	"micro-services/product-server/pkg/config"
)

// PriceRange 枚举类型，表示价格区间
const (
	AllPriceRange      = 0 // 所有价格
	PriceRange0To50    = 1 // 价格区间 0 - 50
	PriceRange51To100  = 2 // 价格区间 51 - 100
	PriceRange101To200 = 3 // 价格区间 101 - 200
	PriceRange201To500 = 4 // 价格区间 201 - 500
	PriceRangeAbove500 = 5 // 价格区间 500 以上
)

// GetProductList 获取商品列表，支持按分类和价格区间过滤
func GetProductList(currentPage int32, pageSize int32, categoryId int32, priceRange int32) (
	[]*pb.ProductListItem, int32, error) {
	// 构建 SQL 查询语句，根据 categoryId 和 priceRange 判断是否过滤
	var query string
	var countQuery string
	var args []interface{}
	// 价格区间条件
	var priceCondition string
	switch priceRange {
	case PriceRange0To50:
		priceCondition = "prices BETWEEN ? AND ?"
		args = append(args, 0, 50)
	case PriceRange51To100:
		priceCondition = "prices BETWEEN ? AND ?"
		args = append(args, 51, 100)
	case PriceRange101To200:
		priceCondition = "prices BETWEEN ? AND ?"
		args = append(args, 101, 200)
	case PriceRange201To500:
		priceCondition = "prices BETWEEN ? AND ?"
		args = append(args, 201, 500)
	case PriceRangeAbove500:
		priceCondition = "prices > ?"
		args = append(args, 500)
	default:
		priceCondition = "" // 不限制价格
	}

	// TODO: 哈皮报错，不知道为什么，已经加了前置判断
	// 根据 categoryId 和 priceRange 构建查询语句
	if categoryId == 0 {
		// 查询所有商品
		if priceCondition == "" {
			query = "SELECT product_id, name, description, cover_url, category_id, prices " +
				"FROM b2c_product.products " +
				"LIMIT ? OFFSET ?"
			args = append(args, pageSize, (currentPage-1)*pageSize)
			countQuery = "SELECT COUNT(*) FROM b2c_product.products"
		} else {
			query = "SELECT product_id, name, description, cover_url, category_id, prices " +
				"FROM b2c_product.products WHERE " + priceCondition +
				" LIMIT ? OFFSET ?"
			args = append(args, pageSize, (currentPage-1)*pageSize)
			countQuery = "SELECT COUNT(*) FROM b2c_product.products WHERE " + priceCondition
		}
	} else {
		// 根据 categoryId 过滤商品
		if priceCondition == "" {
			query = "SELECT product_id, name, description, cover_url, category_id, prices " +
				"FROM b2c_product.products WHERE category_id = ? " +
				"LIMIT ? OFFSET ?"
			args = append(args, categoryId, pageSize, (currentPage-1)*pageSize)
			countQuery = "SELECT COUNT(*) FROM b2c_product.products WHERE category_id = ?"
		} else {
			// 如果 priceCondition 不为空，拼接 `AND` 条件
			query = "SELECT product_id, name, description, cover_url, category_id, prices " +
				"FROM b2c_product.products WHERE category_id = ? AND " + priceCondition +
				" LIMIT ? OFFSET ?"
			args = append(args, categoryId, pageSize, (currentPage-1)*pageSize)
			countQuery = "SELECT COUNT(*) FROM b2c_product.products WHERE category_id = ? AND " + priceCondition
		}
	}

	// 执行商品查询
	row, err := config.MySqlClient.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer row.Close()

	// 获取商品总数
	var totalItems int32
	// 查询总数时不需要分页参数
	if priceCondition == "" {
		if categoryId == 0 {
			err = config.MySqlClient.QueryRow(countQuery).Scan(&totalItems)
		} else {
			err = config.MySqlClient.QueryRow(countQuery, categoryId).Scan(&totalItems)
		}
	} else {
		if categoryId == 0 {
			err = config.MySqlClient.QueryRow(countQuery, args[0], args[1]).Scan(&totalItems)
		} else {
			err = config.MySqlClient.QueryRow(countQuery, args[0], args[1], args[2]).Scan(&totalItems)
		}
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
		var productCover string
		var productCategoryId int32
		var price float64
		err = row.Scan(&productId, &productName, &productDescription, &productCover, &productCategoryId, &price)
		if err != nil {
			return nil, 0, err
		}
		productList = append(productList, &pb.ProductListItem{
			ProductId:         productId,
			ProductName:       productName,
			Description:       productDescription,
			ProductCover:      productCover,
			ProductCategoryId: productCategoryId,
			ProductPrice:      price,
		})
	}

	return productList, totalItems, nil
}
