package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
	"micro-services/user-server/pkg/config"
	"path/filepath"
	"time"
)

//func UploadFileToMinio(file *multipart.FileHeader, fold string, username string, bucketName string) (result bool, msg string, path string) {
//	// 打开上传的文件
//	src, err := file.Open()
//	if err != nil {
//		return false, "文件打开失败", ""
//	}
//	defer src.Close()
//
//	// 定义 MinIO 存储桶和对象名
//	objectName := fold + fmt.Sprintf("%d%s", time.Now().Unix(), filepath.Ext(file.Filename))
//
//	// 创建存储桶（如果不存在）
//	ctx := context.Background()
//	exists, err := config.MinioClient.BucketExists(ctx, bucketName)
//	if err != nil {
//		log.Println(err)
//		return false, "存储桶查找失败", ""
//	}
//	if !exists {
//		err = config.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
//		if err != nil {
//			log.Println(err)
//			return false, "存储桶创建失败", ""
//		}
//	}
//
//	// 上传文件到 MinIO
//	uploadInfo, err := config.MinioClient.PutObject(
//		ctx,
//		bucketName,
//		objectName,
//		src,
//		file.Size,
//		minio.PutObjectOptions{ContentType: "image/jpeg"},
//	)
//	if err != nil {
//		return false, "文件上传失败", ""
//	}
//
//	return true, "上传成功", uploadInfo.Key
//
//}

func UploadFileToMinio(content []byte, filename string, fold string, bucketName string) (result bool, msg string, path string) {
	// 定义 MinIO 存储桶和对象名
	objectName := fold + fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(filename))

	// 创建存储桶（如果不存在）
	ctx := context.Background()
	exists, err := config.MinioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Println(err)
		return false, "存储桶查找失败", ""
	}
	if !exists {
		err = config.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Println(err)
			return false, "存储桶创建失败", ""
		}
	}

	// 上传文件到 MinIO
	uploadInfo, err := config.MinioClient.PutObject(
		ctx,
		bucketName,
		objectName,
		bytes.NewReader(content),
		int64(len(content)),
		minio.PutObjectOptions{ContentType: "image/jpeg"},
	)
	if err != nil {
		return false, "文件上传失败", ""
	}
	fmt.Println(uploadInfo)
	return true, "上传成功", uploadInfo.Key
}
