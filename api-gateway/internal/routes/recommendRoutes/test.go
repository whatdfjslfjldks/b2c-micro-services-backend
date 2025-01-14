package recommendRoutes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/mat"
)

func Test(c *gin.Context) {
	//// 创建原始矩阵 A
	//data := []float64{
	//	5, 4, 0, 0, 1,
	//	4, 0, 0, 2, 3,
	//	1, 0, 0, 5, 4,
	//	0, 3, 4, 0, 5,
	//}
	//// 转换为 gonum 中的矩阵对象
	//A := mat.NewDense(4, 5, data)
	//
	//// 调用 SVD 函数来分解 A
	//var svd mat.SVD
	//ok := svd.Factorize(A, mat.SVDThin)
	//if !ok {
	//	log.Fatal("SVD factorization failed")
	//}
	//
	//// 获取左奇异矩阵 U
	//U := mat.NewDense(A.Dims()) // 创建一个新矩阵来存储 U
	//svd.UTo(U)                  // 提取 U 矩阵
	//
	//// 获取奇异值 Σ，返回的是一个切片
	//Sigma := svd.Values(nil)
	//
	//// 获取右奇异矩阵 V^T
	//Vt := mat.NewDense(A.Dims()) // 创建一个新矩阵来存储 V^T
	//svd.VTo(Vt)                  // 提取 V^T 矩阵
	//
	//// 打印结果
	//fmt.Println("Matrix A:")
	//matPrint(A)
	//
	//fmt.Println("\nMatrix U:")
	//matPrint(U)
	//
	//fmt.Println("\nSingular Values (Σ):")
	//fmt.Println(Sigma)
	//
	//fmt.Println("\nMatrix V^T:")
	//matPrint(Vt)
}

// 打印矩阵的辅助函数
func matPrint(m mat.Matrix) {
	rows, cols := m.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Printf("%.2f ", m.At(i, j))
		}
		fmt.Println()
	}
}
