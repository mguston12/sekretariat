package sekretariat

import (
	"context"
	"sekretariat/internal/entity/sekretariat"
	"sekretariat/pkg/errors"
)

func (s Service) GetAllBanks(ctx context.Context) ([]sekretariat.Bank, error) {
	banks, err := s.data.GetAllBanks(ctx)
	if err != nil {
		return banks, errors.Wrap(err, "[SERVICE][GetAllBanks]")
	}
	return banks, nil
}

func (s Service) GetBankByID(ctx context.Context, id int) (sekretariat.Bank, error) {
	bank, err := s.data.GetBankByID(ctx, id)
	if err != nil {
		return bank, errors.Wrap(err, "[SERVICE][GetBankByID]")
	}
	return bank, nil
}
