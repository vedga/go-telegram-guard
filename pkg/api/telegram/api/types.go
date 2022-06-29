package api

// Location represent user location
type Location struct {
	// Longitude of sender's
	Longitude float32 `json:"longitude"`
	// Latitude of sender's
	Latitude float32 `json:"latitude"`
}

// User represent telegram user or bot
type User struct {
	// ID is unique user or bot identifier
	ID int `json:"message_id"`
	// Name represent first name
	Name string `json:"first_name"`
	// LastName is last name if any
	LastName string `json:"last_name,omitempty"`
	// UserName is user or bot username
	UserName string `json:"username,omitempty"`
}

// Chat represent conversation place
type Chat struct {
	// ID is unique chat ID
	ID int `json:"id"`
	// Type represent conversation place type: “private”, “group”, “supergroup”, “channel”
	Type string `json:"type"`
	// Title represent group or channel title
	Title *string `json:"title,omitempty"`
	// UserName represent user name for some channels
	UserName *string `json:"username,omitempty"`
	// Name represent first name
	Name *string `json:"first_name,omitempty"`
	// LastName is last name if any
	LastName *string `json:"last_name,omitempty"`
	// AllMembersAdministrators true if all members are administators
	AllMembersAdministrators bool `json:"all_members_are_administrators,omitempty"`
}

// MessageEntity represent message entities
type MessageEntity struct {
	// Type of the entity.
	// One of mention (@username), hashtag, bot_command, url, email, bold (bold text), italic (italic text),
	// code (monowidth string), pre (monowidth block), text_link (for clickable text URLs)
	Type string `json:"type"`
	// Offset in UTF-16 code units to the start of the entity
	Offset int `json:"offset"`
	// Length of the entity in UTF-16 code units
	Length int `json:"length"`
	// Link used for “text_link” only, url that will be opened after user taps on the text
	Link string `json:"url,omitempty"`
}

// Message represent conversation message
type Message struct {
	// ID is unique message ID
	ID int `json:"message_id"`
	// From is message sender, may be empty in channels
	From *User `json:"from,omitempty"`
	// UnixTime represent unix time when message is created
	UnixTime int `json:"date,omitempty"`
	// ConversationPlace metadata
	ConversationPlace Chat `json:"chat,omitempty"`
	// ForwardFrom original message issuer if any
	ForwardFrom *User `json:"forward_from,omitempty"`
	// ForwardUnixTime represent unix time when original message is created
	ForwardUnixTime int `json:"forward_date,omitempty"`
	// Reply is reference to original message
	Reply *Message `json:"reply_to_message,omitempty"`
	// Text is message content
	Text *string `json:"text,omitempty"`
	// Entities contain optional message entities
	Entities []MessageEntity `json:"entities,omitempty"`
}

// InlineQuery represent inline query
type InlineQuery struct {
	// ID is query ID
	ID string `json:"id"`
	// From contain query issuer
	From User `json:"from"`
	// Location contain query issuer's location
	Location *Location `json:"location,omitempty"`
	// Query contain text of query
	Query string `json:"query"`
	// Offset of the results to be returned, can be controlled by the bot
	Offset string `json:"offset"`
}

// Update represent conversation update
type Update struct {
	// ID is update unique ID
	ID int `json:"update_id,omitempty"`
	// Message contain conversation message if any
	Message *Message `json:"message,omitempty"`
	// InlineQuery contain incoming inline query
	InlineQuery *InlineQuery `json:"inline_query,omitempty"`
}

// SendMessage represent message to send. On success, the sent Message is returned.
type SendMessage struct {
	// WebHookMessageReply must contain "sendMessage" if used inside webhook reply
	WebHookMessageReply *string `json:"method,omitempty"`
	// Target is unique identifier for the target chat or username of the target channel (in the format @channelusername)
	// Can be string or integer
	Target string `json:"chat_id"`
	// Text is message content
	Text string `json:"text"`
	// ParseMode Send Markdown or HTML, if you want Telegram apps to show bold, italic, fixed-width text or inline URLs in your bot's message.
	ParseMode *string `json:"parse_mode,omitempty"`
	// DisableWebPagePreview disables link previews for links in this message
	DisableWebPagePreview *bool `json:"disable_web_page_preview,omitempty"`
	// DisableNotification Sends the message silently. iOS users will not receive a notification, Android users will receive a notification with no sound.
	DisableNotification *bool `json:"disable_notification,omitempty"`
	// ReplyToMessageID if the message is a reply, ID of the original message
	ReplyToMessageID *int `json:"reply_to_message_id,omitempty"`
	// ReplyMarkup Additional interface options. A JSON-serialized object for an inline keyboard, custom reply keyboard,
	// instructions to hide reply keyboard or to force a reply from the user.
	ReplyMarkup *string `json:"reply_markup,omitempty"`
}
