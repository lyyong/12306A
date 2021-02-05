package service

import "indent/models/indent"

func GetAllIndent (userId int) ([]*indent.Indent, error){
	return indent.GetAllIndent(db, userId)
}
