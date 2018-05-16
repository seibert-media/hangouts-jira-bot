package hangouts

// User struct
type User struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl"`
	Email       string `json:"email"`
}

// Message struct
type Message struct {
	Name       string  `json:"name,omitempty"`
	Sender     User    `json:"sender,omitempty"`
	CreateTime string  `json:"createTime,omitempty"`
	Text       string  `json:"text,omitempty"`
	Thread     *Thread `json:"thread,omitempty"`
	Cards      []*Card `json:"cards,omitempty"`
}

// Thread struct
type Thread struct {
	Name string `json:"name,omitempty"`
}

// Event struct
type Event struct {
	Type  string `json:"type"`
	Space struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"space"`
	Message Message `json:"message,omitempty"`
	User    User    `json:"user"`
}
