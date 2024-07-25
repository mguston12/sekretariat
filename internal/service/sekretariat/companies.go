package sekretariat

import (
	"context"
	"sekretariat/internal/entity/sekretariat"
	"sekretariat/pkg/errors"
)

func (s Service) GetAllCompanies(ctx context.Context) ([]sekretariat.Company, error) {
	companies, err := s.data.GetAllCompanies(ctx)
	if err != nil {
		return companies, errors.Wrap(err, "[SERVICE][GetAllCompanies]")
	}
	return companies, nil
}

func (s Service) GetCompanyByID(ctx context.Context, id int) (sekretariat.Company, error) {
	company, err := s.data.GetCompanyByID(ctx, id)
	if err != nil {
		return company, errors.Wrap(err, "[SERVICE][GetCompanyByID]")
	}
	return company, nil
}
