package domain

// UserSettings описывает сущность настроек пользователя.
type UserSettings struct {
	UserId          int64  `json:"user_id"`
	ActiveCompanyId *int64 `json:"active_company_id"`
}
