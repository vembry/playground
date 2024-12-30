package http

import (
	"app/internal/model"
	"app/internal/module"
	"app/internal/server/http/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	otelchimetric "github.com/riandyrn/otelchi/metric"
	"github.com/segmentio/ksuid"
	"go.opentelemetry.io/otel"
)

type handler struct {
	balanceModule module.IBalance
}

func New(balanceModule module.IBalance) *handler {
	return &handler{
		balanceModule: balanceModule,
	}
}

func (h *handler) GetHandler() http.Handler {
	var serverName = "app-go"

	// define base config for metric middlewares
	baseCfg := otelchimetric.NewBaseConfig(serverName, otelchimetric.WithMeterProvider(otel.GetMeterProvider()))

	// define router
	r := chi.NewRouter()
	r.Use(
		otelchi.Middleware(serverName, otelchi.WithChiRoutes(r)),
		otelchimetric.NewRequestDurationMillis(baseCfg),
		otelchimetric.NewRequestInFlight(baseCfg),
		otelchimetric.NewResponseSizeBytes(baseCfg),
	)

	r.Post("/balance/open", h.Open)
	r.Get("/balance/{balance_id}", h.Get)
	r.Post("/balance/{balance_id}/deposit", h.Deposit)
	r.Post("/balance/{balance_id}/withdraw", h.Withdraw)
	r.Post("/balance/{balance_id_from}/transfer/{balance_id_to}", h.Transfer)

	return r
}

// Open opens new balance. Basically creates new balance entry
func (h *handler) Open(w http.ResponseWriter, r *http.Request) {
	balance, err := h.balanceModule.Open(r.Context())
	if err != nil {
		util.RespondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	util.RespondJson(w, http.StatusOK, balance)
}

// Get gets balance by balance id
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	// get param
	balanceIdRaw := r.PathValue("balance_id")
	balanceId, _ := ksuid.Parse(balanceIdRaw)

	// call service
	balance, err := h.balanceModule.Get(r.Context(), balanceId)
	if err != nil {
		util.RespondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// return
	util.RespondJson(w, http.StatusOK, balance)
}

// Withdraw attempts to withdraw balance
func (h *handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	// read payload
	var in model.WithdrawParam
	err := util.ReadBody(r, in)
	if err != nil {
		util.RespondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// read path value
	balanceIdRaw := r.PathValue("balance_id")
	in.BalanceId, _ = ksuid.Parse(balanceIdRaw)

	// call service
	withdrawal, err := h.balanceModule.Withdraw(r.Context(), &in)
	if err != nil {
		util.RespondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// return
	util.RespondJson(w, http.StatusOK, withdrawal)
}

// Deposit attempts to deposit balance
func (h *handler) Deposit(w http.ResponseWriter, r *http.Request) {
	// read payload
	var in model.DepositParam
	err := util.ReadBody(r, in)
	if err != nil {
		util.RespondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// read path value
	balanceIdRaw := r.PathValue("balance_id")
	in.BalanceId, _ = ksuid.Parse(balanceIdRaw)

	// call service
	deposit, err := h.balanceModule.Deposit(r.Context(), &in)
	if err != nil {
		util.RespondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// return
	util.RespondJson(w, http.StatusOK, deposit)
}

// Transfer attempts to send balance from a balance id to another balance id
func (h *handler) Transfer(w http.ResponseWriter, r *http.Request) {
	// read payload
	var in model.TransferParam
	err := util.ReadBody(r, in)
	if err != nil {
		util.RespondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// read path value for balance-from
	balanceIdFromRaw := r.PathValue("balance_id_from")
	in.BalanceIdFrom, err = ksuid.Parse(balanceIdFromRaw)
	if err != nil {
		util.RespondErrorJson(w, http.StatusBadRequest, "invalid balance id from")
		return
	}

	// read path value for balance-to
	balanceIdToRaw := r.PathValue("balance_id_to")
	in.BalanceIdTo, err = ksuid.Parse(balanceIdToRaw)
	if err != nil {
		util.RespondErrorJson(w, http.StatusBadRequest, "invalid balance id to")
		return
	}

	// call service
	transfer, err := h.balanceModule.Transfer(r.Context(), &in)
	if err != nil {
		util.RespondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// return
	util.RespondJson(w, http.StatusOK, transfer)
}
