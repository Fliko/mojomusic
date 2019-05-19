package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	mojoroutes "github.com/Fliko/mojoMusic/mojoroutes"
	"github.com/bwmarrin/discordgo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

const (
	initialServers = "localhost:8080"
	resolverType   = "manual"
)

var conn *grpc.ClientConn

func init() {
	// Register balancer type
	opts := []grpc.DialOption{grpc.WithBalancerName(roundrobin.Name), grpc.WithInsecure()}
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

	c := mojoroutes.NewRouteClient(conn)
	fmt.Println("Starting the mojo client!", c)
	listenForDiscordEvents(c)
}

var voiceInstances = map[string]*voiceInstance{}

type voiceInstance struct {
	discord     *discordgo.Session
	queue       *[]string
	soundBuffer [][]byte
	channelID   string
	guildID     string
	skip        bool
	stop        bool
	playing     bool
}

func newVoiceInstance(guild *discordgo.Guild) *voiceInstance {
	vi := new(voiceInstance)
	voiceInstances[guild.ID] = vi

	vi.guildID = guild.ID
	vi.queue = &([]string{})
	vi.soundBuffer = make([][]byte, 0)

	vi.discord, _ = discordgo.New("Bot " + "NTYyMzg5MTAxODQwNTY0MjI0.XKKEhw.o_hhe-jNHROcRscH4XhUbgoKx8A")
	return vi
}
func (vi *voiceInstance) playSong(author string) {
	time.Sleep(2 * time.Second)

}
func (vi *voiceInstance) stopSong() {}
func (vi *voiceInstance) skipSong() {}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	c := mojoroutes.NewRouteClient(conn)

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

		switch method {
		case "play":
			// If the server is already in our list assume it is playing and queue song
			if voiceInstances[guild.ID] != nil {
				*voiceInstances[guild.ID].queue = append(*voiceInstances[guild.ID].queue, keyWords)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Queued: %s", keyWords))
			} else {
				voiceInstances[guild.ID] = newVoiceInstance(guild)
				*voiceInstances[guild.ID].queue = append(*voiceInstances[guild.ID].queue, keyWords)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Playing: %s"))
				go voiceInstances[guild.ID].playSong(m.Author.Username)
			}
		case "stop":
			go voiceInstances[guild.ID].stopSong()
		case "skip":
			go voiceInstances[guild.ID].skipSong()
		case "help":
			s.ChannelMessageSend(m.ChannelID, `**!play** <youtube link or query> - Search/Play Youtube link, queues up if another track is playing
**!skip** - Skip current playing track
**!stop** - Stops tracks and clears queue`)
		}
	}
}

func listenForDiscordEvents(c mojoroutes.RouteClient) {
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

	fmt.Println("RUNNIN THIS STREAMIN SHIT")

	req := &mojoroutes.GreetManyTimesRequest{
		Greeting: &mojoroutes.Greeting{
			Name: "Fkin strmin man",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalln("FUCKKKKKKKKK")
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("SREAMD:LKJH:LKJ")
		}
		log.Println(msg.GetResult())
	}
}
