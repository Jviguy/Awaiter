package Awaiter

import "github.com/bwmarrin/discordgo"

type MessageReactionAddAwaiter struct {
	//A slice of Entries belonging to this Awaiter
	Entries []MessageReactionAddEntry

	session *discordgo.Session
}

func (m *MessageReactionAddAwaiter) GetSession() *discordgo.Session {
	return m.session
}

//The function added to the *discordgo.Session called when a reaction is sent
func (m *MessageReactionAddAwaiter) handle(s *discordgo.Session, reaction *discordgo.MessageReactionAdd)  {
	if reaction.UserID == s.State.User.ID {
		return
	}
	for k , entry := range m.Entries {
		if entry.GetMessageId() == reaction.MessageID {
			if !entry.IncludeBots() {
				if user,_ := s.User(reaction.MessageID); user.Bot {
					return
				}
			}
			entry.GetChannel() <- reaction
			m.RemoveEntry(k)
		}
	}
}

func (m *MessageReactionAddAwaiter) RemoveEntry(k int) {
	//Get the Entry from the list
	entry := m.Entries[k]
	//close the channel
	close(entry.GetChannel())
	//remove the entry from the List
	m.Entries[k] = m.Entries[len(m.Entries)-1]
	m.Entries[len(m.Entries)-1] = MessageReactionAddEntry{}
	m.Entries = m.Entries[:len(m.Entries)-1]
}

type MessageReactionAddEntry struct {
	messageId string

	channel chan *discordgo.MessageReactionAdd

	includeBot bool
}

func (m MessageReactionAddEntry) GetMessageId() string {
	return m.messageId
}

func (m MessageReactionAddEntry) IncludeBots() bool {
	return m.includeBot
}

func (m MessageReactionAddEntry) GetChannel() chan *discordgo.MessageReactionAdd {
	return m.channel
}
