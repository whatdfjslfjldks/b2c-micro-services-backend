package repository

import (
	"database/sql"
	"errors"
	"log"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/pkg/utils"
	"micro-services/product-server/pkg/config"
	"micro-services/product-server/pkg/localLog"
	"micro-services/product-server/pkg/model/dto"
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
	// 查询SQL
	query := `
		SELECT 
			pi.image, 
			pt.type_description, 
			ps.stock, 
			ps.sold, 
			pp.price, 
			pp.origin_price
		FROM 
			b2c_product.products p
		LEFT JOIN 
			b2c_product.product_image pi ON pi.product_id = p.id
		LEFT JOIN 
			b2c_product.product_type pt ON pt.product_id = p.id
		LEFT JOIN 
			b2c_product.product_sale ps ON ps.product_id = p.id
		LEFT JOIN 
			b2c_product.product_price pp ON pp.product_id = p.id
		WHERE 
			p.id = ?
	`
	// 执行查询
	rows, err := config.MySqlClient.Query(query, id)
	if err != nil {
		localLog.ProductLog.Error("GetProductDetail--Failed to execute query: " + err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return productDetail, errors.New("GLB-001")
		}
		return productDetail, errors.New("GLB-003")
	}
	defer rows.Close()

	// 处理查询结果
	var imageList []*pb.ProductImg
	var typeList []*pb.ProductType

	for rows.Next() {
		var img pb.ProductImg
		var typ pb.ProductType
		var stock, sold int32
		var price, originPrice float64

		// 使用 sql.NullString 来处理可能的 NULL 值
		var imgUrl sql.NullString
		var typeName sql.NullString

		// 获取每一行的数据
		err := rows.Scan(
			&imgUrl,      // 获取图片路径
			&typeName,    // 获取商品类型描述
			&stock,       // 获取商品库存
			&sold,        // 获取商品已售数量
			&price,       // 获取商品价格
			&originPrice, // 获取商品原始价格
		)
		if err != nil {
			localLog.ProductLog.Error("GetProductDetail--Failed to scan row: " + err.Error())
			if errors.Is(err, sql.ErrNoRows) {
				return productDetail, errors.New("GLB-001")
			}
			return productDetail, errors.New("GLB-003")
		}

		// 检查 imgUrl 和 typeName 是否为 NULL，如果是，则赋予默认值
		if imgUrl.Valid {
			img.ImgUrl = imgUrl.String
		} else {
			img.ImgUrl = "" // 或者可以设置为默认值
		}

		if typeName.Valid {
			typ.TypeName = typeName.String
		} else {
			typ.TypeName = "" // 或者可以设置为默认值
		}

		// 填充返回的结构体
		imageList = append(imageList, &img)
		typeList = append(typeList, &typ)

		// 只取第一个商品的价格、库存等数据
		productDetail.Price = price
		productDetail.OriginPrice = originPrice
		productDetail.Stock = stock
		productDetail.Sold = sold
	}

	// 填充图片列表和分类列表
	productDetail.ImageList = imageList
	productDetail.TypeList = typeList

	// 如果没有查询到任何数据，返回错误
	if len(imageList) == 0 || len(typeList) == 0 {
		if errors.Is(err, sql.ErrNoRows) {
			return productDetail, errors.New("GLB-001")
		}
		return productDetail, errors.New("GLB-003")
	}

	return productDetail, nil
}

func GetSecProductDetail(id int32) (dto.SecProducts, error) {
	// 创建返回值结构
	var productDetail dto.SecProducts
	// 查询SQL
	query := `
		SELECT 
			pi.image, 
			pt.type_description, 
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
			b2c_product.product_image pi ON pi.product_id = p.id
		LEFT JOIN 
			b2c_product.product_type pt ON pt.product_id = p.id
		LEFT JOIN 
			b2c_product.product_sale ps ON ps.product_id = p.id
		LEFT JOIN 
			b2c_product.product_price pp ON pp.product_id = p.id
		LEFT JOIN 
			b2c_product.seckill sk ON p.id = sk.product_id
		WHERE 
			p.id = ?
	`
	// 执行查询
	rows, err := config.MySqlClient.Query(query, id)
	if err != nil {
		localLog.ProductLog.Error("GetProductDetail--Failed to execute query: " + err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return productDetail, errors.New("GLB-001")
		}
		return productDetail, errors.New("GLB-003")
	}
	defer rows.Close()

	// 处理查询结果
	var imageList []*pb.ProductImg
	var typeList []*pb.ProductType

	for rows.Next() {
		var img pb.ProductImg
		var typ pb.ProductType
		var stock, sold int32
		var price, originPrice float64
		var startTime, duration string
		var sessionId int32

		// 使用 sql.NullString 来处理可能的 NULL 值
		var imgUrl sql.NullString
		var typeName sql.NullString

		// 获取每一行的数据
		err := rows.Scan(
			&imgUrl,      // 获取图片路径
			&typeName,    // 获取商品类型描述
			&stock,       // 获取商品库存
			&sold,        // 获取商品已售数量
			&price,       // 获取商品价格
			&originPrice, // 获取商品原始价格
			&startTime,
			&duration,
			&sessionId,
		)
		if err != nil {
			localLog.ProductLog.Error("GetProductDetail--Failed to scan row: " + err.Error())
			if errors.Is(err, sql.ErrNoRows) {
				return productDetail, errors.New("GLB-001")
			}
			return productDetail, errors.New("GLB-003")
		}

		// 检查 imgUrl 和 typeName 是否为 NULL，如果是，则赋予默认值
		if imgUrl.Valid {
			img.ImgUrl = imgUrl.String
		} else {
			img.ImgUrl = "" // 或者可以设置为默认值
		}

		if typeName.Valid {
			typ.TypeName = typeName.String
		} else {
			typ.TypeName = "" // 或者可以设置为默认值
		}

		// 填充返回的结构体
		imageList = append(imageList, &img)
		typeList = append(typeList, &typ)

		// 只取第一个商品的价格、库存等数据
		productDetail.Price = price
		productDetail.OriginPrice = originPrice
		productDetail.Stock = stock
		productDetail.Sold = sold
		productDetail.StartTime = startTime
		productDetail.Duration = duration
		productDetail.SessionId = sessionId
	}

	// 填充图片列表和分类列表
	productDetail.ImageList = imageList
	productDetail.TypeList = typeList

	// 如果没有查询到任何数据，返回错误
	if len(imageList) == 0 || len(typeList) == 0 {
		if errors.Is(err, sql.ErrNoRows) {
			return productDetail, errors.New("GLB-001")
		}
		return productDetail, errors.New("GLB-003")
	}

	return productDetail, nil
}

//	func GetProductById(id int32) (dto.ProductById, error) {
//		return &dto.ProductById{
//
//		},nil
//		//var product dto.ProductById
//		//var cover sql.NullString // 使用 sql.NullString 来处理 cover_url 字段
//		//
//		//// 查询数据库
//		//err := config.MySqlClient.QueryRow("SELECT name, category_id, description, prices, cover_url FROM b2c_product.products WHERE product_id = ?", id).Scan(&product.Name, &product.Category, &product.Description, &product.Price, &cover)
//		//if err != nil {
//		//	log.Printf("Failed to get product by ID-%d: %v", id, err)
//		//	return product, errors.New("GLB-003")
//		//}
//		//
//		//// 如果 cover 是 NULL，将其设为空字符串
//		//if cover.Valid {
//		//	product.Cover = cover.String
//		//} else {
//		//	product.Cover = "" // NULL 处理为空字符串
//		//}
//		//
//		//return product, nil
//	}
//

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
		secProduct.Name, 1, secProduct.CategoryId, secProduct.Description, utils.GetTime())
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
	for _, sType := range secProduct.Type {
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
	for _, sImg := range secProduct.Img {
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

//
//// GetProductList 获取商品列表，支持按分类和价格区间过滤
//func GetProductList(currentPage int32, pageSize int32, categoryId int32, sort int32) (
//	[]*pb.ProductListItem, int32, error) {
//	// 构建 SQL 查询语句，根据 categoryId 和 priceRange 判断是否过滤
//	//log.Printf("sort is %d", sort)
//	var query string
//	var countQuery string
//	var args []interface{}
//
//	// 基本的查询结构
//	query = "SELECT product_id, name, description, cover_url, category_id, prices " +
//		"FROM b2c_product.products WHERE 1=1"
//	// TODO 1=1 避免判断是否为第一个查询，加上1=1可以等加and
//	countQuery = "SELECT COUNT(*) FROM b2c_product.products where 1=1"
//
//	// 根据 categoryId 判断是否过滤分类
//	if categoryId != 0 {
//		query += " AND category_id = ?"
//		countQuery += " AND category_id = ?"
//		args = append(args, categoryId)
//	}
//
//	// 根据 sort 参数构建排序条件
//	switch sort {
//	case 1:
//		// 按价格升序排序
//		query += " ORDER BY prices ASC"
//	case 2:
//		query += " ORDER BY create_at DESC"
//	default:
//		// 默认不排序（或者按默认排序，比如按 ID 排序）
//		query += " ORDER BY product_id ASC"
//	}
//
//	// 加入分页查询
//	query += " LIMIT ? OFFSET ?"
//	args = append(args, pageSize, (currentPage-1)*pageSize)
//
//	// 执行商品查询
//	row, err := config.MySqlClient.Query(query, args...)
//	if err != nil {
//		return nil, 0, err
//	}
//	defer row.Close()
//
//	// 获取商品总数
//	var totalItems int32
//	if categoryId == 0 {
//		err = config.MySqlClient.QueryRow(countQuery).Scan(&totalItems)
//	} else {
//		err = config.MySqlClient.QueryRow(countQuery, categoryId).Scan(&totalItems)
//	}
//
//	if err != nil {
//		return nil, 0, err
//	}
//
//	// 处理查询结果
//	var productList []*pb.ProductListItem
//	for row.Next() {
//		var productId int32
//		var productName string
//		var productDescription string
//		var productCover sql.NullString // 使用 NullString 类型来处理可能的 NULL 值
//		var productCategoryId int32
//		var price float64
//		err = row.Scan(&productId, &productName, &productDescription, &productCover, &productCategoryId, &price)
//		if err != nil {
//			return nil, 0, err
//		}
//
//		// 如果 cover_url 是 NULL，则设置为默认值
//		var coverURL string
//		if productCover.Valid {
//			coverURL = productCover.String
//		} else {
//			coverURL = "" // 使用空字符串作为默认值
//		}
//
//		productList = append(productList, &pb.ProductListItem{
//			ProductId:         productId,
//			ProductName:       productName,
//			Description:       productDescription,
//			ProductCover:      coverURL,
//			ProductCategoryId: productCategoryId,
//			ProductPrice:      price,
//		})
//	}
//
//	//log.Printf("productList is %v", productList)
//
//	return productList, totalItems, nil
//}
//
//// GetSecList 获取秒杀商品列表
//// TODO 已售
//func GetSecList(currentPage int32, pageSize int32, time int32) (
//	[]*pb.SecListItem, error) {
//	// 使用 JOIN 查询一次性获取商品信息及已售数量
//	query := `
//		SELECT
//			p.sec_id,
//			p.name,
//			p.description,
//			p.flash_sale_price,
//			p.original_price,
//			p.stock,
//			p.start_time,
//			p.end_time,
//			p.sec_cover,
//			IFNULL(s.sold, 0) AS sold
//		FROM
//			b2c_product.products_seckill p
//		LEFT JOIN
//			b2c_product.product_seckill_sale s
//		ON
//			p.sec_id = s.sec_id
//		WHERE
//			p.time = ?
//		LIMIT ?
//		OFFSET ?`
//
//	// 执行商品查询
//	rows, err := config.MySqlClient.Query(query, time, pageSize, (currentPage-1)*pageSize)
//	if err != nil {
//		return nil, errors.New("GLB-003")
//	}
//	defer rows.Close()
//
//	// 处理查询结果
//	var productList []*pb.SecListItem
//	for rows.Next() {
//		var secId int32
//		var secName string
//		var secDescription string
//		var secOriginalPrice float64
//		var secPrice float64
//		var secStock int32
//		var secStartTime string
//		var secEndTime string
//		var secCover sql.NullString // 使用 NullString 类型来处理可能的 NULL 值
//		var sold int32
//
//		err = rows.Scan(&secId, &secName, &secDescription, &secPrice, &secOriginalPrice, &secStock, &secStartTime, &secEndTime, &secCover, &sold)
//		if err != nil {
//			return nil, errors.New("GLB-003")
//		}
//
//		productList = append(productList, &pb.SecListItem{
//			SecId:            secId,
//			SecName:          secName,
//			SecOriginalPrice: secOriginalPrice,
//			SecPrice:         secPrice,
//			SecStock:         secStock,
//			SecCover:         secCover.String,
//			SecStartTime:     secStartTime,
//			SecEndTime:       secEndTime,
//			SecSold:          sold,
//		})
//	}
//
//	// 检查查询过程中是否有错误
//	if err := rows.Err(); err != nil {
//		return nil, errors.New("GLB-003")
//	}
//
//	return productList, nil
//}
//
//// GetSecTotalItems 获取秒杀商品总数
//func GetSecTotalItems(time int32) (int32, error) {
//	query := `SELECT COUNT(*) FROM b2c_product.products_seckill WHERE time = ?`
//	var count int32
//	err := config.MySqlClient.QueryRow(query, time).Scan(&count)
//	if err != nil {
//		return 0, errors.New("GLB-003")
//	}
//	return count, nil
//}
