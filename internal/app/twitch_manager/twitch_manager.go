package twitch_manager

import (
	"errors"
	"fmt"
	"time"

	"github.com/bep/debounce"
	"github.com/nicklaw5/helix/v2"
	"github.com/vikpe/qw-demobot/internal/comms/topic"
	"github.com/vikpe/qw-demobot/internal/pkg/zeromq"
	"github.com/vikpe/qw-demobot/internal/pkg/zeromq/message"
)

const rateLimit = 10 * time.Second

func New(clientID, accessToken, broadcasterID, subscriberAddress string) (*zeromq.Subscriber, error) {
	apiClient, err := helix.NewClient(&helix.Options{
		ClientID:       clientID,
		AppAccessToken: accessToken,
	})

	subscriber := zeromq.NewSubscriber(subscriberAddress, topic.TwitchChannelSetTitle)

	if err != nil {
		err := errors.New(fmt.Sprintf("twitch api client error: %s", err))
		return subscriber, err
	}

	debounced := debounce.New(rateLimit)

	subscriber.OnMessage = func(msg message.Message) {
		changeTitle := func() {
			err := SetTitle(apiClient, broadcasterID, msg.Content.ToString())

			if err != nil {
				subscriber.OnError(err)
			}
		}
		debounced(changeTitle)
	}

	return subscriber, nil
}

func SetTitle(apiClient *helix.Client, broadcasterID, title string) error {
	const quakeGameId = "7348"
	fmt.Println("new title", title)

	_, err := apiClient.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID: broadcasterID,
		Title:         title,
		GameID:        quakeGameId,
	})

	return err
}
