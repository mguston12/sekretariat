package sekretariat

import (
	"context"
	"fmt"
	"math"
	"sekretariat/internal/entity/sekretariat"
	"sekretariat/pkg/errors"
)

func (s Service) GetCustomerFiltered(ctx context.Context, company int, keyword string, page, length int) ([]sekretariat.Customer, int, error) {
	limit := length
	offset := (page - 1) * length
	var lastPage int

	if page != 0 && length != 0 {
		customers, count, err := s.data.GetCustomerFilteredCount(ctx, company, keyword)
		if err != nil {
			return customers, lastPage, errors.Wrap(err, "[SERVICE][GetCustomerFiltered][COUNT]")

		}
		lastPage = int(math.Ceil(float64(count) / float64(length)))

		customers, err = s.data.GetCustomerFiltered(ctx, company, keyword, offset, limit)
		if err != nil {
			return customers, lastPage, errors.Wrap(err, "[SERVICE][GetCustomerFiltered]")
		}

		return customers, lastPage, nil
	}

	customers, err := s.data.GetAllCustomers(ctx, company)
	if err != nil {
		return customers, lastPage, errors.Wrap(err, "[SERVICE][GetCustomerFiltered]")
	}

	return customers, lastPage, nil
}

func (s Service) GetCustomerByID(ctx context.Context, id string) (sekretariat.Customer, error) {
	customer, err := s.data.GetCustomerByID(ctx, id)
	if err != nil {
		return customer, errors.Wrap(err, "[SERVICE][GetCustomerByID]")
	}
	return customer, nil
}

func (s Service) GetCustomerID(ctx context.Context, company int) (string, error) {
	com := ""

	switch company {
	case 1:
		com = "PB"
	case 2:
		com = "PBM"
	case 3:
		com = "MMU"
	}

	id, err := s.data.FetchAndIncreaseCounter(ctx, company)
	if err != nil {
		return "", errors.Wrap(err, "[SERVICE][GetCustomerFiltered]")
	}

	return fmt.Sprintf("%s-%06d", com, id), nil
}

func (s Service) CreateCustomer(ctx context.Context, customer sekretariat.Customer) error {
	customerID, err := s.GetCustomerID(ctx, customer.CompanyID)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][CreateCustomer][1]")
	}

	customer.CustomerID = customerID

	err = s.data.CreateCustomer(ctx, customer)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][CreateCustomer][2]")
	}

	return nil
}

func (s Service) UpdateCustomer(ctx context.Context, customer sekretariat.Customer) error {
	err := s.data.UpdateCustomer(ctx, customer)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][UpdateCustomer]")
	}

	return nil
}
