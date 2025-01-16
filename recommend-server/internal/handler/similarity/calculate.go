package similarity

import (
	"log"
	"micro-services/recommend-server/internal/repository"
	"micro-services/recommend-server/internal/service/recommend"
	"micro-services/recommend-server/pkg/config"
	"strconv"
	"time"
)

func CalculateSim() {
	for {
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
				// TODO 数据量少，相似度可以调整低一些，后面可以修改
				if similarity > 0.35 {
					//fmt.Printf("Cosine Similarity between %s and %s: %f\n", userAKey, userBKey, similarity)

					// 存储相似用户对
					similarityMap[userAKey] = append(similarityMap[userAKey], userBKey)
					similarityMap[userBKey] = append(similarityMap[userBKey], userAKey)
				}
			}
		}

		// 将 similarityMap 存储到 Redis 或 MySQL数据库中
		for user, similarUsers := range similarityMap {
			//fmt.Printf("User %s has similar users: %v\n", user, similarUsers)
			e := repository.SaveSimUserInToRedis(user, similarUsers)
			if e != nil {
				log.Printf("Error saving similarity map to Redis: %v\n", e)
			}
		}

		// 每 10 分钟计算一次相似度
		time.Sleep(10 * time.Minute)
	}
}
