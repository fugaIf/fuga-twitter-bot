package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/g8rswimmer/go-twitter"
	"log"
	"net/http"
	"os"
	"strings"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

type tokenStr struct {
	TokenText string `json:"token"`
}

func loadToken() (*tokenStr, error) {
	raw, err := os.ReadFile("./env/token.json")
	if err != nil {
		log.Fatal("loadToken os.Open err:", err)
		return nil, err
	}
	var ts tokenStr

	json.Unmarshal(raw, &ts)
	return &ts, err
}

/**
	In order to run, the user will need to provide the bearer token and the list of ids.
**/
func main() {
	loadToken, err := loadToken()
	token := flag.String("token", loadToken.TokenText, "twitter API token")
	ids := flag.String("ids", "1494700117143658500", "twitter ids")
	flag.Parse()

	tweet := &twitter.Tweet{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	fieldOpts := twitter.TweetFieldOptions{
		Expansions:  []twitter.Expansion{twitter.ExpansionEntitiesMentionsUserName, twitter.ExpansionAuthorID},
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldConversationID, twitter.TweetFieldAttachments},
	}

	lookups, err := tweet.Lookup(context.Background(), strings.Split(*ids, ","), fieldOpts)
	var tweetErr *twitter.TweetErrorResponse
	switch {
	case errors.As(err, &tweetErr):
		printTweetError(tweetErr)
	case err != nil:
		fmt.Println(err)
	default:
		for _, lookup := range lookups {
			printTweetLookup(lookup)
			fmt.Println()
		}
	}
}

func printTweetLookup(lookup twitter.TweetLookup) {
	enc, err := json.MarshalIndent(lookup, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(enc))
}

func printTweetError(tweetErr *twitter.TweetErrorResponse) {
	enc, err := json.MarshalIndent(tweetErr, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(enc))
}
