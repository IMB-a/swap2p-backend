package server

import (
	"net/http"

	"github.com/IMB-a/swap2p-backend/api"
)

func (s Server) GetAllTrades(w http.ResponseWriter, r *http.Request, params api.GetAllTradesParams) {
	// TODO implement me
	panic("implement me")
}

func (s Server) GetPersonalData(w http.ResponseWriter, r *http.Request, chatID api.PChatID) {

}

func (s Server) UpdateState(w http.ResponseWriter, r *http.Request, chatID api.PChatID, params api.UpdateStateParams) {
	// TODO implement me
	panic("implement me")
}

func (s Server) AddWallet(w http.ResponseWriter, r *http.Request, chatID api.PChatID, params api.AddWalletParams) {
	// TODO implement me
	panic("implement me")
}
