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
