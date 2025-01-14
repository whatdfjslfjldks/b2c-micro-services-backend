package recommend

import (
	"fmt"
	"log"
	"micro-services/recommend-server/pkg/config"
	"micro-services/recommend-server/pkg/kafka/model"
)

// GetUserBehavior 获取用户在 click, browse, search, purchase 表中的行为数据
// TODO 可以用来定期同步双rm库数据
func GetUserBehavior(userId int64) ([]model.Recommend, error) {

	// TODO 不需要这个函数了，使用增量算法，每当用户产生click等操作时候，更新数据到redis中
	// 直接再kafka队列消费者里操作
	// 获取点击数据
	clickData, err := getUserActionsFromTable(userId, "click")
	if err != nil {
		return nil, err
	}

	//// 获取搜索数据
	//searchData, err := getUserActionsFromTable(userId, "search")
	//if err != nil {
	//	return nil, err
	//}

	// 获取浏览数据
	browseData, err := getUserActionsFromTable(userId, "browse")
	if err != nil {
		return nil, err
	}

	// 获取购买数据
	purchaseData, err := getUserActionsFromTable(userId, "purchase")
	if err != nil {
		return nil, err
	}

	// 合并所有行为数据
	var allBehaviorData []model.Recommend
	allBehaviorData = append(allBehaviorData, clickData...)
	//allBehaviorData = append(allBehaviorData, searchData...)
	allBehaviorData = append(allBehaviorData, browseData...)
	allBehaviorData = append(allBehaviorData, purchaseData...)

	return allBehaviorData, nil
}

func getUserActionsFromTable(userId int64, table string) ([]model.Recommend, error) {

	// TODO  改成主从数据库模式，从从数据库读取 待完成

	sql := fmt.Sprintf(`SELECT user_id, product_id, status, count, create_at, update_at
						FROM b2c_recommend.%s
						WHERE user_id = ?`, table)

	rows, err := config.MySqlClient.Query(sql, userId)
	if err != nil {
		log.Printf("Failed to get user actions from %s: %v", table, err)
		return nil, err
	}
	defer rows.Close()

	var actions []model.Recommend
	for rows.Next() {
		var action model.Recommend
		if err := rows.Scan(&action.UserId, &action.ProductId, &action.Status, &action.Count, &action.Time, &action.Time); err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		actions = append(actions, action)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Failed to iterate rows: %v", err)
		return nil, err
	}

	return actions, nil
}
