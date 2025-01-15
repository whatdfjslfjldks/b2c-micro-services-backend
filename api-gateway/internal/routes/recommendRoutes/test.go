package recommendRoutes

//
//import (
//	"fmt"
//	"gonum.org/v1/gonum/mat"
//	"math"
//)
//
//// 计算余弦相似度
//func cosineSimilarity(user1, user2 []float64) float64 {
//	var dotProduct, normUser1, normUser2 float64
//	for i := 0; i < len(user1); i++ {
//		dotProduct += user1[i] * user2[i]
//		normUser1 += user1[i] * user1[i]
//		normUser2 += user2[i] * user2[i]
//	}
//	return dotProduct / (math.Sqrt(normUser1) * math.Sqrt(normUser2))
//}
//
//func Test() {
//	// 创建原始用户-商品矩阵 A
//	data := []float64{
//		5, 4, 0, 0, 1, // 用户1
//		4, 0, 0, 2, 3, // 用户2
//		1, 0, 0, 5, 4, // 用户3
//		0, 3, 4, 0, 5, // 用户4
//	}
//	A := mat.NewDense(4, 5, data) // 4 行 5 列
//
//	// 获取 A 的维度
//	//m, n := A.Dims()
//
//	// 计算用户之间的相似度
//	userSimilarities := make([][]float64, m)
//	//for i := 0; i < m; i++ {
//	//	userSimilarities[i] = make([]float64, m)
//	//	for j := 0; j < m; j++ {
//	//		if i != j {
//	//			// 计算用户 i 和 j 的相似度
//	//			userSimilarities[i][j] = cosineSimilarity(A.RowView(i), A.RowView(j))
//	//		} else {
//	//			userSimilarities[i][j] = 1 // 自身与自身的相似度为 1
//	//		}
//	//	}
//	//}
//
//	// 打印用户相似度矩阵
//	fmt.Println("User Similarities Matrix:")
//	for i := 0; i < m; i++ {
//		for j := 0; j < m; j++ {
//			fmt.Printf("%.2f ", userSimilarities[i][j])
//		}
//		fmt.Println()
//	}
//
//	// 预测用户可能喜欢的商品
//	predictRecommendations(A, userSimilarities)
//}
//
//// 预测推荐商品
//func predictRecommendations(A *mat.Dense, userSimilarities [][]float64) {
//	// 获取 A 的维度
//	m, n := A.Dims()
//
//	// 为每个用户预测商品评分
//	for i := 0; i < m; i++ {
//		fmt.Printf("Predicted ratings for User %d: ", i+1)
//
//		// 对每个商品进行评分预测
//		for j := 0; j < n; j++ {
//			if A.At(i, j) == 0 { // 只有在用户没有对该商品评分时进行预测
//				var weightedSum, similaritySum float64
//
//				// 计算基于邻居用户的加权评分
//				for k := 0; k < m; k++ {
//					if k != i && A.At(k, j) > 0 { // 只考虑评分过该商品的用户
//						weightedSum += userSimilarities[i][k] * A.At(k, j)
//						similaritySum += math.Abs(userSimilarities[i][k])
//					}
//				}
//
//				if similaritySum > 0 {
//					// 计算预测评分
//					predictedRating := weightedSum / similaritySum
//					if predictedRating > 2 { // 设置一个阈值，例如评分大于2即为推荐
//						fmt.Printf("推荐商品：%d ", j+1)
//					}
//				}
//			}
//		}
//		fmt.Println()
//	}
//}
