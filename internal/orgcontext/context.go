package orgcontext

import "context"

// OrgContextKey is the request context key for the active organization ID.
type OrgContextKey struct{}

// WithOrgID stores the org ID in the context.
func WithOrgID(ctx context.Context, orgID int64) context.Context {
	return context.WithValue(ctx, OrgContextKey{}, orgID)
}

// OrgIDFromContext returns the org ID from context, if set.
func OrgIDFromContext(ctx context.Context) (int64, bool) {
	value := ctx.Value(OrgContextKey{})
	if value == nil {
		return 0, false
	}
	orgID, ok := value.(int64)
	return orgID, ok
}
