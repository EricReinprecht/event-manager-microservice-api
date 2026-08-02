package auth_service

import (
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"

	appErrors "github.com/reinp/event-platform/backend/internal/appErrors"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/security"
)

type RegistrationService struct {
	users             *repository.UserRepository
	verifications     *repository.EmailVerificationRepository
	verification      *VerificationService
	emailService      EmailSender
	passwordValidator *security.PasswordValidator
	unitOfWork        *repository.AuthUnitOfWork
}

func NewRegistrationService(
	users *repository.UserRepository,
	verifications *repository.EmailVerificationRepository,
	verification *VerificationService,
	emailService EmailSender,
	passwordValidator *security.PasswordValidator,
	unitOfWork *repository.AuthUnitOfWork,
) *RegistrationService {

	return &RegistrationService{
		users:             users,
		verifications:     verifications,
		verification:      verification,
		emailService:      emailService,
		passwordValidator: passwordValidator,
		unitOfWork:        unitOfWork,
	}
}

func (s *RegistrationService) Register(
	ctx context.Context,
	req RegisterRequest,
) (*models.User, error) {

	req.Email = strings.TrimSpace(
		strings.ToLower(req.Email),
	)

	req.Username = strings.TrimSpace(
		req.Username,
	)

	if err := security.ValidateUsername(
		req.Username,
	); err != nil {
		return nil, err
	}

	if _, err := s.users.FindByEmail(
		ctx,
		req.Email,
	); err == nil {
		return nil, appErrors.ErrEmailAlreadyExists
	}

	if _, err := s.users.FindByUsername(
		ctx,
		req.Username,
	); err == nil {
		return nil, appErrors.ErrUsernameAlreadyExists
	}

	if err := s.passwordValidator.Validate(
		req.Password,
		req.Username,
		req.Email,
	); err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Username:     req.Username,
	}

	err = s.unitOfWork.Transaction(
		ctx,
		func(repositories *repository.AuthTransactionRepositories) error {
			if err := repositories.Users.Create(
				ctx,
				user,
			); err != nil {
				return err
			}

			verification, rawToken :=
				s.verification.NewVerification(
					user.ID,
				)

			if err := repositories.Verifications.Create(
				verification,
			); err != nil {
				return err
			}

			return s.emailService.SendVerificationEmail(
				user.Email,
				user.Username,
				rawToken,
			)
		},
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
