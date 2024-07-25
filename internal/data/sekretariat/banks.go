package sekretariat

import (
	"context"
	"sekretariat/internal/entity/sekretariat"

	"sekretariat/pkg/errors"

	"github.com/jmoiron/sqlx"
)

func (d Data) GetAllBanks(ctx context.Context) ([]sekretariat.Bank, error) {
	var (
		rows  *sqlx.Rows
		datas []sekretariat.Bank
		err   error
	)

	rows, err = d.stmt[getAllBanks].QueryxContext(ctx)
	if err != nil {
		return datas, errors.Wrap(err, "[DATA][GetAllBanks]")
	}

	for rows.Next() {
		var data sekretariat.Bank
		err := rows.StructScan(&data)
		if err != nil {
			return datas, errors.Wrap(err, "[DATA][GetAllBanks]")
		}
		datas = append(datas, data)
	}
	defer rows.Close()

	return datas, nil
}

func (d Data) GetBankByID(ctx context.Context, id int) (sekretariat.Bank, error) {
	bank := sekretariat.Bank{}

	if err := d.stmt[getBankByID].QueryRowxContext(ctx, id).StructScan(&bank); err != nil {
		return bank, errors.Wrap(err, "[DATA][GetBankByID]")
	}

	return bank, nil
}