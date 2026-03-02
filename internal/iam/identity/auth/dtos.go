package auth

type VerifyEmailRequestDTO struct {
	Token string `form:"token" binding:"required"`
}

type SendResetPasswordEmailRequestDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type ResendVerificationEmailRequestDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequestDTO struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type RefreshRequestDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequestDTO struct {
	User         RegisterUserRequestDTO         `json:"user" binding:"required"`
	Organization RegisterOrganizationRequestDTO `json:"organization" binding:"required"`
}

type RegisterUserRequestDTO struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterOrganizationRequestDTO struct {
	Cellphone string `json:"cellphone" binding:"required,len=11"`
	Name      string `json:"name" binding:"required,min=2,max=200"`
}

type RegisterInputDTO struct {
	User         RegisterUserRequestDTO         `json:"user" binding:"required"`
	Organization RegisterOrganizationRequestDTO `json:"organization" binding:"required"`
}

type TokenPairDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
