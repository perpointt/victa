package bot

const (
	CommandStart = "start"
	CommandHelp  = "help"
)

const (
	CallbackMainMenu         = "main_menu"
	CallbackDeleteMessage    = "delete_message"
	CallbackClearState       = "clear_state"
	CallbackConfirmOperation = "confirm_operation"
	CallbackBlank            = "blank"
)

const (
	CallbackCreateCompany             = "create_company"
	CallbackUpdateCompany             = "update_company"
	CallbackDeleteCompany             = "delete_company"
	CallbackListCompany               = "list_company"
	CallbackDetailCompany             = "detail_company"
	CallbackCompanyIntegrations       = "company_integrations"
	CallbackBackToDetailCompany       = "back_to_detail_company"
	CallbackUpdateCompanyIntegrations = "update_integrations"
)

const (
	CallbackCreateApp       = "create_app"
	CallbackDeleteApp       = "delete_app"
	CallbackUpdateApp       = "update_app"
	CallbackListApp         = "list_app"
	CallbackDetailApp       = "detail_app"
	CallbackAppIntegrations = "app_integrations"
)

const (
	CallbackInviteUser       = "invite_user"
	CallbackDeleteUser       = "delete_user"
	CallbackListUser         = "list_user"
	CallbackDetailUser       = "detail_user"
	CallbackBackToDetailUser = "back_to_detail_user"
)
