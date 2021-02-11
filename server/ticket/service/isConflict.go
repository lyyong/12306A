package service

import "ticket/models/ticket"

func IsConflict(passengerId *[]int32) (bool, error) {
	return ticket.IsConflict(db,passengerId)
}
