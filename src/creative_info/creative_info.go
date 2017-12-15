package creative_info

type CreativeInfo struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	Type      int    `json:"type"`
	Size      int64  `json:"size"`
	FailTimes int    `json:"fail_times"`
}
