package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages made by this bot
	if s.State.Ready.User.Username == m.Author.Username {
		return
	}

	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	if m.Content[:1] == "!" {
		channel, _ := s.State.Channel(m.ChannelID)
		guild, _ := s.State.Guild(channel.GuildID)
		method := strings.Split(m.Content, " ")[0][1:]
		search := strings.Split(m.Content, " ")[1:]
		keyWords := strings.Join(search[:], " ")
		fmt.Println(channel.ID, guild.ID)
		if method == "play" {
			fmt.Println(keyWords, m.Content, strings.Split(m.Content, " "))
			results := ytSearch(keyWords)

			if voiceInstances[guild.ID] != nil {
				voiceInstances[guild.ID].queueVideo(results.title)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Queued: %s", results.title))
			} else {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Playing: %s", results.title))
				go newVoiceInstance(guild, m.Author.ID, results)
			}
		} else if method == "stop" && voiceInstances[guild.ID] != nil {
			voiceInstances[guild.ID].stopVideo()
		} else if method == "skip" && voiceInstances[guild.ID] != nil {
			voiceInstances[guild.ID].skipVideo()
		} else if method == "help" {
			s.ChannelMessageSend(m.ChannelID, `**!play** <youtube link or query> - Search/Play Youtube link, queues up if another track is playing
**!skip** - Skip current playing track
**!stop** - Stops tracks and clears queue`)
		}
	}
}

const (
	initialServers = "localhost:8080"
	resolverType   = "manual"
)

var conn *grpc.ClientConn

func init() {
	// Register balancer type
	opts := []grpc.DialOption{grpc.WithBalancerName(roundrobin.Name)}
	// Register servers to attach to
	b, _ := manual.GenerateAndRegisterManualResolver()
	addresses := []resolver.Address{}
	for _, addr := range strings.Split(initialServers, ",") {
		addresses = append(addresses, resolver.Address{Addr: addr, Type: resolver.Backend})
	}
	b.InitialAddrs(addresses)
	servers := b
	resolver.Register(servers)
	resolver.SetDefaultScheme(servers.Scheme())
	// Connect
	conn, _ = grpc.Dial("Connected", opts...)
}

func main() {
	defer conn.Close()
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
	// Listen for ^c
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	fmt.Println("Listening")
	<-stop
	fmt.Println("[main] stopping")
}
