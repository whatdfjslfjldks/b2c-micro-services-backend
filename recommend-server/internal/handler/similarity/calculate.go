package similarity

import (
	"fmt"
	"log"
	"micro-services/recommend-server/internal/service/recommend"
	"micro-services/recommend-server/pkg/config"
	"strconv"
)

func CalculateSim() {
	//for {
	// 获取所有 keys
	keys, err := config.RdClient3.Keys(config.Ctx, "*").Result()
	if err != nil {
		log.Println("Error fetching keys:", err)
		//continue
	}

	// 存储所有用户的数据
	users := make(map[string]map[string]int32)

	// 获取所有用户数据
	for _, key := range keys {
		values, err := config.RdClient3.HGetAll(config.Ctx, key).Result()
		if err != nil {
			log.Println("Error fetching values for key", key, ":", err)
			continue
		}

		// 将每个 key 的哈希表值转换成 map[string]int32 类型
		userData := make(map[string]int32)
		for productId, count := range values {
			countInt, _ := strconv.Atoi(count) // 将字符串转换为 int
			userData[productId] = int32(countInt)
		}

		// 保存用户数据
		users[key] = userData
	}

	// 存储相似度大于 5 的用户对
	similarityMap := make(map[string][]string) // 存储每个用户的相似用户列表

	// 计算每对用户之间的余弦相似度，只计算一次
	for i, userAKey := range keys {
		for j := i + 1; j < len(keys); j++ {
			userBKey := keys[j]

			// 获取用户数据
			userA := users[userAKey]
			userB := users[userBKey]

			// 计算余弦相似度
			similarity := recommend.CalculateCosineSimilarity(userA, userB)
			if similarity > 0.5 {
				// 如果相似度大于 0.5，将用户对存储到 similarityMap
				fmt.Printf("Cosine Similarity between %s and %s: %f\n", userAKey, userBKey, similarity)

				// 存储相似用户对
				similarityMap[userAKey] = append(similarityMap[userAKey], userBKey)
				similarityMap[userBKey] = append(similarityMap[userBKey], userAKey)
			}
		}
	}

	// 你可以将 similarityMap 存储到 Redis 或数据库中
	// 假设你要存储在 Redis 中，可以这样做
	// 打印所有相似用户对
	for user, similarUsers := range similarityMap {
		fmt.Printf("User %s has similar users: %v\n", user, similarUsers)
	}

	// 每 10 分钟计算一次相似度
	//	time.Sleep(10 * time.Minute)
	//}
}
