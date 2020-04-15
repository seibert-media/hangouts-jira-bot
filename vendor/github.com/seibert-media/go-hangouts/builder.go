package hangouts

// NewMessage message with builder functions
func NewMessage() *Message {
	return &Message{}
}

// WithText adds text to function
func (m *Message) WithText(t string) *Message {
	m.Text = t
	return m
}

// WithSender adds a User to Message
func (m *Message) WithSender(u User) *Message {
	m.Sender = u
	return m
}

// InThread defines the thread a message should be posted in
func (m *Message) InThread(name string) *Message {
	m.Thread = &Thread{
		Name: name,
	}
	return m
}

// WithCard adds a new card to Message
func (m *Message) WithCard(c *Card) *Message {
	m.Cards = append(m.Cards, c)
	return m
}

// NewCard for being used in Message
func NewCard() *Card {
	return &Card{}
}

// WithHeader adds the defined header values to card
func (c *Card) WithHeader(title, subtitle, imageURL, imageStyle string) *Card {
	c.Header = &Header{
		Title:      title,
		Subtitle:   subtitle,
		ImageURL:   imageURL,
		ImageStyle: imageStyle,
	}
	return c
}

// WithSection adds a new Section to Card
func (c *Card) WithSection(s *Section) *Card {
	c.Sections = append(c.Sections, s)
	return c
}

// WithAction adds a new Action to Card
func (c *Card) WithAction(a *Action) *Card {
	c.Actions = append(c.Actions, a)
	return c
}

// NewSection for being used in Card
func NewSection(header string) *Section {
	return &Section{
		Header: header,
	}
}

// WithWidget adds a new Widget to Section
func (s *Section) WithWidget(w *Widget) *Section {
	s.Widgets = append(s.Widgets, w)
	return s
}

// NewWidget which may **either** contain TextParagraph, KeyValue or Buttons
func NewWidget() *Widget {
	return &Widget{}
}

// WithTextParagraph adds text to widget
func (w *Widget) WithTextParagraph(t string) *Widget {
	w.TextParagraph = &TextParagraph{
		Text: t,
	}
	return w
}

// WithKeyValue adds KeyValue object to widget
func (w *Widget) WithKeyValue(kv *KeyValue) *Widget {
	w.KeyValue = kv
	return w
}

// WithButton adds new Button to Widget
func (w *Widget) WithButton(b *Button) *Widget {
	w.Buttons = append(w.Buttons, b)
	return w
}

// WithImage adds new Image to Widget
func (w *Widget) WithImage(i *Image) *Widget {
	w.Image = i
	return w
}

// NewKeyValue for use in Widget
// TopLabel and Content are required
func NewKeyValue(topLabel, content string, multiline bool) *KeyValue {
	return &KeyValue{
		TopLabel:         topLabel,
		Content:          content,
		ContentMultiline: multiline,
	}
}

// WithBottomLabel for KeyValue
func (kv *KeyValue) WithBottomLabel(l string) *KeyValue {
	kv.BottomLabel = l
	return kv
}

// WithOpenLink adds OnClick containing an OpenLink action to KeyValue
func (kv *KeyValue) WithOpenLink(url string) *KeyValue {
	kv.OnClick = &OnClick{
		Action: nil,
		OpenLink: &OpenLink{
			URL: url,
		},
	}
	return kv
}

// WithAction adds OnClick containing a FormAction to KeyValue
func (kv *KeyValue) WithAction(name string, parameters ...*ActionParameter) *KeyValue {
	kv.OnClick = &OnClick{
		Action: &FormAction{
			MethodName: name,
			Parameters: parameters,
		},
		OpenLink: nil,
	}
	return kv
}

// NewActionParameter for use in KeyValue
func NewActionParameter(key, value string) *ActionParameter {
	return &ActionParameter{
		Key:   key,
		Value: value,
	}
}

// WithButton adds new Button to KeyValue
func (kv *KeyValue) WithButton(b *Button) *KeyValue {
	kv.Button = b
	return kv
}

// NewTextLinkButton containing OpenLink
func NewTextLinkButton(text, url string) *Button {
	return &Button{
		TextButton: &TextButton{
			Text: text,
			OnClick: &OnClick{
				Action: nil,
				OpenLink: &OpenLink{
					URL: url,
				},
			},
		},
	}
}

// NewTextActionButton containing FormAction
func NewTextActionButton(text, action string, parameters ...*ActionParameter) *Button {
	return &Button{
		TextButton: &TextButton{
			Text: text,
			OnClick: &OnClick{
				Action: &FormAction{
					MethodName: action,
					Parameters: parameters,
				},
				OpenLink: nil,
			},
		},
	}
}

// NewImageLinkButton containing OpenLink
func NewImageLinkButton(icon, iconURL, url string) *Button {
	return &Button{
		ImageButton: &ImageButton{
			Icon:    icon,
			IconURL: iconURL,
			OnClick: &OnClick{
				Action: nil,
				OpenLink: &OpenLink{
					URL: url,
				},
			},
		},
	}
}

// NewImageActionButton containing FormAction
func NewImageActionButton(icon, iconURL, action string, parameters ...*ActionParameter) *Button {
	return &Button{
		ImageButton: &ImageButton{
			Icon:    icon,
			IconURL: iconURL,
			OnClick: &OnClick{
				Action: &FormAction{
					MethodName: action,
					Parameters: parameters,
				},
				OpenLink: nil,
			},
		},
	}
}

// NewImage with URL and AspectRatio
func NewImage(url string, ratio int) *Image {
	return &Image{
		ImageURL:    url,
		AspectRatio: ratio,
	}
}

// WithOpenLink adds OnClick with OpenLink to Image
func (i *Image) WithOpenLink(url string) *Image {
	i.OnClick = &OnClick{
		Action: nil,
		OpenLink: &OpenLink{
			URL: url,
		},
	}
	return i
}

// WithAction adds OnClick with FormAction to Image
func (i *Image) WithAction(name string, parameters ...*ActionParameter) *Image {
	i.OnClick = &OnClick{
		Action: &FormAction{
			MethodName: name,
			Parameters: parameters,
		},
		OpenLink: nil,
	}
	return i
}
