package repository

import (
	"fmt"
	"micro-services/pkg/utils"
	"micro-services/risk-server/pkg/config"
)

func SaveLoginInfoInToMysql(userId int64, ip string, agent string, status string) error {
	query := `
        INSERT INTO b2c_risk.risk_login(user_id, login_ip, login_agent, login_status, create_at, count, update_at)
        VALUES(?,?,?,?,?, 1, ?)
        ON DUPLICATE KEY UPDATE 
            login_status = VALUES(login_status),
            count = count + 1,
            update_at = VALUES(update_at),
            create_at = COALESCE(create_at, VALUES(create_at))  -- 仅在首次插入时设置 create_at
    `
	_, err := config.MySqlClient.Exec(query, userId, ip, agent, status, utils.GetTime(), utils.GetTime())
	if err != nil {
		fmt.Printf("could not insert to risk_login: %v", err)
		return err
	}
	return nil
}
