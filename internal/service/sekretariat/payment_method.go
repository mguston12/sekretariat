package sekretariat

import (
	"context"
	"sekretariat/internal/entity/sekretariat"
	"sekretariat/pkg/errors"
)

func (s Service) GetPaymentMethod(ctx context.Context) ([]sekretariat.Pembayaran, error) {
	methods, err := s.data.GetPaymentMethod(ctx)
	if err != nil {
		return methods, errors.Wrap(err, "[SERVICE][GetPaymentMethod]")
	}
	return methods, nil
}
