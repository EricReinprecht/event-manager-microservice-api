package helpers

import (
	"context"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/service"
)

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	FamilyID     uuid.UUID
}

func LoginUser(
	ctx context.Context,
	service *service.AuthService,
	userIdentifier string,
	password string,
) *AuthTokens {

	response, err := service.Login(
		ctx,
		userIdentifier,
		password,
		"test-agent",
		"127.0.0.1",
	)

	if err != nil {
		panic(err)
	}

	// family id is inside JWT,
	// for tests that need it we can decode later
	return &AuthTokens{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}
}
