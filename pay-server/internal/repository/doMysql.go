package repository

import (
	"fmt"
	"micro-services/pay-server/pkg/config"
	"micro-services/pkg/utils"
)

func ReverseOrderStatus(orderId string, orderStatus int32, paymentStatus int32) error {
	result, err := config.MySqlClient.Exec("UPDATE b2c_order.orders SET order_status = ? WHERE order_id = ?", orderStatus, orderId)
	if err != nil {
		return err
	}
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	rs, err := config.MySqlClient.Exec("UPDATE b2c_order.order_payments SET payment_status = ? , payment_date=? WHERE order_id = ?", paymentStatus, utils.GetTime(), orderId)
	if err != nil {
		return err
	}
	if rowsAffected, _ := rs.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}
