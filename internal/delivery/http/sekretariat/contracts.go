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

func (h *Handler) GetAllContractsHeader(w http.ResponseWriter, r *http.Request) {
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
	result, metadata, err = h.sekretariatSvc.GetAllContractsHeader(ctx, company, r.FormValue("keyword"), page, length)

	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}

	resp.Data = result
	resp.Metadata = metadata

	log.Printf("[INFO][%s][%s] %s", r.RemoteAddr, r.Method, r.URL)
}

func (h *Handler) GetDataContractByContractNumber(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.RenderJSON(w, r)

	company, _ := strconv.Atoi(r.FormValue("company"))
	ctx := r.Context()

	result, err := h.sekretariatSvc.GetDataContractByContractNumber(ctx, company, r.FormValue("kontrak"))
	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}

	resp.Data = result

	log.Printf("[INFO][%s][%s] %s", r.RemoteAddr, r.Method, r.URL)
}

func (h *Handler) GetCounterContract(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.RenderJSON(w, r)

	company, _ := strconv.Atoi(r.FormValue("company"))
	ctx := r.Context()

	result, err := h.sekretariatSvc.GetCounterContract(ctx, company)
	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}

	resp.Data = result

	log.Printf("[INFO][%s][%s] %s", r.RemoteAddr, r.Method, r.URL)
}

func (h *Handler) CreateContract(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.RenderJSON(w, r)

	ctx := r.Context()

	header := sekretariat.KontrakHeader{}

	err := json.NewDecoder(r.Body).Decode(&header)
	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}

	err = h.sekretariatSvc.CreateContract(ctx, header)
	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}

	log.Printf("[INFO][%s][%s] %s", r.RemoteAddr, r.Method, r.URL)
}

func (h *Handler) PrintKontrak(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.RenderJSON(w, r)

	company, _ := strconv.Atoi(r.FormValue("company"))
	ctx := r.Context()

	result, err := h.sekretariatSvc.PrintKontrak(ctx, company, r.FormValue("kontrak"))
	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename="+r.FormValue("kontrak")+".pdf")

	_, err = result.WriteTo(w)
	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
		defer resp.RenderJSON(w, r)
		log.Printf("[ERROR][%s][%s] %s | Reason: %s", r.RemoteAddr, r.Method, r.URL, err.Error())
		return
	}

	resp.Data = result

	log.Printf("[INFO][%s][%s] %s", r.RemoteAddr, r.Method, r.URL)
}
