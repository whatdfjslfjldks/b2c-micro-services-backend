package file

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"micro-services/product-server/internal/repository"
	"micro-services/product-server/pkg/enums"
	"strconv"
)

func UploadProduct(file []byte) error {
	// 使用 excelize 打开 excel
	f, err := excelize.OpenReader(bytes.NewReader(file))
	if err != nil {
		log.Printf("打开excel失败:%v", err)
		return errors.New("打开excel失败")
	}

	sheetName := f.GetSheetName(0)
	//fmt.Println("sheetName:", sheetName)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Printf("获取sheet失败:%v", err)
		return errors.New("获取sheet失败")
	}

	// 批量存储产品
	var productBatch []repository.Product

	// 读取每一行
	for index, row := range rows {
		// 跳过第一行表头
		if index == 0 {
			continue
		}

		// 检查行数据是否有效
		if len(row) > 0 {
			// 1. 商品名称
			productName := row[0]
			// 2. 种类
			category := row[1]
			// 对应种类表把种类转化成int32
			c, b := enums.CategoryMap[category]
			if !b {
				info := fmt.Sprintf("种类不存在:%s", category)
				return errors.New(info)
			}
			// 3. 商品描述
			description := row[2]
			// 4. 价格
			price := row[3]
			p, e := strconv.ParseFloat(price, 64)
			if e != nil {
				log.Printf("价格转换失败:%v", e)
				info := fmt.Sprintf("价格转换失败:%v", e)
				return errors.New(info)
			}
			// 5. 库存
			stock := row[4]
			s, e := strconv.ParseInt(stock, 10, 32)
			if e != nil {
				log.Printf("库存转换失败:%v", e)
				info := fmt.Sprintf("库存转换失败:%v", e)
				return errors.New(info)
			}

			// 将解析后的数据存入批量插入数组
			productBatch = append(productBatch, repository.Product{
				Name:        productName,
				Category:    c,
				Description: description,
				Price:       p,
				Stock:       int32(s),
			})

			// 每100条数据批量插入一次
			if len(productBatch) >= 100 {
				err := repository.SaveProductsBatch(productBatch)
				if err != nil {
					return err
				}
				// 清空批量数据
				productBatch = nil
			}
		}
	}

	// 如果剩余数据不足100条，仍然执行批量插入
	if len(productBatch) > 0 {
		err := repository.SaveProductsBatch(productBatch)
		if err != nil {
			return err
		}
	}

	return nil
}
