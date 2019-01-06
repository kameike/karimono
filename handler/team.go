package handler

import (
	"github.com/kameike/karimono/domain"
	"github.com/kameike/karimono/model"
	"github.com/labstack/echo"
	// 	. "github.com/kameike/karimono/error"
)

type borrowingRequest struct {
	Item     string `json:"item"`
	Memo_    string `json:"memo"`
	TeamName string `json:"teamName"`
}

func (r borrowingRequest) ItemName() string {
	return r.Item
}

func (r borrowingRequest) Memo() string {
	return r.Memo_
}

func CreateBorrowing(c echo.Context) error {
	h := createHandler(c)

	var req borrowingRequest
	h.bodyAsJson(&req)

	t, err := h.provider.GetTeamProviderViaTeamName(req.TeamName)

	if err != nil {
		h.renderError(err)
		return nil
	}

	item, err := t.BorrowItem(req)

	if err != nil {
		h.renderError(err)
		return nil
	}

	h.renderJson(item)

	return nil
}

func getTeamBorrowing(t domain.TeamDomain, h *Handler) {
	b, e := t.GetTeamBorrowings()
	if e != nil {
		h.renderError(e)
	}

	type res struct {
		Borrowing []model.Borrowing `json:"borrowings"`
	}

	h.renderJson(res{b})
}

func getTeamMenbers(t domain.TeamDomain, h *Handler) {
	tm, e := t.GetTeamMenbers()
	if e != nil {
		h.renderError(e)
	}

	type res struct {
		Members []model.Account `json:"menbers"`
	}

	h.renderJson(res{tm})
}
