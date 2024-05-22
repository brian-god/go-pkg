package plugin

import (
	"context"
	"github.com/brian-god/go-pkg/hctx"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

const Ignore = "ignore"

type TenantPlugin struct{}

func (t *TenantPlugin) Name() string {
	return "tenant_plugin"
}

func NewTenantPlugin() *TenantPlugin {
	return &TenantPlugin{}
}

func (t *TenantPlugin) Initialize(db *gorm.DB) error {
	if err := db.Callback().Query().Before("gorm:query").Register("tenant_id:before_query", t.beforeQuery); err != nil {
		return err
	}
	if err := db.Callback().Create().Before("gorm:create").Register("tenant_id:before_create", t.beforeCarte); err != nil {
		return err
	}
	return nil
}

// 创建前
func (t *TenantPlugin) beforeCarte(db *gorm.DB) {
	ctx := db.Statement.Context
	tenantID := hctx.GetTenantId(ctx)
	if TenantIDNotNil(tenantID) && !IsIgnoreTenant(ctx) {
		if db.Statement.Schema != nil {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					rv := reflect.Indirect(db.Statement.ReflectValue.Index(i))
					field1 := db.Statement.Schema.FieldsByDBName[hctx.KeyTenantId]
					if field1 != nil {
						err := field1.Set(ctx, rv, tenantID)
						if err != nil {
							err = db.Statement.AddError(err)
							if err != nil {
								return
							}
						}
					}

				}
			case reflect.Struct:
				field := db.Statement.Schema.FieldsByDBName[hctx.KeyTenantId]
				if field != nil {
					db.Statement.SetColumn(hctx.KeyTenantId, tenantID)
				}
			default:
				hlog.Debug("before create tenant nil")
			}
		}
	}
}

// 查询前
func (t *TenantPlugin) beforeQuery(db *gorm.DB) {
	// 一些业务逻辑，拿到 tenantID，可能从 context 中
	ctx := db.Statement.Context
	if !IsIgnoreTenant(ctx) {
		tenantID := hctx.GetTenantId(ctx)
		if TenantIDNotNil(tenantID) {
			db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.Eq{Column: clause.Column{Table: db.Statement.Table, Name: hctx.KeyTenantId}, Value: tenantID},
			}})
		} else {
			db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.Eq{Column: clause.Column{Table: db.Statement.Table, Name: hctx.KeyTenantId}, Value: ""},
			}})
		}
	}
}

// TenantIDNotNil 租户id是否为空
func TenantIDNotNil(tenantID string) bool {
	return tenantID != "" && tenantID != "<nil>" && tenantID != "0"
}

// GetCtxTenantID 获取租户ID
func GetCtxTenantID(ctx context.Context) string {
	tenantID := hctx.GetTenantId(ctx)
	if TenantIDNotNil(tenantID) {
		return tenantID
	}
	return ""
}

// IsIgnoreTenant 判断是否忽略租户
func IsIgnoreTenant(ctx context.Context) bool {
	return hctx.GetTenantId(ctx) == hctx.IgnoreTenantId
}

// AddTenantWhere 添加租户条件
func AddTenantWhere(ctx context.Context, db *gorm.DB, wStr string) *gorm.DB {
	tenantID := hctx.GetTenantId(ctx)
	if TenantIDNotNil(tenantID) {
		db = db.Where(wStr, tenantID)
	} else {
		db = db.Where(wStr, "")
	}
	return db
}
