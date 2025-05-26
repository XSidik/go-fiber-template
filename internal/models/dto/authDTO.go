package dto

type Register struct {
	UserName string `json:"user_name" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=6,max=15,uppercase,lowercase,containsany=!@#$%^&*"`
}

type Login struct {
	UserName string `json:"user_name" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required"`
}

type RefreshToken struct {
	Token string `json:"token" validate:"required"`
}

type Logout struct {
	Token string `json:"token" validate:"required"`
}
