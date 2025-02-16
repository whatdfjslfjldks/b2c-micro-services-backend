package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/google/uuid"
	"log"
	"micro-services/order-server/pkg/config"
	pb "micro-services/pkg/proto/order-server"
	"micro-services/pkg/utils"
	"strings"
	"time"
)

func IsUserIdExist(userId int64) bool {
	var exists bool
	err := config.MySqlClient.QueryRow("SELECT EXISTS(SELECT 1 FROM b2c_user.users WHERE user_id = ?)", userId).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
func CalcTotalPrice(productId []int32, productAmount []int32) (float64, error) {
	if len(productId) != len(productAmount) {
		return 0, fmt.Errorf("the lengths of productId and productAmount must be the same")
	}

	// 生成占位符
	placeholders := make([]string, len(productId))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	// 正确拼接占位符字符串
	placeholderStr := ""
	for i, ph := range placeholders {
		if i > 0 {
			placeholderStr += ", "
		}
		placeholderStr += ph
	}
	query := fmt.Sprintf("SELECT product_id, price FROM b2c_product.product_price WHERE product_id IN (%s)", placeholderStr)

	// 准备参数
	args := make([]interface{}, len(productId))
	for i, id := range productId {
		args[i] = id
	}

	// 执行批量查询
	rows, err := config.MySqlClient.Query(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to query product prices: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}(rows)

	// 存储商品价格的映射
	priceMap := make(map[int32]float64)
	for rows.Next() {
		var id int32
		var price float64
		if err := rows.Scan(&id, &price); err != nil {
			return 0, fmt.Errorf("failed to scan product price: %w", err)
		}
		priceMap[id] = price
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("error occurred while iterating over rows: %w", err)
	}

	// 计算总价
	var totalPrice float64
	for i, id := range productId {
		price, exists := priceMap[id]
		if !exists {
			return 0, fmt.Errorf("price not found for product %d", id)
		}
		totalPrice += price * float64(productAmount[i])
	}

	return totalPrice, nil
}

// CreateOrder 创建订单
func CreateOrder(userId int64, address string, detail string, name string, phone string, note string, productId []int32, typeName []string, productAmount []int32, totalPrice float64) (string, error) {
	// 生成订单编号
	orderId := uuid.New().String()
	// 开启事务
	tx, err := config.MySqlClient.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func(tx *sql.Tx) {
		if p := recover(); p != nil {
			// 回滚事务
			err := tx.Rollback()
			if err != nil {
				log.Printf("Failed to rollback transaction: %v", err)
			}
			panic(p)
		} else if err != nil {
			// 回滚事务
			err := tx.Rollback()
			if err != nil {
				log.Printf("Failed to rollback transaction: %v", err)
			}
		}
	}(tx)
	// 插入订单信息到orders表
	_, err = tx.Exec("INSERT INTO b2c_order.orders (order_id, user_id, order_status, create_at, update_at) VALUES (?, ?, ?, ?, ?)",
		orderId, userId, 0, utils.GetTime(), utils.GetTime())
	if err != nil {
		return "", fmt.Errorf("failed to insert order: %w", err)
	}
	// 插入支付信息到order_payments表
	_, err = tx.Exec("INSERT INTO b2c_order.order_payments (order_id, payment_method, payment_price, payment_status) VALUES (?, ?, ?, ?)",
		orderId, 0, totalPrice, 0)
	if err != nil {
		return "", fmt.Errorf("failed to insert order payment: %w", err)
	}
	// 将切片转换为JSON字符串
	productIdJSON, err := json.Marshal(productId)
	if err != nil {
		return "", fmt.Errorf("failed to marshal productId: %w", err)
	}
	typeNameJSON, err := json.Marshal(typeName)
	if err != nil {
		return "", fmt.Errorf("failed to marshal typeName: %w", err)
	}
	productAmountJSON, err := json.Marshal(productAmount)
	if err != nil {
		return "", fmt.Errorf("failed to marshal productAmount: %w", err)
	}
	// 插入订单项到order_items表
	_, err = tx.Exec("INSERT INTO b2c_order.order_items (order_id, product_ids, type_names, product_amounts, total_price) VALUES (?, ?, ?, ?, ?)",
		orderId, string(productIdJSON), string(typeNameJSON), string(productAmountJSON), totalPrice)
	if err != nil {
		return "", fmt.Errorf("failed to insert order item: %w", err)
	}
	// 插入订单地址信息到order_addresses表
	_, err = tx.Exec("INSERT INTO b2c_order.order_details (order_id, address, detail, name, phone, note, create_at, update_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		orderId, address, detail, name, phone, note, utils.GetTime(), utils.GetTime())
	if err != nil {
		return "", fmt.Errorf("failed to insert order address: %w", err)
	}
	// 提交事务
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}
	return orderId, nil
}

func ReverseOrderStatus(orderId string, orderStatus int32, paymentStatus int32) error {
	result, err := config.MySqlClient.Exec("UPDATE b2c_order.orders SET order_status = ? WHERE order_id = ?", orderStatus, orderId)
	if err != nil {
		return err
	}
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	rs, err := config.MySqlClient.Exec("UPDATE b2c_order.order_payments SET payment_status = ? , payment_date=? WHERE order_id = ?", paymentStatus, utils.GetTime(), orderId)
	if err != nil {
		return err
	}
	if rowsAffected, _ := rs.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func GetOrderDetail(orderId string, userId int64) (*pb.GetOrderDetailResponse, error) {
	// 初始化返回对象
	orderDetail := &pb.GetOrderDetailResponse{}

	// 查询订单信息
	var orderStatus, paymentMethod, paymentStatus int32
	var totalPrice float64
	var orderDate string
	var ud int64
	err := config.MySqlClient.QueryRow(`
		SELECT o.user_id,o.order_status, op.payment_method, op.payment_status, op.payment_price, o.create_at
		FROM b2c_order.orders o
		JOIN b2c_order.order_payments op ON o.order_id = op.order_id
		WHERE o.order_id = ?`, orderId).Scan(&ud, &orderStatus, &paymentMethod, &paymentStatus, &totalPrice, &orderDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query order information: %w", err)
	}
	if userId != ud {
		return nil, fmt.Errorf("userId does not match")
	}

	// 获取订单地址信息
	var name, phone, address, detail, note string
	err = config.MySqlClient.QueryRow(`
		SELECT name, phone, address, detail, note
		FROM b2c_order.order_details
		WHERE order_id = ?`, orderId).Scan(&name, &phone, &address, &detail, &note)
	if err != nil {
		return nil, fmt.Errorf("failed to query order address information: %w", err)
	}

	// 获取订单商品信息
	var productIdJSON, typeNameJSON, productAmountJSON string
	err = config.MySqlClient.QueryRow(`
		SELECT product_ids, type_names, product_amounts
		FROM b2c_order.order_items
		WHERE order_id = ?`, orderId).Scan(&productIdJSON, &typeNameJSON, &productAmountJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items: %w", err)
	}

	// 解析商品相关信息
	var productIds []int32
	var typeNames []string
	var productAmounts []int32

	err = json.Unmarshal([]byte(productIdJSON), &productIds)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal productIds: %w", err)
	}
	err = json.Unmarshal([]byte(typeNameJSON), &typeNames)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal typeNames: %w", err)
	}
	err = json.Unmarshal([]byte(productAmountJSON), &productAmounts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal productAmounts: %w", err)
	}

	// 填充响应结构体
	orderDetail.Code = 200
	orderDetail.StatusCode = "GLB-000"
	orderDetail.Msg = "获取数据成功"
	orderDetail.OrderId = orderId
	orderDetail.OrderDate = orderDate
	orderDetail.OrderStatus = orderStatus
	orderDetail.PaymentMethod = paymentMethod
	orderDetail.PaymentStatus = paymentStatus
	orderDetail.PaymentPrice = totalPrice
	orderDetail.Name = name
	orderDetail.Phone = phone
	orderDetail.Address = address
	orderDetail.Detail = detail
	orderDetail.Note = note
	orderDetail.ProductId = productIds
	orderDetail.TypeName = typeNames
	orderDetail.ProductAmount = productAmounts

	return orderDetail, nil
}

func CheckProductStock(productId []int32, productAmount []int32) error {
	// 加锁
	pool := goredis.NewPool(config.RdClient)
	rs := redsync.New(pool)
	mutex := rs.NewMutex("normal", redsync.WithExpiry(10*time.Second)) // 设置锁的超时时间为 10 秒

	// 尝试获取锁
	if err := mutex.Lock(); err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer func() {
		if _, err := mutex.Unlock(); err != nil {
			log.Printf("failed to release lock: %v", err)
		}
	}()

	if len(productId) != len(productAmount) {
		return fmt.Errorf("the lengths of productId and productAmount must be the same")
	}

	// 生成占位符
	placeholders := make([]string, len(productId))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	placeholderStr := strings.Join(placeholders, ", ")

	// 构建查询语句，检查商品是否上架
	checkQuery := fmt.Sprintf("SELECT id FROM b2c_product.products WHERE kind = 1 AND isListed = 1 AND id IN (%s)", placeholderStr)

	// 准备参数
	args := make([]interface{}, len(productId))
	for i, id := range productId {
		args[i] = id
	}

	// 执行批量查询
	rows, err := config.MySqlClient.Query(checkQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to query product listings: %w", err)
	}
	defer rows.Close()

	// 存储已上架商品的ID
	listedProductIds := make(map[int32]struct{})
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan product listing: %w", err)
		}
		listedProductIds[id] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error occurred while iterating over rows: %w", err)
	}

	// 检查每个商品是否已上架
	for _, id := range productId {
		if _, exists := listedProductIds[id]; !exists {
			return fmt.Errorf("product %d is not listed", id)
		}
	}

	// 构建查询库存的语句
	stockQuery := fmt.Sprintf("SELECT product_id, stock FROM b2c_product.product_sale WHERE product_id IN (%s)", placeholderStr)

	// 执行批量查询库存
	rows, err = config.MySqlClient.Query(stockQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to query product stocks: %w", err)
	}
	defer rows.Close()

	// 存储商品库存的映射
	stockMap := make(map[int32]int32)
	for rows.Next() {
		var id int32
		var stock int32
		if err := rows.Scan(&id, &stock); err != nil {
			return fmt.Errorf("failed to scan product stock: %w", err)
		}
		stockMap[id] = stock
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error occurred while iterating over rows: %w", err)
	}

	// 检查库存是否足够
	var allStocksSufficient bool
	allStocksSufficient = true
	for i, id := range productId {
		stock, exists := stockMap[id]
		if !exists {
			return fmt.Errorf("stock not found for product %d", id)
		}
		if stock < productAmount[i] {
			return fmt.Errorf("insufficient stock for product %d, required: %d, available: %d", id, productAmount[i], stock)
		}
		// 如果库存足够，减去相应的数量
		stockMap[id] -= productAmount[i]
		if stockMap[id] < 0 {
			allStocksSufficient = false
		}
	}

	// 如果所有库存都足够，开始事务更新库存
	if allStocksSufficient {
		tx, err := config.MySqlClient.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}
		defer func() {
			if !allStocksSufficient {
				_ = tx.Rollback()
			}
		}()

		for id, newStock := range stockMap {
			updateQuery := "UPDATE b2c_product.product_sale SET stock = ? WHERE product_id = ?"
			if _, err := tx.Exec(updateQuery, newStock, id); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("failed to update stock for product %d: %w", id, err)
			}
		}

		if err := tx.Commit(); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	} else {
		return fmt.Errorf("not all products have sufficient stock")
	}

	return nil
}
