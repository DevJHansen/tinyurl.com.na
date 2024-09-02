package models

type NewRedirectReqBody struct {
	Target string `json:"target"`
}

type Redirect struct {
	Owner     string `json:"owner"`
	Id        int    `json:"id"`
	Key       string `json:"key"`
	Target    string `json:"target"`
	CreatedAt string `json:"created_at"`
	Deleted   bool   `json:"deleted"`
}

type RedirectsCache map[string]string

type RedirectAnalytic struct {
	Id         int    `json:"id"`
	TargetKey  string `json:"target_key"`
	CreatedAt  string `json:"created_at"`
	UserAgent  string `json:"user_agent"`
	IpAddress  string `json:"ip_address"`
	Referrer   string `json:"referrer"`
	DeviceType string `json:"device_type"`
	Os         string `json:"os"`
	Flagged    bool   `json:"flagged"`
}
