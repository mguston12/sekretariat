package http

import (
	"errors"
	"log"
	"net/http"

	"sekretariat/pkg/response"

	"github.com/gorilla/mux"
)

// Handler will initialize mux router and register handler
func (s *Server) Handler() *mux.Router {
	r := mux.NewRouter()
	// Jika tidak ditemukan, jangan diubah.
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	// Tambahan Prefix di depan API endpoint
	router := r.PathPrefix("").Subrouter()
	// Health Check
	router.HandleFunc("", defaultHandler).Methods("GET")
	router.HandleFunc("/", defaultHandler).Methods("GET")
	// Routes
	// router.HandleFunc("/contracts", s.Sekretariat.GetContractFiltered).Methods("GET")

	contract := r.PathPrefix("/contracts").Subrouter()

	contract.HandleFunc("/detail", s.Sekretariat.GetDataContractByContractNumber).Methods("GET")
	contract.HandleFunc("/create", s.Sekretariat.CreateContract).Methods("POST")
	contract.HandleFunc("/update", s.Sekretariat.UpdateContract).Methods("POST")
	contract.HandleFunc("/print", s.Sekretariat.PrintKontrak).Methods("GET")
	contract.HandleFunc("/counter", s.Sekretariat.GetCounterContract).Methods("GET")
	contract.HandleFunc("/expiredsoon", s.Sekretariat.GetContractExp30Days).Methods("GET")
	contract.HandleFunc("", s.Sekretariat.GetAllContractsHeader).Methods("GET")

	customer := r.PathPrefix("/customers").Subrouter()
	customer.HandleFunc("/create", s.Sekretariat.CreateCustomer).Methods("POST")
	customer.HandleFunc("/update", s.Sekretariat.UpdateCustomer).Methods("PUT")
	customer.HandleFunc("/croncustomers", s.Sekretariat.ImportCustomersFromExcel).Methods("GET")
	customer.HandleFunc("", s.Sekretariat.GetCustomerFiltered).Methods("GET")

	company := r.PathPrefix("/companies").Subrouter()
	company.HandleFunc("", s.Sekretariat.GetAllCompanies).Methods("GET")

	bank := r.PathPrefix("/banks").Subrouter()
	bank.HandleFunc("/create", s.Sekretariat.CreateBank).Methods("POST")
	bank.HandleFunc("/update", s.Sekretariat.UpdateBank).Methods("PUT")
	bank.HandleFunc("/delete", s.Sekretariat.DeleteBankByID).Methods("DELETE")
	bank.HandleFunc("/filter", s.Sekretariat.GetBankByCompanyID).Methods("GET")
	bank.HandleFunc("", s.Sekretariat.GetAllBanks).Methods("GET")

	paymentMethod := r.PathPrefix("/payment-method").Subrouter()
	paymentMethod.HandleFunc("", s.Sekretariat.GetPaymentMethod).Methods("GET")

	return r
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Go Skeleton API"))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	var (
		resp   *response.Response
		err    error
		errRes response.Error
	)
	resp = &response.Response{}
	defer resp.RenderJSON(w, r)

	err = errors.New("404 Not Found")

	if err != nil {
		// Error response handling
		errRes = response.Error{
			Code:   404,
			Msg:    "404 Not Found",
			Status: true,
		}

		log.Printf("[ERROR] %s %s - %v\n", r.Method, r.URL, err)
		resp.StatusCode = 404
		resp.Error = errRes
		return
	}
}
