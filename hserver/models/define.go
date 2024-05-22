package models

type GetIntIdReq struct {
	Id int64 `query:"id" path:"id"` // 通用int类型Id(用于在path或param接受参数)
}

type GeStringIdReq struct {
	Id int64 `query:"id" path:"id"` // 通用string类型Id(用于在path或param接受参数)
}

// Request 分页参数
type Request struct {
	Size    int `query:"size,required"`    // 页码大小，最大100
	Current int `query:"current,required"` // 页码，从1开始
}
