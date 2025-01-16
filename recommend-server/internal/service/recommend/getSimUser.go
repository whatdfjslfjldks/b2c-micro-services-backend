package recommend

import "micro-services/recommend-server/internal/repository"

func GetSimUserId(userId int64) (
	[]string, error) {
	simUserId, err := repository.GetSimUserId(userId)
	if err != nil {
		return nil, err
	}
	return simUserId, nil
}
