package main

import (
	"fmt"
	"github.com/Lukaesebrot/dgc"
	"github.com/bwmarrin/discordgo"
	"github.com/desponda/HollowedBot/pkg/syncer"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	token := os.Getenv("TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	err = session.Open()
	if err != nil {
		panic(err)
	}

	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	}()

	router := dgc.Create(&dgc.Router{
		Prefixes: []string{
			"!",
		},
		IgnorePrefixCase: true,
		BotsAllowed:      false,
		Commands:         []*dgc.Command{},
		Middlewares:      []dgc.Middleware{},
		PingHandler: func(ctx *dgc.Ctx) {
			err = ctx.RespondText("Pong!")
			if err != nil {
				log.Println(err)
			}
		},
	})
	router.RegisterDefaultHelpCommand(session, nil)

	router.RegisterCmd(&dgc.Command{

		Name: "syncmembers",

		Aliases: []string{
			"syncmembers",
		},

		Description: "Responds with the rank of the member",
		Usage:       "syncmembers",
		Example:     "syncmembers",
		Flags:       []string{},
		IgnoreCase:  true,
		SubCommands: []*dgc.Command{},
		RateLimiter: dgc.NewRateLimiter(5*time.Second, 1*time.Second, func(ctx *dgc.Ctx) {
			err = ctx.RespondText("You are being rate limited!")
			log.Println(err)
		}),
		Handler: syncMembers,
	})

	router.RegisterCmd(&dgc.Command{

		Name: "requestrank",

		Aliases: []string{
			"requestrank",
		},

		Description: "Responds with the rank of the member",
		Usage:       "requestrank",
		Example:     "requestrank",
		Flags:       []string{},
		IgnoreCase:  true,
		SubCommands: []*dgc.Command{},
		RateLimiter: dgc.NewRateLimiter(5*time.Second, 1*time.Second, func(ctx *dgc.Ctx) {
			err = ctx.RespondText("You are being rate limited!")
			if err != nil {
				log.Println(err)
			}
		}),
		Handler: requestRank,
	})
	router.Initialize(session)
}

func requestRank(ctx *dgc.Ctx) {
	wSyncer := syncer.NewWoWSyncer()
	var users []syncer.Member
	memberName := ctx.Event.Member.Nick
	if memberName == "" {
		memberName = ctx.Event.Author.Username
	}

	users = append(users, syncer.Member{
		UserName: memberName,
		DiscId:   ctx.Event.Author.ID,
	})

	users = wSyncer.SyncUserRanks(users)
	if users[0].DiscRankName == "" {
		err := ctx.RespondText(memberName + ", I could not find your rank, you are an either an officer or your nickname is not set correctly")
		if err != nil {
			log.Println(err)
		}
	} else {
		err := ctx.RespondText(memberName + " assigning you the " + users[0].DiscRankName + " role")
		if err != nil {
			log.Println(err)
		}
	}
	executeUserSync(ctx, users)
}

func syncMembers(ctx *dgc.Ctx) {
	wSyncer := syncer.NewWoWSyncer()
	ms, err := ctx.Session.GuildMembers(ctx.Event.GuildID, "", 1000)
	var users []syncer.Member
	fmt.Println(err)
	for _, member := range ms {
		user := syncer.Member{}
		memberName := member.Nick
		if memberName == "" {
			memberName = member.User.Username
		}
		user.UserName = memberName
		user.DiscId = member.User.ID
		user.Roles = member.Roles
		users = append(users, user)
	}

	users = wSyncer.SyncUserRanks(users)
	executeUserSync(ctx, users)
}

func executeUserSync(ctx *dgc.Ctx, users []syncer.Member) {

	for _, user := range users {
		for _, rr := range user.RolesToRemove {
			err := ctx.Session.GuildMemberRoleRemove(ctx.Event.GuildID, user.DiscId, rr)
			if err != nil {
				log.Println(err)
			}
		}
		for _, ra := range user.RolesToAdd {
			err := ctx.Session.GuildMemberRoleAdd(ctx.Event.GuildID, user.DiscId, ra)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
