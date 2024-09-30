package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/mohammedfaizan/gator/internal/config"
	"github.com/mohammedfaizan/gator/internal/database"
)

const dbURL = "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

type state struct {
	db *database.Queries
	cfg *config.Config
}

func main()  {
	//creating database link
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("error with database link connection")
	}
	if db == nil {
		log.Fatal("db is nil")
	}

	dbQueries := database.New(db)

	cfg := &config.Config{}
	err = cfg.Read()
	if err != nil {
		fmt.Println(err)
	}

	st := state{
		db: dbQueries,
		cfg: cfg,
	}

	cmds := commands{
		Handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", func(s *state, c command) error {
		return st.HandleLogin(c)
	})
	
	// Register the `register` command
	cmds.register("register", func(s *state, c command) error {
		return registerHandler(s, c.ArgSlice)
	})

	// Register the `reset` command
	cmds.register("reset", func(s *state, c command) error {
		return resetHandler(s)
	})

	cmds.register("users", func(s *state, c command) error {
		return usersHandler(s)
	})

	cmds.register("agg", func(s *state, c command) error {
		return aggHandler(command{
			Name: "url",
			ArgSlice: []string{"https://wagslane.dev/index.xml", "10s"},
		})
	})

	cmds.register("addfeed", middlewareLoggedIn(addFeedHandler))

	cmds.register("feeds", func(s *state, c command) error {
		return feedHandler(s)
	})

	cmds.register("follow", middlewareLoggedIn(followHandler))

	cmds.register("following", middlewareLoggedIn(followingHandler))

	cmds.register("unfollow", middlewareLoggedIn(unfollowHandler))

	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		fmt.Println("Error: Not enough arguments provided")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		Name: cmdName,
		ArgSlice: cmdArgs,
	}

	err = cmds.Run(&st, cmd)
	
	if err != nil {
		fmt.Println("Error running command:", err)
		os.Exit(1)
	}

	
}