package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/security"
)

type AuthService struct {
	users             *repository.UserRepository
	verifications     *repository.EmailVerificationRepository
	refreshTokens     *repository.RefreshTokenRepository
	passwordResets    *repository.PasswordResetTokenRepository
	jwt               *auth.JWT
	emailService      EmailSender
	passwordValidator *security.PasswordValidator

	refreshTokenDuration  time.Duration
	passwordResetCooldown time.Duration
}

func NewAuthService(
	users *repository.UserRepository,
	verifications *repository.EmailVerificationRepository,
	refreshTokens *repository.RefreshTokenRepository,
	passwordResets *repository.PasswordResetTokenRepository,
	jwt *auth.JWT,
	emailService EmailSender,
	passwordValidator *security.PasswordValidator,
	refreshTokenDuration time.Duration,
	passwordResetCooldown time.Duration,
) *AuthService {

	log.Printf("AUTH USERS REPOSITORY: %#v", users)

	if users == nil {
		panic("AUTH SERVICE CREATED WITH NIL USERS REPOSITORY")
	}

	return &AuthService{
		users:                 users,
		verifications:         verifications,
		refreshTokens:         refreshTokens,
		passwordResets:        passwordResets,
		jwt:                   jwt,
		emailService:          emailService,
		passwordValidator:     passwordValidator,
		refreshTokenDuration:  refreshTokenDuration,
		passwordResetCooldown: passwordResetCooldown,
	}
}

type RegisterRequest struct {
	Email    string
	Password string
	Username string
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}

type SessionResponse struct {
	FamilyID  uuid.UUID `json:"familyId"`
	Device    string    `json:"device"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
	Current   bool      `json:"current"`
}

type ForgotPasswordRequest struct {
	Identifier string
}

type ResetPasswordRequest struct {
	Token       string
	NewPassword string
}

func (s *AuthService) Register(
	ctx context.Context,
	req RegisterRequest,
) (*models.User, error) {

	// normalize input
	req.Email = strings.TrimSpace(
		strings.ToLower(req.Email),
	)

	req.Username = strings.TrimSpace(
		req.Username,
	)

	// validate username
	if err := security.ValidateUsername(
		req.Username,
	); err != nil {
		return nil, err
	}

	// check email duplicate
	_, err := s.users.FindByEmail(
		ctx,
		req.Email,
	)

	if err == nil {
		return nil, appErrors.ErrEmailAlreadyExists
	}

	// check username duplicate
	_, err = s.users.FindByUsername(
		ctx,
		req.Username,
	)

	if err == nil {
		return nil, appErrors.ErrUsernameAlreadyExists
	}

	// validate password
	if err := s.passwordValidator.Validate(
		req.Password,
		req.Username,
		req.Email,
	); err != nil {
		return nil, err
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Username:     req.Username,
	}

	token := auth.GenerateToken()

	err = s.users.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			userRepo :=
				repository.NewUserRepository(tx)

			verificationRepo :=
				repository.NewEmailVerificationRepository(tx)

			// create user
			if err := userRepo.Create(
				ctx,
				user,
			); err != nil {
				return err
			}

			// create verification token
			verification := &models.EmailVerification{
				UserID: user.ID,
				Token: auth.HashToken(
					token,
				),
				ExpiresAt: time.Now().Add(
					24 * time.Hour,
				),
			}

			if err := verificationRepo.Create(
				verification,
			); err != nil {
				return err
			}

			// email failure should rollback registration
			if err := s.emailService.SendVerificationEmail(
				user.Email,
				user.Username,
				token,
			); err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	identifier string,
	password string,
	userAgent string,
	ip string,
) (*TokenResponse, error) {

	log.Printf("AUTH SERVICE POINTER: %p", s)
	log.Printf("USERS REPO POINTER: %p", s.users)

	user, err := s.users.FindByIdentifier(
		ctx,
		identifier,
	)

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.VerifiedAt == nil {
		return nil, errors.New(
			"email not verified",
		)
	}

	// create refresh token family
	familyID := uuid.New()

	// create access token with family information
	accessToken, err := s.jwt.Generate(
		user.ID.String(),
		familyID.String(),
	)

	if err != nil {
		return nil, err
	}

	// create refresh token
	refreshToken := auth.GenerateToken()

	refresh := &models.RefreshToken{

		UserID: user.ID,

		TokenHash: auth.HashToken(
			refreshToken,
		),

		FamilyID: familyID,

		UserAgent: userAgent,

		IPAddress: ip,

		DeviceName: "Unknown",

		ExpiresAt: time.Now().Add(
			s.refreshTokenDuration,
		),
	}

	err = s.refreshTokens.Create(
		refresh,
	)

	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Secret() string {
	return s.jwt.Secret
}

func (s *AuthService) ValidateUser(
	ctx context.Context,
	userID uuid.UUID,
) (*models.User, error) {

	return s.users.FindByID(
		ctx,
		userID,
	)
}

func (s *AuthService) VerifyEmail(
	ctx context.Context,
	token string,
) (string, error) {

	hashedToken := auth.HashToken(
		token,
	)

	verification, err := s.verifications.FindByToken(
		ctx,
		hashedToken,
	)

	if err != nil {
		return "", errors.New(
			"invalid verification token",
		)
	}

	if verification.ExpiresAt.Before(time.Now()) {
		return "", errors.New(
			"verification expired",
		)
	}

	user, err := s.users.FindByID(
		ctx,
		verification.UserID,
	)

	if err != nil {
		return "", err
	}

	// already verified
	if user.VerifiedAt != nil {

		familyID := uuid.New()

		return s.jwt.Generate(
			user.ID.String(),
			familyID.String(),
		)
	}

	now := time.Now()

	user.VerifiedAt = &now

	err = s.users.Update(
		ctx,
		user,
	)

	if err != nil {
		return "", err
	}

	err = s.verifications.Delete(
		verification.ID,
	)

	if err != nil {
		return "", err
	}

	familyID := uuid.New()

	return s.jwt.Generate(
		user.ID.String(),
		familyID.String(),
	)
}

func (s *AuthService) Refresh(
	ctx context.Context,
	refreshToken string,
) (*TokenResponse, error) {

	hash := auth.HashToken(
		refreshToken,
	)

	storedToken, err := s.refreshTokens.FindByHash(
		ctx,
		hash,
	)

	if err != nil {
		return nil, errors.New(
			"invalid refresh token",
		)
	}

	if storedToken.RevokedAt != nil {

		// refresh token replay detected
		_ = s.refreshTokens.RevokeFamily(
			ctx,
			storedToken.FamilyID,
		)

		return nil, errors.New(
			"refresh token replay detected",
		)
	}

	if storedToken.ExpiresAt.Before(time.Now()) {
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

	newRefreshToken := auth.GenerateToken()

	refresh := &models.RefreshToken{

		UserID: user.ID,

		TokenHash: auth.HashToken(
			newRefreshToken,
		),

		// keep same device/session
		FamilyID: storedToken.FamilyID,

		UserAgent: storedToken.UserAgent,

		IPAddress: storedToken.IPAddress,

		DeviceName: storedToken.DeviceName,

		ExpiresAt: time.Now().Add(
			s.refreshTokenDuration,
		),
	}

	err = s.refreshTokens.Transaction(
		ctx,
		func(tx database.DBExecutor) error {

			repo := repository.NewRefreshTokenRepository(tx)

			if err := repo.Revoke(
				storedToken.ID,
			); err != nil {
				return err
			}

			return repo.Create(
				refresh,
			)
		},
	)

	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwt.Generate(
		user.ID.String(),
		storedToken.FamilyID.String(),
	)

	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Logout(
	ctx context.Context,
	refreshToken string,
) error {
	hash := auth.HashToken(refreshToken)

	return s.refreshTokens.RevokeByHash(
		ctx,
		hash,
	)
}

func (s *AuthService) Sessions(
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

	seen := make(
		map[uuid.UUID]bool,
	)

	for _, token := range tokens {

		// one session = one family
		if seen[token.FamilyID] {
			continue
		}

		seen[token.FamilyID] = true

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

func (s *AuthService) RevokeSession(
	ctx context.Context,
	userID uuid.UUID,
	familyID uuid.UUID,
) error {

	// security check:
	// only allow deleting own sessions

	tokens, err := s.refreshTokens.FindSessionsByUser(
		ctx,
		userID,
	)

	if err != nil {
		return err
	}

	found := false

	for _, token := range tokens {

		if token.FamilyID == familyID {
			found = true
			break
		}
	}

	if !found {
		return errors.New(
			"session not found",
		)
	}

	return s.refreshTokens.RevokeFamily(
		ctx,
		familyID,
	)
}

func (s *AuthService) CreateSession(
	user *models.User,
	userAgent string,
	ip string,
) (*TokenResponse, error) {

	familyID := uuid.New()

	accessToken, err := s.jwt.Generate(
		user.ID.String(),
		familyID.String(),
	)

	if err != nil {
		return nil, err
	}

	refreshToken := auth.GenerateToken()

	refresh := &models.RefreshToken{
		UserID: user.ID,

		TokenHash: auth.HashToken(
			refreshToken,
		),

		FamilyID: familyID,

		UserAgent: userAgent,

		IPAddress: ip,

		DeviceName: "Unknown",

		ExpiresAt: time.Now().Add(
			s.refreshTokenDuration,
		),
	}

	if err := s.refreshTokens.Create(refresh); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) LogoutAll(
	ctx context.Context,
	userID uuid.UUID,
) error {

	return s.refreshTokens.RevokeAllByUser(
		ctx,
		userID,
	)
}

func (s *AuthService) ForgotPassword(
	ctx context.Context,
	identifier string,
) error {

	user, err := s.users.FindByIdentifier(
		ctx,
		identifier,
	)

	if err != nil {

		// do not reveal if user exists
		return nil
	}

	// invalidate previous reset tokens
	// cooldown check
	latest, err := s.passwordResets.FindLatestByUser(
		ctx,
		user.ID,
	)

	if err == nil {

		if time.Since(latest.CreatedAt) < s.passwordResetCooldown {
			// silently ignore to prevent email spam
			return nil
		}
	}

	// invalidate previous reset tokens
	err = s.passwordResets.InvalidateForUser(
		ctx,
		user.ID,
	)

	if err != nil {
		return err
	}

	rawToken := auth.GenerateToken()

	reset := &models.PasswordResetToken{
		UserID: user.ID,

		TokenHash: auth.HashToken(
			rawToken,
		),

		ExpiresAt: time.Now().Add(
			15 * time.Minute,
		),
	}

	err = s.passwordResets.Create(
		ctx,
		reset,
	)

	if err != nil {
		return err
	}

	err = s.emailService.SendPasswordResetEmail(
		user.Email,
		user.Username,
		rawToken,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResetPassword(
	ctx context.Context,
	token string,
	newPassword string,
) error {

	hash := auth.HashToken(
		token,
	)

	resetToken, err := s.passwordResets.FindByHash(
		ctx,
		hash,
	)

	if err != nil {

		return errors.New(
			"invalid reset token",
		)
	}

	if resetToken.InvalidatedAt != nil {

		return errors.New(
			"reset token invalidated",
		)
	}

	if resetToken.UsedAt != nil {

		return errors.New(
			"reset token already used",
		)
	}

	if resetToken.ExpiresAt.Before(
		time.Now(),
	) {

		return errors.New(
			"reset token expired",
		)
	}

	user, err := s.users.FindByID(
		ctx,
		resetToken.UserID,
	)

	if err != nil {
		return errors.New(
			"cannot reset password for deleted user",
		)
	}

	// validate password rules
	if err := s.passwordValidator.Validate(
		newPassword,
		user.Username,
		user.Email,
	); err != nil {

		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user.PasswordHash = string(passwordHash)

	err = s.users.Update(
		ctx,
		user,
	)

	if err != nil {
		return err
	}

	// mark reset token as used
	err = s.passwordResets.MarkUsed(
		ctx,
		resetToken.ID,
	)

	if err != nil {
		return err
	}

	// logout all devices
	err = s.refreshTokens.RevokeAllByUser(
		ctx,
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResendVerificationEmail(
	ctx context.Context,
	email string,
) error {

	email = strings.TrimSpace(
		strings.ToLower(email),
	)

	user, err := s.users.FindByEmail(
		ctx,
		email,
	)

	if err != nil {
		return errors.New("invalid request")
	}

	if user.DeletedAt.Valid {
		return errors.New("invalid request")
	}

	if user.VerifiedAt != nil {
		return errors.New(
			"email already verified",
		)
	}

	// prevent email spam
	lastVerification, err := s.verifications.FindLatestByUser(
		ctx,
		user.ID,
	)

	if err == nil {

		cooldown := 5 * time.Minute

		if lastVerification.CreatedAt.After(
			time.Now().Add(-cooldown),
		) {
			return errors.New(
				"please wait before requesting another verification email",
			)
		}
	}

	// invalidate previous verification tokens

	err = s.verifications.InvalidateForUser(
		ctx,
		user.ID,
	)

	if err != nil {
		return err
	}

	// create new raw token

	rawToken := auth.GenerateToken()

	verification := &models.EmailVerification{

		UserID: user.ID,

		// store hash only
		Token: auth.HashToken(
			rawToken,
		),

		ExpiresAt: time.Now().Add(
			24 * time.Hour,
		),
	}

	err = s.verifications.Create(
		verification,
	)

	if err != nil {
		return err
	}

	return s.emailService.SendVerificationEmail(
		user.Email,
		user.Username,
		rawToken,
	)
}
