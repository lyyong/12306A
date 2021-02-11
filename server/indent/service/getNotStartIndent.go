package service

import "indent/models/indent"

func GetNotStartIndent (userId int) ([]*indent.Indent, error){
	return indent.GetNotStartIndent(db, userId)
}
