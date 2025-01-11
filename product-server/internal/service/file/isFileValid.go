package file

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
)

// IsFileValid 文件格式检查 只读 sheet1
// 注意： 多出列没关系，但是规定列不能出错，超出的列不会读
func IsFileValid(file []byte) (
	bool, error) {
	// 使用 excelize 打开 excel
	f, err := excelize.OpenReader(bytes.NewReader(file))
	if err != nil {
		log.Printf("打开excel失败:%v", err)
		return false, errors.New("打开excel失败")
	}
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		log.Printf("excel中没有sheet")
		return false, errors.New("excel中没有sheet")
	}
	// 获取第一个 sheet 的所有行
	sheetName := sheetNames[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Printf("获取sheet失败:%v", err)
		return false, errors.New("获取sheet失败")
	}
	expectedHeaders := []string{"商品名称", "种类", "商品描述", "价格", "库存"}
	if len(rows) > 0 {
		headers := rows[0]
		for i, expectedHeaders := range expectedHeaders {
			if headers[i] != expectedHeaders {
				log.Printf("excel中第1行第%d列的值为%s,预期为%s", i+1, headers[i], expectedHeaders)
				s := fmt.Sprintf("excel中第1行第%d列的值为%s,预期为%s", i+1, headers[i], expectedHeaders)
				return false, errors.New(s)
			}
		}
	}
	// TODO 暂时不做数据类型的校验，因为如果数据量大的话响应会很慢，如果需要校验可以单独写一个接口提供校验
	return true, nil
}
