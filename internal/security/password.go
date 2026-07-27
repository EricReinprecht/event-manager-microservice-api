package security

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	zxcvbn "github.com/nbutton23/zxcvbn-go"
)

type PasswordValidator struct {
	client *http.Client
}

func NewPasswordValidator() *PasswordValidator {

	return &PasswordValidator{
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (p *PasswordValidator) Validate(
	password string,
	username string,
	email string,
) error {

	var validationErrors []string

	passwordNormalized := normalize(password)
	lower := strings.ToLower(password)

	/*
		1. Length
	*/

	if len(password) < 12 {

		validationErrors = append(
			validationErrors,
			"password must be at least 12 characters",
		)
	}

	if len(password) > 128 {

		validationErrors = append(
			validationErrors,
			"password cannot exceed 128 characters",
		)
	}

	/*
		2. Character rules
	*/

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {

		switch {

		case unicode.IsUpper(char):
			hasUpper = true

		case unicode.IsLower(char):
			hasLower = true

		case unicode.IsDigit(char):
			hasNumber = true

		case unicode.IsPunct(char) ||
			unicode.IsSymbol(char):

			hasSpecial = true
		}
	}

	if !hasUpper {

		validationErrors = append(
			validationErrors,
			"password must contain an uppercase letter",
		)
	}

	if !hasLower {

		validationErrors = append(
			validationErrors,
			"password must contain a lowercase letter",
		)
	}

	if !hasNumber {

		validationErrors = append(
			validationErrors,
			"password must contain a number",
		)
	}

	if !hasSpecial {

		validationErrors = append(
			validationErrors,
			"password must contain a special character",
		)
	}

	/*
		3. Personal information protection
	*/

	if containsPersonalInfo(
		passwordNormalized,
		username,
	) {

		validationErrors = append(
			validationErrors,
			"password cannot contain username",
		)
	}

	emailName, _, _ := strings.Cut(
		email,
		"@",
	)

	if containsPersonalInfo(
		passwordNormalized,
		emailName,
	) {

		validationErrors = append(
			validationErrors,
			"password cannot contain email",
		)
	}

	/*
		4. Forbidden words
	*/

	for _, word := range forbiddenPasswordWords {

		if strings.Contains(
			passwordNormalized,
			word,
		) {

			validationErrors = append(
				validationErrors,
				"password contains forbidden words",
			)

			break
		}
	}

	/*
		5. Pattern checks
	*/

	for _, pattern := range keyboardPatterns {

		if strings.Contains(
			lower,
			pattern,
		) {

			validationErrors = append(
				validationErrors,
				"password contains a common pattern",
			)

			break
		}
	}

	for _, pattern := range sequentialPatterns {

		if strings.Contains(
			lower,
			pattern,
		) {

			validationErrors = append(
				validationErrors,
				"password contains sequential characters",
			)

			break
		}
	}

	/*
		6. Repeated characters
	*/

	if hasRepeatedCharacters(password) {

		validationErrors = append(
			validationErrors,
			"password contains too many repeated characters",
		)
	}

	/*
		7. Breached password check
	*/

	compromised, err := p.IsCompromised(
		password,
	)

	if err != nil {

		log.Printf(
			"password breach check failed: %v",
			err,
		)

	} else if compromised {

		validationErrors = append(
			validationErrors,
			"password has appeared in a data breach",
		)
	}

	/*
		8. Password strength
	*/

	result := zxcvbn.PasswordStrength(
		password,
		[]string{
			username,
			email,
			"event",
			"platform",
			"ticket",
			"party",
		},
	)

	if result.Score < 2 {

		validationErrors = append(
			validationErrors,
			"password is too weak",
		)
	}

	if len(validationErrors) > 0 {

		return &PasswordValidationError{
			Field:  "password",
			Errors: validationErrors,
		}
	}

	return nil
}

type PasswordValidationError struct {
	Field  string
	Errors []string
}

func (e *PasswordValidationError) Error() string {

	return strings.Join(
		e.Errors,
		", ",
	)
}

func (p *PasswordValidator) IsCompromised(
	password string,
) (bool, error) {

	hash := sha1.Sum(
		[]byte(password),
	)

	hashString := strings.ToUpper(
		hex.EncodeToString(
			hash[:],
		),
	)

	prefix := hashString[:5]

	suffix := hashString[5:]

	resp, err := p.client.Get(
		"https://api.pwnedpasswords.com/range/" + prefix,
	)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(
		resp.Body,
	)

	if err != nil {
		return false, err
	}

	for _, line := range strings.Split(
		string(body),
		"\n",
	) {

		parts := strings.Split(
			line,
			":",
		)

		if len(parts) != 2 {
			continue
		}

		if strings.TrimSpace(
			parts[0],
		) == suffix {

			return true, nil
		}
	}

	return false, nil
}

func containsPersonalInfo(
	password string,
	value string,
) bool {

	if value == "" {
		return false
	}

	parts := strings.FieldsFunc(
		value,
		func(r rune) bool {

			return r == '.' ||
				r == '-' ||
				r == '_' ||
				r == ' ' ||
				r == '@'
		},
	)

	for _, part := range parts {

		part = normalize(part)

		if len(part) >= 3 &&
			strings.Contains(
				password,
				part,
			) {

			return true
		}
	}

	return false
}

func hasRepeatedCharacters(
	password string,
) bool {

	runes := []rune(password)

	count := 1

	for i := 1; i < len(runes); i++ {

		if runes[i] == runes[i-1] {

			count++

			if count >= 4 {
				return true
			}

		} else {

			count = 1
		}
	}

	return false
}

func normalize(
	value string,
) string {

	value = strings.ToLower(
		value,
	)

	replacer := strings.NewReplacer(
		"@", "a",
		"$", "s",
		"!", "i",
		"0", "o",
		"1", "i",
		"3", "e",
		"5", "s",
		"7", "t",
	)

	value = replacer.Replace(
		value,
	)

	value = strings.ReplaceAll(
		value,
		".",
		"",
	)

	value = strings.ReplaceAll(
		value,
		"-",
		"",
	)

	value = strings.ReplaceAll(
		value,
		"_",
		"",
	)

	value = strings.ReplaceAll(
		value,
		" ",
		"",
	)

	return value
}

var keyboardPatterns = []string{

	"qwerty",
	"asdfgh",
	"zxcvbn",

	"123456",
	"234567",

	"abcdef",
	"bcdefg",
}

var sequentialPatterns = []string{

	"abcdefghijklmnopqrstuvwxyz",

	"0123456789",
}

var forbiddenPasswordWords = []string{

	"password",
	"pass",
	"secret",
	"admin",
	"welcome",
	"letmein",
	"changeme",
	"login",
}
