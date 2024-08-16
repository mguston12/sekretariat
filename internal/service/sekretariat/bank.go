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

func (s Service) GetBankByCompanyID(ctx context.Context, id int) ([]sekretariat.Bank, error) {
	bank, err := s.data.GetBankByCompanyID(ctx, id)
	if err != nil {
		return bank, errors.Wrap(err, "[SERVICE][GetBankByCompanyID]")
	}
	return bank, nil
}

func (s Service) CreateBank(ctx context.Context, bank sekretariat.Bank) error {
	err := s.data.CreateBank(ctx, bank)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][CreateBank]")
	}
	return nil
}

func (s Service) UpdateBank(ctx context.Context, bank sekretariat.Bank) error {
	err := s.data.UpdateBank(ctx, bank)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][UpdateBank]")
	}
	return nil
}

func (s Service) DeleteBankByID(ctx context.Context, id int) error {
	err := s.data.DeleteBankByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][DeleteBankByID]")
	}
	return nil
}
