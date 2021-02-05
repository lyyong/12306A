package service

import "ticket/models/ticket"

func IsConflict(passengerId *[]int32) bool {
	return ticket.IsConflict(db,passengerId)
}
