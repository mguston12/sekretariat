package sekretariat

import (
	"context"
	"sekretariat/internal/entity/sekretariat"
)

// Data ...
// Masukkan function dari package data ke dalam interface ini
type Data interface {
	// Contracts
	GetAllContractsHeader(ctx context.Context, company int, keyword string) ([]sekretariat.KontrakHeader, error)
	GetAllContractsHeaderPage(ctx context.Context, company int, keyword string, offset, limit int) ([]sekretariat.KontrakHeader, error)
	GetAllContractsHeaderCount(ctx context.Context, company int) ([]sekretariat.KontrakHeader, int, error)
	GetCounterContract(ctx context.Context, company int) (int, error)
	GetContractsHeaderByContractNumber(ctx context.Context, company int, no_kontrak string) (sekretariat.KontrakHeader, error)
	GetContractDetailsByContractNumber(ctx context.Context, no_kontrak string) ([]sekretariat.KontrakDetail, error)

	CreateContractHeader(ctx context.Context, header sekretariat.KontrakHeader) error
	CreateContractDetail(ctx context.Context, detail sekretariat.KontrakDetail) error
	IncreaseCounterContract(ctx context.Context, company int) error

	// Customer
	GetCustomerFiltered(ctx context.Context, company int, keyword string, offset, limit int) ([]sekretariat.Customer, error)
	GetCustomerFilteredCount(ctx context.Context, company int, keyword string) ([]sekretariat.Customer, int, error)
	CreateCustomer(ctx context.Context, customer sekretariat.Customer) error
	GetCustomerByID(ctx context.Context, id string) (sekretariat.Customer, error)
	FetchAndIncreaseCounter(ctx context.Context, company int) (int, error)
	UpdateCustomer(ctx context.Context, customer sekretariat.Customer) error

	// Company
	GetAllCompanies(ctx context.Context) ([]sekretariat.Company, error)
	GetCompanyByID(ctx context.Context, id int) (sekretariat.Company, error)

	// Bank
	GetAllBanks(ctx context.Context) ([]sekretariat.Bank, error)
	GetBankByID(ctx context.Context, id int) (sekretariat.Bank, error)
}

// Service ...
// Tambahkan variable sesuai banyak data layer yang dibutuhkan
type Service struct {
	data Data
}

// New ...
// Tambahkan parameter sesuai banyak data layer yang dibutuhkan
func New(data Data) Service {
	// Assign variable dari parameter ke object
	return Service{
		data: data,
	}
}