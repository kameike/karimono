package handler

import (
	"github.com/kameike/karimono/domain"
	"github.com/kameike/karimono/error"
	"github.com/kameike/karimono/model"
)

type createTeamRequest struct {
	Name string `json:"name"`
	Pass string `json:"password"`
}

type createTeamResponse struct {
	Team model.Team `json:"team"`
}

func (r createTeamRequest) TeamName() string {
	return r.Name
}

func (r createTeamRequest) Password() string {
	return r.Pass
}

func createTeam(a domain.AccountDomain, h *Handler) {
	var req createTeamRequest
	h.bodyAsJson(&req)

	team, err := a.CreateTeam(req)
	if err != nil {
		h.renderError(err)
		return
	}

	h.renderJson(createTeamResponse{
		Team: *team,
	})
}

func returnBorrowing(a domain.AccountDomain, h *Handler) {
	req := h.context.Param("idHash")

	if req == "" {
		h.renderError(apperror.ApplicationError{apperror.ErrorRequestFormat})
		return
	}

	b, er := a.RetrunBorrowingWithHash(req)
	if er != nil {
		h.renderError(er)
		return
	}

	var res struct {
		Borrwoing model.Borrowing `json:"borrowing"`
	}
	res.Borrwoing = *b

	h.renderJson(res)
}

func joinTeam(a domain.AccountDomain, h *Handler) {
	var req createTeamRequest
	h.bodyAsJson(&req)

	team, err := a.JoinTeam(req)
	if err != nil {
		h.renderError(err)
		return
	}

	h.renderJson(createTeamResponse{
		Team: *team,
	})
}

type getHistoryResponse struct {
	histories []model.Hisotry
}

func getHistory(a domain.AccountDomain, h *Handler) {
	his, err := a.GetHistory()
	if err != nil {
		h.renderError(err)
		return
	}

	h.renderJson(getHistoryResponse{
		histories: his,
	})
}

type getBorrowingsResponse struct {
	borrowings []model.Borrowing
}

func getBorrowings(a domain.AccountDomain, h *Handler) {
	b, err := a.GetBorrowings()
	if err != nil {
		h.renderError(err)
		return
	}

	h.renderJson(getBorrowingsResponse{
		borrowings: b,
	})
}

type getTeamsResponse struct {
	Teams []model.Team `json:"teams"`
}

func getTeams(a domain.AccountDomain, h *Handler) {
	t, err := a.GetTeams()

	if err != nil {
		h.renderError(err)
		return
	}

	h.renderJson(getTeamsResponse{
		Teams: t,
	})
}
