package repository

import (
	"errors"
	"micro-services/pay-server/pkg/config"
)

func DeleteAliPayQRCode(orderId string) error {
	result := config.RdClient.Del(config.Ctx, orderId)
	if result.Err() != nil {
		return result.Err()
	}
	if result.Val() == 0 {
		return errors.New("订单不存在！")
	}
	return nil
}
