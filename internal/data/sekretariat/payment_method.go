package sekretariat

import (
	"context"
	"sekretariat/internal/entity/sekretariat"

	"sekretariat/pkg/errors"

	"github.com/jmoiron/sqlx"
)

func (d Data) GetPaymentMethod(ctx context.Context) ([]sekretariat.Pembayaran, error) {
	var (
		rows  *sqlx.Rows
		datas []sekretariat.Pembayaran
		err   error
	)

	rows, err = d.stmt[getPaymentMethod].QueryxContext(ctx)
	if err != nil {
		return datas, errors.Wrap(err, "[DATA][GetPaymentMethod]")
	}

	for rows.Next() {
		var data sekretariat.Pembayaran
		err := rows.StructScan(&data)
		if err != nil {
			return datas, errors.Wrap(err, "[DATA][GetPaymentMethod]")
		}
		datas = append(datas, data)
	}
	defer rows.Close()

	return datas, nil
}
