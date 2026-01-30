package model

// User 用户模型
type User struct {
	BaseModel
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Nickname string `gorm:"type:varchar(50)" json:"nickname"`
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"`
	Status   Status `gorm:"type:tinyint;default:1" json:"status"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname" binding:"max=50"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname *string `json:"nickname" binding:"omitempty,max=50"`
	Avatar   *string `json:"avatar" binding:"omitempty,url"`
	Status   *Status `json:"status" binding:"omitempty,oneof=0 1"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      *User  `json:"user"`
}

// UserListQuery 用户列表查询参数
type UserListQuery struct {
	PageQuery
	SortQuery
	Username string  `form:"username"`
	Email    string  `form:"email"`
	Status   *Status `form:"status"`
}
