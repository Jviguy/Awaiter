package Awaiter

import "github.com/bwmarrin/discordgo"

type Awaiter interface {
	GetSession() *discordgo.Session
	RemoveEntry(k int)
}
