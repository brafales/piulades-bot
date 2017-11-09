package twitter

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"io/ioutil"

	twitterAPI "github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

//Client allows you to communicate with the Twitter API
type Client struct {
	client *twitterAPI.Client
}

//Tweet holds information about a specific twitter status
type Tweet struct {
	apiTweet        *twitterAPI.Tweet
	imageDownloader ImageDownloader
}

//ImageDownloader provides a means to download an image from a URL
type ImageDownloader interface {
	Get(string) ([]byte, error)
}

//DefaultImageDownloader implements a basic ImageDownloader
type DefaultImageDownloader struct{}

//Get downloads a URL and returns it as a []byte
func (d DefaultImageDownloader) Get(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return []byte{}, nil
	}
	if statusCode := response.StatusCode; statusCode != 200 {
		return []byte{}, fmt.Errorf("Bad response code: %d", statusCode)
	}
	return ioutil.ReadAll(response.Body)
}

//NewClient returns a new Client instance.
//apiKey is the twitter API key and apiSecret is the
//twitter API secret
func NewClient(apiKey string, apiSecret string) *Client {
	config := &clientcredentials.Config{ClientID: apiKey,
		ClientSecret: apiSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token"}
	httpClient := config.Client(oauth2.NoContext)
	newClient := Client{client: twitterAPI.NewClient(httpClient)}
	return &newClient
}

//NewTweet returns a new instance of Tweet from a low level API response
func NewTweet(apiTweet *twitterAPI.Tweet, imageDownloader ImageDownloader) Tweet {
	return Tweet{apiTweet: apiTweet, imageDownloader: imageDownloader}
}

//GetTwit returns a Tweet associated to the statusID
func (t Client) GetTwit(statusID int64) (Tweet, error) {
	params := &twitterAPI.StatusShowParams{
		TweetMode: "extended",
	}
	apiTweet, _, err := t.client.Statuses.Show(statusID, params)
	if err == nil {
		return NewTweet(apiTweet, DefaultImageDownloader{}), nil
	}
	return Tweet{}, err
}

//PrintableText returns a string with the tweet information
func (t Tweet) PrintableText(user string) string {
	return fmt.Sprintf("Tweet enviat per [%s]: https://twitter.com/%s/status/%d - %s",
		user,
		t.apiTweet.User.ScreenName,
		t.apiTweet.ID,
		t.apiTweet.Text)
}

//URL returns the original tweet URL
func (t Tweet) URL() string {
	return fmt.Sprintf("https://twitter.com/%s/status/%d",
		t.apiTweet.User.ScreenName,
		t.apiTweet.ID)
}

//ExtendedEntities returns a slice with all the raw images of the Tweet
func (t Tweet) ExtendedEntities() ([][]byte, error) {
	if t.apiTweet.ExtendedEntities != nil && len(t.apiTweet.ExtendedEntities.Media) > 1 {
		entities := make([][]byte, len(t.apiTweet.ExtendedEntities.Media), len(t.apiTweet.ExtendedEntities.Media))
		for i, entity := range t.apiTweet.ExtendedEntities.Media {
			entityData, err := t.imageDownloader.Get(entity.MediaURL)
			if err != nil {
				return entities, err
			}
			entities[i] = entityData
		}
		return entities, nil
	}
	return make([][]byte, 0, 0), nil
}

//GetStatusID gives you the twitter status id for a given text
func GetStatusID(text string) (int64, error) {
	r := regexp.MustCompile(`https?://(www\.)?twitter\.com/\w+/status/(?P<id>\d+)`)
	matches := r.FindStringSubmatch(text)
	if matches != nil {
		id, err := strconv.ParseInt(matches[len(matches)-1], 10, 64)
		if err == nil {
			return id, nil
		}
	}
	return 0, errors.New("Could not find a twitter status id")
}
