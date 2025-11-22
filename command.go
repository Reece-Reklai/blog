package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Reece-Reklai/blog/internal"
	"github.com/Reece-Reklai/blog/internal/database"
	"github.com/google/uuid"
)

type configState struct {
	db  *database.Queries
	cfg *internal.Config
}

type command struct {
	name        string
	commandArgs []string
}

type cli struct {
	cliCommand map[string]func(*configState, command) error
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (cliHandler *cli) run(userState *configState, command command) error {
	runHandler, ok := cliHandler.cliCommand[command.name]
	if ok == false {
		return errors.New("unknown command")
	}
	err := runHandler(userState, command)
	if err != nil {
		return err
	}
	return nil
}

func (cliHandler *cli) register(name string, currentFunc func(*configState, command) error) error {
	_, ok := cliHandler.cliCommand[name]
	if ok == true {
		return errors.New("command exists")
	}
	cliHandler.cliCommand[name] = currentFunc
	return nil
}

func handlerLogin(loginState *configState, loginCommand command) error {
	user, err := loginState.db.GetUser(context.Background(), loginCommand.commandArgs[0])
	if user.Name != loginCommand.commandArgs[0] {
		return errors.New("user does not exist")
	}
	if user.Name == loginState.cfg.CurrentUserName {
		fmt.Println("user is already login")
		return nil
	}
	err = loginState.cfg.SetUser(loginCommand.commandArgs[0])
	if err != nil {
		fmt.Println("user was not set")
		return err
	}
	fmt.Println("user has login successfully")
	return nil
}

func handlerRegister(registerState *configState, registerCommand command) error {
	user, err := registerState.db.GetUser(context.Background(), registerCommand.commandArgs[0])
	if user.Name == registerCommand.commandArgs[0] {
		return errors.New("user is already registered")
	}
	createdUser, err := registerState.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: registerCommand.commandArgs[0]})
	if err != nil {
		fmt.Println(err)
		return errors.New("user was not created")
	}
	registerState.cfg.SetUser(createdUser.Name)
	fmt.Println("user has been set successfully after being registered into the database")
	return nil
}

func handlerUsers(userState *configState, _ command) error {
	users, err := userState.db.GetAllUsers(context.Background())
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to retrieve users from database")
	}
	for _, value := range users {
		if userState.cfg.CurrentUserName == value.Name {
			fmt.Printf("%s (current)\n", value.Name)
			continue
		}
		fmt.Println(value.Name)
	}
	return nil
}

func handlerReset(userState *configState, _ command) error {
	err := userState.db.DeleteAllUsers(context.Background())
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to remove all users")
	}
	userState.cfg.SetUser("")
	return nil
}

func handleUserFeed(userState *configState, _ command) error {
	userFeed := make(map[string][]string)
	feeds, err := userState.db.GetAllFeeds(context.Background())
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to retrieve feeds from database")
	}
	users, err := userState.db.GetAllUsers(context.Background())
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to retrieve users from database")
	}
	for _, valUser := range users {
		for _, valFeed := range feeds {
			if valUser.ID == valFeed.UserID {
				userFeed[valUser.Name] = append(userFeed[valUser.Name], valFeed.Name)
			}
		}
		_, ok := userFeed[valUser.Name]
		if ok == false {
			userFeed[valUser.Name] = append(userFeed[valUser.Name], "")
		}
	}
	for key, val := range userFeed {
		if val[0] == "" {
			fmt.Println(key)
			continue
		}
		for _, blog := range val {
			fmt.Println(blog)
		}
		fmt.Println(key)
	}
	return nil
}

func handlerAddFeed(userState *configState, cmd command, user database.User) error {
	userID := user.ID
	nameFeed := cmd.commandArgs[0]
	url := cmd.commandArgs[1]
	_, err := userState.db.CreateFeed(context.Background(), database.CreateFeedParams{Name: nameFeed, Url: url, UserID: userID})
	if err != nil {
		fmt.Println("failed to create feed(feed url needs to be unique)")
		return err
	}
	feed, err := userState.db.GetFeedURL(context.Background(), url)
	if err != nil {
		fmt.Println("failed to retrieve a feed")
		return err
	}
	_, err = userState.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		fmt.Println("failed to retrieve the feed followers")
		return err
	}
	return nil
}

func handlerFollow(userState *configState, cmd command, user database.User) error {
	feedURL := cmd.commandArgs[0]
	feed, err := userState.db.GetFeedURL(context.Background(), feedURL)
	if err != nil {
		fmt.Println("failed to retrieve a feed")
		return err
	}
	_, err = userState.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		fmt.Println("failed to retrieve the feed followers")
		return err
	}
	return nil
}

func handlerFollowing(userState *configState, cmd command, user database.User) error {
	following, err := userState.db.GetFeedFollowForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Println("failed to retrieve feed associated with user")
		return nil
	}
	fmt.Printf("Current User: %s\n", user.Name)
	for _, value := range following {
		fmt.Println(value.FeedName)
	}
	return nil
}

func handlerUnFollow(userState *configState, cmd command, user database.User) error {
	feedURL := cmd.commandArgs[0]
	feed, err := userState.db.GetFeedURL(context.Background(), feedURL)
	if err != nil {
		fmt.Println("failed to retrieve a feed")
		return err
	}
	err = userState.db.DeletFeedFollowByURL(context.Background(), database.DeletFeedFollowByURLParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		fmt.Println("failed to delete feed associated with user")
		return err
	}
	return nil
}

func middlewareLoggedIn(handler func(userState *configState, cmd command, user database.User) error) func(*configState, command) error {
	return func(userState *configState, cmd command) error {
		userName := userState.cfg.CurrentUserName
		user, err := userState.db.GetUser(context.Background(), userName)
		if err != nil {
			fmt.Println("failed to retrieve user")
			return err
		}
		handler(userState, cmd, user)
		return nil
	}
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var rss RSSFeed
	request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmt.Println("failed to create a request")
		return nil, err
	}
	request.Header.Set("User-Agent", "gator")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Failed to retrieve response")
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("failed to read from response body into bytes")
		return nil, err
	}
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		fmt.Println("failed to unmarshal the xml into the rss struct")
		return nil, err
	}
	defer response.Body.Close()
	return &rss, nil
}

func scrapeFeeds(userState *configState, cmd command) error {
	currentFeed, err := userState.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("failed to fetch next feed")
		return err
	}
	fmt.Println(currentFeed.Url)
	userState.db.MarkFeedFetched(context.Background(), currentFeed.ID)
	rss, err := fetchFeed(context.Background(), currentFeed.Url)
	if err != nil {
		fmt.Println("failed to fetch rss feed")
		fmt.Println(err)
	}
	channel := rss.Channel
	item := channel.Item
	// fmt.Println("--- channel ---")
	// fmt.Println(html.UnescapeString(channel.Title))
	// fmt.Println(channel.Link)
	// fmt.Println("--- items ---")
	for _, value := range item {
		fmt.Println(value.Title)
		break
	}
	return nil
}
