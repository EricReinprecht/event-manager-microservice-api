package party_member_repository

import (
	"github.com/reinp/event-platform/backend/internal/database"
)

type Facade struct {
	Repository *PartyMemberRepository
	Write      *PartyMemberWriteRepository
}

func NewFacade(
	db database.DBExecutor,
	transactionManager *database.TransactionManager,
) *Facade {

	return &Facade{
		Repository: NewPartyMemberRepository(
			db,
		),

		Write: NewPartyMemberWriteRepository(
			transactionManager,
		),
	}
}
