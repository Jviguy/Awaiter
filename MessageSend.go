package Awaiter

import "github.com/bwmarrin/discordgo"

//An Awaiter for awaiting messages to be sent in a said channel
type MessageSendAwaiter struct {

	//A slice of Entries belonging to this Awaiter
	Entries []MessageSendEntry

	//The *discordgo.Session this awaiter belongs to
	session *discordgo.Session
}

//A Entry for the MessageSendAwaiter
type MessageSendEntry struct {
	//the channelId the Entry is for
	channelId string
	//The channel to return to when a message appears in the specified channel named by channelId
	channel chan *discordgo.MessageCreate
	//if true the if bot check will be ignored
	includeBots bool
}

//Returns wether to include bots in the await
func (m MessageSendEntry) IncludeBots() bool {
	return m.includeBots
}

//Returns the channel that will be returned too
func (m MessageSendEntry) GetChannel() chan *discordgo.MessageCreate {
	return m.channel
}

//inherited from Entry and returns the ChannelId used in the Entry
func (m MessageSendEntry) GetChannelId() string {
	return m.channelId
}

//Adds a Entry to the MessageSendAwaiter and has to be manually handled.
func (m *MessageSendAwaiter) Await(entry MessageSendEntry) {
	m.Entries = append(m.Entries,entry)
}

//Adds a Entry to the MessageSendAwaiter and returns the message when it is received.
func (m *MessageSendAwaiter) AwaitMessage(channelId string,IncludeBots bool) *discordgo.MessageCreate {
	//Make the channel
	channel := make(chan *discordgo.MessageCreate)
	//Form the Entry
	entry := MessageSendEntry{channelId,channel,IncludeBots}
	//Add the Entry
	m.Await(entry)
	//return the recieved message from the channel
	return <-channel
}

//Returns the *discordgo.Session that the Awaiter has been added to.
func (m *MessageSendAwaiter) GetSession() *discordgo.Session {
	return m.session
}

//Removes a Entry from a MessageSendAwaiter
func (m *MessageSendAwaiter) RemoveEntry(k int) {
	//Get the Entry from the list
	entry := m.Entries[k]
	//close the channel
	close(entry.GetChannel())
	//remove the entry from the List
	m.Entries[k] = m.Entries[len(m.Entries)-1]
	m.Entries[len(m.Entries)-1] = MessageSendEntry{}
	m.Entries = m.Entries[:len(m.Entries)-1]
}

//The function added to the *discordgo.Session called when a message is sent
func (m *MessageSendAwaiter) handle(s *discordgo.Session,msg *discordgo.MessageCreate)  {
	if msg.Author.ID == s.State.User.ID {
		return
	}
	for k , entry := range m.Entries {
		if entry.GetChannelId() == msg.ID{
			if !entry.IncludeBots() {
				if msg.Author.Bot {
					return
				}
			}
			entry.GetChannel() <- msg
			m.RemoveEntry(k)
		}
	}
}

//Initializes a new MessageSendAwaiter ready for use
func NewMessageSendAwaiter(s *discordgo.Session) *MessageSendAwaiter{
	awaiter := &MessageSendAwaiter{session: s,Entries: make([]MessageSendEntry,0)}
	s.AddHandler(awaiter.handle)
	return awaiter
}