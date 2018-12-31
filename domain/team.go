package domain

import (
	"github.com/kameike/karimono/model"
)

// Interfaces
type TeamDescriable interface {
	TeamNameProvider
	TeamIdProvider
}

type TeamPasswordUpdateRequester interface {
	TeamPasswordProvider
}

type TeamNameUpdateRequester interface {
	TeamNameProvider
}

type AccountDescriable interface {
	AccountIdProvider
}

type TeamDomain interface {
	UpdateTeamPassword(TeamPasswordUpdateRequester) (model.Team, error)
	UpdateTeamName(TeamNameUpdateRequester) (model.Team, error)
	KickAccount(AccountDescriable)
	GetTeamInfo() (model.Team, error)
	GetHistories() ([]model.Hisotry, error)
}
