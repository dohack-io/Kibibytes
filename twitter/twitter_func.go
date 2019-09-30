package twitter

import (
	"Kibibytes/utils/secrets"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var consumerKey = secrets.Get("TWITTER_CONSUMER_KEY")
var consumerSecret = secrets.Get("TWITTER_CONSUMER_SECRET")
var accessToken = secrets.Get("TWITTER_ACCESS_TOKEN")
var accessSecret = secrets.Get("TWITTER_ACCESS_SECRET")

const userToFind = "polizei_nrw_do"

//noinspection GoUnusedParameter
func Command(args string) string {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	userTimelineParams := &twitter.UserTimelineParams{ScreenName: userToFind, Count: 5}
	tweets, resp, _ := client.Timelines.UserTimeline(userTimelineParams)
	if resp.StatusCode == 404 {
		return "Something went wrong! Try again :)"
	}

	m := ""
	for _, tweet := range tweets {
		m += tweet.Text + "\n"
	}

	return m
}
