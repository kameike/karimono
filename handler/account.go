package handler

import (
	"github.com/kameike/karimono/domain"
	"github.com/kameike/karimono/model"
)

type createTeamRequest struct {
	name     string
	password string
}

type createTeamResponse struct {
	team model.Team
}

func (r createTeamRequest) TeamName() string {
	return r.name
}

func (r createTeamRequest) Password() string {
	return r.password
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
		team: *team,
	})
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
		team: *team,
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
