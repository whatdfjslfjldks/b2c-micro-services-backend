package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"micro-services/order-server/pkg/config"
	"micro-services/pkg/utils"
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

func ReverseOrderStatus(orderId string, orderStatus int32) error {
	result, err := config.MySqlClient.Exec("UPDATE b2c_order.orders SET order_status = ? WHERE order_id = ?", orderStatus, orderId)
	if err != nil {
		return err
	}
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}
