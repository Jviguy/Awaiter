package Awaiter

import "github.com/bwmarrin/discordgo"

//An Awaiter for awaiting messages to be sent in a said channel
type MessageDeleteAwaiter struct {
	//A slice of Entries belonging to this Awaiter
	session *discordgo.Session
	//The *discordgo.Session this awaiter belongs to
	Entries []MessageDeleteEntry
}

//A Entry for the MessageSendAwaiter
type MessageDeleteEntry struct {
	//the channelId the Entry is for
	channelId string
	//The channel to return to when a message appears in the specified channel named by channelId
	channel chan *discordgo.Message
}

//Returns the channel that will be returned too
func (m *MessageDeleteEntry) GetChannel() chan *discordgo.Message {
	return m.channel
}

//inherited from Entry and returns the ChannelId used in the Entry
func (m *MessageDeleteEntry) GetChannelId() string {
	return m.channelId
}

//Adds a Entry to the MessageDeleteAwaiter and returns the message when it is received.
func (m *MessageDeleteAwaiter) AwaitDeletedMessage(channelId string) *discordgo.Message {
	//Make the channel
	channel := make(chan *discordgo.Message)
	//Form the Entry
	entry := MessageDeleteEntry{channelId,channel}
	//Add the Entry
	m.Await(entry)
	//return the recieved message from the channel
	return <-channel
}

//Adds a Entry to the MessageDeleteAwaiter and has to be manually handled.
func (m *MessageDeleteAwaiter) Await(entry MessageDeleteEntry)  {
	m.Entries = append(m.Entries,entry)
}

//Returns the *discordgo.Session that the Awaiter has been added to.
func (m *MessageDeleteAwaiter) GetSession() *discordgo.Session {
	return m.session
}

//Removes a Entry from a MessageDeleteAwaiter
func (m *MessageDeleteAwaiter) RemoveEntry(k int) {
	//Get the Entry from the list
	entry := m.Entries[k]
	//close the channel
	close(entry.GetChannel())
	//remove the entry from the List
	m.Entries[k] = m.Entries[len(m.Entries)-1]
	m.Entries[len(m.Entries)-1] = MessageDeleteEntry{}
	m.Entries = m.Entries[:len(m.Entries)-1]
}

//The function added to the *discordgo.Session called when a message is sent
func (m *MessageDeleteAwaiter) handle(s *discordgo.Session,msg *discordgo.MessageDelete) {
	for k,entry := range m.Entries {
		if entry.GetChannelId() == msg.ChannelID{
			entry.channel <- msg.BeforeDelete
			m.RemoveEntry(k)
		}
	}
}

