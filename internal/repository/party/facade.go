package party_repository

import (
	"github.com/reinp/event-platform/backend/internal/database"
	baseRepository "github.com/reinp/event-platform/backend/internal/repository"
	partyCategoryRepository "github.com/reinp/event-platform/backend/internal/repository/party_category"
	ticketCategoryRepository "github.com/reinp/event-platform/backend/internal/repository/ticket_category"
)

type Facade struct {
	Repository *PartyRepository
	Write      *PartyWriteRepository
}

func NewFacade(
	db database.DBExecutor,
	transactionManager *database.TransactionManager,
	partyImages *baseRepository.PartyImageRepository,
	partyCategories *partyCategoryRepository.PartyCategoryWriteRepository,
	ticketCategories *ticketCategoryRepository.TicketCategoryWriteRepository,
) *Facade {

	return &Facade{
		Repository: NewPartyRepository(
			db,
		),

		Write: NewPartyWriteRepository(
			transactionManager,
			partyImages,
			partyCategories,
			ticketCategories,
		),
	}
}
