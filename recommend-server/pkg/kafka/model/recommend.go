package model

type Recommend struct {
	UserId    int64  `json:"user_id"`
	ProductId int32  `json:"product_id"`
	Status    string `json:"status"` // click,purchase,browse,search...
	Time      string `json:"time"`   // 更新时间,如果create_at存在就是更新时间，不存在就是创建时间
}
