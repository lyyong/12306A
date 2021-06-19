package cache

import "fmt"

type CandidateCache struct {
}

func (cc CandidateCache) GetKeyByTrainIDAndDate(TrainID uint, Date string) string {
	return fmt.Sprintf("candidate-%d-%s", TrainID, Date)
}

func (cc CandidateCache) GetKeyByUserID(UserID uint) string {
	return fmt.Sprintf("candidate-%d", UserID)
}

func (cc CandidateCache) GetKeyByOrderIDUnPay(OrderID string) string {
	return fmt.Sprintf("candidate-unpay-%s", OrderID)
}

func (cc CandidateCache) GetTrainIDSCacheKey() string {
	return "CAN-trainIDS"
}
