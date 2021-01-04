package main

type TgUpdateResult struct {
	Ok bool `json:"ok"`
	Result []TgUpdate `json:"result"`
}

type TgUpdate struct {
	UpdateId int64 `json:"update_id"`
	Message TgMessage `json:"message"`
}

type TgMessage struct {
	MessageId int64 `json:"message_id"`
	From TgFrom `json:"from"`
	Chat TgChat `json:"chat"`
	Date int64 `json:"date"`
	Text string `json:"text"`
	Entities []TgEntities `json:"entities"`
}

type TgFrom struct {
	Id int64 `json:"id"`
	IsBot bool `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type TgChat struct {
	Id int64 `json:"id"`
	FirstName string `json:"first_name"`
	Username string `json:"username"`
	Type string `json:"type"`
}

type TgEntities struct {
	Offset int64 `json:"offset"`
	Length int64 `json:"length"`
	Type string `json:"type"`
}


