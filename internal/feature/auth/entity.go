package auth

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var (
	upperRegex   = regexp.MustCompile(`[A-Z]`)
	specialRegex = regexp.MustCompile(`[!@#$%^&*()\-_=+\[\]{};':"\\|,.<>\/?]`)
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,strong_password"`
	Name     string `json:"name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type GetNewTokensRequest struct {
	Token string `json:"token" validate:"required"`
}

type Tokens struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}
type VerifyAccountRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  int64  `json:"code" validate:"required"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}
type ConfirmResetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Code        int64  `json:"code" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,strong_password"`
}
type ResendVerifyCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return upperRegex.MatchString(password) && specialRegex.MatchString(password)
	})
	return v
}
