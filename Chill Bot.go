package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	//"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

const xp_multiplier = 1.0

const defaultprefix string = "cb!"

// Append this string with error messages
var error_string string

//bot id
var BotID string

var VC *discordgo.VoiceConnection
var Vc discordgo.VoiceConnection

// something that isn't used for blackjack
type deck struct {
	cards [13]bool
}

// the blackjack info struct
type blackjack_game struct {
	active     bool
	playerID   string
	starttime  int // not used intended to be like a timer for the game
	bet        int
	dealerhand int
	playerhand int
}

type DICTIONARYENTRY struct {
	Keyword    string
	Definition string
}

type CHANNEL struct {
	Id          string
	Blacklisted bool
	Whitelisted bool
	About       string
}

type SERVER struct {
	Id         string
	Prefix     string
	Dictionary []DICTIONARYENTRY
	Channels   []CHANNEL
}

type MEMBER struct {
	server string

	xpi int
	xps string

	leveli int
	levels string

	xpuntills   string
	untillwhats string

	slicesi int
	slicess string

	IAM string
}

var SERVERS []SERVER

var GLOBALDICTIONARY []DICTIONARYENTRY

type GLOBALMEMBER struct {
	xpi int
	xps string

	leveli int
	levels string

	xpuntills   string
	untillwhats string

	slicesi int
	slicess string

	IAM string
}

var blackjack blackjack_game

const Pride_Factory_s string = "575078840804835358"
const Roles_r string = "632070676576075776"

const Blue_r string = "575085615305981997"
const Red_r string = "575085647019376683"
const Yellow_r string = "575085699791978497"
const Green_r string = "575085736920088598"
const Purple_r string = "575085818092191756"
const Orange_r string = "575085772114493491"
const Gray_r string = "632079267395534848"
const Pink_r string = "601632152827985941"
const Pinker_r string = "632076195218849792"

const Friends_r string = "632082165911257089"
const VC_r string = "632082190964097044"

const BotSpam_c string = "632734731863064597"

const kazka_u string = "340665281791918092"

var CBPresence *discordgo.Presence

var state discordgo.State

func main() {

	b, err := ioutil.ReadFile("TOKEN.txt") // b has type []byte
	if err != nil {
		log.Fatal(err)
	}
	token := string(b)

	dg, err := discordgo.New(token)
	if err != nil {
		fmt.Println(err.Error() + " Error 1")
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		fmt.Println(err.Error() + " Error 2")
		return
	} else {
		fmt.Println(u)
	}

	BotID = u.ID

	dg.AddHandler(messageHandler)
	dg.AddHandler(addreactionHandler)
	dg.AddHandler(subreactionHandler)

	dg.AddHandler(MemberJoinHandler)
	dg.AddHandler(MemberLeaveHandler)

	err = dg.Open()

	if err != nil {
		fmt.Println(err.Error() + " Error 3")
		return
	}

	LoadServersConfig()
	LoadGlobalDictionary()

	state = *discordgo.NewState()
	//load / create state
	/*	{
		state.MaxMessageCount = 50
		for i, currentserver := range SERVERS {
			hmm, _ := s.Guild(currentserver.Id)
			fmt.Println(hmm, i)
		}

	}*/

	fmt.Println("Bot is running POGGERS")

	<-make(chan struct{})
	return
}

//add server

func SaveServersConfig() {
	savedata, err := json.Marshal(SERVERS)
	if err != nil {
		fmt.Println(err.Error() + " Error SaveServersConfig 1")
		return
	}

	err = ioutil.WriteFile("SERVERS.config", savedata, 0644)
	if err != nil {
		fmt.Println(err.Error() + " Error SaveServersConfig 2")
		return
	}
}

func LoadServersConfig() {

	savedata, err := ioutil.ReadFile("SERVERS.config")
	if err != nil {
		fmt.Println(err.Error() + " Error LoadServersConfig 1")
		return
	}
	json.Unmarshal(savedata, &SERVERS)
	if err != nil {
		fmt.Println(err.Error() + " Error LoadServersConfig 2")
		return
	}
}

func AddGlobalDictionaryEntry(keyword string, definition string) {

	var dictionaryentry DICTIONARYENTRY
	dictionaryentry.Keyword = keyword
	dictionaryentry.Definition = definition

	GLOBALDICTIONARY = append(GLOBALDICTIONARY, dictionaryentry)

}

func SaveGlobalDictionary() {

	savedata, err := json.Marshal(GLOBALDICTIONARY)
	if err != nil {
		fmt.Println(err.Error() + " Error SaveGlobalDictionary 1")
		return
	}

	err = ioutil.WriteFile("GLOBALDICTIONARY.config", savedata, 0644)
	if err != nil {
		fmt.Println(err.Error() + " Error SaveGlobalDictionary 2")
		return
	}
}

func LoadGlobalDictionary() {

	savedata, err := ioutil.ReadFile("GLOBALDICTIONARY.config")
	if err != nil {
		fmt.Println(err.Error() + " Error LoadGlobalDictionary 1")
		return
	}
	json.Unmarshal(savedata, &GLOBALDICTIONARY)
	if err != nil {
		fmt.Println(err.Error() + " Error LoadGlobalDictionary 2")
		return
	}
}

func VoiceSpeakingUpdateHandler(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
	VCRECIEVE := <-vc.OpusRecv
	fmt.Println(VCRECIEVE)
	vc.Speaking(true)
	vc.OpusSend <- VCRECIEVE.Opus
	vc.Speaking(false)

}

func subreactionHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.ChannelID == Roles_r {
		switch r.Emoji.Name {
		case "💙":
			s.GuildMemberRoleRemove(Pride_Factory_s, r.UserID, Blue_r)
		case "❤":
			s.GuildMemberRoleRemove(Pride_Factory_s, r.UserID, Red_r)
		case "💛":
			s.GuildMemberRoleRemove(Pride_Factory_s, r.UserID, Yellow_r)
		case "💚":
			s.GuildMemberRoleRemove(Pride_Factory_s, r.UserID, Green_r)
		case "💜":
			s.GuildMemberRoleRemove(Pride_Factory_s, r.UserID, Purple_r)
		case "📙":
			s.GuildMemberRoleRemove(Pride_Factory_s, r.UserID, Orange_r)
		case "🌝":
			s.GuildMemberRoleRemove(Pride_Factory_s, r.UserID, Gray_r)
		case "🌸":
			s.GuildMemberRoleRemove(Pride_Factory_s, r.UserID, Pink_r)
		case "🌷":
			s.GuildMemberRoleRemove(Pride_Factory_s, r.UserID, Pinker_r)
		case "🎙":
			s.GuildMemberRoleRemove(Pride_Factory_s, r.UserID, VC_r)
		case "👋":
			s.GuildMemberRoleRemove(Pride_Factory_s, r.UserID, Friends_r)

		}
	}

}

func addreactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {

	if r.ChannelID == Roles_r {
		switch r.Emoji.Name {
		case "💙":
			s.GuildMemberRoleAdd(Pride_Factory_s, r.UserID, Blue_r)
		case "❤":
			s.GuildMemberRoleAdd(Pride_Factory_s, r.UserID, Red_r)
		case "💛":
			s.GuildMemberRoleAdd(Pride_Factory_s, r.UserID, Yellow_r)
		case "💚":
			s.GuildMemberRoleAdd(Pride_Factory_s, r.UserID, Green_r)
		case "💜":
			s.GuildMemberRoleAdd(Pride_Factory_s, r.UserID, Purple_r)
		case "📙":
			s.GuildMemberRoleAdd(Pride_Factory_s, r.UserID, Orange_r)
		case "🌝":
			s.GuildMemberRoleAdd(Pride_Factory_s, r.UserID, Gray_r)
		case "🌸":
			s.GuildMemberRoleAdd(Pride_Factory_s, r.UserID, Pink_r)
		case "🌷":
			s.GuildMemberRoleAdd(Pride_Factory_s, r.UserID, Pinker_r)
		case "🎙":
			s.GuildMemberRoleAdd(Pride_Factory_s, r.UserID, VC_r)
		case "👋":
			s.GuildMemberRoleAdd(Pride_Factory_s, r.UserID, Friends_r)

		}
	}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return
	} else {
		if m.ChannelID != "547188977053204521" {
			activity_tracker(m)
		}
		if keywords(s, m) == true {
			// terminate
		}
		if commands(s, m) == true {
			// terminate
		}
		if m.ChannelID == BotSpam_c || test_room(m.ChannelID) || m.Author.ID == kazka_u {
			if casino(s, m) == true {
				// terminate
			}
		}
		if m.ChannelID == "436669514931765279" || test_room(m.ChannelID) || m.Author.ID == kazka_u {
			if store(s, m) == true {
				// terminate
			}
		}
		//testshit needs to move

		for i, currentserver := range SERVERS {
			if m.GuildID == currentserver.Id {
				//SERVERADMIN
				{

					/*					memberid, _ := s.GuildMember(m.GuildID, m.Author.ID)
										for i, currentrole_s := range memberid.Roles {
											currentrole, _ := state.Role(m.GuildID, currentrole_s)


										}*/

				}
				if m.Content == SERVERS[i].Prefix+"d" {
					var gdk_s string
					var ldk_s string
					for _, currententry := range GLOBALDICTIONARY {
						gdk_s += "\t" + currententry.Keyword
					}
					for _, currententry := range currentserver.Dictionary {
						ldk_s += "\t" + currententry.Keyword
					}
					s.ChannelMessageSend(m.ChannelID, "Now do the same command, but with one of these keywords following: "+gdk_s+ldk_s)
				} else if strings.HasPrefix(m.Content, SERVERS[i].Prefix+"d") {
					trim := strings.TrimPrefix(m.Content, SERVERS[i].Prefix+"d")
					tolower := strings.ToLower(trim)
					var definitions_s string
					for _, currententry := range GLOBALDICTIONARY {
						if strings.Contains(tolower, currententry.Keyword) {
							definitions_s += currententry.Definition + "\n"
						}

					}
					for _, currententry := range currentserver.Dictionary {
						if strings.Contains(tolower, currententry.Keyword) {
							definitions_s += currententry.Definition + "\n"
						}
					}
					if definitions_s == "" {
						definitions_s = "That shit dont mean shit, but you're still valid"
					}
					s.ChannelMessageSend(m.ChannelID, definitions_s)
				}
			}
		}
		//kazkashit needs to move
		if m.Author.ID == kazka_u {
			for _, currentserver := range SERVERS {
				if m.GuildID == currentserver.Id {

					if strings.HasPrefix(m.Content, "kazka!"+"state") {
						fmt.Println("here you go ", state.Guilds)
					}

					if strings.HasPrefix(m.Content, "kazka!"+"help") {
						s.ChannelMessageSend(m.ChannelID, "kazka!help,\tkazka!ADDSERVER,\tkazka!globaldefine")
					}

					if strings.HasPrefix(m.Content, "kazka!"+"ADDSERVER") {
						s.ChannelMessageSend(m.ChannelID, "This server is already added dickhead.")
					}

					if strings.HasPrefix(m.Content, "kazka!"+"globaldefine") {
						trim := strings.TrimPrefix(m.Content, "kazka!"+"globaldefine")
						trim = strings.TrimSpace(trim)
						split := strings.SplitN(trim, " ", 2)
						var newentry DICTIONARYENTRY
						newentry.Keyword = strings.ToLower(split[0])
						newentry.Definition = split[1]

						for i, currententry := range GLOBALDICTIONARY {
							if currententry.Keyword == newentry.Keyword {
								if currententry.Definition == newentry.Definition {
									s.ChannelMessageSend(m.ChannelID, "Kazka, FR this shit already defined dumb ass.")
								} else {
									if strings.ToLower(newentry.Definition) == "delete" {
										GLOBALDICTIONARY = append(GLOBALDICTIONARY[:i], GLOBALDICTIONARY[i+1:]...)
										s.ChannelMessageSend(m.ChannelID, "Succsessfully undefined: \""+newentry.Keyword+"\"")
										SaveGlobalDictionary()
									} else {
										GLOBALDICTIONARY[i].Definition = newentry.Definition
										s.ChannelMessageSend(m.ChannelID, "Redefined \""+currententry.Keyword+"\" from \""+currententry.Definition+"\"to \""+newentry.Definition+"\"")
										SaveGlobalDictionary()
									}
								}
								return
							}
						}
						if strings.ToLower(newentry.Definition) == "delete" {
							s.ChannelMessageSend(m.ChannelID, "That already doesnt exist 🤔")
						} else {
							s.ChannelMessageSend(m.ChannelID, "Successfully Globally Defined!\nKeyword: "+split[0]+"\nDefinition: "+split[1])
							GLOBALDICTIONARY = append(GLOBALDICTIONARY, newentry)
						}
						SaveGlobalDictionary()
					}

					return
				}
			}
			if m.Content == "kazka!"+"ADDSERVER" {
				var newserver SERVER
				newserver.Id = m.GuildID
				newserver.Prefix = defaultprefix
				SERVERS = append(SERVERS, newserver)
				s.ChannelMessageSend(m.ChannelID, "Server Added, use "+defaultprefix+"help for a list of commands")
				SaveServersConfig()
			}
		}
	}
}

func keywords(s *discordgo.Session, m *discordgo.MessageCreate) (err bool) {

	if strings.Contains(m.Content, "monkaS") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":monkaS:432277630692360194")
	}

	if strings.Contains(m.Content, "haHAA") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":HAhaa:434488099083649049")
	}

	if strings.Contains(m.Content, "POGGERS") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":POGGERS:434360204147032074")
	}

	if strings.Contains(m.Content, "FeelsGoodMan") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":FeelsGoodMan:434360842985668608")
	}

	if strings.Contains(m.Content, "FeelsBadMan") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":FeelsBadMan:434360842377625631")
	}

	if strings.Contains(m.Content, "D:") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":CapitalDColon:434488099255615488")
	}

	if strings.Contains(m.Content, "ecksdee") {
		s.ChannelMessageSend(m.ChannelID, "xD")
	}

	if strings.Contains(m.Content, "DansGame") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":DansGame:438526394331561993")
	}

	if strings.Contains(strings.ToLower(m.Content), "dab") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":HAhaa:434488099083649049")
		s.ChannelMessageSend(m.ChannelID, "\\ <:HAhaa:434488099083649049> >")
	}

	if strings.Contains(strings.ToLower(m.Content), "good bot") {
		s.ChannelMessageSend(m.ChannelID, ":)")
	}

	if strings.Contains(m.Content, "OMEGALUL") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":OMEGALUL:434488099083386890")
	}

	if strings.Contains(m.Content, "gachiBASS") || strings.Contains(strings.ToLower(m.Content), " ass ") || strings.Contains(strings.ToLower(m.Content), " dick ") || strings.Contains(strings.ToLower(m.Content), " cock ") || strings.Contains(strings.ToLower(m.Content), " penis ") || strings.HasPrefix(strings.ToLower(m.Content), "ass") || strings.HasPrefix(strings.ToLower(m.Content), "dick") || strings.HasPrefix(strings.ToLower(m.Content), "cock") || strings.HasPrefix(strings.ToLower(m.Content), "penis") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "a:gachiBASS:434488099163078667")
	}

	if strings.Contains(m.Content, "Kreygasm") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":Kreygasm:465040254840340481")
	}

	if strings.Contains(strings.ToLower(m.Content), "bitch") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":CapitalDColon:434488099255615488")
		s.MessageReactionAdd(m.ChannelID, m.ID, "🇧")
	}

	if strings.Contains(strings.ToLower(m.Content), "cunt") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":CapitalDColon:434488099255615488")
		s.MessageReactionAdd(m.ChannelID, m.ID, "🇨")
	}

	if strings.Contains(strings.ToLower(m.Content), "faggot") || strings.Contains(strings.ToLower(m.Content), "fag") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":CapitalDColon:434488099255615488")
		s.MessageReactionAdd(m.ChannelID, m.ID, "🇫")
	}

	if strings.Contains(strings.ToLower(m.Content), "retard") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":CapitalDColon:434488099255615488")
		s.MessageReactionAdd(m.ChannelID, m.ID, "🇷")
	}

	if strings.Contains(m.Content, "Clap") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "a:Clap:434360842620895253")
	}

	if m.GuildID == "" {
		DMCHANNEL, _ := s.UserChannelCreate(m.Author.ID)
		s.ChannelMessageSend(DMCHANNEL.ID, "Hello "+m.Author.Username+".")
	}

	return false
}

func echo(v *discordgo.VoiceConnection) {

	recv := make(chan *discordgo.Packet, 2)
	go dgvoice.ReceivePCM(v, recv)

	send := make(chan []int16, 2)
	go dgvoice.SendPCM(v, send)

	v.Speaking(true)
	defer v.Speaking(false)

	for {

		p, ok := <-recv
		if !ok {
			return
		}

		send <- p.PCM
	}
}

func commands(s *discordgo.Session, m *discordgo.MessageCreate) (err bool) {
	// general commands
	if m.ChannelID == BotSpam_c || m.Author.ID == kazka_u {

		//!who
		if m.Content == "!who" {
			s.ChannelMessageSend(m.ChannelID, "I am a program running off Kazka's computer that he is currently working on for Asexual Discords. If you want to see what I can do, you can type ``!help``")
		}

		if strings.HasPrefix(m.Content, "!xp") || strings.HasPrefix(m.Content, "!xpc") || strings.HasPrefix(m.Content, "!slices") || strings.HasPrefix(m.Content, "!slicesc") {
			s.ChannelMessageSend(m.ChannelID, "Depreciated, Try ``!check`` or ``!check [@]``")
		}

		if strings.HasPrefix(m.Content, "!iam") {
			if len(m.Content)-5 > 140 {
				s.ChannelMessageSend(m.ChannelID, "Yea, ok that's too long.")
			} else {
				iam_update(m.Author.ID, strings.TrimPrefix(m.Content, "!iam"))
				s.ChannelMessageSend(m.ChannelID, "You are now \""+strings.TrimPrefix(m.Content, "!iam ")+"\"")
			}
		}

		if strings.HasPrefix(m.Content, "!check") {
			if m.Content == "!check" {
				member := member_get(s, m.Author.ID, true)
				embed := &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{
						IconURL: m.Author.AvatarURL(""),
					},
					Color:       0x8a72da, // purple
					Title:       "Here are your stats " + m.Author.Username + "!",
					Description: member.IAM,
					Fields: []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:   "Total XP",
							Value:  member.xps,
							Inline: true,
						},
						&discordgo.MessageEmbedField{
							Name:   "Current Level",
							Value:  member.levels,
							Inline: true,
						},
						&discordgo.MessageEmbedField{
							Name:   member.untillwhats,
							Value:  member.xpuntills,
							Inline: true,
						},
						&discordgo.MessageEmbedField{
							Name:   "Total Slices",
							Value:  member.slicess,
							Inline: true,
						},
					},
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: m.Author.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
				}
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
			} else {
				mentions := m.Mentions
				howmany := len(mentions)
				if howmany > 0 {
					for i := 0; i < howmany; i++ {
						member := member_get(s, mentions[i].ID, mentions[i].ID == m.Author.ID)
						embed := &discordgo.MessageEmbed{
							Author: &discordgo.MessageEmbedAuthor{
								IconURL: mentions[i].AvatarURL(""),
							},
							Color:       0x8ad0da, // purple
							Title:       "...Checking " + mentions[i].Username + "'s Stats!",
							Description: member.IAM,
							Fields: []*discordgo.MessageEmbedField{
								&discordgo.MessageEmbedField{
									Name:   "Total XP",
									Value:  member.xps,
									Inline: true,
								},
								&discordgo.MessageEmbedField{
									Name:   "Current Level",
									Value:  member.levels,
									Inline: true,
								},
								&discordgo.MessageEmbedField{
									Name:   member.untillwhats,
									Value:  member.xpuntills,
									Inline: true,
								},
								&discordgo.MessageEmbedField{
									Name:   "Total Slices",
									Value:  member.slicess,
									Inline: true,
								},
							},
							Thumbnail: &discordgo.MessageEmbedThumbnail{
								URL: mentions[i].AvatarURL(""),
							},
							Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
						}
						s.ChannelMessageSendEmbed(m.ChannelID, embed)
					}

				} else {
					s.ChannelMessageSend(m.ChannelID, "????????")
					return
				}
			}

		}

		//!help
		if m.Content == "!help" {
			s.ChannelMessageSend(m.ChannelID, `<#`+BotSpam_c+`> Commands:
!help - gives a list of commands
!keywords - list of words that chillbot reacts to
!who - who am I, you ask?
!help blackjack - help for asexual blackjack
!help roulette - roulette help
!check - check your CAD stats
!check [@mentions]- check others CAD stats
!iam [Identity / Status / Pronouns / idc] - update your IAM
!transfer [@mention] [number of slices] - transfer your slices to another user
!buy - list of things you can buy with slices

Anywhere Commands:
!roll [number] - roll a die
!SWAT - when u need to SWAT someone
!hug - when u wanna hug someone
!d / !define / !dictionary - gives you a list of words we have defined
!info - gives room info (please use passive aggressively)

Staff Commands:
!mute [@mention] [reason]
!kick [@mention] [reason]
!ban [@mention] [reason]
!xpa [@mention] [amount] - for adjusting XP
!slicesa [@mention] [amount] - for asjusting slices
!terminate - turns the bot off
`)
		}

		//!keywords
		if m.Content == "!keywords" {
			s.ChannelMessageSend(m.ChannelID, "monkaS, FeelsGoodMan, FeelsBadMan, Clap, haHAA, DansGame, POGGERS, ecksdee, good bot, D:, Dab, OMEGALUL, gachiBASS")
		}

		//!help blackjack
		if m.Content == "!help blackjack" {
			s.ChannelMessageSend(m.ChannelID, "Type ``!blackjack`` or ``!bj`` followed by your bet to start the game. For example you can bet 10 slices like so: ```!bj 10```")
			s.ChannelMessageSend(m.ChannelID, "In Asexual BlackJack Kings and Queens are worth ``1`` point, and Aces are the most at ``11`` points.")
			s.ChannelMessageSend(m.ChannelID, "You win by hitting 21 points first, not going above 21 first, or staying when you have more points.")
			s.ChannelMessageSend(m.ChannelID, "Draw another card with ``!hit``, both you and the dealer draw a card.")
			s.ChannelMessageSend(m.ChannelID, "Compare cards with ``!stay``, if you're feeling confident that is.")
			s.ChannelMessageSend(m.ChannelID, "Force end a blackjack session with ``!bjCLEAR`` **DO NOT ABUSE**")
		}

		//!help roulette
		if m.Content == "!help roulette" {
			s.ChannelMessageSend(m.ChannelID, "Bruh it's roulette. >.>")
			s.ChannelMessageSend(m.ChannelID, "Type ``!roulette`` followed by your bet to play, for example: ```!roulette 10```")
		}

		//!transfer [@] [amount]
		if strings.HasPrefix(m.Content, "!transfer ") {
			mentions := m.Mentions
			if len(mentions) == 1 {
				SLICEStransfer_s := strings.TrimPrefix(m.Content, "!transfer <@"+mentions[0].ID+"> ")
				if strings.HasPrefix(SLICEStransfer_s, "!transfer") {
					SLICEStransfer_s = strings.TrimPrefix(m.Content, "!transfer <@!"+mentions[0].ID+"> ")
				}
				SLICEStransfer_i, _ := strconv.Atoi(SLICEStransfer_s)
				if SLICEStransfer_i < 1 {
					s.ChannelMessageSend(m.ChannelID, "WTF, no???")
					return false
				}
				SLICEShas, _ := strconv.Atoi(currency_get(m.Author.ID))
				if SLICEStransfer_i > SLICEShas {
					s.ChannelMessageSend(m.ChannelID, "Uh... You might wanna double check your bank account...")
					return false
				} else {
					currency_adjust(m.ChannelID, -SLICEStransfer_i, m.Author.ID)
					currency_adjust(m.ChannelID, SLICEStransfer_i, mentions[0].ID)
					s.ChannelMessageSend(m.ChannelID, "User <@"+m.Author.ID+"> has transfered ``"+SLICEStransfer_s+"`` to user <@"+mentions[0].ID+">!")
					return true
				}
			}
		}
	}

	// anywhere commands

	//dictionary
	if strings.HasPrefix(strings.ToLower(m.Content), "!dictionary") || strings.HasPrefix(strings.ToLower(m.Content), "!define") || strings.HasPrefix(strings.ToLower(m.Content), "!d") {
		dtolower := strings.ToLower(m.Content)

		if dtolower == "!dictionary" || dtolower == "!define" || dtolower == "!d" {
			s.ChannelMessageSend(m.ChannelID, "Do ``!d`` followed by one of our defined words.\nWords we have defined are:\n🔻\nAspec\nAsexual\nAllo\nRomantic\nAromantic\nAlterous\nPlatonic\nLith/Akio/Akoi\nCupio\nSapio\nDemi\nBi/Homo/Hetero\nPan\nSensual\nAesthetic\nGSRM\nAutorchoris\nLibido\nQPR\nSensual\nRecip")
		}

		if strings.Contains(dtolower, "asexual") || strings.Contains(dtolower, "ace") {
			s.ChannelMessageSend(m.ChannelID, "```Asexual``````A person without sexual attraction.```")
		}

		if strings.Contains(dtolower, "acespec") || strings.Contains(dtolower, "aspec") {
			s.ChannelMessageSend(m.ChannelID, "```Ace Spectrum``````A person who is some degree of asexual and/or aromantic.```")
		}

		if strings.Contains(dtolower, "allo") {
			s.ChannelMessageSend(m.ChannelID, "```Allo (prefix)``````A person who experiences attraction (ex. sexual, romantic, etc.). Typically allo is short for allosexual, or someone who is not asexual.```")
		}

		if strings.Contains(dtolower, " ro ") || strings.Contains(dtolower, "-ro ") || strings.Contains(dtolower, "romo") || strings.Contains(dtolower, " romantic") || strings.HasSuffix(m.Content, "-ro") || strings.HasSuffix(m.Content, " ro") {
			s.ChannelMessageSend(m.ChannelID, "```Romantic``````A person with romantic attraction.```")
		}

		if strings.Contains(dtolower, "aro ") || strings.Contains(dtolower, "aromantic") {
			s.ChannelMessageSend(m.ChannelID, "```Aromantic``````A person without romantic attraction.```")
		}

		if strings.Contains(dtolower, "cupio") {
			s.ChannelMessageSend(m.ChannelID, "```Cupio (prefix)``````A person without attraction that has a desire for such a relationship.```")
		}

		if strings.Contains(dtolower, "akio") || strings.Contains(dtolower, "lith") || strings.Contains(dtolower, "akoi") {
			s.ChannelMessageSend(m.ChannelID, "```Akio (AKA Lith & Akoi) (prefix)``````A person that loses attraction upon reciprocation.```")
		}

		if strings.Contains(dtolower, "demi") {
			s.ChannelMessageSend(m.ChannelID, "```Demi (prefix)``````A person who gains attraction after they feel a sufficient bond has been formed.```")
		}

		if strings.Contains(dtolower, "sapio") {
			s.ChannelMessageSend(m.ChannelID, "```Sapio (prefix)``````A person who is attracted to percieved intellegence.```")
		}

		if strings.Contains(dtolower, "platonic") || strings.Contains(dtolower, "friend") {
			s.ChannelMessageSend(m.ChannelID, "```Platonic & Platonic Attraction``````A nonsexual or romantic relationship, friendship.``````An attraction that is purely a desire to become friends.```")
		}

		if strings.Contains(dtolower, "alterous") || strings.Contains(dtolower, "friend++") {
			s.ChannelMessageSend(m.ChannelID, "```Alterous Attraction``````An attraction that is described as \"more than platonic\" and \"less than romantic\" or an \"emotional attraction.\" Generally the relationship between individuals in a QPR. ```")
		}

		if strings.Contains(dtolower, "grey") || strings.Contains(dtolower, "gray") || strings.Contains(dtolower, "grace") {
			s.ChannelMessageSend(m.ChannelID, "```Grey (prefix)``````An attraction that is seemingly random, or does not seem to follow a known pattern.```")
		}
		//
		if strings.Contains(dtolower, " bi") || strings.Contains(dtolower, "bi-") {
			s.ChannelMessageSend(m.ChannelID, "```Bi (prefix)``````An attraction towards the both binary genders.```")
		}

		if strings.Contains(dtolower, " pan") || strings.Contains(dtolower, "pan-") {
			s.ChannelMessageSend(m.ChannelID, "```Pan (prefix)``````No preference for gender in attraction.```")
		}

		if strings.Contains(dtolower, "homo") || strings.Contains(dtolower, "homo-") || strings.Contains(dtolower, "gay") {
			s.ChannelMessageSend(m.ChannelID, "```Homo (prefix)``````An attraction to the same gender.```")
		}

		if strings.Contains(dtolower, " het") || strings.Contains(dtolower, "hetero") || strings.Contains(dtolower, "straight") {
			s.ChannelMessageSend(m.ChannelID, "```Hetero (prefix)``````An attraction the other gender.```")
		}

		if strings.Contains(dtolower, "aesthetic") {
			s.ChannelMessageSend(m.ChannelID, "```Aesthetic Attraction``````Non-sexual/nonromantic attraction to the way someone looks. Described as \"Like looking at a nice painting.\" (CGB)```")
		}

		if strings.Contains(dtolower, "autochoris") || strings.Contains(dtolower, "idea") {
			s.ChannelMessageSend(m.ChannelID, "```Autochoris (prefix)``````a disconnection between oneself and a sexual target/object of arousal; may involve sexual fantasies or arousal in response to erotica or pornography, but lacking any desire to be a participant in the sexual activities therein. (Anthony Bogaert)```")
		}

		if strings.Contains(dtolower, "gsrm") || strings.Contains(dtolower, "gsm") || strings.Contains(dtolower, "lgbt") {
			s.ChannelMessageSend(m.ChannelID, "```Gender Sexual & Romantic Minorities (GSRM)``````Much better than the Alphabet Soup.```")
		}

		if strings.Contains(dtolower, "libido") || strings.Contains(dtolower, "sex drive") {
			s.ChannelMessageSend(m.ChannelID, "```Libido``````Often reffered to as a \"Sex Drive.\" Better described as a persons sexual energy or capacity.\nAn Ace may or may not have a libido just like any other sexuality.```")
		}

		if strings.Contains(dtolower, "queer platonic relationship") || strings.Contains(dtolower, "qpr") {
			s.ChannelMessageSend(m.ChannelID, "```Queer Platonic Relationship (QPR)``````A relationship that is not romantic but involves a close emotional connection (Alterous) beyond what most people consider friendship.```")
		}

		if strings.Contains(dtolower, "sensual") {
			s.ChannelMessageSend(m.ChannelID, "```Sensual Attraction``````An attraction based on one of the senses, Touch, Sight, Sound, Smell and or Taste, can be non-sexual, or sexual.```")
		}

		if strings.Contains(dtolower, "reciprocate") || strings.Contains(dtolower, "recip") {
			s.ChannelMessageSend(m.ChannelID, "```Recip (prefix)``````A reciprocated attraction. Like the mirror of attraction, the opposite of Akoi.```")
		}
		/*
			if strings.Contains(dtolower, "aro ") || strings.Contains(dtolower, "aromantic") {
				s.ChannelMessageSend(m.ChannelID, "```Aromantic``````A person without romantic attraction.```")
			}

			if strings.Contains(dtolower, "aro ") || strings.Contains(dtolower, "aromantic") {
				s.ChannelMessageSend(m.ChannelID, "```Aromantic``````A person without romantic attraction.```")
			}

			if strings.Contains(dtolower, "aro ") || strings.Contains(dtolower, "aromantic") {
				s.ChannelMessageSend(m.ChannelID, "```Aromantic``````A person without romantic attraction.```")
			}

			if strings.Contains(dtolower, "aro ") || strings.Contains(dtolower, "aromantic") {
				s.ChannelMessageSend(m.ChannelID, "```Aromantic``````A person without romantic attraction.```")
			}
		*/
	}
	{
		//!hug
		if strings.HasPrefix(m.Content, "!hug") {

			mentions := m.Mentions
			if len(mentions) == 1 {
				if m.Author.ID == "163691732565753857" {
					if mentions[0].ID == kazka_u {
						s.ChannelMessageSend(m.ChannelID, "***HUGS <@"+kazka_u+">!***")
					} else {
						s.ChannelMessageSend(m.ChannelID, "Please don't.")
					}
				} else if m.Author.ID == kazka_u {
					if mentions[0].ID == "163691732565753857" {
						s.ChannelMessageSend(m.ChannelID, "***HUGS <@163691732565753857>!***")
					} else {
						s.ChannelMessageSend(m.ChannelID, "Please don't.")
					}
				} else {
					s.ChannelMessageSend(m.ChannelID, "Please don't.")
				}
			} else if strings.HasSuffix(m.Content, "daddy") || strings.HasSuffix(m.Content, "Daddy") {
				s.ChannelMessageSend(m.ChannelID, "<a:gachiBASS:434488099163078667>")
			} else {
				s.ChannelMessageSend(m.ChannelID, "Please don't.")
			}
		}

		//!d
		if strings.HasPrefix(m.Content, "!roll ") {
			trimmed_s := strings.TrimPrefix(m.Content, "!roll ")
			trimmed_i, _ := strconv.Atoi(trimmed_s)
			if trimmed_i > 0 {
				s.ChannelMessageSend(m.ChannelID, strconv.Itoa((rand.Int()%(trimmed_i))+1))
			}
		}

		//!SWAT
		if strings.HasPrefix(m.Content, "!SWAT") {
			s.ChannelMessageSend(m.ChannelID, "PUT UR HANDS UP CRIMINAL\nhttps://www.youtube.com/watch?v=6SC-KleaUBc")
		}

	}

	//support commands
	/*
	   	if hard_support(s, m.Author.ID) || hard_admin(s, m.Author.ID) || m.Author.ID == kazka_u {

	   		//!squirrel [reason]
	   		if strings.HasPrefix(m.Content, "!squirrel ") {
	   			if hard_admin(s, m.Author.ID) || hard_support(s, m.Author.ID) || m.Author.ID == "200254471845052416" {
	   				var ii int
	   				var squirel string
	   				for ii < 20 {
	   					ii++
	   					squirel += "Hey what's that? --->\n\n\n\n"
	   					if ii == 10 {
	   						squirel += "🐿\n\n\n\n️"
	   					}
	   				}
	   				s.ChannelMessageSend(m.ChannelID, squirel)
	   				moderation_log(s, m.Author.ID, 1, m.ChannelID, "", strings.TrimPrefix(m.Content, "!squirrel"))
	   			} else {
	   				s.ChannelMessageSend(m.ChannelID, "There aint nuthin there...")
	   			}
	   		}
	   	}
	   /*
	   	//staff commands
	   	/*if hard_admin(s, m.Author.ID) || m.Author.ID == kazka_u {

	   	//!mute [reason]
	   	if strings.HasPrefix(m.Content, "!mute") {
	   		if hard_admin(s, m.Author.ID) {
	   			mentions := m.Mentions
	   			if len(mentions) == 1 {
	   				s.GuildMemberRoleAdd(CAD_s, m.Mentions[0].ID, "454807029207269376")
	   				moderation_log(s, m.Author.ID, 2, m.ChannelID, m.Mentions[0].ID, strings.TrimPrefix(m.Content, "!mute "))
	   			}
	   		} else {
	   			s.ChannelMessageSend(m.ChannelID, "Excuse me, who do you think you are? I should mute you...")
	   		}
	   	}

	   	//!kick [reason]
	   	if strings.HasPrefix(m.Content, "!kick") {
	   		if hard_admin(s, m.Author.ID) {
	   			mentions := m.Mentions
	   			if len(mentions) == 1 {
	   				s.GuildMemberDeleteWithReason(CAD_s, m.Mentions[0].ID, strings.TrimPrefix(m.Content, "!kick"))
	   				moderation_log(s, m.Author.ID, 4, m.ChannelID, m.Mentions[0].ID, strings.TrimPrefix(m.Content, "!kick "))
	   			}
	   		} else {
	   			s.ChannelMessageSend(m.ChannelID, "Excuse me, who do you think you are? I should kick you...")
	   		}
	   	}

	   	//!ban [reason]
	   	if strings.HasPrefix(m.Content, "!ban") {
	   		if hard_admin(s, m.Author.ID) {
	   			mentions := m.Mentions
	   			if len(mentions) == 1 {
	   				s.GuildBanCreateWithReason(CAD_s, m.Mentions[0].ID, strings.TrimPrefix(m.Content, "!ban "), 0)
	   				moderation_log(s, m.Author.ID, 3, m.ChannelID, m.Mentions[0].ID, strings.TrimPrefix(m.Content, "!ban "))
	   			} else if len(mentions) == 0 {
	   				s.GuildBanCreateWithReason(CAD_s, m.Mentions[0].ID, strings.TrimPrefix(m.Content, "!ban "), 0)
	   				moderation_log(s, m.Author.ID, 3, m.ChannelID, m.Mentions[0].ID, strings.TrimPrefix(m.Content, "!ban "))
	   			}
	   		} else {
	   			s.ChannelMessageSend(m.ChannelID, "Excuse me, who do you think you are? I should ban you...")
	   		}
	   	}
	*/
	//!xpa [@] [amount]
	if strings.HasPrefix(m.Content, "!xpa") {
		if m.Author.ID == kazka_u {
			mentions := m.Mentions
			if len(mentions) == 1 {
				XPadjust_s := strings.TrimPrefix(m.Content, "!xpa <@"+mentions[0].ID+"> ")
				if strings.HasPrefix(XPadjust_s, "!xpa") {
					XPadjust_s = strings.TrimPrefix(m.Content, "!xpa <@!"+mentions[0].ID+"> ")
				}
				XPadjust_i, _ := strconv.Atoi(XPadjust_s)
				xp_adjust(m.ChannelID, XPadjust_i, mentions[0].ID)
				s.ChannelMessageSend(m.ChannelID, "XP Adjusted")
			}
		}
	}

	//!slicesa [@] [amount]
	if strings.HasPrefix(m.Content, "!slicesa") {
		if m.Author.ID == kazka_u {
			mentions := m.Mentions
			if len(mentions) == 1 {
				SLICESadjust_s := strings.TrimPrefix(m.Content, "!slicesa <@"+mentions[0].ID+"> ")
				if strings.HasPrefix(SLICESadjust_s, "!slicesa") {
					SLICESadjust_s = strings.TrimPrefix(m.Content, "!slicesa <@!"+mentions[0].ID+"> ")
				}
				SLICESadjust_i, _ := strconv.Atoi(SLICESadjust_s)
				currency_adjust(m.ChannelID, SLICESadjust_i, mentions[0].ID)
				s.ChannelMessageSend(m.ChannelID, "FREE SLICES! <:POGGERS:434360204147032074>")
			} else {
				s.ChannelMessageSend(m.ChannelID, "Uh... no, this my cake <:CapitalDColon:434488099255615488> B!")
			}
		}
	}

	//kazka commands
	if m.Author.ID == kazka_u {

		if strings.HasPrefix(m.Content, "!kiss") {
			if m.Content == "!kiss" {
				s.ChannelMessageSend(m.ChannelID, "*Kazka kisses goatboy*")
			} else {

				kisses, _ := strconv.Atoi(strings.TrimPrefix(m.Content, "!kiss "))

				var i int
				var ks string

				for i < kisses {
					i++
					ks += " kis"
				}
				s.ChannelMessageSend(m.ChannelID, "*Kazka"+ks+"s'es goatboy*")

			}
		}

		//!test

		/*		type MessageEmbed struct {
				URL         string                 `json:"url,omitempty"`
				Type        string                 `json:"type,omitempty"`
				Title       string                 `json:"title,omitempty"`
				Description string                 `json:"description,omitempty"`
				Timestamp   string                 `json:"timestamp,omitempty"`
				Color       int                    `json:"color,omitempty"`
				Footer      *MessageEmbedFooter    `json:"footer,omitempty"`
				Image       *MessageEmbedImage     `json:"image,omitempty"`
				Thumbnail   *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
				Video       *MessageEmbedVideo     `json:"video,omitempty"`
				Provider    *MessageEmbedProvider  `json:"provider,omitempty"`
				Author      *MessageEmbedAuthor    `json:"author,omitempty"`
				Fields      []*MessageEmbedField   `json:"fields,omitempty"`
			}*/

		if strings.HasPrefix(m.Content, "!idban") {
			err := s.GuildBanCreate(Pride_Factory_s, strings.TrimPrefix(m.Content, "!idban "), 0)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "!idban error: "+err.Error())
			} else {
				s.ChannelMessageSend(m.ChannelID, "I THINK I BANNED HIM DADDY KAZKA")
			}

			//			s.ChannelMessageSend(m.ChannelID, )

			//			s.GuildBanCreateWithReason(CAD_s, strings.TrimPrefix(m.Content, "!idban "), "IDBAN", 0)
		}

		if m.Content == "!test" {
			DMCHANNEL, _ := s.UserChannelCreate(m.Author.ID)
			s.ChannelMessageSend(DMCHANNEL.ID, "Hello Dave.")
		}

		if strings.Contains(m.Content, "!join") {
			//if join bot channel
			mute := false
			deaf := false
			var err error
			VC, err = s.ChannelVoiceJoin(Pride_Factory_s, "524267098977992709", mute, deaf)
			if err != nil {
				fmt.Println(err.Error() + " Error !test")
			} else {
				if VC.Ready == true {
					s.ChannelMessageSend(m.ChannelID, "VC joined and ready 2 go!")
					//VC.AddHandler(VoiceSpeakingUpdateHandler)

					recv := make(chan *discordgo.Packet, 2)
					go dgvoice.ReceivePCM(VC, recv)

					send := make(chan []int16, 2)
					go dgvoice.SendPCM(VC, send)

					VC.Speaking(true)
					defer VC.Speaking(false)

					for {

						p, ok := <-recv
						if !ok {
							return false
						}

						send <- p.PCM
					}

					//					VoiceStateUpdate

				}
			}
		}

		if strings.Contains(m.Content, "!leave") {
			// if bot channel empty
			VC.Disconnect()
			s.ChannelMessageSend(m.ChannelID, "C ya VC!")
		}

		if strings.HasPrefix(m.Content, "!talking") {
			if VC.Ready == true {
				VC.Speaking(true)
			} else {
				s.ChannelMessageSend(m.ChannelID, "I'm not in VC you fucking mong")
			}
			return
		}

		if strings.HasPrefix(m.Content, "!stoptalking") {
			if VC.Ready == true {
				VC.Speaking(false)
			} else {
				s.ChannelMessageSend(m.ChannelID, "I'm not in VC you fucking mong")
			}
			return
		}
		if strings.HasPrefix(m.Content, "!beep") {
			if VC.Ready == true {

				/////////
				///////// opus / dgvoice
				/////////
				/*
					c_dgp := make(chan *discordgo.Packet)


					dgvoice.ReceivePCM(VC, c_dgp)


					dgphandler := <-c_dgp
					c_i16 := make(chan []int16)
					c_i16 <- dgphandler.PCM

					VC.Speaking(true)
					dgvoice.SendPCM(VC, c_i16)
					VC.Speaking(false)
					/**/

				stop := make(chan bool)
				dgvoice.PlayAudioFile(VC, "sound.wav", stop)

				/////////
				///////// opus / dgvoice
				/////////
			}
		}

		//!playing
		if strings.HasPrefix(m.Content, "!playing") {

			var usd discordgo.UpdateStatusData

			var cbgame discordgo.Game
			cbgame.Name = string(strings.TrimPrefix(m.Content, "!playing"))
			usd.Game = &cbgame

			s.UpdateStatusComplex(usd)

			s.ChannelMessageSend(m.ChannelID, "I am now playing "+strings.TrimPrefix(m.Content, "!playing"))

		}

		//!terminate
		if m.Content == "!terminate" {
			if m.Author.ID == kazka_u {
				s.ChannelMessageSend(m.ChannelID, "A'ight, peace out yall.")
				os.Exit(-1)
			} else {
				s.ChannelMessageSend(m.ChannelID, "The FUCK do you think you are???")
			}
		}

	}

	return false
}

func casino(s *discordgo.Session, m *discordgo.MessageCreate) (err bool) {
	var BJ []string
	BJ = strings.SplitN(m.Content, " ", 2)
	if BJ[0] == "!blackjack" || BJ[0] == "!bj" {
		if blackjack.active == false {
			blackjack.playerID = m.Author.ID
			blackjack.starttime = 0
			if len(BJ) > 1 {
				bet_int, _ := strconv.Atoi(BJ[1])
				if bet_int < 1 {
					s.ChannelMessageSend(m.ChannelID, "You must bet atleast 1 slice!")
					return false
				}
				blackjack.bet = bet_int
			} else {
				blackjack.bet = 1
			}
			currency := currency_get(m.Author.ID)
			currency_int, _ := strconv.Atoi(currency)
			if currency_int < blackjack.bet {
				s.ChannelMessageSend(m.ChannelID, "You only have "+currency+" you scamming fuck! >:(")
			} else if blackjack.bet > 1000000000000 {
				s.ChannelMessageSend(m.ChannelID, "Look pal, Tapioca and Kazka broke the bank, you can't bet over a tril any more :/ sorry...")
			} else {
				bet_string := strconv.Itoa(blackjack.bet)
				s.ChannelMessageSend(m.ChannelID, "You have bet "+bet_string+" in a game of BlackJack!")
				pcard1, pcard1n := pickacardbtw(rand.Int())
				pcard2, pcard2n := pickacardbtw(rand.Int())
				dcard1, _ := pickacardbtw(rand.Int())
				dcard2, _ := pickacardbtw(rand.Int())

				pcard1_string := strconv.Itoa(pcard1)
				pcard2_string := strconv.Itoa(pcard2)
				blackjack.playerhand = pcard1 + pcard2
				ptotal := strconv.Itoa(blackjack.playerhand)
				blackjack.dealerhand = dcard1 + dcard2
				s.ChannelMessageSend(m.ChannelID, "You have drawn the cards "+pcard1n+" & "+pcard2n+" with a value of "+pcard1_string+" + "+pcard2_string+" for a total of "+ptotal+" points.")

				if blackjack.dealerhand > 21 && blackjack.playerhand > 21 {
					s.ChannelMessageSend(m.ChannelID, "[Loss] We both bust, what the hell?!?")
					currency_adjust(m.ChannelID, -blackjack.bet, m.Author.ID)
				} else if blackjack.dealerhand > 21 {
					s.ChannelMessageSend(m.ChannelID, "[Win] Im bust <:CapitalDColon:434488099255615488>")
					currency_adjust(m.ChannelID, blackjack.bet, m.Author.ID)
				} else if blackjack.playerhand > 21 {
					s.ChannelMessageSend(m.ChannelID, "[Loss] You're bust! <:FeelsGoodMan:434360842985668608>")
					currency_adjust(m.ChannelID, -blackjack.bet, m.Author.ID)
				} else if blackjack.dealerhand == 21 && blackjack.playerhand == 21 {
					s.ChannelMessageSend(m.ChannelID, "[Loss] Wait... we both blackjacked on draw? Well fuck, that's a tie.")
				} else if blackjack.dealerhand == 21 {
					s.ChannelMessageSend(m.ChannelID, "[Loss] HA! I got blackjack!")
					currency_adjust(m.ChannelID, -blackjack.bet, m.Author.ID)
				} else if blackjack.playerhand == 21 {
					s.ChannelMessageSend(m.ChannelID, "[Win] You won first draw? You counting cards or something buddy?")
					currency_adjust(m.ChannelID, blackjack.bet, m.Author.ID)
				} else {
					blackjack.active = true
				}
			}
		} else {
			if blackjack.playerID == m.Author.ID {
				s.ChannelMessageSend(m.ChannelID, "Smh... You're already playing idiot; type !hit or !stay.")
			} else {
				s.ChannelMessageSend(m.ChannelID, "WAIT YOUR TURN!")
			}
		}
	}

	if m.Content == "!hit" || m.Content == "!hitme" || m.Content == "!hitmedad" {
		if blackjack.active == true {
			if blackjack.playerID == m.Author.ID {
				if m.Content == "!hitmedad" {
					s.ChannelMessageSend(m.ChannelID, "OK <a:gachiBASS:434488099163078667> <a:Clap:434360842620895253> \n IM GONNA BUST!!! ")
				}
				blackjack.active = false
				pcard, pcardn := pickacardbtw(rand.Int())
				dcard, _ := pickacardbtw(rand.Int())

				blackjack.playerhand += pcard
				blackjack.dealerhand += dcard

				pcard_string := strconv.Itoa(pcard)
				ptotal_string := strconv.Itoa(blackjack.playerhand)

				s.ChannelMessageSend(m.ChannelID, "You have drawn card "+pcardn+" of value "+pcard_string+" for a total of "+ptotal_string+" points.")

				if blackjack.dealerhand > 21 && blackjack.playerhand > 21 {
					s.ChannelMessageSend(m.ChannelID, "[Loss] You and I have bust!")
					currency_adjust(m.ChannelID, -blackjack.bet, m.Author.ID)
				} else if blackjack.dealerhand > 21 {
					s.ChannelMessageSend(m.ChannelID, "[Win] Im bust <:FeelsBadMan:434360842377625631>")
					currency_adjust(m.ChannelID, blackjack.bet, m.Author.ID)
				} else if blackjack.playerhand > 21 {
					s.ChannelMessageSend(m.ChannelID, "[Loss] YOU BUST! I Win <:POGGERS:434360204147032074>")
					currency_adjust(m.ChannelID, -blackjack.bet, m.Author.ID)
				} else if blackjack.dealerhand == 21 && blackjack.playerhand == 21 {
					s.ChannelMessageSend(m.ChannelID, "[Loss] We... both... won? Uh ok...")
					currency_adjust(m.ChannelID, -blackjack.bet, m.Author.ID)
				} else if blackjack.dealerhand == 21 {
					s.ChannelMessageSend(m.ChannelID, "[Loss] Woo blackjack!")
					currency_adjust(m.ChannelID, -blackjack.bet, m.Author.ID)
				} else if blackjack.playerhand == 21 {
					s.ChannelMessageSend(m.ChannelID, "[Win] Damn, you got blackjack")
					currency_adjust(m.ChannelID, blackjack.bet, m.Author.ID)
				} else {
					blackjack.active = true
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "You stupid or something pal?")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "We... we aren't playing a game... right now... I think you meant !bj...")
		}
	}

	if m.Content == "!stay" {
		if blackjack.active == true {
			if blackjack.playerID == m.Author.ID {
				blackjack.active = false
				ptotal := strconv.Itoa(blackjack.playerhand)
				dtotal := strconv.Itoa(blackjack.dealerhand)
				s.ChannelMessageSend(m.ChannelID, "You have chosen to stay! You have "+ptotal+" points, and I have "+dtotal+" points.")
				if blackjack.playerhand > blackjack.dealerhand {
					s.ChannelMessageSend(m.ChannelID, "[Win] You got lucky")
					currency_adjust(m.ChannelID, blackjack.bet, m.Author.ID)
				} else if blackjack.playerhand < blackjack.dealerhand {
					s.ChannelMessageSend(m.ChannelID, "[Loss] Woo, I got a mean poker face")
					currency_adjust(m.ChannelID, -blackjack.bet, m.Author.ID)
				} else {
					s.ChannelMessageSend(m.ChannelID, "[Loss] What are the chances?")
					currency_adjust(m.ChannelID, -blackjack.bet, m.Author.ID)
				}
			}
		}
	}

	if BJ[0] == "!roulette" || BJ[0] == "!r" {
		if len(BJ) > 1 {
			bet_s := BJ[1]
			bet_i, _ := strconv.Atoi(bet_s)
			if bet_i < 1 {
				s.ChannelMessageSend(m.ChannelID, "You must bet atleast 1 Slice!")
			} else {
				currency := currency_get(m.Author.ID)
				currency_int, _ := strconv.Atoi(currency)
				if currency_int < bet_i {
					s.ChannelMessageSend(m.ChannelID, "You only have "+currency+" you scamming fuck! >:(")
				} else {
					if rand.Int()%2 == 1 {
						s.ChannelMessageSend(m.ChannelID, "[Win] You have won "+bet_s+" Slices!")
						currency_adjust(m.ChannelID, bet_i, m.Author.ID)
					} else {
						s.ChannelMessageSend(m.ChannelID, "[Loss] You have lost "+bet_s+" Slices!")
						currency_adjust(m.ChannelID, -bet_i, m.Author.ID)
					}
				}
			}
		}
	}

	if m.Content == "!bjCLEAR" {
		blackjack.active = false
	}

	return false
}

func store(s *discordgo.Session, m *discordgo.MessageCreate) (err bool) {

	/*	if m.ChannelID == BotSpam_c || m.Author.ID == kazka_u {

		//!buy
		if m.Content == "!buy" {
			s.ChannelMessageSend(m.ChannelID, "!buy color [HEX color] - 1,000,000 slices, changes the color of the color role")
		}
		//!buy color [color]
		if strings.HasPrefix(m.Content, "!buy color ") || strings.HasPrefix(m.Content, "!buy colour ") {
			var color_s string
			if strings.HasPrefix(m.Content, "!buy color ") {
				color_s = strings.TrimPrefix(m.Content, "!buy color ")
			} else if strings.HasPrefix(m.Content, "!buy colour ") {
				color_s = strings.TrimPrefix(m.Content, "!buy colour ")
			}
			member, err := s.State.Member("409907314045353984", m.Author.ID)
			if err != nil {
				fmt.Println(err.Error() + " buy color error")
				return false
			}
			for _, RoleID := range member.Roles {
				if RoleID == color_r {
					slices_s := currency_get(m.Author.ID)
					slices_b, _ := strconv.Atoi(slices_s)
					if slices_b >= 1000000 {
						currency_adjust(m.ChannelID, -1000000, m.Author.ID)
						h2d, _ := strconv.ParseInt("0x"+color_s, 0, 64)
						s.GuildRoleEdit(CAD_s, color_r, "Color", int(h2d), false, 0, false) //(st *Role, err error)

						s.ChannelMessageSend(m.ChannelID, "You have bought a color!")
						return true
					} else {
						s.ChannelMessageSend(m.ChannelID, "OOF, You only have "+slices_s+" slices!")
						return true
					}
				}
			}
			s.ChannelMessageSend(m.ChannelID, "You must be in the color role to change the color of the color role.")
			return true
		}
	}*/
	return true
}

func test_room(CHANNELID string) bool {
	/*	if CHANNELID == "442493156584587265" || CHANNELID == "410522839548952596" || CHANNELID == "403460796106932225" {
			return true
		} else {
			return false
		}*/
	return false
}

/*
func moderation_log(s *discordgo.Session, who string, what int, where string, towho string, why string) {
	//	410522993102422026 #moderationlog ID
	if what == 1 { // squirrel
		s.ChannelMessageSend(modlog_c, "User \"<@"+who+">\" has squirreled in channel \"<#"+where+">\" because \"``"+why+"``\".")
	} else if what == 2 { // chloroform
		s.ChannelMessageSend(modlog_c, "User \"<@"+who+">\" has chloroformed \"<@"+towho+">\" in channel \"<#"+where+">\" because \"``"+why+"``\".")
	} else if what == 3 { // ban
		s.ChannelMessageSend(modlog_c, "User \"<@"+who+">\" has banned \"<@"+towho+">\" in channel \"<#"+where+">\" because \"``"+why+"``\".")
		s.ChannelMessageSend(announcements_c, "User \"<@"+who+">\" has banned \"<@"+towho+">\" in channel \"<#"+where+">\" because \"``"+why+"``\".")
	} else if what == 4 { // kick
		s.ChannelMessageSend(modlog_c, "User \"<@"+who+">\" has kicked \"<@"+towho+">\" in channel \"<#"+where+">\" because \"``"+why+"``\".")
		s.ChannelMessageSend(announcements_c, "User \"<@"+who+">\" has kicked \"<@"+towho+">\" in channel \"<#"+where+">\" because \"``"+why+"``\".")
	}

}

func hard_admin(s *discordgo.Session, AUTHORID string) bool {
	if AUTHORID == kazka_u {
		return true
	}
	member, err := s.State.Member("409907314045353984", AUTHORID)
	if err != nil {
		fmt.Println(err.Error() + " Error 4")
		return false
	}
	for _, RoleID := range member.Roles {
		if RoleID == staff_r || RoleID == admin_r || RoleID == moderator_r {
			return true
		}
	}
	return false
}

func hard_support(s *discordgo.Session, AUTHORID string) bool {
	if AUTHORID == kazka_u {
		return true
	}
	member, err := s.State.Member("409907314045353984", AUTHORID)
	if err != nil {
		fmt.Println(err.Error() + " Error 4")
		return false
	}
	for _, RoleID := range member.Roles {
		if RoleID == support_r || RoleID == staff_r || RoleID == admin_r || RoleID == moderator_r {
			return true
		}
	}
	return false
}
*/

func activity_tracker(m *discordgo.MessageCreate) {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		fmt.Println(err.Error() + " Error 4")
	}

	db.Update(func(tx *bolt.Tx) error {
		xp_bucket, err := tx.CreateBucketIfNotExists([]byte("xp"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		xp_byte := xp_bucket.Get([]byte(m.Author.ID))

		level_int, _ := strconv.Atoi(string(xp_byte))

		if (strings.Count(" ", m.Content) + 1) < 10 {
			level_int += strings.Count(m.Content, " ") + (1 * xp_multiplier)

		} else {
			level_int += (10 * xp_multiplier)
		}

		xp_string := strconv.Itoa(level_int)

		xp_byte = []byte(xp_string)

		err = xp_bucket.Put([]byte(m.Author.ID), xp_byte)

		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		c_bucket, err := tx.CreateBucketIfNotExists([]byte("currency"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		xp_byte := c_bucket.Get([]byte(m.Author.ID))

		level_int, _ := strconv.Atoi(string(xp_byte))

		if (strings.Count(" ", m.Content) + 1) < 10 {
			level_int += strings.Count(m.Content, " ") + 1
		} else {
			level_int += 10
		}

		xp_string := strconv.Itoa(level_int)

		xp_byte = []byte(xp_string)

		err = c_bucket.Put([]byte(m.Author.ID), xp_byte)

		return nil
	})

	defer db.Close()
}

func xp_till(ID string) (xp string, till string) {

	xp = xp_get(ID)

	xp_i, _ := strconv.Atoi(xp)
	if xp_i < 500 {
		xp = strconv.Itoa(500 - xp_i)
		till = "XP till Verified"
	} else if xp_i < 40500 {
		xp = strconv.Itoa(40500 - xp_i)
		till = "XP till Chill Squad"
	} else {
		return "", ""
	}
	return xp, till
}

func xp_get(ID string) (xp string) {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		fmt.Println(err.Error() + " Error 4")
	}

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("xp"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		get_b := b.Get([]byte(ID))

		get_i, _ := strconv.Atoi(string(get_b))

		get_s := strconv.Itoa(get_i)

		xp = get_s

		return nil
	})
	defer db.Close()

	return xp
}

func iam_update(ID string, IAM string) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		fmt.Println(err.Error() + " IAM error")
	}
	db.Update(func(tx *bolt.Tx) error {
		iam_bucket, err := tx.CreateBucketIfNotExists([]byte("iam"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = iam_bucket.Put([]byte(ID), []byte(IAM))
		return nil
	})
	defer db.Close()
}

func xp_adjust(ChannelID string, adjustment_value int, ID string) {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		fmt.Println(err.Error() + " Error 4")
	}

	db.Update(func(tx *bolt.Tx) error {
		c_bucket, err := tx.CreateBucketIfNotExists([]byte("xp"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		xp_byte := c_bucket.Get([]byte(ID))

		level_int, _ := strconv.Atoi(string(xp_byte))

		level_int += adjustment_value

		xp_string := strconv.Itoa(level_int)

		xp_byte = []byte(xp_string)

		err = c_bucket.Put([]byte(ID), xp_byte)

		return nil
	})

	defer db.Close()
}

func member_get(s *discordgo.Session, ID string, roles bool) (member MEMBER) {

	{ //currency check
		db, err := bolt.Open("my.db", 0600, nil)
		if err != nil {
			fmt.Println(err.Error() + " currency error")
		}
		// get slices
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("currency"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			get_byte := b.Get([]byte(ID))
			member.slicesi, _ = strconv.Atoi(string(get_byte))
			member.slicess = strconv.Itoa(member.slicesi)
			return nil
		})
		// get xp
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("xp"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			get_byte := b.Get([]byte(ID))
			member.xpi, _ = strconv.Atoi(string(get_byte))
			member.xps = strconv.Itoa(member.xpi)
			return nil
		})
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("iam"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			member.IAM = string(b.Get([]byte(ID)))
			return nil
		})
		defer db.Close()
		if member.IAM == "" {
			member.IAM = "Hey! Hey Listen!"

		}
	}

	{ //level get
		i := member.xpi
		nextlevel := 0
		for i >= nextlevel {
			member.leveli++
			nextlevel = (500 * (member.leveli * member.leveli))
		}
		member.levels = strconv.Itoa(member.leveli)
		if member.leveli < 2 {
			member.untillwhats = "XP till Verified"
			member.xpuntills = strconv.Itoa((500 * (1)) - member.xpi)
		} else if member.leveli < 10 {
			member.untillwhats = "XP till Chill Squad"
			member.xpuntills = strconv.Itoa((500 * (81)) - member.xpi)
		} else {
			member.untillwhats = "XP till next LVL"
			member.xpuntills = strconv.Itoa(nextlevel - member.xpi)
		}
	}

	//assign activity roles
	if roles == true {
		if member.leveli > 1 {
			s.GuildMemberRoleAdd("409907314045353984", ID, "410520872676360193")
		}
		if member.leveli >= 10 {
			s.GuildMemberRoleAdd("409907314045353984", ID, "446112365055049729")
		}
	}

	{ //iam get

	}

	return member
}

func level_get(xp string) (Level string) {
	var Level_i int
	Level_i = 0
	level_i, _ := strconv.Atoi(string(xp))
	for level_i >= (500 * (Level_i * Level_i)) {
		Level_i++
	}

	Level = strconv.Itoa(Level_i)

	return Level
}

func currency_get(ID string) (currency string) {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		fmt.Println(err.Error() + " Error 4")
	}

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("currency"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		get_byte := b.Get([]byte(ID))

		get_int, _ := strconv.Atoi(string(get_byte))

		get_string := strconv.Itoa(get_int)

		currency = get_string

		return nil
	})
	defer db.Close()
	return currency
}

func currency_adjust(ChannelID string, adjustment_value int, ID string) {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		fmt.Println(err.Error() + " Error 4")
	}

	db.Update(func(tx *bolt.Tx) error {
		c_bucket, err := tx.CreateBucketIfNotExists([]byte("currency"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		currency_byte := c_bucket.Get([]byte(ID))

		currency_int, _ := strconv.Atoi(string(currency_byte))

		currency_int += adjustment_value

		currency_string := strconv.Itoa(currency_int)

		currency_byte = []byte(currency_string)

		err = c_bucket.Put([]byte(ID), currency_byte)

		return nil
	})

	defer db.Close()
}

func pickacardbtw(seed int) (value int, name string) {
	switch seed % 12 {
	case 0:
		return 11, "Ace"
	case 1:
		return 2, "Two"
	case 2:
		return 3, "Three"
	case 3:
		return 4, "Four"
	case 4:
		return 5, "Five"
	case 5:
		return 6, "Six"
	case 6:
		return 7, "Seven"
	case 7:
		return 8, "Eight"
	case 8:
		return 9, "Nine"
	case 9:
		return 10, "Ten"
	case 10:
		return 10, "Jack"
	case 11:
		return 1, "Queen"
	case 12:
		return 1, "King"
	default:
		return seed % 13, "somethings fucky, KAZKA!!! :<"
	}
}

func MemberJoinHandler(s *discordgo.Session, ma *discordgo.GuildMemberAdd) {

	histring := "New Fag"
	switch rand.Int() % 5 {
	case 0:
		histring = ("Welcome to the Chill Asexual Discord " + ma.Member.User.Username + "!")
	case 1:
		histring = ("Welcome to Chili's " + ma.Member.User.Username + "!")
	case 2:
		histring = ("Khajiit " + ma.Member.User.Username + " has wares if you have SLICES...")
	case 3:
		histring = ("Bathroom is the 4th door on the left " + ma.Member.User.Username + ".")
	case 4:
		histring = (ma.Member.User.Username + " is here, Please clap.")
	}

	embed := &discordgo.MessageEmbed{
		Color:       0x00ff00,
		Title:       histring,
		Description: ma.Member.User.ID,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: ma.Member.User.AvatarURL(""),
		},
		//		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
	}
	s.ChannelMessageSendEmbed("409907314045353986", embed)

}

func MemberLeaveHandler(s *discordgo.Session, ma *discordgo.GuildMemberRemove) {

	byestring := "Bye Fag"
	switch rand.Int() % 11 {
	case 0:
		byestring = ("I'm sorry " + ma.Member.User.Username + ", You are not the biggest loser.")
	case 1:
		byestring = ("I'm sorry " + ma.Member.User.Username + ", You have been voted off the island.")
	case 2:
		byestring = ("*psst* I heard " + ma.Member.User.Username + " wasn't even chill anyway :/.")
	case 3:
		byestring = ("Don't let the door hit you on the ass on your way out " + ma.Member.User.Username + ".")
	case 4:
		byestring = ("You are the weakest link. Goodbye " + ma.Member.User.Username + ".")
	case 5:
		byestring = ("We didn't like you anyway <@" + ma.Member.User.Username + ">. 🖕🖕🖕")
	case 6:
		byestring = (ma.Member.User.Username + ", didn't even curtesy flush.")
	case 7:
		byestring = (ma.Member.User.Username + "'s CIS HET Otherkin Otaku Senpai is in another castle!")
	case 8:
		byestring = ("Let me guess " + ma.Member.User.Username + "... someone stole your cake.")
	case 9:
		byestring = ("Your cold never bothered us anyway " + ma.Member.User.Username + "...")
	case 10:
		byestring = "Piss off " + ma.Member.User.Username + "!"
	}

	embed := &discordgo.MessageEmbed{
		Color:       0xff0000,
		Title:       byestring,
		Description: ma.Member.User.ID,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: ma.Member.User.AvatarURL(""),
		},
	}
	s.ChannelMessageSendEmbed("436669514931765279", embed)

}

// ideas
// daily pot
// bingo
// update help
// play VC
