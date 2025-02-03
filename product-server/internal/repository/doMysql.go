package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/pkg/utils"
	"micro-services/product-server/pkg/config"
	"micro-services/product-server/pkg/localLog"
	"micro-services/product-server/pkg/model/dto"
	"strings"
)

// SaveProductsBatch 只批量上传普通商品
func SaveProductsBatch(products []dto.Product) error {
	// 使用事务
	tx, err := config.MySqlClient.Begin()
	if err != nil {
		localLog.ProductLog.Error("SaveProductsBatch--Failed to begin transaction: " + err.Error())
		return errors.New("GLB-003")
	}
	defer tx.Rollback() // 如果发生错误，回滚事务

	// 插入顶层products表
	sql := "INSERT INTO b2c_product.products (name, category, kind, description, create_time) VALUES "
	var productArgs []interface{}
	for _, p := range products {
		sql += "(?, ?, ?, ?, ?),"
		productArgs = append(productArgs, p.Name, p.Category, 1, p.Description, utils.GetTime())
	}
	// 去掉最后一个逗号
	sql = sql[:len(sql)-1]
	// 执行批量插入
	result, err := tx.Exec(sql, productArgs...)
	if err != nil {
		localLog.ProductLog.Error("SaveProductsBatch--Failed to execute batch insert: " + err.Error())
		return errors.New("GLB-003")
	}

	// 获取最后插入的自增ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		localLog.ProductLog.Error("SaveProductsBatch--Failed to get last insert ID: " + err.Error())
		return errors.New("GLB-003")
	}

	// 获取批量插入的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		localLog.ProductLog.Error("SaveProductsBatch--Failed to get rows affected: " + err.Error())
		return errors.New("GLB-003")
	}

	// 创建一个切片存储每个插入的 product_id
	productIds := make([]int64, 0, rowsAffected)
	for i := int64(0); i < rowsAffected; i++ {
		productIds = append(productIds, lastInsertID+i)
	}
	// 初始化sale表
	sqlSale := "INSERT INTO b2c_product.product_sale (product_id, stock, sold) VALUES "
	var saleArgs []interface{}
	for i, p := range products {
		sqlSale += "(?, ?, ?),"
		saleArgs = append(saleArgs, productIds[i], p.Stock, 0)
	}
	sqlSale = sqlSale[:len(sqlSale)-1] // 去掉最后一个逗号
	_, err = tx.Exec(sqlSale, saleArgs...)
	if err != nil {
		localLog.ProductLog.Error("SaveProductsBatch--Failed to execute batch insert: " + err.Error())
		return errors.New("GLB-003")
	}

	// 初始化price表
	sqlPrice := "INSERT INTO b2c_product.product_price (product_id, price, origin_price, price_change_date) VALUES "
	var priceArgs []interface{}
	for i, p := range products {
		// 通过 productIds[i] 获取对应的 product_id
		sqlPrice += "(?, ?, ?, ?),"
		priceArgs = append(priceArgs, productIds[i], p.Price, p.Price, utils.GetTime())
	}
	sqlPrice = sqlPrice[:len(sqlPrice)-1] // 去掉最后一个逗号
	_, err = tx.Exec(sqlPrice, priceArgs...)
	if err != nil {
		localLog.ProductLog.Error("SaveProductsBatch--Failed to execute batch insert: " + err.Error())
		return errors.New("GLB-003")
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		localLog.ProductLog.Error("SaveProductsBatch--Failed to commit transaction: " + err.Error())
		return errors.New("GLB-003")
	}
	return nil
}

func GetInfoFromProducts(id int32) (
	dto.Products, error) {
	query := "SELECT name, kind, category, description, create_time FROM b2c_product.products WHERE id=?"
	var products dto.Products
	products.ID = id
	err := config.MySqlClient.QueryRow(query, id).Scan(&products.Name, &products.Kind, &products.Category, &products.Description, &products.CreateTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return products, errors.New("GLB-001")
		}
		return products, errors.New("GLB-003")
	}
	return products, nil
}

func GetProductDetail(id int32) (dto.NormalProducts, error) {
	// 创建返回值结构
	var productDetail dto.NormalProducts

	// 查询商品基本信息、价格、库存等
	baseQuery := `
        SELECT 
            ps.stock, 
            ps.sold, 
            pp.price, 
            pp.origin_price
        FROM 
            b2c_product.products p
        LEFT JOIN 
            b2c_product.product_sale ps ON ps.product_id = p.id
        LEFT JOIN 
            b2c_product.product_price pp ON pp.product_id = p.id
        WHERE 
            p.id = ?
    `
	row := config.MySqlClient.QueryRow(baseQuery, id)
	var stock, sold int32
	var price, originPrice float64
	err := row.Scan(&stock, &sold, &price, &originPrice)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			localLog.ProductLog.Error("GetProductDetail--Product not found: " + err.Error())
			return productDetail, nil
		}
		localLog.ProductLog.Error("GetProductDetail--Failed to execute base query: " + err.Error())
		return productDetail, errors.New("GLB-003")
	}
	productDetail.Stock = stock
	productDetail.Sold = sold
	productDetail.Price = price
	productDetail.OriginPrice = originPrice

	// 查询商品图片
	imgQuery := `
        SELECT 
            image
        FROM 
            b2c_product.product_image
        WHERE 
            product_id = ?
    `
	imgRows, err := config.MySqlClient.Query(imgQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			localLog.ProductLog.Error("GetProductDetail--Product not found: " + err.Error())
			return productDetail, nil
		}
		localLog.ProductLog.Error("GetProductDetail--Failed to execute image query: " + err.Error())
		return productDetail, errors.New("GLB-003")
	}
	defer imgRows.Close()

	var imageList []*pb.PImg
	for imgRows.Next() {
		var imgUrl sql.NullString
		err := imgRows.Scan(&imgUrl)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				localLog.ProductLog.Error("GetProductDetail--Product not found: " + err.Error())
				return productDetail, nil
			}
			localLog.ProductLog.Error("GetProductDetail--Failed to scan image row: " + err.Error())
			return productDetail, errors.New("GLB-003")
		}
		if imgUrl.Valid {
			img := pb.PImg{ImgUrl: imgUrl.String}
			imageList = append(imageList, &img)
		}
	}
	productDetail.ImageList = imageList

	// 查询商品类型
	typeQuery := `
        SELECT 
            type_description
        FROM 
            b2c_product.product_type
        WHERE 
            product_id = ?
    `
	typeRows, err := config.MySqlClient.Query(typeQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			localLog.ProductLog.Error("GetProductDetail--Product not found: " + err.Error())
			return productDetail, nil
		}
		localLog.ProductLog.Error("GetProductDetail--Failed to execute type query: " + err.Error())
		return productDetail, errors.New("GLB-003")
	}
	defer typeRows.Close()

	var typeList []*pb.PType
	for typeRows.Next() {
		var typeName sql.NullString
		err := typeRows.Scan(&typeName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				localLog.ProductLog.Error("GetProductDetail--Product not found: " + err.Error())
				return productDetail, nil
			}
			localLog.ProductLog.Error("GetProductDetail--Failed to scan type row: " + err.Error())
			return productDetail, errors.New("GLB-003")
		}
		if typeName.Valid {
			typ := pb.PType{TypeName: typeName.String}
			typeList = append(typeList, &typ)
		}
	}
	productDetail.TypeList = typeList

	// 如果没有查询到任何数据，返回错误
	if len(imageList) == 0 && len(typeList) == 0 {
		return productDetail, nil
	}

	return productDetail, nil
}

func GetSecProductDetail(id int32) (dto.SecProducts, error) {
	// 创建返回值结构
	var productDetail dto.SecProducts

	// 查询商品基本信息、价格、库存、秒杀信息
	baseQuery := `
        SELECT 
            ps.stock, 
            ps.sold, 
            pp.price, 
            pp.origin_price,
            sk.start_time,
            sk.duration,
            sk.session_id
        FROM 
            b2c_product.products p
        LEFT JOIN 
            b2c_product.product_sale ps ON ps.product_id = p.id
        LEFT JOIN 
            b2c_product.product_price pp ON pp.product_id = p.id
        LEFT JOIN 
            b2c_product.seckill sk ON p.id = sk.product_id
        WHERE 
            p.id = ?
    `
	row := config.MySqlClient.QueryRow(baseQuery, id)
	var stock, sold, sessionId int32
	var price, originPrice float64
	var startTime, duration string
	err := row.Scan(&stock, &sold, &price, &originPrice, &startTime, &duration, &sessionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			localLog.ProductLog.Error("GetSecProductDetail--Product not found: " + err.Error())
			return productDetail, errors.New("GLB-001")
		}
		localLog.ProductLog.Error("GetSecProductDetail--Failed to execute base query: " + err.Error())
		return productDetail, errors.New("GLB-003")
	}
	productDetail.Stock = stock
	productDetail.Sold = sold
	productDetail.Price = price
	productDetail.OriginPrice = originPrice
	productDetail.StartTime = startTime
	productDetail.Duration = duration
	productDetail.SessionId = sessionId

	// 查询商品图片
	imgQuery := `
        SELECT 
            image
        FROM 
            b2c_product.product_image
        WHERE 
            product_id = ?
    `
	imgRows, err := config.MySqlClient.Query(imgQuery, id)
	if err != nil {
		localLog.ProductLog.Error("GetSecProductDetail--Failed to execute image query: " + err.Error())
		return productDetail, errors.New("GLB-003")
	}
	defer imgRows.Close()

	var imageList []*pb.PImg
	for imgRows.Next() {
		var imgUrl sql.NullString
		err := imgRows.Scan(&imgUrl)
		if err != nil {
			localLog.ProductLog.Error("GetSecProductDetail--Failed to scan image row: " + err.Error())
			return productDetail, errors.New("GLB-003")
		}
		if imgUrl.Valid {
			img := pb.PImg{ImgUrl: imgUrl.String}
			imageList = append(imageList, &img)
		}
	}
	productDetail.ImageList = imageList

	// 查询商品类型
	typeQuery := `
        SELECT 
            type_description
        FROM 
            b2c_product.product_type
        WHERE 
            product_id = ?
    `
	typeRows, err := config.MySqlClient.Query(typeQuery, id)
	if err != nil {
		localLog.ProductLog.Error("GetSecProductDetail--Failed to execute type query: " + err.Error())
		return productDetail, errors.New("GLB-003")
	}
	defer typeRows.Close()

	var typeList []*pb.PType
	for typeRows.Next() {
		var typeName sql.NullString
		err := typeRows.Scan(&typeName)
		if err != nil {
			localLog.ProductLog.Error("GetSecProductDetail--Failed to scan type row: " + err.Error())
			return productDetail, errors.New("GLB-003")
		}
		if typeName.Valid {
			typ := pb.PType{TypeName: typeName.String}
			typeList = append(typeList, &typ)
		}
	}
	productDetail.TypeList = typeList

	// 如果没有查询到任何数据，返回错误
	if len(imageList) == 0 && len(typeList) == 0 {
		return productDetail, nil
	}

	return productDetail, nil
}

// UploadSecProduct 插入秒杀商品数据
func UploadSecProduct(secProduct *pb.UploadSecKillProductRequest) error {
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
	if err != nil {
		return errors.New("PRT-003")
	}

	// 插入秒杀商品数据
	result, err := tx.Exec("INSERT INTO b2c_product.products (name, kind, category, description, create_time)  VALUES (?, ?, ?, ?, ?)",
		secProduct.Name, 2, secProduct.CategoryId, secProduct.Description, utils.GetTime())
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
	query := "INSERT INTO b2c_product.product_type (product_id, type_description) VALUES"
	for _, sType := range secProduct.PType {
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
	query2 := "INSERT INTO b2c_product.product_image (product_id, image) VALUES"
	for _, sImg := range secProduct.PImg {
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
	_, err = tx.Exec("INSERT INTO b2c_product.product_sale (product_id,stock, sold) VALUES (?,?, ?)", id, secProduct.Stock, 0)
	if err != nil {
		log.Printf("Failed to insert product_sale: %v", err)
		return errors.New("GLB-003")
	}
	// seckill表
	_, err = tx.Exec("INSERT INTO b2c_product.seckill(product_id, start_time, duration, session_id) VALUES (?,?,?,?)", id, s, secProduct.Duration, secProduct.SessionId)
	if err != nil {
		log.Printf("Failed to insert seckill: %v", err)
		return errors.New("GLB-003")
	}
	// price表
	_, err = tx.Exec("INSERT INTO b2c_product.product_price(product_id,origin_price,price) VALUES (?,?,?)", id, secProduct.OriginalPrice, secProduct.Price)
	if err != nil {
		log.Printf("Failed to insert price: %v", err)
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

// GetProductList 获取商品列表
func GetProductList(currentPage int32, pageSize int32, categoryId int32, kind int32, sort int32) ([]*pb.ProductListItem, int32, error) {
	// 创建返回值结构
	var productList []*pb.ProductListItem

	// 构建基础查询语句和参数
	baseQuery, args := buildBaseQuery(currentPage, pageSize, categoryId, kind, sort)

	// 执行基本信息查询
	rows, err := config.MySqlClient.Query(baseQuery, args...)
	if err != nil {
		logError("Failed to execute base query", err)
		return nil, 0, err
	}
	defer rows.Close()

	countQuery, countArgs := buildCountQuery(categoryId, kind)
	var totalItems int32
	err = config.MySqlClient.QueryRow(countQuery, countArgs...).Scan(&totalItems)
	if err != nil {
		logError("Failed to execute count query", err)
		return nil, 0, err
	}

	// 存储商品 ID 列表和商品映射
	productIDs, productMap, err := processBaseQuery(rows)
	if err != nil {
		return nil, 0, err
	}

	// 查询每个商品的图片
	if len(productIDs) > 0 {
		err = queryImages(productIDs, productMap)
		if err != nil {
			return nil, 0, err
		}
	}

	// 查询每个商品的类型
	if len(productIDs) > 0 {
		err = queryTypes(productIDs, productMap)
		if err != nil {
			return nil, 0, err
		}
	}

	// 按照原始顺序将商品添加到列表中
	for _, id := range productIDs {
		if product, ok := productMap[id]; ok {
			productList = append(productList, product)
		}
	}

	// 如果没有查询到任何数据，返回错误
	if len(productList) == 0 {
		return nil, 0, err
	}

	return productList, totalItems, nil
}

// buildCountQuery 构建统计满足条件商品数量的查询语句和参数
func buildCountQuery(categoryId int32, kind int32) (string, []interface{}) {
	countQuery := `
        SELECT 
            COUNT(*)
        FROM 
            b2c_product.products p
        WHERE 
            1=1
    `
	var args []interface{}

	// 如果传入了 categoryId，则添加过滤条件
	if categoryId != 0 {
		countQuery += " AND p.category = ?"
		args = append(args, categoryId)
	}

	// 如果传入了 kind，则添加过滤条件
	if kind != 0 {
		countQuery += " AND p.kind = ?"
		args = append(args, kind)
	}

	return countQuery, args
}

// buildBaseQuery 构建基础查询语句和参数
func buildBaseQuery(currentPage int32, pageSize int32, categoryId int32, kind int32, sort int32) (string, []interface{}) {
	baseQuery := `
        SELECT 
            p.id, 
            p.name,
            p.description,
            p.create_time,
            ps.stock, 
            ps.sold, 
            pp.price, 
            pp.origin_price
        FROM 
            b2c_product.products p
        LEFT JOIN 
            b2c_product.product_sale ps ON ps.product_id = p.id
        LEFT JOIN 
            b2c_product.product_price pp ON pp.product_id = p.id
        WHERE 
            1=1
    `
	var args []interface{}

	// 如果传入了 categoryId，则添加过滤条件
	if categoryId != 0 {
		baseQuery += " AND p.category = ?"
		args = append(args, categoryId)
	}

	// 如果传入了 kind，则添加过滤条件
	if kind != 0 {
		baseQuery += " AND p.kind = ?"
		args = append(args, kind)
	}

	// 根据排序标准添加排序条件
	switch sort {
	case 1:
		// 按价格升序排序
		baseQuery += " ORDER BY pp.price ASC"
	case 2:
		// 按创建时间降序排序
		baseQuery += " ORDER BY p.create_time DESC"
	default:
		// 默认按产品 ID 升序排序
		baseQuery += " ORDER BY p.id ASC"
	}

	// 添加分页查询的条件
	offset := (currentPage - 1) * pageSize
	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)

	return baseQuery, args
}

// processBaseQuery 处理基础查询结果
func processBaseQuery(rows *sql.Rows) ([]int32, map[int32]*pb.ProductListItem, error) {
	var productIDs []int32
	productMap := make(map[int32]*pb.ProductListItem)

	for rows.Next() {
		var product pb.ProductListItem
		var stock sql.NullInt32
		var sold sql.NullInt32
		var price sql.NullFloat64
		var originPrice sql.NullFloat64
		var createTime sql.NullString

		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&createTime,
			&stock,
			&sold,
			&price,
			&originPrice,
		)
		if err != nil {
			logError("Failed to scan base row", err)
			return nil, nil, err
		}

		if createTime.Valid {
			product.CreateTime = createTime.String
		}
		if stock.Valid {
			product.Stock = stock.Int32
		} else {
			product.Stock = 0
		}
		if sold.Valid {
			product.Sold = sold.Int32
		} else {
			product.Sold = 0
		}
		if price.Valid {
			product.Price = price.Float64
		} else {
			product.Price = 0
		}
		if originPrice.Valid {
			product.OriginalPrice = originPrice.Float64
		} else {
			product.OriginalPrice = 0
		}

		productMap[product.Id] = &product
		productIDs = append(productIDs, product.Id)
	}

	return productIDs, productMap, nil
}

// queryImages 查询商品图片
func queryImages(productIDs []int32, productMap map[int32]*pb.ProductListItem) error {
	imgQuery := `
        SELECT 
            product_id, 
            image
        FROM 
            b2c_product.product_image
        WHERE 
            product_id IN (?` + strings.Repeat(",?", len(productIDs)-1) + `)
    `
	imgArgs := make([]interface{}, len(productIDs))
	for i, id := range productIDs {
		imgArgs[i] = id
	}
	imgRows, err := config.MySqlClient.Query(imgQuery, imgArgs...)
	if err != nil {
		logError("Failed to execute image query", err)
		return err
	}
	defer imgRows.Close()

	for imgRows.Next() {
		var productID int32
		var imgUrl sql.NullString

		err := imgRows.Scan(&productID, &imgUrl)
		if err != nil {
			logError("Failed to scan image row", err)
			return err
		}

		if imgUrl.Valid {
			img := pb.PImg{ImgUrl: imgUrl.String}
			if product, ok := productMap[productID]; ok {
				product.PImg = append(product.PImg, &img)
			}
		}
	}

	return nil
}

// queryTypes 查询商品类型
func queryTypes(productIDs []int32, productMap map[int32]*pb.ProductListItem) error {
	typeQuery := `
        SELECT 
            product_id, 
            type_description
        FROM 
            b2c_product.product_type
        WHERE 
            product_id IN (?` + strings.Repeat(",?", len(productIDs)-1) + `)
    `
	typeArgs := make([]interface{}, len(productIDs))
	for i, id := range productIDs {
		typeArgs[i] = id
	}
	typeRows, err := config.MySqlClient.Query(typeQuery, typeArgs...)
	if err != nil {
		logError("Failed to execute type query", err)
		return err
	}
	defer typeRows.Close()

	for typeRows.Next() {
		var productID int32
		var typeName sql.NullString

		err := typeRows.Scan(&productID, &typeName)
		if err != nil {
			logError("Failed to scan type row", err)
			return err
		}

		if typeName.Valid {
			typ := pb.PType{TypeName: typeName.String}
			if product, ok := productMap[productID]; ok {
				product.PType = append(product.PType, &typ)
			}
		}
	}

	return nil
}

// GetSecList 获取秒杀商品列表
func GetSecList(currentPage int32, pageSize int32, sessionId int32) ([]*pb.SecListItem, int32, error) {
	// 创建返回值结构
	var productList []*pb.SecListItem

	// 构建基础查询语句和参数
	baseQuery, args := buildSecBaseQuery(currentPage, pageSize, sessionId, 2)

	//log.Printf("baseQuery: %s, args: %v", baseQuery, args)

	// 执行基本信息查询
	rows, err := config.MySqlClient.Query(baseQuery, args...)
	if err != nil {
		logError("Failed to execute base query", err)
		return nil, 0, err
	}
	defer rows.Close()

	countQuery, countArgs := buildSecCountQuery(2, sessionId)

	//log.Printf("countQuery: %s, countArgs: %v", countQuery, countArgs)

	var totalItems int32
	err = config.MySqlClient.QueryRow(countQuery, countArgs...).Scan(&totalItems)
	if err != nil {
		logError("Failed to execute count query", err)
		return nil, 0, err
	}

	// 存储商品 ID 列表和商品映射
	productIDs, productMap, err := processSecBaseQuery(rows)
	if err != nil {
		return nil, 0, err
	}

	// 查询每个商品的图片
	if len(productIDs) > 0 {
		err = querySecImages(productIDs, productMap)
		if err != nil {
			return nil, 0, err
		}
	}

	// 查询每个商品的类型
	if len(productIDs) > 0 {
		err = querySecTypes(productIDs, productMap)
		if err != nil {
			return nil, 0, err
		}
	}

	// 按照原始顺序将商品添加到列表中
	for _, id := range productIDs {
		if product, ok := productMap[id]; ok {
			productList = append(productList, product)
		}
	}

	// 如果没有查询到任何数据，返回错误
	if len(productList) == 0 {
		return nil, 0, err
	}

	return productList, totalItems, nil
}

// processSecBaseQuery 处理基础查询结果
func processSecBaseQuery(rows *sql.Rows) ([]int32, map[int32]*pb.SecListItem, error) {
	var productIDs []int32
	productMap := make(map[int32]*pb.SecListItem)

	for rows.Next() {
		var product pb.SecListItem
		var stock sql.NullInt32
		var sold sql.NullInt32
		var price sql.NullFloat64
		var originPrice sql.NullFloat64
		var createTime sql.NullString

		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&createTime,
			&stock,
			&sold,
			&price,
			&originPrice,
			&product.SessionId,
		)
		if err != nil {
			logError("Failed to scan base row", err)
			return nil, nil, err
		}

		if createTime.Valid {
			product.CreateTime = createTime.String
		}
		if stock.Valid {
			product.Stock = stock.Int32
		} else {
			product.Stock = 0
		}
		if sold.Valid {
			product.Sold = sold.Int32
		} else {
			product.Sold = 0
		}
		if price.Valid {
			product.Price = price.Float64
		} else {
			product.Price = 0
		}
		if originPrice.Valid {
			product.OriginalPrice = originPrice.Float64
		} else {
			product.OriginalPrice = 0
		}

		productMap[product.Id] = &product
		productIDs = append(productIDs, product.Id)
	}

	return productIDs, productMap, nil
}

// queryImages 查询商品图片
func querySecImages(productIDs []int32, productMap map[int32]*pb.SecListItem) error {
	imgQuery := `
        SELECT 
            product_id, 
            image
        FROM 
            b2c_product.product_image
        WHERE 
            product_id IN (?` + strings.Repeat(",?", len(productIDs)-1) + `)
    `
	imgArgs := make([]interface{}, len(productIDs))
	for i, id := range productIDs {
		imgArgs[i] = id
	}
	imgRows, err := config.MySqlClient.Query(imgQuery, imgArgs...)
	if err != nil {
		logError("Failed to execute image query", err)
		return err
	}
	defer imgRows.Close()

	for imgRows.Next() {
		var productID int32
		var imgUrl sql.NullString

		err := imgRows.Scan(&productID, &imgUrl)
		if err != nil {
			logError("Failed to scan image row", err)
			return err
		}

		if imgUrl.Valid {
			img := pb.PImg{ImgUrl: imgUrl.String}
			if product, ok := productMap[productID]; ok {
				product.PImg = append(product.PImg, &img)
			}
		}
	}

	return nil
}

// queryTypes 查询商品类型
func querySecTypes(productIDs []int32, productMap map[int32]*pb.SecListItem) error {
	typeQuery := `
        SELECT 
            product_id, 
            type_description
        FROM 
            b2c_product.product_type
        WHERE 
            product_id IN (?` + strings.Repeat(",?", len(productIDs)-1) + `)
    `
	typeArgs := make([]interface{}, len(productIDs))
	for i, id := range productIDs {
		typeArgs[i] = id
	}
	typeRows, err := config.MySqlClient.Query(typeQuery, typeArgs...)
	if err != nil {
		logError("Failed to execute type query", err)
		return err
	}
	defer typeRows.Close()

	for typeRows.Next() {
		var productID int32
		var typeName sql.NullString

		err := typeRows.Scan(&productID, &typeName)
		if err != nil {
			logError("Failed to scan type row", err)
			return err
		}

		if typeName.Valid {
			typ := pb.PType{TypeName: typeName.String}
			if product, ok := productMap[productID]; ok {
				product.PType = append(product.PType, &typ)
			}
		}
	}

	return nil
}

// buildSecCountQuery 构建统计满足条件商品数量的查询语句和参数
func buildSecCountQuery(kind int32, sessionId int32) (string, []interface{}) {
	countQuery := `
       SELECT
            COUNT(*)
        FROM
            b2c_product.products p
        LEFT JOIN
            b2c_product.seckill sk 
        ON 
            sk.product_id = p.id AND sk.session_id = ?
        WHERE
            p.kind = ?;
    `
	var args []interface{}
	args = append(args, sessionId)
	args = append(args, kind)

	return countQuery, args
}

// buildBaseQuery 构建基础查询语句和参数
func buildSecBaseQuery(currentPage int32, pageSize int32, sessionId int32, kind int32) (string, []interface{}) {
	baseQuery := `
        SELECT 
            p.id, 
            p.name,
            p.description,
            p.create_time,
            ps.stock, 
            ps.sold, 
            pp.price, 
            pp.origin_price,
            sk.session_id
        FROM 
            b2c_product.products p
        LEFT JOIN 
            b2c_product.product_sale ps ON ps.product_id = p.id
        LEFT JOIN 
            b2c_product.product_price pp ON pp.product_id = p.id
        LEFT JOIN 
            b2c_product.seckill sk ON sk.product_id = p.id AND sk.session_id = ?
        WHERE
            p.kind = ?
    `
	args := []interface{}{sessionId, kind}

	// 添加分页查询的条件，使用占位符
	baseQuery += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, (currentPage-1)*pageSize)

	return baseQuery, args
}

// logError 记录错误日志
func logError(message string, err error) {
	localLog.ProductLog.Error("GetList--" + message + ": " + err.Error())
}
