package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mohammedfaizan/gator/internal/database"
)

type command struct {
	Name string
	ArgSlice []string
}

type commands struct {
	Handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error)  {
	_, exists := c.Handlers[name]
	if exists {
		fmt.Printf("%s command already exists", name)
		return
	} 

	c.Handlers[name] = f

}




func (c *commands) Run(s *state, cmd command) error {
	_, exists := c.Handlers[cmd.Name]
	if !exists {
		return errors.New("command doesn't exist")
	}

	err := c.Handlers[cmd.Name](s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (s *state) HandleLogin(cmd command) error {
	if len(cmd.ArgSlice) == 0 {
		return errors.New("no commands given")
	}

	name := cmd.ArgSlice[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(name, s.cfg.DbUrl)
	if err != nil {
		return errors.New("couldn't login user")
	}

	log.Printf("User %s has been set", cmd.ArgSlice[0])
	return nil
}

func registerHandler(state *state, args []string) error {
	if len(args) < 1 {
		return errors.New("please provide a name")
	}
	name := args[0]

	dbUser, err := state.db.GetUser(context.Background(), name)
	if err != nil {
		log.Println("user exists macha")
	}

	if dbUser.Name == name {

		log.Println("user exists already")
		os.Exit(1)
	}

	newUser, err := state.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.NewString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
	})

	if err != nil {
		log.Fatal(err)
		return errors.New("user already exists")
	}

	
	err = state.cfg.SetUser(newUser.Name, dbURL)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully:")
	
	return nil
}

func resetHandler(state *state) error {
	
	

	err := state.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unexpected error during records deletion")
	}

	
	fmt.Println("Records were reset successfully")
	return nil
}