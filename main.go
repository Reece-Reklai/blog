package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Reece-Reklai/blog/internal"
	"github.com/Reece-Reklai/blog/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cliHandler := cli{cliCommand: make(map[string]func(*configState, command) error)}
	cliHandler.register("login", handlerLogin)
	cliHandler.register("register", handlerRegister)
	cliHandler.register("users", handlerUsers)
	cliHandler.register("reset", handlerReset)
	cliHandler.register("feeds", handleUserFeed)
	cliHandler.register("agg", scrapeFeeds)
	cliHandler.register("browse", middlewareLoggedIn(handlerBrowse))
	cliHandler.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cliHandler.register("follow", middlewareLoggedIn(handlerFollow))
	cliHandler.register("following", middlewareLoggedIn(handlerFollowing))
	cliHandler.register("unfollow", middlewareLoggedIn(handlerUnFollow))
	currentConfig, err := internal.Read()
	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Println("requires a command argument")
		os.Exit(1)
	}
	userCommand := command{name: arguments[1]}
	for index, value := range arguments {
		if index > 1 {
			userCommand.commandArgs = append(userCommand.commandArgs, value)
		}
	}
	if err != nil {
		fmt.Println("failed to read json, made its way to main file")
	}
	appState := configState{cfg: &currentConfig}
	db, err := sql.Open("postgres", appState.cfg.DbURL)
	if err != nil {
		fmt.Println(err)
	}
	dbQuery := database.New(db)
	appState.db = dbQuery
	switch userCommand.name {
	case "login":
		if len(arguments) != 3 {
			fmt.Println("only command and username should be given")
			os.Exit(1)
		}
		err := cliHandler.run(&appState, userCommand)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "register":
		if len(arguments) != 3 {
			fmt.Println("only command and username should be given")
			os.Exit(1)
		}
		err := cliHandler.run(&appState, userCommand)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "users":
		if len(arguments) != 2 {
			fmt.Println("requires only command")
			os.Exit(1)
		}
		err := cliHandler.run(&appState, userCommand)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "reset":
		if len(arguments) != 2 {
			fmt.Println("requires only command")
			os.Exit(1)
		}
		err := cliHandler.run(&appState, userCommand)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "browse":
		if len(arguments) != 2 {
			fmt.Println("requires only command ... for ... now")
			os.Exit(1)
		}
		err := cliHandler.run(&appState, userCommand)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "agg":
		if len(arguments) != 3 {
			fmt.Println("requires only command and duration string (1s, 1m, 1h) ")
			os.Exit(1)
		}
		timeBetweenRequests, err := time.ParseDuration(userCommand.commandArgs[0])
		if err != nil {
			fmt.Println("Invalid duration string")
			os.Exit(1)
		}
		fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
		ticker := time.NewTicker(timeBetweenRequests)
		defer ticker.Stop()
		for ; ; <-ticker.C {
			err = cliHandler.run(&appState, userCommand)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	case "addfeed":
		if len(arguments) != 4 {
			fmt.Println("requires a command, feed name, and feed url")
			os.Exit(1)
		}
		err := cliHandler.run(&appState, userCommand)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	case "feeds":
		if len(arguments) != 2 {
			fmt.Println("requires only command")
			os.Exit(1)
		}
		err := cliHandler.run(&appState, userCommand)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "follow":
		if len(arguments) != 3 {
			fmt.Println("requires a command and feed url to follow")
			os.Exit(1)
		}
		err := cliHandler.run(&appState, userCommand)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "following":
		if len(arguments) != 2 {
			fmt.Println("requires only command")
			os.Exit(1)
		}
		err := cliHandler.run(&appState, userCommand)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "unfollow":
		if len(arguments) != 3 {
			fmt.Println("requires a command and a feed url to unfollow")
			os.Exit(1)
		}
		err := cliHandler.run(&appState, userCommand)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		fmt.Println("unknown comannd")
	}
	currentConfig, err = internal.Read()
	if err != nil {
		fmt.Println(err)
	}
}
