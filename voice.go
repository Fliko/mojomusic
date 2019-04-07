package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/bwmarrin/discordgo"
)

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

func (vi *voiceInstance) queueVideo(fileName string) {
	tmp := append((*vi.queue), fileName+".dca")
	vi.queue = &tmp
}
func (vi *voiceInstance) skipVideo() {}
func (vi *voiceInstance) stopVideo() {}

// loads a dca file from disk into instance's soundBuffer
func (vi *voiceInstance) loadSound() error {
	file, err := os.Open((*vi.queue)[0])
	if err != nil {
		logger(err, "Failed to open music file")
	}
	defer file.Close()

	var frameLen int16

	for {
		// Read frame length
		err = binary.Read((io.Reader)(file), binary.LittleEndian, &frameLen)

		// If eof return
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}

		if err != nil {
			logger(err, "failed to load sound file length")
		}

		inBuf := make([]byte, frameLen)
		err = binary.Read(file, binary.LittleEndian, &inBuf)
		if err != nil {
			logger(err, "failed  to load sound file")
		}
		vi.soundBuffer = append(vi.soundBuffer, inBuf)
	}
}

func (vi *voiceInstance) playSound() error {
	err := vi.discord.Open()
	if err != nil {
		logger(err, "Failed to connect to discord server")
	}
	defer vi.discord.Close()
	fmt.Println(vi.guildID, vi.channelID)
	chn, err := vi.discord.ChannelVoiceJoin(vi.guildID, vi.channelID, false, true)
	fmt.Println("CHANNEL", chn, err)
	if err != nil {
		logger(err, "Failed to connect Voice")
	}
	// Notify discord that you are about to speak
	//vi.ChannelVoiceJoin(vi.guildID, vi. chnnelID, false, true)
	chn.Speaking(true)
	defer chn.Speaking(false)

	// speak some
	for _, buff := range vi.soundBuffer {
		chn.OpusSend <- buff
	}

	tmp := (*vi.queue)[1:]
	vi.queue = &tmp
	fmt.Println("FUCKN ELL")
	return nil
}

func newVoiceInstance(guild *discordgo.Guild, author string, result ytResult) {
	vi := new(voiceInstance)
	voiceInstances[guild.ID] = vi

	vi.guildID = guild.ID
	vi.queue = &([]string{})
	vi.soundBuffer = make([][]byte, 0)

	vi.discord, _ = discordgo.New("Bot " + "NTYyMzg5MTAxODQwNTY0MjI0.XKKEhw.o_hhe-jNHROcRscH4XhUbgoKx8A")
	vi.queueVideo(result.title + "-" + result.id)
	vi.loadSound()
	// Look for the message sender in that guild's current voice states.
	for _, vs := range guild.VoiceStates {
		if vs.UserID == author {
			vi.channelID = vs.ChannelID
			err := vi.playSound()
			if err != nil {
				fmt.Println("Error playing sound:", err)
			}

			return
		}
	}
	//vi.playSound()
}
