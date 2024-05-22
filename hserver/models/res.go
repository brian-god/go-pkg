package models

type PageRes[T any] struct {
	Total int64 `json:"total"`
	List  []*T  `json:"list"`
}

func NewPageRes[T any](total int64, list []*T) *PageRes[T] {
	return &PageRes[T]{
		Total: total,
		List:  list,
	}
}
