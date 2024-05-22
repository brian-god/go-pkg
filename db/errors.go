package db

import (
	"errors"
	"gorm.io/gorm"
)

func IfErrorNotNotFound(err error) bool {
	return err != nil && !errors.Is(err, gorm.ErrRecordNotFound)
}
func ErrorIsRecordNotFound(err error) bool {
	return err != nil && !errors.Is(err, gorm.ErrRecordNotFound)
}
func IfErrorNotFound(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrRecordNotFound)
}

// IsUniqueIndexError ， 判断是否为索引错误
// 参数：
//
//	err ： desc
//
// 返回值：
//
//	bool ：desc
//func IsUniqueIndexError(err error) bool {
//	errType := reflect.TypeOf(err).String()
//	if errType == "*mysql.MySQLError" && err.(*mysql.MySQLError).Number == 1062 {
//		return true
//	}
//	return false
//}
