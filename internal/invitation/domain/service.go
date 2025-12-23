package domain

import "context"

type InviteRequest struct {
	OrganizationID string `json:"organization_id"`
	Email string `json:"email"`
	Role string `json:"role"`
}

type BatchInviteRequest struct {
	Invitations []InviteRequest `json:"invitations"`
}

type VerifyRequest struct {
	OrganizationID string `json:"organization_id"`
	Code string `json:"code"`
}

type Service interface {
	BatchInvite(context.Context, BatchInviteRequest) error
	Invite(context.Context, InviteRequest) error
	Verify(context.Context, VerifyRequest) error
}