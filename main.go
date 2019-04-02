package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages made by this bot
	if s.State.Ready.User.Username == m.Author.Username {
		return
	}

	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	if m.Content[:1] == "!" {
		channel, _ := s.Channel(m.ChannelID)
		serverID := channel.GuildID
		method := strings.Split(m.Content, " ")[0][1:]

		if method == "play" {
			results := ytSearch(strings.Split(m.Content, " ")[1:])

			if voiceInstances[serverID] != nil {
				voiceInstances[serverID].queueVideo(results.title)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Queued: %s", results.title))
			} else {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Playing: %s", results.title))
				go newVoiceInstance(serverID, channel.ID, results)
			}
		} else if method == "stop" && voiceInstances[serverID] != nil {
			voiceInstances[serverID].stopVideo()
		} else if method == "skip" && voiceInstances[serverID] != nil {
			voiceInstances[serverID].skipVideo()
		} else if method == "help" {
			s.ChannelMessageSend(m.ChannelID, `**!play** <youtube link or query> - Search/Play Youtube link, queues up if another track is playing
**!skip** - Skip current playing track
**!stop** - Stops tracks and clears queue`)
		}
	}
}

func main() {
	// NTYyMzg5MTAxODQwNTY0MjI0.XKKEhw.o_hhe-jNHROcRscH4XhUbgoKx8A
	discord, err := discordgo.New("Bot " + "NTYyMzg5MTAxODQwNTY0MjI0.XKKEhw.o_hhe-jNHROcRscH4XhUbgoKx8A")
	if err != nil {
		fmt.Println(os.Getenv("DISCORD_TOKEN"))
	}

	discord.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Listening...")
	lock := make(chan int)
	<-lock

	//search := os.Args[1:]
	//ytSearch(search)
}
