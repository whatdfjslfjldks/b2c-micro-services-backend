package repository

import (
	"context"
	"errors"
	"micro-services/order-server/pkg/config"
	"time"
)

func SaveAliPayQRCode(orderId string, qrCode string) error {
	result := config.RdClient.Set(context.Background(), orderId, qrCode, 35*time.Minute)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func GetAliPayQRCode(orderId string) (string, error) {
	//fmt.Println("订单编号: ", orderId)
	result := config.RdClient.Get(context.Background(), orderId)
	if result.Err() != nil {
		return "", result.Err()
	}
	return result.Val(), nil
}

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
