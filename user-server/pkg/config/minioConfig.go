package config

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

var MinioClient *minio.Client

// InitMinio 初始化 MinIO 客户端
func InitMinio() {
	// 创建 MinIO 客户端实例
	var err error
	MinioClient, err = minio.New("127.0.0.1:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""), // 使用静态凭证
		Secure: false,                                                   // 设置为 false 表示使用 HTTP，true 表示使用 HTTPS
	})
	if err != nil {
		log.Fatalln("Failed to initialize MinIO client:", err)
	}
}
