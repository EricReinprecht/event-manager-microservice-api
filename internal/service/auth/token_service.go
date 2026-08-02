package auth_service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/clock"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
	refreshTokenRepository "github.com/reinp/event-platform/backend/internal/repository/refresh_token"
)

const defaultDeviceName = "Unknown"

type TokenService struct {
	users              *repository.UserRepository
	refreshTokens      *repository.RefreshTokenRepository
	refreshTokenWriter *refreshTokenRepository.RefreshTokenWriteRepository
	jwt                *auth.JWT
	clock              clock.Clock

	refreshTokenDuration time.Duration
}

func NewTokenService(
	users *repository.UserRepository,
	refreshTokens *repository.RefreshTokenRepository,
	refreshTokenWriter *refreshTokenRepository.RefreshTokenWriteRepository,
	jwt *auth.JWT,
	clock clock.Clock,
	refreshTokenDuration time.Duration,
) *TokenService {

	return &TokenService{
		users:                users,
		refreshTokens:        refreshTokens,
		refreshTokenWriter:   refreshTokenWriter,
		jwt:                  jwt,
		clock:                clock,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (s *TokenService) now() time.Time {
	return s.clock.Now()
}

func (s *TokenService) Secret() string {
	return s.jwt.Secret
}

func (s *TokenService) GenerateAccessToken(
	userID uuid.UUID,
	familyID uuid.UUID,
) (string, error) {

	return s.jwt.Generate(
		userID.String(),
		familyID.String(),
	)
}

func (s *TokenService) CreateRefreshToken(
	userID uuid.UUID,
	familyID uuid.UUID,
	userAgent string,
	ipAddress string,
	deviceName string,
) (*models.RefreshToken, string) {

	rawToken := auth.GenerateToken()

	refreshToken := &models.RefreshToken{
		UserID: userID,

		TokenHash: auth.HashToken(
			rawToken,
		),

		FamilyID: familyID,

		UserAgent: userAgent,

		IPAddress: ipAddress,

		DeviceName: deviceName,

		ExpiresAt: s.now().Add(
			s.refreshTokenDuration,
		),
	}

	return refreshToken, rawToken
}

func (s *TokenService) CreateSessionTokens(
	user *models.User,
	userAgent string,
	ipAddress string,
) (*TokenResponse, error) {

	familyID := uuid.New()

	accessToken, err := s.GenerateAccessToken(
		user.ID,
		familyID,
	)

	if err != nil {
		return nil, err
	}

	refreshToken, rawRefreshToken :=
		s.CreateRefreshToken(
			user.ID,
			familyID,
			userAgent,
			ipAddress,
			defaultDeviceName,
		)

	if err := s.refreshTokens.Create(
		refreshToken,
	); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
	}, nil
}

func (s *TokenService) Refresh(
	ctx context.Context,
	rawRefreshToken string,
) (*TokenResponse, error) {

	tokenHash := auth.HashToken(
		rawRefreshToken,
	)

	storedToken, err := s.refreshTokens.FindByHash(
		ctx,
		tokenHash,
	)

	if err != nil {
		return nil, errors.New(
			"invalid refresh token",
		)
	}

	if storedToken.RevokedAt != nil {

		_ = s.refreshTokens.RevokeFamily(
			ctx,
			storedToken.FamilyID,
		)

		return nil, errors.New(
			"refresh token replay detected",
		)
	}

	if storedToken.ExpiresAt.Before(s.now()) {
		return nil, errors.New(
			"refresh token expired",
		)
	}

	user, err := s.users.FindByID(
		ctx,
		storedToken.UserID,
	)

	if err != nil {
		return nil, errors.New(
			"invalid user",
		)
	}

	newRefreshToken, rawNewRefreshToken :=
		s.CreateRefreshToken(
			user.ID,
			storedToken.FamilyID,
			storedToken.UserAgent,
			storedToken.IPAddress,
			storedToken.DeviceName,
		)

	err = s.refreshTokenWriter.Rotate(
		ctx,
		storedToken,
		newRefreshToken,
	)

	if err != nil {
		return nil, err
	}

	accessToken, err := s.GenerateAccessToken(
		user.ID,
		storedToken.FamilyID,
	)

	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: rawNewRefreshToken,
	}, nil
}
