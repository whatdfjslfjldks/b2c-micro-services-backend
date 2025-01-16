package recommend

import "micro-services/recommend-server/internal/repository"

func GetSimProductId(targetUserId int64) (
	[]string, error) {
	simProductId, err := repository.GetSimProductId(targetUserId)
	if err != nil {
		return nil, err
	}
	return simProductId, nil
}
