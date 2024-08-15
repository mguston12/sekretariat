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

func (d Data) CreateBank(ctx context.Context, bank sekretariat.Bank) error {
	_, err := d.stmt[createBank].ExecContext(ctx,
		bank.Name,
		bank.Norek,
		bank.AtasNama,
	)

	if err != nil {
		return errors.Wrap(err, "[DATA][CreateBank]")
	}
	return nil
}

func (d Data) UpdateBank(ctx context.Context, bank sekretariat.Bank) error {
	_, err := d.stmt[updateBank].ExecContext(ctx,
		bank.Name,
		bank.Norek,
		bank.AtasNama,
		bank.ID)
	if err != nil {
		return errors.Wrap(err, "[DATA][UpdateBank]")
	}
	return nil
}

func (d Data) DeleteBankByID(ctx context.Context, id int) error {
	_, err := d.stmt[deleteBankByID].ExecContext(ctx, id)

	if err != nil {
		return errors.Wrap(err, "[DATA][DeleteBankByID]")
	}
	return nil
}
