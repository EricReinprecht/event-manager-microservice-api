package party_repository

import (
	"github.com/reinp/event-platform/backend/internal/database"
	partyCategoryRepository "github.com/reinp/event-platform/backend/internal/repository/party_category"
	partyImageRepository "github.com/reinp/event-platform/backend/internal/repository/party_image"
	ticketCategoryRepository "github.com/reinp/event-platform/backend/internal/repository/ticket_category"
)

type Facade struct {
	Repository *PartyRepository
	Write      *PartyWriteRepository
}

func NewFacade(
	db database.DBExecutor,
	transactionManager *database.TransactionManager,
	partyImages *partyImageRepository.PartyImageRepository,
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
