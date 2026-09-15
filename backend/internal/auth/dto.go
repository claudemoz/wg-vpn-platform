package auth

import "backend/internal/user"

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required,min=2,max=80"`
	LastName  string `json:"last_name"  binding:"required,min=2,max=80"`
	Email     string `json:"email"      binding:"required,email,max=160"`
	Phone     string `json:"phone"      binding:"omitempty,max=20"`
	Password  string `json:"password"   binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	User user.UserResponse `json:"user"`
}

type LoginResponse struct {
	Token string            `json:"token"`
	User  user.UserResponse `json:"user"`
}
