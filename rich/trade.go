package rich

import (
	"time"

	"github.com/xuender/goutils"
)

type Trade struct {
	Id     goutils.ID  `json:"id"`     // 主键
	Status TradeStatus `json:"status"` // 状态
	Ca     time.Time   `json:"ca"`     // 创建时间
	Pa     time.Time   `json:"pa"`     // 付款时间
	Orders []Order     `json:"orders"` // 订单详情
	Total  int64       `json:"total"`  // 总额,单位分
}

// 创建订单
func NewTrade(orders []Order) *Trade {
	t := &Trade{
		Id:     goutils.NewId(TradeIdPrefix),
		Ca:     time.Now(),
		Status: 预订,
		Orders: orders,
		Total:  0,
	}
	for _, o := range orders {
		t.Total += o.Price * o.Num
	}
	return t
}
