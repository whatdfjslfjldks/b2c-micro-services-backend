package enums

// CategoryMap 定义商品类型
var CategoryMap = map[string]int32{
	"运动户外": 1,
	"馋嘴零食": 2,
	"潮电数码": 3,
	"服饰时尚": 4,
	"家装建材": 5,
	"办公文具": 6,
	"家居生活": 7,
	"健康美容": 8,
	"母婴用品": 9,
	"书籍音像": 10,
}

// KindMap 定义商品销售类型 普通，秒杀，预售
var KindMap = map[string]int32{
	"普通": 1,
	"秒杀": 2,
	"预售": 3,
}
