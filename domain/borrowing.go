package domain

import (
	"github.com/kameike/karimono/model"
)

// Interfaces
type BorrowItemProvider interface {
	ItemName() string
}

type BorrowItemRequester interface {
	BorrowItemProvider
}

type ReturnItemRequester interface {
	BorrowItemProvider
}

type BorrowingDomain interface {
	GetBorrowings() ([]model.Hisotry, error)
	BorrowItem(BorrowItemRequester) (model.Borrowing, error)
	ReturnItem() (model.Borrowing, error)
}

type TeamDomainRequester interface {
	TeamIdProvider
}
