package request

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
	Dob      string `json:"dob" binding:"required"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
	Dob      string `json:"dob" binding:"required"`
}
