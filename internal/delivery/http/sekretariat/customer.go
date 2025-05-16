package sekretariat

import (
	"encoding/json"
	"log"
	"net/http"
	httpHelper "sekretariat/internal/delivery/http"
	"sekretariat/internal/entity/sekretariat"
	"sekretariat/pkg/response"
	"strconv"
)

func (h *Handler) GetCustomerFiltered(w http.ResponseWriter, r *http.Request) {
	var (
		result   interface{}
		metadata interface{}
		err      error
		resp     response.Response
	)
	defer resp.RenderJSON(w, r)

	company, _ := strconv.Atoi(r.FormValue("company"))

	page, _ := strconv.Atoi(r.FormValue("page"))
	length, _ := strconv.Atoi(r.FormValue("length"))

	ctx := r.Context()
	result, metadata, err = h.sekretariatSvc.GetCustomerFiltered(ctx, company, r.FormValue("keyword"), page, length)

	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}

	resp.Data = result
	resp.Metadata = metadata

	log.Printf("[INFO][%s][%s] %s", r.RemoteAddr, r.Method, r.URL)
}

func (h *Handler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	var (
		result interface{}
		err    error
		resp   response.Response
	)
	defer resp.RenderJSON(w, r)

	ctx := r.Context()
	result, err = h.sekretariatSvc.GetCustomer(ctx, r.FormValue("keyword"))

	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}

	resp.Data = result

	log.Printf("[INFO][%s][%s] %s", r.RemoteAddr, r.Method, r.URL)
}

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.RenderJSON(w, r)

	ctx := r.Context()

	customer := sekretariat.Customer{}

	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}

	err = h.sekretariatSvc.CreateCustomer(ctx, customer)
	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}

	log.Printf("[INFO][%s][%s] %s", r.RemoteAddr, r.Method, r.URL)
}

func (h *Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.RenderJSON(w, r)

	ctx := r.Context()

	customer := sekretariat.Customer{}

	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}

	err = h.sekretariatSvc.UpdateCustomer(ctx, customer)
	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}

	log.Printf("[INFO][%s][%s] %s", r.RemoteAddr, r.Method, r.URL)
}

func (h *Handler) ImportCustomersFromExcel(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		resp response.Response
	)
	defer resp.RenderJSON(w, r)

	ctx := r.Context()
	err = h.sekretariatSvc.ImportCustomersFromExcel(ctx)

	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}

	log.Printf("[INFO][%s][%s] %s", r.RemoteAddr, r.Method, r.URL)
}
