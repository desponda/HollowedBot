package syncer

import (
	"context"
	"fmt"
	"github.com/FuzzyStatic/blizzard/v2"
	"github.com/FuzzyStatic/blizzard/v2/wowp"
	"os"
	"strconv"
)

type WoWSyncer struct {
	client *blizzard.Client
}

//TODO: Make rankMappings dynamic
var rankMappings = map[int]RankMapping{
	8: {
		RoleId:   "770760519266336769",
		RoleName: "Grunt",
	},
	7: {
		RoleId:   "770760519266336769",
		RoleName: "Trial",
	},
	6: {
		RoleId:   "772898890642751489",
		RoleName: "Sorcerer",
	},
	5: {
		RoleId:   "772899591749500939",
		RoleName: "Gladiator",
	},
	4: {
		RoleId:   "788985953561739286",
		RoleName: "Champion",
	},
}

func NewWoWSyncer() *WoWSyncer {
	blizz := blizzard.NewClient(os.Getenv("BLIZZARD_CLIENT_ID"), os.Getenv("BLIZZARD_CLIENT_SECRET"), blizzard.US, blizzard.EnUS)

	err := blizz.AccessTokenRequest(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &WoWSyncer{
		client: blizz,
	}
}

func (ws *WoWSyncer) SyncUserRanks(users []Member) []Member {
	//TODO: make server/guild dynamic
	gr, _, _ := ws.client.WoWGuildRoster(context.Background(), "illidan", "hollowed")

	for n, _ := range users {
		gameRank := ws.findMember(*gr, users[n].UserName)
		users[n].DiscRankId = rankMappings[gameRank].RoleId
		users[n].DiscRankName = rankMappings[gameRank].RoleName
		users[n] = ws.setSyncActions(users[n])
	}

	return users
}

func (ws *WoWSyncer) findMember(roster wowp.GuildRoster, member string) int {
	members := roster.Members
	for _, m := range members {
		if m.Character.Name == member {
			fmt.Println(m.Character.Name + " rank is " + strconv.Itoa(m.Rank))
			return m.Rank
		}
	}
	return -1
}

func (ws *WoWSyncer) setSyncActions(user Member) Member {
	for _, m := range rankMappings {
		if user.DiscRankId == m.RoleId {
			user.RolesToAdd = append(user.Roles, m.RoleId)
		} else {
			user.RolesToRemove = append(user.Roles, m.RoleId)
		}
	}
	return user
}
