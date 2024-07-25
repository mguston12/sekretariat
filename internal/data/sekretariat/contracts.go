package sekretariat

import (
	"context"
	"sekretariat/internal/entity/sekretariat"

	"sekretariat/pkg/errors"

	"github.com/jmoiron/sqlx"
)

func (d Data) GetAllContractsHeader(ctx context.Context, company int, keyword string) ([]sekretariat.KontrakHeader, error) {
	var (
		rows  *sqlx.Rows
		datas []sekretariat.KontrakHeader
		err   error
	)

	_keyword := "%" + keyword + "%"

	rows, err = d.stmt[getAllContractsHeader].QueryxContext(ctx, company, _keyword, _keyword, _keyword, _keyword)
	if err != nil {
		return datas, errors.Wrap(err, "[DATA][GetAllContractsHeader]")
	}

	for rows.Next() {
		var data sekretariat.KontrakHeader
		err := rows.StructScan(&data)
		if err != nil {
			return datas, errors.Wrap(err, "[DATA][GetAllContractsHeader]")
		}
		datas = append(datas, data)
	}
	defer rows.Close()

	return datas, nil
}

func (d Data) GetAllContractsHeaderPage(ctx context.Context, company int, keyword string, offset, limit int) ([]sekretariat.KontrakHeader, error) {
	headers := []sekretariat.KontrakHeader{}

	_keyword := "%" + keyword + "%"

	rows, err := d.stmt[getAllContractsHeaderPage].QueryxContext(ctx, company, _keyword, _keyword, _keyword, _keyword, offset, limit)
	if err != nil {
		return headers, errors.Wrap(err, "[DATA][GetAllContractsHeaderPage1]")
	}

	defer rows.Close()

	for rows.Next() {
		header := sekretariat.KontrakHeader{}
		if err = rows.StructScan(&header); err != nil {
			return headers, errors.Wrap(err, "[DATA][GetAllContractsHeaderPage2]")
		}
		headers = append(headers, header)
	}

	return headers, nil
}

func (d Data) GetAllContractsHeaderCount(ctx context.Context, company int) ([]sekretariat.KontrakHeader, int, error) {
	headers := []sekretariat.KontrakHeader{}
	var count int

	if err := d.stmt[getAllContractsHeaderCount].QueryRowxContext(ctx, company).Scan(&count); err != nil {
		return headers, count, errors.Wrap(err, "[DATA][GetAllContractsHeaderCount]")
	}

	return headers, count, nil
}

func (d Data) GetContractsHeaderByContractNumber(ctx context.Context, company int, no_kontrak string) (sekretariat.KontrakHeader, error) {
	header := sekretariat.KontrakHeader{}
	_no_kontrak := no_kontrak + "%"

	if err := d.stmt[getContractHeaderByContractNumber].QueryRowxContext(ctx, company, _no_kontrak).StructScan(&header); err != nil {
		return header, errors.Wrap(err, "[DATA][GetContractsHeaderByContractNumber]")
	}

	return header, nil
}

func (d Data) GetContractDetailsByContractNumber(ctx context.Context, no_kontrak string) ([]sekretariat.KontrakDetail, error) {
	var (
		rows  *sqlx.Rows
		datas []sekretariat.KontrakDetail
		err   error
	)

	_no_kontrak := no_kontrak + "%"

	rows, err = d.stmt[getContractDetailsByContractNumber].QueryxContext(ctx, _no_kontrak)
	if err != nil {
		return datas, errors.Wrap(err, "[DATA][GetContractDetailsByContractNumber]")
	}

	for rows.Next() {
		var data sekretariat.KontrakDetail
		err := rows.StructScan(&data)
		if err != nil {
			return datas, errors.Wrap(err, "[DATA][GetContractDetailsByContractNumber]")
		}
		datas = append(datas, data)
	}
	defer rows.Close()

	return datas, nil
}

func (d Data) CreateContractHeader(ctx context.Context, header sekretariat.KontrakHeader) error {
	_, err := d.stmt[createContractHeader].ExecContext(ctx,
		header.NoKontrak,
		header.TanggalBuat,
		header.CompanyID,
		header.CustomerID,
		header.Bank.ID,
		header.ActiveYN,
		header.UpdatedBy,
	)

	if err != nil {
		return errors.Wrap(err, "[DATA][CreateContractHeader]")
	}
	return nil
}

func (d Data) CreateContractDetail(ctx context.Context, detail sekretariat.KontrakDetail) error {
	_, err := d.stmt[createContractDetail].ExecContext(ctx,
		detail.NoKontrak,
		detail.Quantity,
		detail.Tipe,
		detail.Speed,
		detail.Harga,
		detail.FreeCopy,
		detail.OverCopy,
		detail.FreeCopyColor,
		detail.OverCopyColor,
		detail.PeriodeAwal,
		detail.PeriodeAkhir,
		detail.Penempatan,
		detail.ActiveYN,
		detail.UpdatedBy,
	)

	if err != nil {
		return errors.Wrap(err, "[DATA][CreateContractDetail]")
	}
	return nil
}

// func (d Data) UpdateStorage(ctx context.Context, storage purchasing.Storage) error {
// 	_, err := d.stmt[updateStorage].ExecContext(ctx,
// 		storage.Name,
// 		storage.ActiveYN,
// 		storage.UserID,
// 		storage.Code)
// 	if err != nil {
// 		return errors.Wrap(err, "[DATA][UpdateStorage]")
// 	}
// 	return nil
// }

// func (d Data) DeleteStorageByCode(ctx context.Context, stocode int) error {
// 	_, err := d.stmt[deleteStorageByCode].ExecContext(ctx, stocode)

// 	if err != nil {
// 		return errors.Wrap(err, "[DATA][DeleteStorageByCode]")
// 	}
// 	return nil
// }

func (d Data) GetCounterContract(ctx context.Context, company int) (int, error) {
	counter := 0
	tx, err := d.db.Beginx()
	if err != nil {
		return counter, errors.Wrap(err, "[DATA][GetCounterContract]")
	}
	if err := tx.QueryRowxContext(ctx, `SELECT count FROM counter_kontrak WHERE company = ?`, company).Scan(&counter); err != nil {
		return counter, errors.Wrap(err, "[DATA][GetCounterContract][A]")
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return counter, errors.Wrap(err, "[DATA][GetCounterContract]")
	}
	return counter, nil
}

func (d Data) IncreaseCounterContract(ctx context.Context, company int) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "[DATA][FetchAndIncreaseCounter]")
	}
	if _, err := tx.ExecContext(ctx, `UPDATE counter_kontrak SET count = count + 1 WHERE company = ?`, company); err != nil {
		return errors.Wrap(err, "[DATA][FetchAndIncreaseCounter][B]")
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "[DATA][FetchAndIncreaseCounter]")
	}
	return nil
}
