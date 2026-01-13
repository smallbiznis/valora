package domain

import "time"

// Inbox View (Unassigned / Needs Action)

type InboxRequest struct {
	Limit int `json:"limit" form:"limit"`
}

type InboxItem struct {
	EntityType   string     `json:"entity_type"` // "invoice" | "customer"
	EntityID     string     `json:"entity_id"`
	EntityName   string     `json:"entity_name"`   // invoice_number or customer_name
	RiskCategory string     `json:"risk_category"` // "overdue" | "failed_payment" | "high_exposure"
	RiskScore    int        `json:"risk_score"`    // For sorting
	AmountDue    int64      `json:"amount_due"`
	Currency     string     `json:"currency"`
	DaysOverdue  int        `json:"days_overdue,omitempty"`
	LastAttempt  *time.Time `json:"last_attempt,omitempty"`
	PublicToken  string     `json:"public_token,omitempty"`
}

type InboxResponse struct {
	Items    []InboxItem `json:"items"`
	Currency string      `json:"currency"`
}

// My Work View (Claimed by Me)

type MyWorkRequest struct {
	Limit int `json:"limit" form:"limit"`
}

type MyWorkItem struct {
	AssignmentID string `json:"assignment_id"`
	EntityType   string `json:"entity_type"`
	EntityID     string `json:"entity_id"`
	EntityName   string `json:"entity_name"`

	// Snapshot values at claim time (stable)
	AmountDueAtClaim   int64 `json:"amount_due_at_claim"`
	DaysOverdueAtClaim int   `json:"days_overdue_at_claim"`

	// Current values (optional, non-sorting)
	CurrentAmountDue   int64 `json:"current_amount_due,omitempty"`
	CurrentDaysOverdue int   `json:"current_days_overdue,omitempty"`

	Currency      string     `json:"currency"`
	ClaimedAt     time.Time  `json:"claimed_at"`
	AssignmentAge string     `json:"assignment_age"` // "2h 15m"
	Status        string     `json:"status"`         // "claimed" | "in_progress"
	LastActionAt  *time.Time `json:"last_action_at,omitempty"`
	PublicToken   string     `json:"public_token,omitempty"`
}

type MyWorkResponse struct {
	Items    []MyWorkItem `json:"items"`
	Currency string       `json:"currency"`
}

// Recently Resolved View

type RecentlyResolvedRequest struct {
	Limit int `json:"limit" form:"limit"`
}

type ResolvedItem struct {
	AssignmentID string    `json:"assignment_id"`
	EntityType   string    `json:"entity_type"`
	EntityID     string    `json:"entity_id"`
	EntityName   string    `json:"entity_name"`
	Status       string    `json:"status"` // "resolved" | "released" | "escalated"
	ResolvedAt   time.Time `json:"resolved_at"`
	ResolvedBy   string    `json:"resolved_by"`
	Reason       string    `json:"reason,omitempty"`
	ClaimedAt    time.Time `json:"claimed_at"`
	Duration     string    `json:"duration"` // "3h 45m"
}

type RecentlyResolvedResponse struct {
	Items []ResolvedItem `json:"items"`
}

// Team View (Manager Only)

type TeamViewRequest struct{}

type TeamMemberWorkload struct {
	UserID             string `json:"user_id"`
	ActiveAssignments  int    `json:"active_assignments"`
	AvgAssignmentAge   string `json:"avg_assignment_age"` // "1h 30m"
	TotalExposureOwned int64  `json:"total_exposure_owned"`
	EscalationCount    int    `json:"escalation_count"`
}

type TeamViewResponse struct {
	Members  []TeamMemberWorkload `json:"members"`
	Currency string               `json:"currency"`
}
