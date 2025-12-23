package domain

import "context"

type Service interface {
	GenerateInvoice(context.Context)
}
