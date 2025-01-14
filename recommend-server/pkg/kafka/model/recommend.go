package model

type Recommend struct {
	UserId    int64  `json:"user_id"`
	ProductId int32  `json:"product_id"` // 利用指针可以显式表明是否有值，无值nil
	Keyword   string `json:"keyword"`
	Status    string `json:"status"` // click,purchase,browse,search...
	Time      string `json:"time"`   // 更新时间,如果create_at存在就是更新时间，不存在就是创建时间

	// 数量，但是不用用户传入，测试用
	Count int32 `json:"count"`
}
