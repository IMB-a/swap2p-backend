package server

import (
	"net/http"

	"github.com/IMB-a/swap2p-backend/api"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

func (s *Server) InitPersonalData(w http.ResponseWriter, r *http.Request, chatID api.PChatID) {
	ctx := r.Context()
	err := s.db.UpsertPerson(ctx, string(chatID))
	if err != nil {
		respond(w, r, errors.Wrap(err, "upsert person"))
		return
	}
}

func (s Server) GetAllTrades(w http.ResponseWriter, r *http.Request, params api.GetAllTradesParams) {
	// TODO implement me
	panic("implement me")
}

func (s Server) GetPersonalData(w http.ResponseWriter, r *http.Request, chatID api.PChatID) {
	ctx := r.Context()
	pd, err := s.db.GetPersonalData(ctx, string(chatID))
	if err != nil {
		respond(w, r, err)
		return
	}
	respond(w, r, pd)
}

func (s Server) UpdateState(w http.ResponseWriter, r *http.Request, chatID api.PChatID, params api.UpdateStateParams) {
	// TODO implement me
	panic("implement me")
}

func (s Server) AddWallet(w http.ResponseWriter, r *http.Request, chatID api.PChatID, params api.AddWalletParams) {
	// TODO implement me
	panic("implement me")
}

func errHandler(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	render.DefaultResponder(w, r, api.ErrorResp{Error: err.Error()})
}

func respond(w http.ResponseWriter, r *http.Request, payload interface{}) {
	if err, ok := payload.(error); ok {
		errHandler(w, r, err)
		return
	}
	render.DefaultResponder(w, r, payload)
}
