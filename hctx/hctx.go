package hctx

import (
	"context"
	"fmt"
	"github.com/brian-god/go-pkg/constant"
	"github.com/google/uuid"
)

const (
	keyAccessToken = "access_token"
	KeyUserId      = "userId"
	KeyPlatform    = "platform"
	KeyToken       = "token"
	KeyRole        = "role"
	KeyTenantId    = "tenant_id"
	DeviceId       = "deviceId"
	DeviceName     = "deviceName"
	IpAddress      = "ipAddress"
	IgnoreTenantId = "ignore_tenant_Id"
	AdminRole      = "Admin"
)

func WithUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, KeyUserId, userId)
}

func GetUserId(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(KeyUserId))
}
func WithPlatform(ctx context.Context, platform string) context.Context {
	return context.WithValue(ctx, KeyPlatform, platform)
}

func GetPlatform(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(KeyPlatform))
}
func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, KeyToken, token)
}

func GetToken(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(KeyToken))
}

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, KeyRole, role)
}

func GetRole(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(KeyRole))
}
func WithTenantId(ctx context.Context, tenantId string) context.Context {
	return context.WithValue(ctx, KeyTenantId, tenantId)
}

func GetTenantId(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(KeyTenantId))
}

func WithDeviceId(ctx context.Context, deviceId string) context.Context {
	return context.WithValue(ctx, DeviceId, deviceId)
}

func GetDeviceId(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(DeviceId))
}
func WithDeviceName(ctx context.Context, deviceName string) context.Context {
	return context.WithValue(ctx, DeviceName, deviceName)
}

func GetDeviceName(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(DeviceName))
}

func WithIpAddress(ctx context.Context, addr string) context.Context {
	return context.WithValue(ctx, IpAddress, addr)
}

func GetIpAddress(ctx context.Context) string {
	return fmt.Sprintf("%v", ctx.Value(IpAddress))
}

// BuildIgnoreTenantCtx 构建忽略租户的ctx
func BuildIgnoreTenantCtx(ctx context.Context) context.Context {
	return WithTenantId(ctx, IgnoreTenantId)
}

func IsAdmin(ctx context.Context) bool {
	role := GetRole(ctx)
	switch role {
	case constant.RoleAdmin:
		return true
	case constant.RoleUser:
		return false
	default:
		return false
	}
}

// WithOperationID 配合调用openim服务的链路id
func WithOperationID(ctx context.Context) context.Context {
	return context.WithValue(ctx, constant.OperationID, uuid.New().String())
}

// WithOpUserID 操作ID
func WithOpUserID(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, constant.OpUserID, userId)
}
