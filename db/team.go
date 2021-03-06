package db

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fatih/structs"
	"github.com/slack-go/slack"
)

type ApiForTeam struct {
	Team                  Team
	SlackBotUserApiClient SlackApiInterface
}

type SlackApiInterface interface {
	// This interface is meant to make a *slack.Client mockable easily
	GetUsers() ([]slack.User, error)
	GetUserInfo(user string) (*slack.User, error)
	PostEphemeral(channelID, userID string, options ...slack.MsgOption) (string, error)
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
	PublishView(userID string, view slack.HomeTabViewRequest, hash string) (*slack.ViewResponse, error)
}

type IncomingWebhook struct {
	Channel          string `json:"channel"`
	ChannelID        string `json:"channel_id"`
	ConfigurationURL string `json:"configuration_url"`
	URL              string `json:"url"`
}

type Team struct {
	SlackTeamId            string          `json:"slack_team_id"`
	SlackBotUserToken      string          `json:"slack_bot_user_oauth_token"`
	SlackRegularOauthToken string          `json:"slack_regular_oauth_token"`
	IncomingWebhook        IncomingWebhook `json:"incoming_webhook"`
	Updated                time.Time       `json:"updated"`
}

func FindOrCreateTeamById(id string) (*Team, error) {
	// CreateTableIfNotCreated(tableName, "slack_team_id")
	team, findErr := FindTeamById(id)
	if findErr != nil {
		// TODO: check the error string, I wasn't able to make sure of this one
		if !strings.Contains(findErr.Error(), "Item does not exist") {
			return createTeamById(id)
		} else {
			log.Printf("%s", findErr)
			return nil, findErr
		}
	}
	log.Printf("Found team of ID %s called %s", team.SlackTeamId, "TODO")
	return team, nil
}

func createTeamById(id string) (*Team, error) {
	tableName := os.Getenv("DYNAMODB_TABLE_PREFIX") + "teams"
	team := &Team{SlackTeamId: id}
	if Store(tableName, structs.New(team).Map()) {
		return team, nil
	} else {
		log.Printf("Error creating a team of ID %s", id)
		return nil, errors.New("Could not create the team in DynamoDB")
	}
}

func FindTeamById(id string) (*Team, error) {
	tableName := os.Getenv("DYNAMODB_TABLE_PREFIX") + "teams"
	out, err := Get(tableName,
		map[string]*dynamodb.AttributeValue{
			"slack_team_id": {
				S: aws.String(id),
			},
		},
	)
	if err != nil {
		log.Printf("Error retrieving an item ID")
		return nil, err
	}

	var team Team
	unMarshallErr := dynamodbattribute.UnmarshalMap(out.Item, &team)
	if err != nil {
		log.Printf("Error unmarshalling a DynamoDB item ID as a Team")
		return nil, err
	}

	return &team, unMarshallErr
}

func GetTeams() ([]*Team, error) {
	tableName := os.Getenv("DYNAMODB_TABLE_PREFIX") + "teams"

	out, err := GetAll(tableName)

	if err != nil {
		log.Println(err)
		return []*Team{}, err
	}
	var teams []*Team
	unmarshalErr := dynamodbattribute.UnmarshalListOfMaps(out.Items, &teams)

	if unmarshalErr != nil {
		log.Printf("Error retrieving all teams")
		return nil, unmarshalErr
	}
	return teams, nil
}
