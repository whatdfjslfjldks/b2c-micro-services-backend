package repository

import (
	"log"
	"micro-services/recommend-server/pkg/config"
	"micro-services/recommend-server/pkg/kafka/model"
)

// SaveClickMsgIntoMysql 调用通用的 SaveMsgIntoMysql 函数
func SaveClickMsgIntoMysql(message model.Recommend) error {
	return SaveMsgIntoMysql("click", message)
}

// SaveBrowseMsgIntoMysql 调用通用的 SaveMsgIntoMysql 函数
func SaveBrowseMsgIntoMysql(message model.Recommend) error {
	return SaveMsgIntoMysql("browse", message)
}

// SavePurchaseMsgIntoMysql 调用通用的 SaveMsgIntoMysql 函数
func SavePurchaseMsgIntoMysql(message model.Recommend) error {
	return SaveMsgIntoMysql("purchase", message)
}

// SaveSearchMsgIntoMysql 调用通用的 SaveMsgIntoMysql 函数
func SaveSearchMsgIntoMysql(message model.Recommend) error {
	sql := `insert into b2c_recommend.search (user_id, keyword, status, count, create_at, update_at) 
            values(?, ?, ?, 1, ?, ?)
            ON DUPLICATE KEY UPDATE 
                user_id = VALUES(user_id),
                keyword = VALUES(keyword),
                status = VALUES(status),
                count = count + 1,
                update_at = VALUES(update_at),
                create_at = COALESCE(create_at, VALUES(create_at))`

	// 执行 SQL 语句
	_, err := config.MySqlClient.Exec(sql, message.UserId, message.Keyword, message.Status, message.Time, message.Time)
	if err != nil {
		log.Printf("Failed to insert into MySQL table search: %v", err)
		return err
	}
	return nil
}

func SaveMsgIntoMysql(tableName string, message model.Recommend) error {
	sql := `insert into b2c_recommend.` + tableName + `(user_id, product_id, status, count, create_at, update_at) 
            values(?, ?, ?, 1, ?, ?)
            ON DUPLICATE KEY UPDATE 
                user_id = VALUES(user_id),
                product_id = VALUES(product_id),
                status = VALUES(status),
                count = count + 1,
                update_at = VALUES(update_at),
                create_at = COALESCE(create_at, VALUES(create_at))`

	// 执行 SQL 语句
	_, err := config.MySqlClient.Exec(sql, message.UserId, message.ProductId, message.Status, message.Time, message.Time)
	if err != nil {
		log.Printf("Failed to insert into MySQL table %s: %v", tableName, err)
		return err
	}
	return nil
}
