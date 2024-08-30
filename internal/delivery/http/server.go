package http

import (
	"net/http"

	"sekretariat/pkg/grace"

	"github.com/rs/cors"
)

type SekretariatHandler interface {
	// Contract
	GetAllContractsHeader(w http.ResponseWriter, r *http.Request)
	GetContractExp30Days(w http.ResponseWriter, r *http.Request)
	GetDataContractByContractNumber(w http.ResponseWriter, r *http.Request)
	GetCounterContract(w http.ResponseWriter, r *http.Request)

	CreateContract(w http.ResponseWriter, r *http.Request)
	UpdateContract(w http.ResponseWriter, r *http.Request)
	PrintKontrak(w http.ResponseWriter, r *http.Request)

	// Company
	GetAllCompanies(w http.ResponseWriter, r *http.Request)

	// Customer
	GetCustomerFiltered(w http.ResponseWriter, r *http.Request)
	CreateCustomer(w http.ResponseWriter, r *http.Request)
	UpdateCustomer(w http.ResponseWriter, r *http.Request)
	ImportCustomersFromExcel(w http.ResponseWriter, r *http.Request)

	// Bank
	GetAllBanks(w http.ResponseWriter, r *http.Request)
	GetBankByCompanyID(w http.ResponseWriter, r *http.Request)
	CreateBank(w http.ResponseWriter, r *http.Request)
	UpdateBank(w http.ResponseWriter, r *http.Request)
	DeleteBankByID(w http.ResponseWriter, r *http.Request)

	// Payment Method
	GetPaymentMethod(w http.ResponseWriter, r *http.Request)
}

// Server ...
type Server struct {
	Sekretariat SekretariatHandler
}

// Serve is serving HTTP gracefully on port x ...
func (s *Server) Serve(port string) error {
	handler := cors.AllowAll().Handler(s.Handler())
	return grace.Serve(port, handler)
}
