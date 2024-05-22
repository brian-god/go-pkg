package baserepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/brian-god/go-pkg/db"
	"gorm.io/gorm"
	"reflect"
	"time"
)

// IBaseRepo [T interface{}] ， 基础数据层方法
type IBaseRepo[T IModel] interface {
	FindPage(context.Context, *db.ListRequest) (int64, []*T, error) // 上下文/请求参数
	FindById(context.Context, interface{}) (*T, error)              // 根据 id 获取模型
	DelById(context.Context, interface{}) error                     // 根据 id 删除
	DelByIds(context.Context, interface{}) error                    // 根据 id 批量删除
	DelByIdUnScoped(context.Context, interface{}) error             // 根据 id 物理删除（可单个可批量）
	EditById(context.Context, interface{}, interface{}) error       // 上下文/id/需要更新的数据模型或者map
	Add(context.Context, *T) (*T, error)                            // 创建并返回模型
	Count(context.Context, []*db.WhereOption) (int64, error)        // 统计数量
}

type IModel interface {
	GetPrimaryKey() string // 获取主键
}

// BaseRepo [T interface{}] ， 基础数据层方法
type BaseRepo[T IModel] struct {
	Model T        // 模型
	DB    *gorm.DB // 数据库连接
}

// FindPage ， 获取模型列表
// 参数：
//
//	ctx ： 上下文
//	params ： desc
//
// 返回值：
//
//	int64 ：desc
//	[]*T ：desc
//	error ：desc
func (r *BaseRepo[T]) FindPage(ctx context.Context, params *db.ListRequest) (int64, []*T, error) {
	var res []*T

	// 查询用户
	resDb := r.DB.WithContext(ctx).
		Model(r.Model)
	//Order("sort")

	// 查询记录条数
	var count int64
	countDb := r.DB.WithContext(ctx).
		Model(&r.Model)

	if params.Wheres != nil && len(params.Wheres) > 0 {
		for _, v := range params.Wheres {
			resDb.Where(v.Sql, v.Params...)
			countDb.Where(v.Sql, v.Params...)
		}
	}

	if params.Select != nil && len(params.Select) > 0 {
		resDb.Select(params.Select)
	}

	if params.No > 0 && params.Size > 0 {
		resDb.Offset(int((params.No - 1) * params.Size)).
			Limit(int(params.Size))

		if err := countDb.Count(&count).Error; err != nil {
			return 0, nil, err
		}
	}

	if params.Order != nil && len(params.Order) > 0 {
		for _, v := range params.Order {
			resDb.Order(v)
		}
	}

	if err := resDb.Find(&res).Error; err != nil {
		return 0, nil, err
	}

	return count, res, nil
}

// FindById ， 根据 id 获取模型
// 参数：
//
//	ctx ： 上下文
//	id ： 模型 id
//
// 返回值：
//
//	*T ：desc
//	error ：desc
func (r *BaseRepo[T]) FindById(ctx context.Context, id interface{}) (*T, error) {
	var res T
	resDb := r.DB.WithContext(ctx).
		Model(r.Model)

	v := reflect.ValueOf(r.Model)
	if v.FieldByName("DeletedAt").IsValid() {
		resDb.Where("deleted_at = 0")
	}

	//根据id查询
	if err := resDb.First(&res, id).Error; err != nil {
		//if gorm.ErrRecordNotFound == err {
		//	return nil, nil
		//}
		return nil, err
	}

	return &res, nil
}

// DelById ， 根据 id 删除
// 参数：
//
//	ctx ： 上下文
//	id ： 模型 id
//
// 返回值：
//
//	error ：desc
func (r *BaseRepo[T]) DelById(ctx context.Context, id interface{}) error {

	//if r.PK != "" {
	//	// todo 可以用反射
	//	return gorm.ErrInvalidField
	//}

	//db :=r.DB.WithContext(ctx).Model(r.Model).Where(fmt.Sprintf("%v = ?",r.PK),id)

	db := r.DB.WithContext(ctx)

	v := reflect.ValueOf(r.Model)
	if v.FieldByName("DeletedAt").IsValid() {

		if r.Model.GetPrimaryKey() == "" {
			// todo 可以用反射
			return errors.New("base repo model pk is not defined")
		}

		db.Model(r.Model).Where(fmt.Sprintf("%v = ?", r.Model.GetPrimaryKey()), id).Update("deleted_at", time.Now().Unix())
	} else {
		db.Delete(&r.Model, id)
	}

	if err := db.Error; err != nil {
		return err
	}

	return nil
}

// DelByIdUnScoped ， 根据 id 删除(物理删除)
// 参数：
//
//	ctx ： 上下文
//	id ： 模型 id
//
// 返回值：
//
//	error ：desc
func (r *BaseRepo[T]) DelByIdUnScoped(ctx context.Context, id interface{}) error {
	return r.DB.WithContext(ctx).Unscoped().Delete(&r.Model, id).Error
}

// DelByIds ， 根据 id 批量删除
// 参数：
//
//	ctx ： 上下文
//	ids ： 模型 id 数组
//
// 返回值：
//
//	error ：desc
func (r *BaseRepo[T]) DelByIds(ctx context.Context, ids interface{}) error {

	db := r.DB.WithContext(ctx)

	v := reflect.ValueOf(r.Model)
	if v.FieldByName("DeletedAt").IsValid() {

		if r.Model.GetPrimaryKey() == "" {
			return errors.New("base repo model pk is not defined")
		}

		db.Model(r.Model).Where(fmt.Sprintf("%v in ?", r.Model.GetPrimaryKey()), ids).Update("deleted_at", time.Now().Unix())
	} else {
		db.Delete(&r.Model).Where(fmt.Sprintf("%v in ?", r.Model.GetPrimaryKey()), ids)
	}

	if err := db.Error; err != nil {
		return err
	}

	return nil
}

// EditById ， 根据 id 更新 模型
// 参数：
//
//	ctx ： desc
//	id ： desc
//	data ： desc
//
// 返回值：
//
//	error ：desc
func (r *BaseRepo[T]) EditById(ctx context.Context, id interface{}, data interface{}) error {
	if r.Model.GetPrimaryKey() == "" {
		// todo 可以用反射
		return errors.New("base repo model pk is not defined")
	}
	newCtx, cancelFunc := context.WithTimeout(ctx, 30*time.Second)
	defer func() {
		cancelFunc()
	}()
	db := r.DB.WithContext(newCtx).
		Model(r.Model)

	v := reflect.ValueOf(r.Model)
	if v.FieldByName("DeletedAt").IsValid() {
		db.Where("deleted_at = 0")
	}
	v = reflect.ValueOf(data)
	if v.Kind() != reflect.Map {
		updated := v.Elem().FieldByName("UpdatedAt")
		if updated.IsValid() {
			updated.SetInt(time.Now().Unix())
		}
	}
	if err := db.Where(fmt.Sprintf("%v = ?", r.Model.GetPrimaryKey()), id).
		Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// Add ， 创建模型
// 参数：
//
//	ctx ： 上下文
//	data ： 模型数据
//
// 返回值：
//
//	error ：desc
func (r *BaseRepo[T]) Add(ctx context.Context, data *T) (*T, error) {
	v := reflect.ValueOf(data)
	created := v.Elem().FieldByName("CreatedAt")
	if created.IsValid() {
		created.SetInt(time.Now().Unix())
	}
	if err := r.DB.WithContext(ctx).Create(data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// Count , 获取记录总数通用方法
// 参数：
//
//	     ctx ： 上下文
//	     db ： gorm.DB
//	     wheres ： 筛选条件
//			model: 模型类
//
// 返回值：
//
//	int64 ：记录总数
//	error ：错误
func (r *BaseRepo[T]) Count(ctx context.Context, wheres []*db.WhereOption) (int64, error) {
	var count int64
	countDb := r.DB.WithContext(ctx).
		Model(r.Model)

	if len(wheres) > 0 {
		for _, v := range wheres {
			countDb.Where(v.Sql, v.Params...)
		}
	}

	if err := countDb.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
