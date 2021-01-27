package Awaiter
//The base Entry type which is "overloaded" for each independent awaiter
type Entry interface {

	GetChannelId() string

	IncludeBots() bool
}

type MessageEntry interface {
	GetMessageId() string

	IncludeBots() bool
}
