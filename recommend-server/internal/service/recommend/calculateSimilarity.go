package recommend

import "math"

// CalculateCosineSimilarity 计算两个用户的余弦相似度
// TODO 相似度计算，具体计算方法可以修改
func CalculateCosineSimilarity(userA, userB map[int32]int32) float64 {
	var dotProduct, magnitudeA, magnitudeB float64

	// 计算点积
	for productId, countA := range userA {
		if countB, exists := userB[productId]; exists {
			dotProduct += float64(countA * countB)
		}
	}

	// 计算 A 和 B 的模
	for _, countA := range userA {
		magnitudeA += float64(countA * countA)
	}
	for _, countB := range userB {
		magnitudeB += float64(countB * countB)
	}

	magnitudeA = math.Sqrt(magnitudeA)
	magnitudeB = math.Sqrt(magnitudeB)

	// 返回余弦相似度
	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}
	return dotProduct / (magnitudeA * magnitudeB)
}
