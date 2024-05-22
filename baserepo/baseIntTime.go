package baserepo

type BaseIntTime struct {
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;not null;default:0;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;not null;default:0;comment:更新时间"`
	DeletedAt int64 `json:"deletedAt" gorm:"column:deleted_at;not null;default:0;comment:删除时间"`
}
