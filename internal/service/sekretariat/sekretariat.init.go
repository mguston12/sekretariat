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
	GetAllContractsHeaderCount(ctx context.Context, company int, keyword string) ([]sekretariat.KontrakHeader, int, error)
	GetCounterContract(ctx context.Context, company int) (int, error)
	GetContractsHeaderByContractNumber(ctx context.Context, company int, no_kontrak string) (sekretariat.KontrakHeader, error)
	GetContractDetailsByContractNumber(ctx context.Context, no_kontrak string) ([]sekretariat.KontrakDetail, error)
	GetContractByNumber(ctx context.Context, kontrak string) (bool, error)
	GetContractExp30Days(ctx context.Context, company int) ([]sekretariat.KontrakDetail, error)

	CreateContractHeader(ctx context.Context, header sekretariat.KontrakHeader) error
	CreateContractDetail(ctx context.Context, detail sekretariat.KontrakDetail) error
	UpdateContractHeader(ctx context.Context, header sekretariat.KontrakHeader) error
	IncreaseCounterContract(ctx context.Context, company int) error
	DeleteContractDetail(ctx context.Context, kontrak string) error

	// Customer
	GetCustomer(ctx context.Context, keyword string) ([]sekretariat.Customer, error)
	GetAllCustomers(ctx context.Context, company int) ([]sekretariat.Customer, error)
	GetCustomerIDByNameAndAddress(ctx context.Context, name, address string) (string, error)
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
	GetBankByCompanyID(ctx context.Context, id int) ([]sekretariat.Bank, error)
	CreateBank(ctx context.Context, bank sekretariat.Bank) error
	UpdateBank(ctx context.Context, bank sekretariat.Bank) error
	DeleteBankByID(ctx context.Context, id int) error

	// Payment Method
	GetPaymentMethod(ctx context.Context) ([]sekretariat.Pembayaran, error)
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
