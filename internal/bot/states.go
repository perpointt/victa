package bot

type ChatState int

const (
	StateWaitingCreateCompany ChatState = iota
	StateWaitingUpdateCompany
	StateWaitingConfirmDeleteCompany
	StateWaitingConfirmDeleteUser
	StateWaitingUpdateCompanyIntegration
)
