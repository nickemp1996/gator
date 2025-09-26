package main

import (
	"fmt"
	"os"
	"database/sql"
	"github.com/nickemp1996/gator/internal/config"
	"github.com/nickemp1996/gator/internal/database"
)

import _ "github.com/lib/pq"

type state struct {
	cfg		*config.Config
	queries	*database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config file: ", err)
		os.Exit(1)
	}

	s := state{cfg: &cfg}

	db, err := sql.Open("postgres", s.cfg.URL)
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)
	s.queries = dbQueries

	c := commands{handlers: make(map[string]func(*state, command) error)}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
	c.register("agg", handlerAgg)
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))	
	c.register("feeds", handlerFeeds)
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("following", middlewareLoggedIn(handlerFollowing))
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	c.register("browse", middlewareLoggedIn(handlerBrowse))

	allArgs := os.Args

	if len(allArgs) < 2 {
		fmt.Println("At least two arguments are required!")
		os.Exit(1)
	}

	var cmd command
	cmd.name = allArgs[1]
	cmd.args = allArgs[2:]

	err = c.run(&s, cmd)
	if err != nil {
		fmt.Println("Error running command:", err)
		os.Exit(1)
	}
}