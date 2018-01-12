package creative_info

type CreativeInfo struct {
	Id        string `json:"id,omitempty"`
	Url       string `json:"url,omitempty"`
	Type      int    `json:"type,omitempty"`
	Size      int64  `json:"size,omitempty"`
	FailTimes int    `json:"fail_times,omitempty"`
}
