package consts

type Status int

const (
	InOut Status = iota
	InEnum
	InStruct
	InConst
)
