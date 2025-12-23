package tenantctx

import "context"

type keyType string

const (
	TenantIDKey keyType = "tenant_id"
)

func TenantID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(TenantIDKey).(int64)
	return id, ok
}
