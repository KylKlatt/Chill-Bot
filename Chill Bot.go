package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const missingClientSecretsMessage = `
Please configure OAuth 2.0`

// Append this string with error messages
var error_string string

//bot id
var BotID string

var VC *discordgo.VoiceConnection

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

var blackjack blackjack_game

const CAD_s string = "409907314045353984"
const modlog_c string = "410522993102422026"
const announcements_c string = "435758751782535169"
const eventsnbots_c string = "436669514931765279"
const casino_c string = "451255642087358464"
const botroom_c string = "442493156584587265"
const kazka_u string = "340665281791918092"
const pram_u string = "398643515363688448"
const staff_r string = "410522026868998146"
const admin_r string = "410521789782032384"
const moderator_r string = "410521251304570882"
const support_r string = "458036299732090892"
const color_r string = "534412981694627864"

var CBPresence *discordgo.Presence

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func channelsListByUsername(service *youtube.Service, part string, forUsername string) {
	call := service.Channels.List(part)
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	handleError(err, "")
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
}

func main() {
	//////////////////////////////////////
	//////////////////////////////////////
	//////////////////////////////////////
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	channelsListByUsername(service, "snippet,contentDetails,statistics", "GoogleDevelopers")
	//////////////////////////////////////
	//////////////////////////////////////
	//////////////////////////////////////

	b, err = ioutil.ReadFile("TOKEN.txt") // b has type []byte
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
	//u.Assets.SmallText = "Asexual BlackJack"
	BotID = u.ID

	//////
	/*
		file := xlsx.NewFile()
		sheet, err := file.AddSheet("Sheet1")
		if err != nil {
			fmt.Printf(err.Error())
		}
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.Value = BotID
		err = file.Save("CAD DATA.xlsx")
		if err != nil {
			fmt.Printf(err.Error())
		}
	*/
	///////

	dg.AddHandler(messageHandler)
	dg.AddHandler(addreactionHandler)
	dg.AddHandler(subreactionHandler)

	dg.AddHandler(presenceHandler)
	// this if for telling if people join or leave, not used rn
	dg.AddHandler(MemberJoinHandler)
	dg.AddHandler(MemberLeaveHandler)

	//	dg.AddHandler(messageDeleteHandler)

	err = dg.Open()

	if err != nil {
		fmt.Println(err.Error() + " Error 3")
		return
	}

	fmt.Println("Bot is running POGGERS")

	<-make(chan struct{})
	return
}

/*
func messageDeleteHandler(s *discordgo.Session, m *discordgo.MessageDelete) {

	//s.ChannelMessageSend("565246284684984320", "User "+m.Author.Username+", <@"+m.Author.ID+"> has deleted the following message in channel <#"+m.ChannelID+">.")

	s.ChannelMessageSend("565246284684984320", ":) why does this cause panic??????")

}*/

func presenceHandler(s *discordgo.Session, p *discordgo.PresenceUpdate) {

	if p.Presence.User.Username != "" {
		//logme := p.Presence.User.Username + " Status " + p.Presence.Status

		//fmt.Println(logme)
	}
}

func subreactionHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.ChannelID == "556647040633798677" {
		switch r.Emoji.Name {
		case "🥃":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "424210978746400768")
		case "🍑":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "504815202320121857")
		case "📝":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "458282517062221859")
		case "🎰":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "451256712658288643")
		case "🎮":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "455495472874782720")
		case "🗾":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "455495577581387794")
		case "🦊":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "465350963251904512")
		case "🎨":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "465350853847416843")
		case "🚗":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "467166719908249602")
		case "💖":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "560497354822647808")
		case "💚":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "560497403845804032")
		case "🌈":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "560056946439225354")

		case "🎙":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "516841560953061378")
		case "🖱":
			s.GuildMemberRoleRemove("409907314045353984", r.UserID, "513557311592202250")
		}

	}

	if r.MessageID == "560684099006758912" {
		switch r.Emoji.Name {
		case "💩":
			{
				s.GuildMemberRoleRemove(CAD_s, r.UserID, color_r)
			}
		}
	}

}

func addreactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	/*
		type MessageReaction struct {
			UserID    string `json:"user_id"`
			MessageID string `json:"message_id"`
			Emoji     Emoji  `json:"emoji"`
			ChannelID string `json:"channel_id"`
			GuildID   string `json:"guild_id,omitempty"`
		}

		type Emoji struct {
			ID            string   `json:"id"`
			Name          string   `json:"name"`
			Roles         []string `json:"roles"`
			Managed       bool     `json:"managed"`
			RequireColons bool     `json:"require_colons"`
			Animated      bool     `json:"animated"`
		}
	*/

	if r.MessageID == "557223964603187203" || r.MessageID == "560525971661258773" || r.MessageID == "560526178432319488" {
		s.MessageReactionAdd(r.ChannelID, r.MessageID, r.Emoji.Name)
		if r.UserID != BotID {
			switch r.Emoji.Name {
			case "🥃":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "424210978746400768")
			case "🍑":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "504815202320121857")
			case "📝":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "458282517062221859")
			case "🎰":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "451256712658288643")
			case "🎮":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "455495472874782720")
			case "🗾":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "455495577581387794")
			case "🦊":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "465350963251904512")
			case "🎨":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "465350853847416843")
			case "🚗":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "467166719908249602")
			case "💖":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "560497354822647808")
			case "💚":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "560497403845804032")
			case "🌈":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "560056946439225354")

			case "🎙":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "516841560953061378")
			case "🖱":
				s.GuildMemberRoleAdd("409907314045353984", r.UserID, "513557311592202250")
			}

		}

	}

	if r.MessageID == "560684099006758912" {
		switch r.Emoji.Name {
		case "💩":
			{
				member, err := s.State.Member(CAD_s, r.UserID)
				if err != nil {
					fmt.Println(err.Error() + " buy color error")
					return
				}
				for _, RoleID := range member.Roles {
					if RoleID == color_r {
						return
					}
				}

				slices_s := currency_get(r.UserID)
				slices_b, _ := strconv.Atoi(slices_s)
				if slices_b >= 100000000 {
					currency_adjust("", -100000000, r.UserID)
					s.GuildMemberRoleAdd(CAD_s, r.UserID, color_r)
				} else {
					//EDIT A LOG MESSAGE
					//s.MessageReactionRemove(channelID, messageID, emojiID, userID string) error
					s.MessageReactionRemove("556647040633798677", "560684099006758912", r.Emoji.Name, r.UserID)
				}
			}

		}

	}
	/*
		if r.MessageID == "560690684898836494" {
			s.MessageReactionAdd("556647040633798677", "560690684898836494", r.Emoji.Name)
		}
	*/

	if r.MessageID == "560690684898836494" {
		s.MessageReactionRemove("556647040633798677", "560690684898836494", r.Emoji.Name, r.UserID)
		member, err := s.State.Member(CAD_s, r.UserID)
		if err != nil {
			fmt.Println(err.Error() + " buy color error")
			return
		}
		for _, RoleID := range member.Roles {
			if RoleID == color_r {
				slices_s := currency_get(r.UserID)
				slices_b, _ := strconv.Atoi(slices_s)
				if slices_b >= 1000000 {
					currency_adjust(r.ChannelID, -1000000, r.UserID)
					switch r.Emoji.Name {
					case "💙":
						s.GuildRoleEdit(CAD_s, color_r, "Color", 6139372, false, 0, false)
					case "❤":
						s.GuildRoleEdit(CAD_s, color_r, "Color", 12458289, false, 0, false)
					case "💜":
						s.GuildRoleEdit(CAD_s, color_r, "Color", 11177686, false, 0, false)
					case "💚":
						s.GuildRoleEdit(CAD_s, color_r, "Color", 7909721, false, 0, false)
					case "💛":
						s.GuildRoleEdit(CAD_s, color_r, "Color", 16632664, false, 0, false)
					case "🖤":
						s.GuildRoleEdit(CAD_s, color_r, "Color", 1, false, 0, false)
					case "🔶":
						s.GuildRoleEdit(CAD_s, color_r, "Color", 16755763, false, 0, false)
					case "⚪":
						s.GuildRoleEdit(CAD_s, color_r, "Color", 15132648, false, 0, false)
					case "💩":
						s.GuildRoleEdit(CAD_s, color_r, "Color", 12544338, false, 0, false)
					}
				}
			}
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
		if m.ChannelID == "451255642087358464" || test_room(m.ChannelID) || hard_admin(s, m.Author.ID) {
			if casino(s, m) == true {
				// terminate
			}
		}
		if m.ChannelID == "436669514931765279" || test_room(m.ChannelID) || hard_admin(s, m.Author.ID) {
			if store(s, m) == true {
				// terminate
			}
		}
		if strings.Contains(m.Content, "bread") || strings.Contains(m.Content, "🍞") || strings.Contains(m.Content, "🥖") {

		}
		if m.Author.ID == "" {

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

	if strings.Contains(m.Content, "Dab") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":HAhaa:434488099083649049")
		s.ChannelMessageSend(m.ChannelID, "\\ <:HAhaa:434488099083649049> >")
	}

	if strings.Contains(m.Content, "good bot") {
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

	if strings.Contains(m.Content, "bitch") || strings.Contains(m.Content, "Bitch") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":CapitalDColon:434488099255615488")
		s.MessageReactionAdd(m.ChannelID, m.ID, "🇧")
	}

	if strings.Contains(m.Content, "cunt") || strings.Contains(m.Content, "Cunt") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":CapitalDColon:434488099255615488")
		s.MessageReactionAdd(m.ChannelID, m.ID, "🇨")
	}

	if strings.Contains(m.Content, "Faggot") || strings.Contains(m.Content, "Fag") || strings.Contains(m.Content, "faggot") || strings.Contains(m.Content, "fag") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":CapitalDColon:434488099255615488")
		s.MessageReactionAdd(m.ChannelID, m.ID, "🇫")
	}

	if strings.Contains(m.Content, "Retarded") || strings.Contains(m.Content, "Retard") || strings.Contains(m.Content, "retarded") || strings.Contains(m.Content, "retard") {
		s.MessageReactionAdd(m.ChannelID, m.ID, ":CapitalDColon:434488099255615488")
		s.MessageReactionAdd(m.ChannelID, m.ID, "🇷")
	}

	if strings.Contains(m.Content, "Clap") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "a:Clap:434360842620895253")
	}

	if strings.Contains(m.Content, "https://discord.gg/") || strings.Contains(m.Content, "discord.gg/") {
		if hard_admin(s, m.Author.ID) == false {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			s.ChannelMessageSend(m.ChannelID, "OOF sorry, you cant link to discord channels. You can @ a staff member and ask them to link to it if its ok!")
			s.ChannelMessageSend(modlog_c, "Deleted message from <@"+m.Author.ID+"> in channel <#"+m.ChannelID+"> that contained a discord link.: ``"+m.Content+"``")
		}
	}

	return false
}

func commands(s *discordgo.Session, m *discordgo.MessageCreate) (err bool) {
	// general commands
	if m.ChannelID == eventsnbots_c || m.ChannelID == casino_c || hard_admin(s, m.Author.ID) || hard_support(s, m.Author.ID) || m.Author.ID == kazka_u || m.ChannelID == botroom_c {

		// info
		if m.Content == "!info" {
			switch m.ChannelID {
			case "435758751782535169":
				s.ChannelMessageSend(m.ChannelID, "welcome to spam room, fucking kazka's spam room holy shit dude calm down no one even reads this shit :rolling eyes:")
				//RNI

			case "409907314045353986":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#409907314045353986>, this is the entry point for newbs, and also the primary channel for general convos. <:Chill:425981027735961600>")
				//general

			case "454125920664682516":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#454125920664682516>, this is like an overflow channel, for general conversations. <:Chill:425981027735961600>")
				//generalbutdifferent

			case "410522736952082433":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#410522736952082433>, channel for shitposting, memes, or things that otherwise don't have a better place.")
				//offtopic

			case "436669514931765279":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#436669514931765279>, THIS IS MY CHANNEL FOOL!!! MUHAHAHAHAH!!")
				//bot room

			case "455992881237327873":
				s.ChannelMessageSend(m.ChannelID, "Zoinks!, where are we Scoob?")
				//illegal seagul

			case "547188977053204521":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#547188977053204521>,\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nWelcome to <#547188977053204521>,\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nWelcome to <#547188977053204521>,\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nWelcome to <#547188977053204521>,\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
				//spam

			case "543741127023525889":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#543741127023525889>, personal aspec identity questions here, y'kno beginer stuff <:AmIAce:518559063575887953>")
				//am i ace?

			case "442795706416365578":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#442795706416365578>, for more general A/Sexual or A/Romantic Spectrum conversations, and topics.")
				//aspec

			case "530818485530394629":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#530818485530394629>, personal developement topics, and questions, school, work, that kinda thing.")
				//life101

			case "421047453530324992":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#421047453530324992>, dissagreements of any caliber, if it gets heated might need to @support.")
				//debate

			case "455586955367940106":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#455586955367940106>, something on your mind, something you need to get off your chest? Here you go. (also pro tip, not everyone is looking for advice)")
				//ventroom

			case "556647040633798677":
				s.ChannelMessageSend(m.ChannelID, "What the fuck, you like... stupid or something?")
				//optin

			case "513563819012784143":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#513563819012784143>, other non-aspec identities and/or how they interact. Check the pinned messages for more info. (ie gender, autism)\nMANAGED BY CUTEGOATBOY AND GNER0")
				//intersections

			case "560523418559184896":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#560523418559184896>, Talking about GROSS ROMANCE THINGS.")
				//romanace

			case "560495867170390036":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#560495867170390036>, Talking about aro things.")
				//aromanace

			case "438921045345435648":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#438921045345435648>, for conversation oriented, shitposting potentially explicit things is allowed so long as you aren't like interupting a convo.\nMANAGED BY GNER0")
				//explicit1

			case "504815410655395840":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#504815410655395840>, more for shit posting NSFW stuff, but you can use it as a backup for <#438921045345435648> if it's in use, Just remeber to check the DISCORD TOS (linked in #rulesandinfo) or hit up the room manager if you're not sure if something would be allowed here or not.\nMANAGED BY GNER0")
				//explicit2

			case "458282694661505044":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#458282694661505044>, finding, sharing, and taking silly tests! (Tests are pinned). If you have a new test, feel free to hit up one of the managers.\nMANAGED BY CUTEGOATBOY AND ZEBB9")
				//test room

			case "451255642087358464":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#451255642087358464>, you can use !help blackjack or !help roulette to understand the game.")
				//casino

			case "455495913243410434":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#455495913243410434>, I'd recomend checking the pinned messages if you want to find others / share your gaming contact info!\nMANAGED BY  UNKNOWNSOLDIER86")
				//gameroom

			case "455495861036646421":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#455495861036646421>, general weeb asia culture that kinda thing.\nMANAGED BY GOGO")
				//japan

			case "465639181037731840":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#465639181037731840>, SFW channel for talking about and sharing furry stuff (please link). For sharing stuff you've made, consider posting it in <#465639104982417419>!\nMANAGED BY KAZKA")
				//furry

			case "465639104982417419":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#465639104982417419>, for talking about things you are making / working on / made and sharing them :)\nMANAGED BY STARFIERY")
				//maker

			case "467174723558834177":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#467174723558834177>, for all your mechanical interests!\nMANAGED BY 🐐 AND ILIKECARSMORETHANWOMEN")
				//garage

			case "457586853155962880":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#457586853155962880>, Channel for VC ")
				//no context

			case "457820278764863489":
				s.ChannelMessageSend(m.ChannelID, "Welcome to <#457820278764863489>, share music er smthn idc this is the last room im doing :shrug:")
				//music

			}
		}

		//!who
		if m.Content == "!who" {
			s.ChannelMessageSend(m.ChannelID, "I am a program running off Kazka's computer that he is currently working on for our discord. If you want to see what i can do, you can type ``!help``")
		}

		//!xp
		if m.Content == "!xp" || m.Content == "!fuckyoukio" {
			xp := xp_get(m.Author.ID)
			s.ChannelMessageSend(m.ChannelID, "You've got "+xp+" experience!")

			level_s := level_get(xp)
			level_i, _ := strconv.Atoi(level_s)
			if level_i > 1 {
				s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "410520872676360193")
			}
			if level_i >= 10 {
				s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "446112365055049729")
			}

			if level_i < 3 {
				s.ChannelMessageSend(m.ChannelID, "Oh, and you're only level "+level_s+". <:OMEGALUL:434488099083386890>")
			} else if level_i < 10 {
				s.ChannelMessageSend(m.ChannelID, "You're also level "+level_s+". <:HAhaa:434488099083649049>")
			} else if level_i >= 10 {
				s.ChannelMessageSend(m.ChannelID, "Congrats you're level "+level_s+". <:FeelsGoodMan:434360842985668608> <a:Clap:434360842620895253>")
			}

			xp_i, _ := strconv.Atoi(xp)
			if xp_i < 500 {
				till := strconv.Itoa(500 - xp_i)
				s.ChannelMessageSend(m.ChannelID, "You have "+till+" xp more till you are verified at 500 xp (lvl 2).")
			} else if xp_i < 40500 {
				till := strconv.Itoa(40500 - xp_i)
				s.ChannelMessageSend(m.ChannelID, "You need "+till+" more xp till you reach Chill Squad at 40500xp (lvl 10).")
			}

		}

		//!slices
		if m.Content == "!slices" {
			slices := currency_get(m.Author.ID)
			s.ChannelMessageSend(m.ChannelID, "You have "+slices+" slices!")
		}

		if m.Content == "!overflow" {

		}

		//!help
		if m.Content == "!help" {
			s.ChannelMessageSend(m.ChannelID, `<#`+eventsnbots_c+`> Commands:
!help - gives a list of commands
!keywords - list of words that chillbot reacts to
!who - who am I, you ask?
!help blackjack - help for asexual blackjack
!help roulette - roulette help
!xp - lists experience points & level
!xpc [@mention] - check another user's XP & level, however, maybe just ask them
!slices - lists slices
!slicesc [@mention] - check another user's slices, however, maybe just ask them
!transfer [@mention] [number of slices] - transfer your slices to another user
!buy - list of things you can buy with slices

Anywhere Commands:
!roll [number] - roll a die
!SWAT - when u need to SWAT someone
!hug - when u wanna hug someone
!d / !define / !dictionary - gives you a list of words we have defined
!info - gives room info (please use passive aggressively)

Support Commands
!squirrel [reason] - used to 'reset' a chat room, reason required

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

		//!slicesc [@]
		if strings.HasPrefix(m.Content, "!slicesc ") {
			mentions := m.Mentions
			if len(mentions) == 1 {
				slices := currency_get(mentions[0].ID)
				if slices != " " {
					s.ChannelMessageSend(m.ChannelID, "User <@"+mentions[0].ID+"> has "+slices+" slices.")
				} else {
					s.ChannelMessageSend(m.ChannelID, "Looks I had an issue finding that user. <:FeelsBadMan:434360842377625631>")
				}
			}
		}

		//!xpc [@]
		if strings.HasPrefix(m.Content, "!xpc ") {
			mentions := m.Mentions
			if len(mentions) == 1 {
				xp := xp_get(mentions[0].ID)
				lvl := level_get(xp)

				if xp != " " {
					s.ChannelMessageSend(m.ChannelID, "User <@"+mentions[0].ID+"> has "+xp+" xp, and is level "+lvl+".")
				} else {
					s.ChannelMessageSend(m.ChannelID, "Looks I had an issue finding that user. <:FeelsBadMan:434360842377625631>")
				}
			}
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

		if strings.Contains(dtolower, "akio") || strings.Contains(dtolower, "akio") || strings.Contains(dtolower, "akoi") {
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

		if strings.Contains(dtolower, "Sensual") {
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
			} else if strings.HasSuffix(m.Content, "daddy") || strings.HasSuffix(m.Content, "Daddy") || m.ChannelID == botroom_c {
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
	if hard_support(s, m.Author.ID) || hard_admin(s, m.Author.ID) || m.Author.ID == kazka_u || m.ChannelID == botroom_c {

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

	//staff commands
	if hard_admin(s, m.Author.ID) || m.Author.ID == kazka_u || m.ChannelID == botroom_c {

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

		//!xpa [@] [amount]
		if strings.HasPrefix(m.Content, "!xpa") {
			if hard_admin(s, m.Author.ID) {
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
			if hard_admin(s, m.Author.ID) {
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

	}

	//kazka commands
	if m.Author.ID == kazka_u || m.ChannelID == botroom_c {

		//!test
		if strings.Contains(m.Content, "!join") {
			//if join bot channel
			mute := false
			deaf := false
			var err error
			VC, err = s.ChannelVoiceJoin(CAD_s, "524267098977992709", mute, deaf)
			if err != nil {
				fmt.Println(err.Error() + " Error !test")
			} else {
				if VC.Ready == true {
					s.ChannelMessageSend(m.ChannelID, "VC joined and ready 2 go!")
				}
			}
		}

		if strings.Contains(m.Content, "!leave") {
			// if bot channel empty
			VC.Disconnect()
			s.ChannelMessageSend(m.ChannelID, "C ya VC!")
		}

		if strings.HasPrefix(m.Content, "!talking") {
			VC.Speaking(true)
			return
		}
		VC.Speaking(false)

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
			if hard_admin(s, m.Author.ID) {
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

	if BJ[0] == "!roulette" {
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

	if m.ChannelID == eventsnbots_c || m.ChannelID == casino_c || hard_admin(s, m.Author.ID) || hard_support(s, m.Author.ID) || m.Author.ID == kazka_u || m.ChannelID == botroom_c {

		//!buy
		if m.Content == "!buy" {
			s.ChannelMessageSend(m.ChannelID, "!buy color [HEX color] - 1,000,000 slices, changes the color of the color role")
		}
		/*
			//!buy colorrole
			if m.Content == "!buy colorrole" || m.Content == "!buy colourrole" {
				slices_s := currency_get(m.Author.ID)
				slices_b, _ := strconv.Atoi(slices_s)
				if slices_b >= 100000000 {
					currency_adjust(m.ChannelID, -100000000, m.Author.ID)
					s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, color_r)
					s.ChannelMessageSend(m.ChannelID, "You have opted into the color role!")
				} else {
					s.ChannelMessageSend(m.ChannelID, "OOF, You only have "+slices_s+" slices!")
				}
			}
		*/
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
	}
	/*
		member, err := s.State.Member("409907314045353984", AUTHORID)
		if err != nil {
			fmt.Println(err.Error() + " Error 4")
		}
		for _, RoleID := range member.Roles {
			if RoleID == support_r || RoleID == staff_r || RoleID == admin_r || RoleID == moderator_r {
				return true
			}
		}
		return false
	*/
	return true
}

func test_room(CHANNELID string) bool {
	if CHANNELID == "442493156584587265" || CHANNELID == "410522839548952596" || CHANNELID == "403460796106932225" {
		return true
	} else {
		return false
	}
}

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
			level_int += strings.Count(m.Content, " ") + 1

		} else {
			level_int += 10
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

	switch rand.Int() % 5 {
	case 0:
		s.ChannelMessageSend("409907314045353986", "Welcome to the Chill Asexual Discord <@"+ma.Member.User.ID+">!")
	case 1:
		s.ChannelMessageSend("409907314045353986", "Welcome to Chili's <@"+ma.Member.User.ID+">!")
	case 2:
		s.ChannelMessageSend("409907314045353986", "Khajiit <@"+ma.Member.User.ID+"> has wares if you have SLICES...")
	case 3:
		s.ChannelMessageSend("409907314045353986", "Bathroom is the 4th door on the left <@"+ma.Member.User.ID+">.")
	case 4:
		s.ChannelMessageSend("409907314045353986", "<@"+ma.Member.User.ID+"> is here, Please clap.")
	}
}

func MemberLeaveHandler(s *discordgo.Session, ma *discordgo.GuildMemberRemove) {

	switch rand.Int() % 10 {
	case 0:
		s.ChannelMessageSend("409907314045353986", "I'm sorry <@"+ma.Member.User.ID+">, You are not the biggest loser.")
	case 1:
		s.ChannelMessageSend("409907314045353986", "I'm sorry <@"+ma.Member.User.ID+">, You have been voted off the island.")
	case 2:
		s.ChannelMessageSend("409907314045353986", "*psst* ||I heard <@"+ma.Member.User.ID+"> wasn't even chill anyway :/.||")
	case 3:
		s.ChannelMessageSend("409907314045353986", "Don't let the door hit you on the ass on your way out <@"+ma.Member.User.ID+">.")
	case 4:
		s.ChannelMessageSend("409907314045353986", "You are the weakest link. Goodbye <@"+ma.Member.User.ID+">.")
	case 5:
		s.ChannelMessageSend("409907314045353986", "We didn't like you anyway <@"+ma.Member.User.ID+">. 🖕🖕🖕")
	case 6:
		s.ChannelMessageSend("409907314045353986", "<@"+ma.Member.User.ID+">, didn't even curtesy flush.")
	case 7:
		s.ChannelMessageSend("409907314045353986", "<@"+ma.Member.User.ID+">'s CIS HET Otherkin Otaku Senpai is in another castle!")
	case 8:
		s.ChannelMessageSend("409907314045353986", "Let me guess <@"+ma.Member.User.ID+">... someone stole your cake.")
	case 9:
		s.ChannelMessageSend("409907314045353986", "Your cold never bothered us anyway <@"+ma.Member.User.ID+">...")
	}
}

// ideas
// daily pot
// bingo
// hard code Kazka
// update help
// play VC
// playing... command
