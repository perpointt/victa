package domain

type GitlabWebhook struct {
	ObjectKind string `json:"object_kind"`
	User       struct {
		Name string `json:"name"`
	} `json:"user"`
	Project struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Homepage  string `json:"homepage"`
	} `json:"project"`
	ObjectAttributes Attributes `json:"object_attributes"`
	Changes          struct {
		CreatedAt   *DateTimeChange `json:"created_at"`
		UpdatedAt   *DateTimeChange `json:"updated_at"`
		ClosedAt    *DateTimeChange `json:"closed_at"`
		Description *struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"description"`
	} `json:"changes"`
	Issue Attributes `json:"issue"`
}

type Attributes struct {
	IID         int    `json:"iid"`
	Title       string `json:"title"`
	Note        string `json:"note"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Action      string `json:"action"`
	State       string `json:"state"`
}
