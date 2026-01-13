package providers

import (
	"github.com/smallbiznis/valora/internal/payment"
	"github.com/smallbiznis/valora/internal/providers/email"
	"github.com/smallbiznis/valora/internal/providers/pdf"
	"go.uber.org/fx"
)

var Module = fx.Module("providers",
	email.Module,
	payment.Module,
	pdf.Module,
)
