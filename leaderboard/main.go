package leaderboard

import (
	"fmt"
	"log"
	"sort"

	"dont-slack-evil/apphome"
	dsedb "dont-slack-evil/db"

	"github.com/slack-go/slack"
)

// SendLeaderboardNotification sends the leaderboard notification
func SendLeaderboardNotification() (int, error) {
	notificationsSent := 0
	teams, teamsErr := dsedb.GetTeams()
	if teamsErr != nil {
		log.Println(teamsErr)
		return 0, teamsErr
	}
	for _, team := range teams {
		// FIXME: remove
		if (team.SlackTeamId != "TU7KB9FB9") {
			continue
		}
		slackBotUserApiClient := slack.New(team.SlackBotUserToken)
		users, err := slackBotUserApiClient.GetUsers()
		if err != nil {
			log.Printf("Could not instantiate bot client for team %v", team.SlackTeamId)
			continue
		}
		type UserScore struct {
			ID string
			Name string
			Good int
			Total int
			Score float64
		}
		var userScores []UserScore;
		for _, user := range users {
			// This is the best way I found to distinguish bots from real users
			// Note that user.IsBot doesn't work because it's false even for bot users...
			log.Println(user.RealName)
			log.Println(user.Profile.BotID)
			if (len(user.Profile.BotID) == 0) {
				good, total := apphome.GetWeeklyStats(user.ID)
				var score float64;
				if (total > 0) {
					score = float64(good) / float64(total)
				} else {
					score = 0
				}
				userScore := UserScore{
					ID: user.ID,
					Name: user.RealName,
					Good: good,
					Total: total,
					Score: score,
				}
				userScores = append(userScores, userScore)
			}
		}
		// Sort by positivity scores
		sort.Slice(userScores, func(i, j int) bool {
			return userScores[i].Score > userScores[j].Score
		})
		var text = "*Weekly positivity rankings:*"
		if (len(userScores) > 0) {
			text += fmt.Sprintf(
				"\n\nCongratulations to <@%s> for being the most positive person this week :tada:",
				userScores[0].ID,
			)
			text += "\n\nHere are the standings:"
			text += fmt.Sprintf(
				"\n:first_place_medal: <@%s> with a %.2f score (%d / %d)",
				userScores[0].ID,
				userScores[0].Score * 100,
				userScores[0].Good,
				userScores[0].Total,
			)
		}
		if (len(userScores) > 1) {
			text += fmt.Sprintf(
				"\n:second_place_medal: <@%s> with a %.2f score (%d / %d)",
				userScores[1].ID,
				userScores[1].Score * 100,
				userScores[1].Good,
				userScores[1].Total,
			)
		}
		if (len(userScores) > 2) {
			text += fmt.Sprintf(
				"\n:third_place_medal: <@%s> with a %.2f score (%d / %d)",
				userScores[2].ID,
				userScores[2].Score * 100,
				userScores[2].Good,
				userScores[2].Total,
			)
		}
	}

	return notificationsSent, nil
}