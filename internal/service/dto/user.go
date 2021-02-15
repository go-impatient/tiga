package dto

// UserRequest 请求数据结果
type UserRequest struct {
	ID       uint64  `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// UserResponse 返回数据结构
type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

