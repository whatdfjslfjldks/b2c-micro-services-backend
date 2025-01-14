package recommend

import (
	"micro-services/recommend-server/pkg/kafka/model"
)

// BuildUserBehaviorVector 构建用户的行为向量
// TODO 可以用来定期同步双rm库数据
func BuildUserBehaviorVector(userBehaviorData []model.Recommend) map[int32]int32 {
	// 使用一个 map 来表示每个用户对每个商品的行为强度
	userBehaviorVector := make(map[int32]int32)

	for _, behavior := range userBehaviorData {
		// 点击赋予权重 1，浏览赋予权重 2，购买赋予权重 3
		//fmt.Println("behavior: ", behavior.Status)
		weight := 0
		if behavior.Status == "CLICK" {
			weight = 1
		} else if behavior.Status == "BROWSE" {
			weight = 2
		} else if behavior.Status == "PURCHASE" {
			weight = 3
		}

		// 计算每个商品的行为强度
		userBehaviorVector[behavior.ProductId] += behavior.Count * int32(weight)
	}

	return userBehaviorVector
}
