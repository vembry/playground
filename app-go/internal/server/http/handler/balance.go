package handler

import (
	"app/internal/model"
	"app/internal/module"
	"net/http"

	"github.com/segmentio/ksuid"
)

type balance struct {
	balanceModule module.IBalance
}

func NewBalance(balanceModule module.IBalance) *balance {
	return &balance{
		balanceModule: balanceModule,
	}
}

func (b *balance) GetMux() *http.ServeMux {
	balancemux := http.NewServeMux()

	balancemux.HandleFunc("POST /open", b.Open)
	balancemux.HandleFunc("GET /:balance_id", b.Get)
	balancemux.HandleFunc("POST /:balance_id/deposit", b.Deposit)
	balancemux.HandleFunc("POST /:balance_id/withdraw", b.Withdraw)
	balancemux.HandleFunc("POST /:balance_id_from/transfer/:balance_id_to", b.Transfer)

	return groupMux("/balance", balancemux)
}

// Open opens new balance. Basically creates new balance entry
func (b *balance) Open(w http.ResponseWriter, r *http.Request) {
	balance, err := b.balanceModule.Open(r.Context())
	if err != nil {
		respondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJson(w, http.StatusOK, balance)
}

// Get gets balance by balance id
func (b *balance) Get(w http.ResponseWriter, r *http.Request) {
	// get param
	balanceIdRaw := r.PathValue("balance_id")
	balanceId, _ := ksuid.Parse(balanceIdRaw)

	// call service
	balance, err := b.balanceModule.Get(r.Context(), balanceId)
	if err != nil {
		respondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// return
	respondJson(w, http.StatusOK, balance)
}

// Withdraw attempts to withdraw balance
func (b *balance) Withdraw(w http.ResponseWriter, r *http.Request) {
	// read payload
	var in model.WithdrawParam
	err := readBody(r, in)
	if err != nil {
		respondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// read path value
	balanceIdRaw := r.PathValue("balance_id")
	in.BalanceId, _ = ksuid.Parse(balanceIdRaw)

	// call service
	withdrawal, err := b.balanceModule.Withdraw(r.Context(), &in)
	if err != nil {
		respondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// return
	respondJson(w, http.StatusOK, withdrawal)
}

// Deposit attempts to deposit balance
func (b *balance) Deposit(w http.ResponseWriter, r *http.Request) {
	// read payload
	var in model.DepositParam
	err := readBody(r, in)
	if err != nil {
		respondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// read path value
	balanceIdRaw := r.PathValue("balance_id")
	in.BalanceId, _ = ksuid.Parse(balanceIdRaw)

	// call service
	deposit, err := b.balanceModule.Deposit(r.Context(), &in)
	if err != nil {
		respondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// return
	respondJson(w, http.StatusOK, deposit)
}

// Transfer attempts to send balance from a balance id to another balance id
func (b *balance) Transfer(w http.ResponseWriter, r *http.Request) {
	// read payload
	var in model.TransferParam
	err := readBody(r, in)
	if err != nil {
		respondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// read path value for balance-from
	balanceIdFromRaw := r.PathValue("balance_id_from")
	in.BalanceIdFrom, err = ksuid.Parse(balanceIdFromRaw)
	if err != nil {
		respondErrorJson(w, http.StatusBadRequest, "invalid balance id from")
		return
	}

	// read path value for balance-to
	balanceIdToRaw := r.PathValue("balance_id_to")
	in.BalanceIdTo, err = ksuid.Parse(balanceIdToRaw)
	if err != nil {
		respondErrorJson(w, http.StatusBadRequest, "invalid balance id to")
		return
	}

	// call service
	transfer, err := b.balanceModule.Transfer(r.Context(), &in)
	if err != nil {
		respondErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	// return
	respondJson(w, http.StatusOK, transfer)
}
