package twitter_test

import (
	"fmt"
	"testing"

	"errors"

	"bytes"

	"github.com/brafales/piulades-bot/twitter"
	twitterAPI "github.com/dghubble/go-twitter/twitter"
)

type goodTestImageDownloader struct{}

func (g goodTestImageDownloader) Get(url string) ([]byte, error) {
	return []byte{}, nil
}

type badTestImageDownloader struct{}

func (b badTestImageDownloader) Get(url string) ([]byte, error) {
	return []byte{}, errors.New("Error")
}

//GetStatusID
func TestGetsStatusIDFromHTTPTwitterLink(t *testing.T) {
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod http://twitter.com/user/status/864174389435729920 tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. httpExcepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	act, err := twitter.GetStatusID(text)
	if err != nil {
		t.Error(err)
	}
	var exp int64 = 864174389435729920
	if act != exp {
		t.Errorf("Expected %v, actual result was %v", exp, act)
	}
}

func TestGetsStatusIDFromHTTPSTwitterLink(t *testing.T) {
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod https://twitter.com/user/status/864174389435729920 tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. httpExcepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	act, err := twitter.GetStatusID(text)
	if err != nil {
		t.Error(err)
	}
	var exp int64 = 864174389435729920
	if act != exp {
		t.Errorf("Expected %v, actual result was %v", exp, act)
	}
}

func TestGetsStatusIDFromHTTPAndWWWTwitterLink(t *testing.T) {
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod http://www.twitter.com/user/status/864174389435729920 tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. httpExcepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	act, err := twitter.GetStatusID(text)
	if err != nil {
		t.Error(err)
	}
	var exp int64 = 864174389435729920
	if act != exp {
		t.Errorf("Expected %v, actual result was %v", exp, act)
	}
}

func TestGetsStatusIDFromHTTPSAndWWWTwitterLink(t *testing.T) {
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod https://www.twitter.com/user/status/864174389435729920 tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. httpExcepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	act, err := twitter.GetStatusID(text)
	if err != nil {
		t.Error(err)
	}
	var exp int64 = 864174389435729920
	if act != exp {
		t.Errorf("Expected %v, actual result was %v", exp, act)
	}
}

func TestReturnsErrorIfNoTwitterLink(t *testing.T) {
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. httpExcepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	_, err := twitter.GetStatusID(text)
	if err == nil {
		t.Error("Expected an error")
	}
}

//ExtendedEntities
func TestReturnsEntities(t *testing.T) {
	urlOne := "1"
	urlTwo := "2"
	entities := twitterAPI.ExtendedEntity{
		Media: []twitterAPI.MediaEntity{
			twitterAPI.MediaEntity{
				MediaURL: urlOne,
			},
			twitterAPI.MediaEntity{
				MediaURL: urlTwo,
			},
		},
	}

	apiTweet := twitterAPI.Tweet{
		ExtendedEntities: &entities,
	}

	tweet := twitter.NewTweet(&apiTweet, goodTestImageDownloader{})

	exp := make([][]byte, 2, 2)

	act, err := tweet.ExtendedEntities()
	if err != nil {
		t.Error(err)
	}
	if len(exp) != len(act) && (!bytes.Equal(exp[0], act[0]) || !bytes.Equal(exp[1], act[1])) {
		t.Errorf("Expected %v, got %v", exp, act)
	}
}

//URL
func TestReturnsTweetURL(t *testing.T) {
	name := "name"
	var ID int64 = 12345
	user := twitterAPI.User{
		ScreenName: name,
	}
	apiTweet := twitterAPI.Tweet{
		User: &user,
		ID:   ID,
	}

	tweet := twitter.NewTweet(&apiTweet, goodTestImageDownloader{})

	exp := fmt.Sprintf("https://twitter.com/%s/status/%d", name, ID)
	if act := tweet.URL(); act != exp {
		t.Errorf("Expected %s, actual result was %s", exp, act)
	}
}

//PrintableText
func TestReturnsTweetPrintableText(t *testing.T) {
	name := "name"
	text := "text"
	username := "user"
	var ID int64 = 12345
	user := twitterAPI.User{
		ScreenName: name,
	}
	apiTweet := twitterAPI.Tweet{
		User:     &user,
		ID:       ID,
		FullText: text,
	}

	tweet := twitter.NewTweet(&apiTweet, goodTestImageDownloader{})

	exp := fmt.Sprintf("Tweet enviat per [%s]: https://twitter.com/%s/status/%d\n%s",
		username,
		name,
		ID,
		text)
	if act := tweet.PrintableText(username); act != exp {
		t.Errorf("Expected %s, actual result was %s", exp, act)
	}
}
