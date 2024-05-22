package db

import "gorm.io/gorm"

// Page 分页参数
type Page struct {
	Size int64 `json:"size" query:"size"` // 页码大小，最大100
	No   int64 `json:"no" query:"no"`     // 页码，从1开始
}

// Response 分页响应
type Response struct {
	Page
	Total int64       `json:"total"` // 总条数
	Data  interface{} `json:"data"`  // 列表数据
}

type ListRequest struct {
	Page
	Wheres []*WhereOption `json:"wheres"` // 条件
	Order  []string       `json:"order"`  // 排序
	Select []string       `json:"select"` // 需要返回的字段，不传则返回全部
}

type WhereOption struct {
	Sql    string        `json:"sql"`
	Params []interface{} `json:"params"`
}

func NewWhereOption(sql string, params ...interface{}) *WhereOption {
	return &WhereOption{
		Sql:    sql,
		Params: params,
	}
}

func (s *ListRequest) WithWhere(sql string, params ...interface{}) *ListRequest {
	s.Wheres = append(s.Wheres, &WhereOption{
		Sql:    sql,
		Params: params,
	})

	return s
}

func (s *ListRequest) WithSelect(params ...string) *ListRequest {
	s.Select = append(s.Select, params...)
	return s
}

func (s *ListRequest) WithOrder(params ...string) *ListRequest {
	s.Order = append(s.Order, params...)
	return s
}

func (s *ListRequest) WithEqual(col string, val interface{}) *ListRequest {
	return s.WithWhere(col+" = ?", val)
}

func (s *ListRequest) WithLike(column string, param string) *ListRequest {
	return s.WithWhere(column+" like ?", "%"+param+"%")
}

func (s *ListRequest) WithNotDelete() *ListRequest {
	return s.WithWhere("deleted_at = 0")
}

func (r *Page) Fix() {
	if r.No <= 0 {
		r.No = 1
	}

	if r.Size <= 0 {
		r.Size = 10
	} else if r.Size > 100 {
		r.Size = 100
	}
}
func Operation(r *Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		r.Fix()
		offset := int((r.No - 1) * r.Size)
		return db.Offset(offset).Limit(int(r.Size))
	}
}
