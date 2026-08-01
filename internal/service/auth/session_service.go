package auth_service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
)

type SessionService struct {
	users         *repository.UserRepository
	refreshTokens *repository.RefreshTokenRepository
	tokens        *TokenService
}

func NewSessionService(
	users *repository.UserRepository,
	refreshTokens *repository.RefreshTokenRepository,
	tokens *TokenService,
) *SessionService {

	return &SessionService{
		users:         users,
		refreshTokens: refreshTokens,
		tokens:        tokens,
	}
}

func (s *SessionService) Login(
	ctx context.Context,
	identifier string,
	password string,
	userAgent string,
	ipAddress string,
) (*TokenResponse, error) {

	identifier = strings.TrimSpace(identifier)

	user, err := s.users.FindByIdentifier(
		ctx,
		identifier,
	)

	if err != nil {
		return nil, errors.New(
			"invalid credentials",
		)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, errors.New(
			"invalid credentials",
		)
	}

	if user.VerifiedAt == nil {
		return nil, errors.New(
			"email not verified",
		)
	}

	return s.tokens.CreateSessionTokens(
		user,
		userAgent,
		ipAddress,
	)
}

func (s *SessionService) Logout(
	ctx context.Context,
	refreshToken string,
) error {

	tokenHash := auth.HashToken(
		refreshToken,
	)

	return s.refreshTokens.RevokeByHash(
		ctx,
		tokenHash,
	)
}

func (s *SessionService) LogoutAll(
	ctx context.Context,
	userID uuid.UUID,
) error {

	return s.refreshTokens.RevokeAllByUser(
		ctx,
		userID,
	)
}

func (s *SessionService) Sessions(
	ctx context.Context,
	userID uuid.UUID,
	currentFamily uuid.UUID,
) ([]SessionResponse, error) {

	tokens, err := s.refreshTokens.FindSessionsByUser(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	sessions := make(
		[]SessionResponse,
		0,
		len(tokens),
	)

	seenFamilies := make(
		map[uuid.UUID]struct{},
	)

	for _, token := range tokens {

		if _, exists := seenFamilies[token.FamilyID]; exists {
			continue
		}

		seenFamilies[token.FamilyID] = struct{}{}

		sessions = append(
			sessions,
			SessionResponse{
				FamilyID:  token.FamilyID,
				Device:    token.UserAgent,
				IP:        token.IPAddress,
				CreatedAt: token.CreatedAt,
				Current:   token.FamilyID == currentFamily,
			},
		)
	}

	return sessions, nil
}

func (s *SessionService) RevokeSession(
	ctx context.Context,
	userID uuid.UUID,
	familyID uuid.UUID,
) error {

	tokens, err := s.refreshTokens.FindSessionsByUser(
		ctx,
		userID,
	)

	if err != nil {
		return err
	}

	for _, token := range tokens {

		if token.FamilyID != familyID {
			continue
		}

		return s.refreshTokens.RevokeFamily(
			ctx,
			familyID,
		)
	}

	return errors.New(
		"session not found",
	)
}

func (s *SessionService) CreateSession(
	user *models.User,
	userAgent string,
	ipAddress string,
) (*TokenResponse, error) {

	return s.tokens.CreateSessionTokens(
		user,
		userAgent,
		ipAddress,
	)
}
