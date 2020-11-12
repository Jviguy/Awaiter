# Awaiter
A Discordgo helper allowing for the use of AwaiteMessages(channelId) like thing
# Usage
```go
package main

import (
  "github.com/bwmarrin/discordgo"
  "github.com/Jviguy/Awaiter"
)

var awaiter Awaiter.MessageSendAwaiter

func main() {
  	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
  
  awaiter = Awaiter.NewMessageSendAwaiter(dg)
  
	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "!awaitMessage" log the next message
	if m.Content == "!awaitMessage" {
    //await a new message
    msg := awaiter.AwaitMessage(m.ChannelId)
    //print it
    fmt.Println(msg.Content)
  }
}
```
