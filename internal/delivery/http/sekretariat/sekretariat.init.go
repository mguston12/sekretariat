package sekretariat

import (
	"bytes"
	"context"
	sekretariat "sekretariat/internal/entity/sekretariat"
)

type SekretariatSvc interface {
	// Contract
	GetAllContractsHeader(ctx context.Context, company int, keyword string, page, length int) ([]sekretariat.KontrakHeader, int, error)
	GetDataContractByContractNumber(ctx context.Context, company int, no_kontrak string) (sekretariat.KontrakHeader, error)
	GetCounterContract(ctx context.Context, company int) (string, error)

	// GetContractFiltered(ctx context.Context, company int, id, name string) ([]sekretariat.Kontrak, error)

	CreateContract(ctx context.Context, header sekretariat.KontrakHeader) error
	PrintKontrak(ctx context.Context, company int, no_kontrak string) (bytes.Buffer, error)

	// Company
	GetAllCompanies(ctx context.Context) ([]sekretariat.Company, error)

	// Customer
	GetCustomerFiltered(ctx context.Context, company int, keyword string, page, length int) ([]sekretariat.Customer, int, error)
	CreateCustomer(ctx context.Context, customer sekretariat.Customer) error
	UpdateCustomer(ctx context.Context, customer sekretariat.Customer) error

	// Bank
	GetAllBanks(ctx context.Context) ([]sekretariat.Bank, error)
}

type Handler struct {
	sekretariatSvc SekretariatSvc
}

func New(is SekretariatSvc) *Handler {
	return &Handler{
		sekretariatSvc: is,
	}
}
