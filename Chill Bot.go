package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/bwmarrin/discordgo"
)

// Append this string with error messages
var error_string string

var BotID string

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

func main() {

	b, err := ioutil.ReadFile("TOKEN.txt") // b has type []byte
	if err != nil {
		log.Fatal(err)
	}
	token := string(b)

	dg, err := discordgo.New("Bot " + token)
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
	// this if for telling if people join or leave, not used rn
	//dg.AddHandler(memberHandler)

	err = dg.Open()

	if err != nil {
		fmt.Println(err.Error() + " Error 3")
		return
	}

	fmt.Println("Bot is running POGGERS")

	<-make(chan struct{})
	return
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return
	} else {

		activity_tracker(m)
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

	if strings.Contains(m.Content, "gachiBASS") {
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

	return false
}

func commands(s *discordgo.Session, m *discordgo.MessageCreate) (err bool) {
	if m.Content == "!who" {
		s.ChannelMessageSend(m.ChannelID, "I am a program running off Kazka's computer that he is currently working on for our discord. I cant do much right now, but I can react to certain keywords and bad words. If you want to see my commands, you can type ``!help``")
	}

	if m.Content == "!xp" /*|| m.Content == "!fuckyoukio" */ {
		xp := xp_get(m.Author.ID)
		s.ChannelMessageSend(m.ChannelID, "You've got "+xp+" experience!")

		xp_i, _ := strconv.Atoi(level_get(xp))
		if xp_i > 1 {
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "410520872676360193")
		}
		if xp_i >= 10 {
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "446112365055049729")
		}

		if xp_i < 3 {
			s.ChannelMessageSend(m.ChannelID, "Oh, and you're only level "+level_get(xp)+". <:OMEGALUL:434488099083386890>")
		} else if xp_i < 10 {
			s.ChannelMessageSend(m.ChannelID, "You're also level "+level_get(xp)+". <:HAhaa:434488099083649049>")
		} else if xp_i >= 10 {
			s.ChannelMessageSend(m.ChannelID, "Congrats you're level "+level_get(xp)+". <:FeelsGoodMan:434360842985668608> <a:Clap:434360842620895253>")
		}
	}

	if m.Content == "!squirrel" {
		if hard_admin(s, m.Author.ID) || hard_support(s, m.Author.ID) {
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
		} else {
			s.ChannelMessageSend(m.ChannelID, "There aint nuthin there...")
		}
	}

	if m.Content == "!slices" {
		slices := currency_get(m.Author.ID)
		s.ChannelMessageSend(m.ChannelID, "You have "+slices+" slices!")
	}

	if m.Content == "!overflow" {

	}

	if m.Content == "!help" {
		s.ChannelMessageSend(m.ChannelID, "!help - gives a list of commands \n!help keywords - list of keywords \n!xp - lists experience points & level \n!terminate - turns the bot off \n!who - who am I, you ask?\n!help blackjack - help for asexual blackjack\n!optin & !optout - used for gaining access and leaving our optin channels\n!help roulette - roulette help")
	}

	if m.Content == "!help keywords" {
		s.ChannelMessageSend(m.ChannelID, "monkaS, FeelsGoodMan, FeelsBadMan, Clap, haHAA, DansGame, POGGERS, ecksdee, good bot, D:, Dab, OMEGALUL, gachiBASS")
	}

	if m.Content == "!help blackjack" {
		s.ChannelMessageSend(m.ChannelID, "Type ``!blackjack`` or ``!bj`` followed by your bet to start the game. For example you can bet 10 slices like so: ```!bj 10```")
		s.ChannelMessageSend(m.ChannelID, "In Asexual BlackJack Kings and Queens are worth ``1`` point, and Aces are the most at ``11`` points.")
		s.ChannelMessageSend(m.ChannelID, "You win by hitting 21 points first, not going above 21 first, or staying when you have more points.")
		s.ChannelMessageSend(m.ChannelID, "Draw another card with ``!hit``, both you and the dealer draw a card.")
		s.ChannelMessageSend(m.ChannelID, "Compare cards with ``!stay``, if you're feeling confident that is.")
		s.ChannelMessageSend(m.ChannelID, "Force end a blackjack session with ``!bjCLEAR`` **DO NOT ABUSE**")
	}
	if m.Content == "!help roulette" {
		s.ChannelMessageSend(m.ChannelID, "Bruh it's roulette. >.>")
		s.ChannelMessageSend(m.ChannelID, "Type ``!roulette`` followed by your bet to play, for example: ```!roulette 10```")
	}

	if m.Content == "!terminate" {
		if hard_admin(s, m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, "A'ight, peace out yall.")
			os.Exit(-1)
		} else {
			s.ChannelMessageSend(m.ChannelID, "The FUCK do you think you are???")
		}
	}

	if strings.HasPrefix(m.Content, "!xpc ") {
		trim := strings.TrimPrefix(m.Content, "!xpc ")

		xp := xp_get(trim)
		lvl := level_get(xp)

		if xp != " " {
			s.ChannelMessageSend(m.ChannelID, "ID "+trim+" has "+xp+" xp "+" and is level "+lvl+".")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Looks I had an issue finding that ID. <:FeelsBadMan:434360842377625631>")
		}
	}

	if strings.HasPrefix(m.Content, "!optin") {
		if strings.Contains(m.Content, "explicit2") {
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "504815202320121857")
			s.ChannelMessageSend(m.ChannelID, "You have opted in to \"definitely-explicit\"")
		} else if strings.Contains(m.Content, "explicit") {
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "424210978746400768")
			s.ChannelMessageSend(m.ChannelID, "You have opted in to \"potentially-explicit\"")
		}
		if strings.Contains(m.Content, "playpen") {
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "503075855619063818")
			s.ChannelMessageSend(m.ChannelID, "You have opted in to \"playpen\"")
		}
		if strings.Contains(m.Content, "casino") {
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "451256712658288643")
			s.ChannelMessageSend(m.ChannelID, "You have opted in to \"casino\"")
		}

		if strings.Contains(m.Content, "japan") {
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "455495577581387794")
			s.ChannelMessageSend(m.ChannelID, "You have opted in to \"japan\"")
		}

		if strings.Contains(m.Content, "gameroom") {
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "455495472874782720")
			s.ChannelMessageSend(m.ChannelID, "You have opted in to \"gameroom\"")
		}

		if strings.Contains(m.Content, "testroom") {
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "458282517062221859")
			s.ChannelMessageSend(m.ChannelID, "You have opted in to \"testroom\"")
		}

		if strings.Contains(m.Content, "makerroom") {
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "465350853847416843")
			s.ChannelMessageSend(m.ChannelID, "You have opted in to \"makerroom\"")
		}

		if strings.Contains(m.Content, "furryroom") {
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "465350963251904512")
			s.ChannelMessageSend(m.ChannelID, "You have opted in to \"furryroom\"")
		}

		if strings.Contains(m.Content, "garage") {
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "467166719908249602")
			s.ChannelMessageSend(m.ChannelID, "You have opted in to \"garage\"")
		}

		if strings.Contains(m.Content, "all") {
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "424210978746400768")
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "451256712658288643")
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "455495577581387794")
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "455495472874782720")
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "458282517062221859")
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "465350853847416843")
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "465350963251904512")
			s.GuildMemberRoleAdd("409907314045353984", m.Author.ID, "467166719908249602")

			s.ChannelMessageSend(m.ChannelID, "You have opted in to all channels")
		}

	}

	if strings.HasPrefix(m.Content, "!optout") {
		if strings.Contains(m.Content, "explicit") {
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "424210978746400768")
			s.ChannelMessageSend(m.ChannelID, "You have opted out of \"explicit\"")
		}

		if strings.Contains(m.Content, "playpen") {
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "503075855619063818")
			s.ChannelMessageSend(m.ChannelID, "You have opted out of \"playpen\"")
		}

		if strings.Contains(m.Content, "casino") {
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "451256712658288643")
			s.ChannelMessageSend(m.ChannelID, "You have opted out of \"casino\"")
		}

		if strings.Contains(m.Content, "japan") {
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "455495577581387794")
			s.ChannelMessageSend(m.ChannelID, "You have opted out of \"japan\"")
		}

		if strings.Contains(m.Content, "gameroom") {
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "455495472874782720")
			s.ChannelMessageSend(m.ChannelID, "You have opted out of \"gameroom\"")
		}

		if strings.Contains(m.Content, "testroom") {
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "458282517062221859")
			s.ChannelMessageSend(m.ChannelID, "You have opted out of \"testroom\"")
		}

		if strings.Contains(m.Content, "makerroom") {
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "465350853847416843")
			s.ChannelMessageSend(m.ChannelID, "You have opted out of \"makerroom\"")
		}

		if strings.Contains(m.Content, "furryroom") {
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "465350963251904512")
			s.ChannelMessageSend(m.ChannelID, "You have opted out of \"furryroom\"")
		}

		if strings.Contains(m.Content, "garage") {
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "467166719908249602")
			s.ChannelMessageSend(m.ChannelID, "You have opted out of \"garage\"")
		}

		if strings.Contains(m.Content, "all") {
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "424210978746400768")
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "451256712658288643")
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "455495577581387794")
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "455495472874782720")
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "458282517062221859")
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "465350853847416843")
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "465350963251904512")
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "467166719908249602")
			s.GuildMemberRoleRemove("409907314045353984", m.Author.ID, "503075855619063818")
			s.ChannelMessageSend(m.ChannelID, "You have opted out of all channels")
		}

	}

	if strings.HasPrefix(m.Content, "!hug") {

		if strings.HasSuffix(m.Content, "daddy") || strings.HasSuffix(m.Content, "Daddy") {
			s.ChannelMessageSend(m.ChannelID, "<a:gachiBASS:434488099163078667>")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Please don't.")
		}
	}

	// can probably make this a single command with the !d prefix and a number suffix, I made it in like 5 mins lol
	if strings.HasPrefix(m.Content, "!d20") {
		s.ChannelMessageSend(m.ChannelID, strconv.Itoa((rand.Int()%19)+1))
	}

	if strings.HasPrefix(m.Content, "!d6") {
		s.ChannelMessageSend(m.ChannelID, strconv.Itoa((rand.Int()%5)+1))
	}

	if strings.HasPrefix(m.Content, "!d10") {
		s.ChannelMessageSend(m.ChannelID, strconv.Itoa((rand.Int()%9)+1))
	}

	if strings.HasPrefix(m.Content, "!d4") {
		s.ChannelMessageSend(m.ChannelID, strconv.Itoa((rand.Int()%3)+1))
	}

	if strings.HasPrefix(m.Content, "!d12") {
		s.ChannelMessageSend(m.ChannelID, strconv.Itoa((rand.Int()%11)+1))
	}

	if strings.HasPrefix(m.Content, "!d8") {
		s.ChannelMessageSend(m.ChannelID, strconv.Itoa((rand.Int()%7)+1))
	}

	if strings.HasPrefix(m.Content, "!givea") {
		if hard_admin(s, m.Author.ID) {
			var Slices []string
			var SlicesID string
			Slices = strings.SplitN(m.Content, " ", 3)
			if Slices[0] == "!givea" {
				if len(Slices) > 1 {
					SlicesID = Slices[1]
					if len(Slices) > 2 {
						Slicesadjust, _ := strconv.Atoi(Slices[2])
						currency_adjust(m.ChannelID, Slicesadjust, SlicesID)
						s.ChannelMessageSend(m.ChannelID, "FREE SLICES! <:POGGERS:434360204147032074>")
					}
				}
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Uh... no, this my cake <:CapitalDColon:434488099255615488> B!")
		}
	}

	if strings.HasPrefix(m.Content, "!xpa") {
		if hard_admin(s, m.Author.ID) {
			var XP []string
			var XPID string
			XP = strings.SplitN(m.Content, " ", 3)
			if XP[0] == "!xpa" {
				if len(XP) > 1 {
					XPID = XP[1]
					if len(XP) > 2 {
						XPadjust, _ := strconv.Atoi(XP[2])
						xp_adjust(m.ChannelID, XPadjust, XPID)
						s.ChannelMessageSend(m.ChannelID, "XP Adjusted")
					}
				}
			}
		}
	}

	if strings.Contains(m.Content, "!test") {

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
				}
				blackjack.bet = bet_int
			} else {
				blackjack.bet = 1
			}
			currency := currency_get(m.Author.ID)
			currency_int, _ := strconv.Atoi(currency)
			if currency_int < blackjack.bet {
				s.ChannelMessageSend(m.ChannelID, "You only have "+currency+" you scamming fuck! >:(")
			} else if blackjack.bet > 1000000 {
				s.ChannelMessageSend(m.ChannelID, "Look pal, Tapioca and Kazka broke the bank, you can't bet over a mil any more :/ sorry...")
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

	return true
}

func test_room(CHANNELID string) bool {
	if CHANNELID == "442493156584587265" || CHANNELID == "410522839548952596" || CHANNELID == "403460796106932225" {
		return true
	} else {
		return false
	}
}

func hard_admin(s *discordgo.Session, AUTHORID string) bool {
	member, _ := s.State.Member("409907314045353984", AUTHORID)
	for _, RoleID := range member.Roles {
		if RoleID == "410522026868998146" {
			return true
		}
	}
	return false
}

func hard_support(s *discordgo.Session, AUTHORID string) bool {
	member, _ := s.State.Member("409907314045353984", AUTHORID)
	for _, RoleID := range member.Roles {
		if RoleID == "458036299732090892" {
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

		xp_int, _ := strconv.Atoi(string(xp_byte))

		if (strings.Count(" ", m.Content) + 1) < 10 {
			xp_int += strings.Count(m.Content, " ") + 1

		} else {
			xp_int += 10
		}

		xp_string := strconv.Itoa(xp_int)

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

		xp_int, _ := strconv.Atoi(string(xp_byte))

		if (strings.Count(" ", m.Content) + 1) < 10 {
			xp_int += strings.Count(m.Content, " ") + 1
		} else {
			xp_int += 10
		}

		xp_string := strconv.Itoa(xp_int)

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

		xp_int, _ := strconv.Atoi(string(xp_byte))

		xp_int += adjustment_value

		xp_string := strconv.Itoa(xp_int)

		xp_byte = []byte(xp_string)

		err = c_bucket.Put([]byte(ID), xp_byte)

		return nil
	})

	defer db.Close()
}

func level_get(xp string) (Level string) {
	var Level_i int
	Level_i = 0
	xp_i, _ := strconv.Atoi(string(xp))
	for xp_i >= (500 * (Level_i * Level_i)) {
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

/*
func memberHandler(s *discordgo.Session, ma *discordgo.GuildMemberAdd,) {
	s.ChannelMessageSend(m., "")


}
*/
// game ideas
// black jack
// daily pot
// bingo
// kreygasm
// a command to mute pram
