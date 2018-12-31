package domain

import (
	"github.com/kameike/karimono/model"
)

// Interfaces
type TeamIdProvider interface {
	TeamId() string
}

type TeamNameProvider interface {
	TeamName() string
}

type TeamPasswordProvider interface {
	Password() string
}

type TeamCreateRequester interface {
	TeamNameProvider
	TeamPasswordProvider
	TeamIdProvider
}

type AccountDomain interface {
	CreateTeam(TeamCreateRequester) (model.Team, error)
	JoinTeam() (model.Team, error)
	LeaveTeam() error
	GetHistory() ([]model.Hisotry, error)
	GetBorrowings() ([]model.Hisotry, error)
}
