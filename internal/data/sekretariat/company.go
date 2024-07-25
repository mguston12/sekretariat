package sekretariat

import (
	"context"
	"sekretariat/internal/entity/sekretariat"

	"sekretariat/pkg/errors"

	"github.com/jmoiron/sqlx"
)

func (d Data) GetAllCompanies(ctx context.Context) ([]sekretariat.Company, error) {
	var (
		rows  *sqlx.Rows
		datas []sekretariat.Company
		err   error
	)

	rows, err = d.stmt[getAllCompanies].QueryxContext(ctx)
	if err != nil {
		return datas, errors.Wrap(err, "[DATA][GetAllCompanies]")
	}

	for rows.Next() {
		var data sekretariat.Company
		err := rows.StructScan(&data)
		if err != nil {
			return datas, errors.Wrap(err, "[DATA][GetAllCompanies]")
		}
		datas = append(datas, data)
	}
	defer rows.Close()

	return datas, nil
}

func (d Data) GetCompanyByID(ctx context.Context, id int) (sekretariat.Company, error) {
	company := sekretariat.Company{}

	if err := d.stmt[getCompanyByID].QueryRowxContext(ctx, id).StructScan(&company); err != nil {
		return company, errors.Wrap(err, "[DATA][GetCompanyByID]")
	}

	return company, nil
}
