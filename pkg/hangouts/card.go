package hangouts

// Card struct
type Card struct {
	Header   *Header    `json:"header,omitempty"`
	Sections []*Section `json:"sections"`
	Actions  []*Action  `json:"cardActions,omitempty"`
}

// Header struct
type Header struct {
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle,omitempty"`
	ImageURL   string `json:"imageUrl,omitempty"`
	ImageStyle string `json:"imageStyle,omitempty"`
}

// Section struct
type Section struct {
	Header  string    `json:"header,omitempty"`
	Widgets []*Widget `json:"widgets,omitempty"`
}

// Widget struct
type Widget struct {
	TextParagraph *TextParagraph `json:"textParagraph,omitempty"`
	KeyValue      *KeyValue      `json:"keyValue,omitempty"`
}

// ActionWidget struct
type ActionWidget struct {
	Buttons []Button `json:"buttons"`
}

// TextParagraph struct
type TextParagraph struct {
	Text string `json:"text"`
}

// KeyValue struct
type KeyValue struct {
	TopLabel         string   `json:"topLabel,omitempty"`
	Content          string   `json:"content,omitempty"`
	ContentMultiline bool     `json:"contentMultiline,omitempty"`
	BottomLabel      string   `json:"bottomLabel,omitempty"`
	OnClick          *OnClick `json:"onClick,omitempty"`
	Icon             string   `json:"icon,omitempty"`
	Button           *Button  `json:"button,omitempty"`
}

// Button struct
type Button struct {
	TextButton  *TextButton  `json:"textButton,omitempty"`
	ImageButton *ImageButton `json:"imageButton,omitempty"`
}

// TextButton struct
type TextButton struct {
	Text    string   `json:"text,omitempty"`
	OnClick *OnClick `json:"onClick,omitempty"`
}

// ImageButton struct
type ImageButton struct {
	IconURL string   `json:"iconUrl,omitempty"`
	Icon    string   `json:"icon,omitempty"`
	OnClick *OnClick `json:"onClick,omitempty"`
}

// Action struct
type Action struct {
	Label   string   `json:"actionLabel,omitempty"`
	OnClick *OnClick `json:"onClick,omitempty"`
}

// OnClick struct
type OnClick struct {
	OpenLink *OpenLink `json:"openLink,omitempty"`
}

// OpenLink struct
type OpenLink struct {
	URL string `json:"url,omitempty"`
}

// Image struct
type Image struct {
	ImageURL string   `json:"imageUrl,omitempty"`
	OnClick  *OnClick `json:"onClick,omitempty"`
}
