package dto

type UserUpdateRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"`
}

type UserQueryRequest struct {
	Page   int    `form:"page" binding:"omitempty,numeric,gte=1"`
	Limit  int    `form:"limit" binding:"omitempty,numeric,gte=1"`
	Search string `form:"search"`
	Sort   string `form:"sort" binding:"omitempty,oneof=full_name_asc full_name_desc email_asc email_desc"`
}

