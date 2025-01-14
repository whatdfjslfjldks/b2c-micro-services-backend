package recommend

import (
	"micro-services/recommend-server/pkg/kafka/model"
)

// RecommendProducts 根据相似用户推荐商品
func RecommendProducts(targetUserId int64) ([]model.Recommend, error) {

	//TODO 待修改

	return nil, nil
	// 获取目标用户的行为数据
	//targetUserBehaviorData, err := GetUserBehavior(targetUserId)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// 构建目标用户的行为向量
	//targetUserBehaviorVector := BuildUserBehaviorVector(targetUserBehaviorData)
	//
	//// 获取所有用户的行为数据（可以根据实际情况限制查询的用户数量）
	//sql := `SELECT user_id FROM b2c_recommend.click UNION
	//        SELECT user_id FROM b2c_recommend.search UNION
	//        SELECT user_id FROM b2c_recommend.browse UNION
	//        SELECT user_id FROM b2c_recommend.purchase`
	//
	//rows, err := config.MySqlClient.Query(sql)
	//if err != nil {
	//	log.Printf("Failed to get users' data: %v", err)
	//	return nil, err
	//}
	//defer rows.Close()
	//
	//var allUsers []int64
	//for rows.Next() {
	//	var userId int64
	//	if err := rows.Scan(&userId); err != nil {
	//		log.Printf("Failed to scan row: %v", err)
	//		return nil, err
	//	}
	//	allUsers = append(allUsers, userId)
	//}
	//
	//// 根据相似度推荐商品
	//var recommendedProducts []model.Recommend
	//for _, userId := range allUsers {
	//	if userId != targetUserId {
	//		// 获取相似用户的行为数据
	//		userBehaviorData, err := GetUserBehavior(userId)
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		// 构建相似用户的行为向量
	//		userBehaviorVector := BuildUserBehaviorVector(userBehaviorData)
	//
	//		// 计算与目标用户的相似度
	//		similarity := CalculateCosineSimilarity(targetUserBehaviorVector, userBehaviorVector)
	//		if similarity > 0.8 { // 相似度阈值
	//			// 推荐该用户的商品
	//			for productId := range userBehaviorVector {
	//				// 你可以根据需求筛选商品，比如排除已经推荐过的商品
	//				recommendedProducts = append(recommendedProducts, model.Recommend{
	//					UserId:    targetUserId,
	//					ProductId: productId,
	//					Status:    "recommended",
	//					Time:      utils.GetTime(),
	//				})
	//			}
	//		}
	//	}
	//}
	//
	//return recommendedProducts, nil
}
