package sekretariat

import (
	"context"
	"sekretariat/internal/entity/sekretariat"

	"sekretariat/pkg/errors"

	"github.com/jmoiron/sqlx"
)

func (d Data) GetAllCustomers(ctx context.Context, company int) ([]sekretariat.Customer, error) {
	var (
		rows  *sqlx.Rows
		datas []sekretariat.Customer
		err   error
	)

	rows, err = d.stmt[getAllCustomers].QueryxContext(ctx, company)
	if err != nil {
		return datas, errors.Wrap(err, "[DATA][GetAllCustomers]")
	}

	for rows.Next() {
		var data sekretariat.Customer
		err := rows.StructScan(&data)
		if err != nil {
			return datas, errors.Wrap(err, "[DATA][GetAllCustomers]")
		}
		datas = append(datas, data)
	}
	defer rows.Close()

	return datas, nil
}

func (d Data) GetCustomerFiltered(ctx context.Context, company int, keyword string, offset, limit int) ([]sekretariat.Customer, error) {
	var (
		rows  *sqlx.Rows
		datas []sekretariat.Customer
		err   error
	)

	_keyword := "%" + keyword + "%"

	rows, err = d.stmt[getCustomerFiltered].QueryxContext(ctx, company, _keyword, _keyword, _keyword, _keyword, offset, limit)
	if err != nil {
		return datas, errors.Wrap(err, "[DATA][GetCustomerFiltered]")
	}

	for rows.Next() {
		var data sekretariat.Customer
		err := rows.StructScan(&data)
		if err != nil {
			return datas, errors.Wrap(err, "[DATA][GetCustomerFiltered]")
		}
		datas = append(datas, data)
	}
	defer rows.Close()

	return datas, nil
}

func (d Data) GetCustomerFilteredCount(ctx context.Context, company int, keyword string) ([]sekretariat.Customer, int, error) {
	customers := []sekretariat.Customer{}
	var count int

	if err := d.stmt[getCustomerFilteredCount].QueryRowxContext(ctx, company, keyword, keyword, keyword, keyword).Scan(&count); err != nil {
		return customers, count, errors.Wrap(err, "[DATA][GetCustomerFilteredCount]")
	}

	return customers, count, nil
}

func (d Data) GetCustomerByID(ctx context.Context, id string) (sekretariat.Customer, error) {
	customer := sekretariat.Customer{}

	if err := d.stmt[getCustomerByID].QueryRowxContext(ctx, id).StructScan(&customer); err != nil {
		return customer, errors.Wrap(err, "[DATA][GetCustomerByID]")
	}

	return customer, nil
}

func (d Data) CreateCustomer(ctx context.Context, customer sekretariat.Customer) error {
	_, err := d.stmt[createCustomer].ExecContext(ctx,
		customer.CustomerID,
		customer.CompanyID,
		customer.CustomerName,
		customer.Address,
		customer.PIC,
		customer.PenandaTangan,
		customer.Jabatan,
		customer.NoTelp,
		customer.UpdatedBy,
	)
	if err != nil {
		return errors.Wrap(err, "[DATA][CreateCustomer]")
	}
	return nil
}

func (d Data) UpdateCustomer(ctx context.Context, customer sekretariat.Customer) error {
	_, err := d.stmt[updateCustomer].ExecContext(ctx,
		customer.CustomerName,
		customer.Address,
		customer.PIC,
		customer.PenandaTangan,
		customer.Jabatan,
		customer.NoTelp,
		customer.UpdatedBy,
		customer.CustomerID,
	)
	if err != nil {
		return errors.Wrap(err, "[DATA][UpdateCustomer]")
	}
	return nil
}

func (d Data) FetchAndIncreaseCounter(ctx context.Context, company int) (int, error) {
	counter := 0
	tx, err := d.db.Beginx()
	if err != nil {
		return counter, errors.Wrap(err, "[DATA][FetchAndIncreaseCounter]")
	}
	if err := tx.QueryRowxContext(ctx, `SELECT count FROM counter WHERE company = ?`, company).Scan(&counter); err != nil {
		return counter, errors.Wrap(err, "[DATA][FetchAndIncreaseCounter][A]")
	}
	if _, err := tx.ExecContext(ctx, `UPDATE counter SET count = count + 1 WHERE company = ?`, company); err != nil {
		return counter, errors.Wrap(err, "[DATA][FetchAndIncreaseCounter][B]")
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return counter, errors.Wrap(err, "[DATA][FetchAndIncreaseCounter]")
	}
	return counter, nil
}
