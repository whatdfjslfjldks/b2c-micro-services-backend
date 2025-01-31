package repository

import (
	"database/sql"
	"errors"
	"log"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/pkg/utils"
	"micro-services/product-server/pkg/config"
	"micro-services/product-server/pkg/model/dto"
)

func SaveProductsBatch(products []dto.Product) error {
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
	result, err := config.MySqlClient.Exec(sql, args...)
	if err != nil {
		return errors.New("GLB-003")
	}
	// 获取最后插入的自增ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return errors.New("GLB-003")
	}
	// 获取批量插入的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("GLB-003")
	}

	// 创建一个切片存储每个插入的 product_id
	productIds := make([]int64, 0, rowsAffected)

	// 生成每个插入的 product_id
	for i := int64(0); i < rowsAffected; i++ {
		productIds = append(productIds, lastInsertID+i)
	}

	//log.Printf("Inserted %d products with IDs: %v", rowsAffected, productIds)
	// 初始化sale表
	err = initSale(productIds)
	if err != nil {
		log.Printf("Failed to init sale table: %v", err)
		return errors.New("GLB-003")
	}
	return nil
}

func initSale(ids []int64) error {
	// 生成批量插入的SQL语句
	sql := "INSERT INTO b2c_product.product_sale (product_id, sold) VALUES "
	var args []interface{} // 用于存储插入的参数
	// 构建SQL语句和参数
	for _, id := range ids {
		sql += "(?, 0),"
		args = append(args, id)
	}
	// 去掉最后一个逗号
	sql = sql[:len(sql)-1]
	// 执行批量插入
	_, err := config.MySqlClient.Exec(sql, args...)
	if err != nil {
		log.Printf("Error inserting product sale records: %v", err)
		return errors.New("GLB-003")
	}
	// 成功插入
	//fmt.Println("Sale records inserted successfully.")
	return nil
}

func GetProductById(id int32) (dto.ProductById, error) {
	var product dto.ProductById
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
func GetProductDetailById(id int32) (dto.ProductDetail, error) {
	var product dto.ProductDetail
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

	for row.Next() {
		var pType pb.ProductType // 直接使用结构体，而不是指针
		err = row.Scan(&pType.TypeName, &pType.TypeId)
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

	for rowImg.Next() {
		var img pb.ProductImg // 直接使用结构体，而不是指针
		err = rowImg.Scan(&img.ImgUrl, &img.ImgId)
		if err != nil {
			log.Printf("Failed to get product image by ID-%d: %v", id, err)
			return product, errors.New("GLB-003")
		}
		product.ProductImg = append(product.ProductImg, &img)
	}
	// 获取商品已售数量
	sold, err := getProductSold(id)
	if err != nil {
		log.Printf("Failed to get product sold by ID-%d: %v", id, err)
		return product, errors.New("GLB-003")
	}
	product.Sold = sold

	return product, nil
}

// getProductSold 获取企业商品已售数量
func getProductSold(id int32) (int32, error) {
	var sold int32
	err := config.MySqlClient.QueryRow("SELECT sold FROM b2c_product.product_sale WHERE product_id = ?", id).Scan(&sold)
	if err != nil {
		log.Printf("Failed to get product sold by ID-%d: %v", id, err)
		return sold, errors.New("GLB-003")
	}
	return sold, nil
}

// UploadSecProduct 插入秒杀商品数据
func UploadSecProduct(secProduct dto.SecKillProduct) error {
	// 开始一个事务, 保证事务可回滚
	tx, err := config.MySqlClient.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return errors.New("GLB-003")
	}

	// 确保在函数结束时进行回滚（如果出现错误）
	defer func() {
		if err != nil {
			tx.Rollback() // 出现错误时回滚事务
		}
	}()

	// 转换时间
	s, err := utils.ConvertTime(secProduct.StartTime)
	e, err := utils.ConvertTime(secProduct.EndTime)
	if err != nil {
		return errors.New("PRT-003")
	}

	// 插入秒杀商品数据
	result, err := tx.Exec("INSERT INTO b2c_product.products_seckill (name, description, flash_sale_price, original_price, stock, start_time, end_time, create_at, update_at,time,sec_cover) VALUES (?, ?, ?, ?, ?, ?,?,?,?,?,?)",
		secProduct.SecName, secProduct.SecDescription, secProduct.SecOriginalPrice, secProduct.SecPrice, secProduct.Stock, s, e, utils.GetTime(), utils.GetTime(), secProduct.Time, secProduct.SecImg[0].ImgUrl)
	if err != nil {
		log.Printf("Failed to insert sec_kill_product: %v", err)
		return errors.New("GLB-003")
	}

	// 获取最后插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last insert ID: %v", err)
		return errors.New("GLB-003")
	}

	// TODO 避免在循环体操作数据库，优化为批量插入
	// 插入秒杀商品类型数据
	var args []interface{}
	query := "INSERT INTO b2c_product.product_seckill_type (sec_id, sec_type) VALUES"
	for _, sType := range secProduct.SecType {
		query += "(?,?),"
		args = append(args, id, sType.TypeName)
	}
	query = query[:len(query)-1]
	_, err = tx.Exec(query, args...)
	if err != nil {
		log.Printf("Failed to insert sec_kill_product_type: %v", err)
		return errors.New("GLB-003")
	}
	// 插入秒杀商品图片数据
	var args2 []interface{}
	query2 := "INSERT INTO b2c_product.product_seckill_image (sec_id,sec_img) VALUES"
	for _, sImg := range secProduct.SecImg {
		query2 += "(?,?),"
		args2 = append(args2, id, sImg.ImgUrl)
	}
	query2 = query2[:len(query2)-1]
	_, err = tx.Exec(query2, args2...)
	if err != nil {
		log.Printf("Failed to insert sec_kill_product_image: %v", err)
		return errors.New("GLB-003")
	}
	// 初始化sale表
	_, err = tx.Exec("INSERT INTO b2c_product.product_seckill_sale (sec_id, sold) VALUES (?, ?)", id, 0)
	if err != nil {
		log.Printf("Failed to insert product_sale: %v", err)
		return errors.New("GLB-003")
	}
	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return errors.New("GLB-003")
	}

	return nil
}

// GetProductList 获取商品列表，支持按分类和价格区间过滤
func GetProductList(currentPage int32, pageSize int32, categoryId int32, sort int32) (
	[]*pb.ProductListItem, int32, error) {
	// 构建 SQL 查询语句，根据 categoryId 和 priceRange 判断是否过滤
	//log.Printf("sort is %d", sort)
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

// GetSecList 获取秒杀商品列表
// TODO 已售
func GetSecList(currentPage int32, pageSize int32, time int32) (
	[]*pb.SecListItem, error) {
	// 使用 JOIN 查询一次性获取商品信息及已售数量
	query := `
		SELECT
			p.sec_id,
			p.name,
			p.description,
			p.flash_sale_price,
			p.original_price,
			p.stock,
			p.start_time,
			p.end_time,
			p.sec_cover,
			IFNULL(s.sold, 0) AS sold
		FROM
			b2c_product.products_seckill p
		LEFT JOIN
			b2c_product.product_seckill_sale s
		ON
			p.sec_id = s.sec_id
		WHERE
			p.time = ? 
		LIMIT ? 
		OFFSET ?`

	// 执行商品查询
	rows, err := config.MySqlClient.Query(query, time, pageSize, (currentPage-1)*pageSize)
	if err != nil {
		return nil, errors.New("GLB-003")
	}
	defer rows.Close()

	// 处理查询结果
	var productList []*pb.SecListItem
	for rows.Next() {
		var secId int32
		var secName string
		var secDescription string
		var secOriginalPrice float64
		var secPrice float64
		var secStock int32
		var secStartTime string
		var secEndTime string
		var secCover sql.NullString // 使用 NullString 类型来处理可能的 NULL 值
		var sold int32

		err = rows.Scan(&secId, &secName, &secDescription, &secPrice, &secOriginalPrice, &secStock, &secStartTime, &secEndTime, &secCover, &sold)
		if err != nil {
			return nil, errors.New("GLB-003")
		}

		productList = append(productList, &pb.SecListItem{
			SecId:            secId,
			SecName:          secName,
			SecOriginalPrice: secOriginalPrice,
			SecPrice:         secPrice,
			SecStock:         secStock,
			SecCover:         secCover.String,
			SecStartTime:     secStartTime,
			SecEndTime:       secEndTime,
			SecSold:          sold,
		})
	}

	// 检查查询过程中是否有错误
	if err := rows.Err(); err != nil {
		return nil, errors.New("GLB-003")
	}

	return productList, nil
}

// GetSecTotalItems 获取秒杀商品总数
func GetSecTotalItems(time int32) (int32, error) {
	query := `SELECT COUNT(*) FROM b2c_product.products_seckill WHERE time = ?`
	var count int32
	err := config.MySqlClient.QueryRow(query, time).Scan(&count)
	if err != nil {
		return 0, errors.New("GLB-003")
	}
	return count, nil
}
