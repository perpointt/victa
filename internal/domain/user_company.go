package domain

// UserCompany описывает сущность настроек пользователя.
type UserCompany struct {
	UserID    int64 `json:"user_id"`
	CompanyID int64 `json:"company_id"`
	RoleID    int64 `json:"role_id"`
}
