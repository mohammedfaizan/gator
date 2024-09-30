package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

func usersHandler(state *state) error {
	users, err := state.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving the users from db")
	}

	for _, dbUser := range users {
		printUser(dbUser, state.cfg.CurrentUserName)
	}

	return nil
}

func aggHandler(c command) error {

	if len(c.ArgSlice) != 2 {
		return fmt.Errorf("2 args required")
	}

	urlString := c.ArgSlice[0]
	durationStr := c.ArgSlice[1]

	timeBetweenReqs, err := time.ParseDuration(durationStr)
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := fetchAndPrintFeed(urlString)
			if err != nil {
				fmt.Println("there was an error:", err)
			}
		}
	}
	
}

func fetchAndPrintFeed(urlString string) error {
	feed, err := fetchFeed(context.Background(), urlString)
	if err != nil {
		fmt.Println("there was an error")
		return err
	}

	fmt.Printf("Title: %v\n", feed.Channel.Title)
	fmt.Printf("Link: %v\n", feed.Channel.Link)
	fmt.Printf("Description: %v\n", feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		fmt.Printf("Item Title: %v", item.Title)
		fmt.Printf("Item Link: %v", item.Link)
		fmt.Printf("Item Description: %v", item.Description)
		fmt.Printf("Item Pubdate: %v", item.PubDate)
	}
	
	return nil
}

func addFeedHandler(s *state, c command, user database.User) error {

	if len(c.ArgSlice) != 2 {
		return fmt.Errorf("addFeed requires two parameters/args")
	}

	url := c.ArgSlice[1]
	currentUser, err := s.db.GetUser(context.Background(),s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("user doesn't exist")
	}
	nullUserId := sql.NullString {
		String: currentUser.ID,
		Valid: true,
	}

	newFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.NewString(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: c.ArgSlice[0],
		Url: url,
		UserID: nullUserId,
	})
	if err != nil {
		return fmt.Errorf("error retrieving feed")
	}

	nullFeedId := sql.NullString {
		String: newFeed.ID,
		Valid: true,
	}

	_ , err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.NewString(),
		UserID: nullUserId,
		FeedID: nullFeedId,
	})
	if err != nil {
		return fmt.Errorf("error creating feedfollow table")
	}
	log.Println("Feed follow also initialised")

	printFeed(newFeed)

	return nil

}

func feedHandler(s *state) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error with feeds")
	}

	for _, dbFeed := range feeds {
		printFeed(dbFeed)
		feedUser, err := s.db.GetUserById(context.Background(), dbFeed.UserID.String)
		if err != nil {
			return fmt.Errorf("error fetching user of feed")
		}
		fmt.Println("User Name: ", feedUser.Name)
	}

	return nil
}

func followHandler(s *state, c command, user database.User) error {
	if len(c.ArgSlice) != 1 {
		return fmt.Errorf("follow takes one arg 'url'")
	}

	url := c.ArgSlice[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't retrieve feed")
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		fmt.Println("err1", err)
		return err
	}


	user_id := sql.NullString {
		String: currentUser.ID,
		Valid: true,
	}
	feed_id := sql.NullString{
		String: feed.ID,
		Valid: true,
	}

	feedFollow, err := s.db.GetFeedByIDs(context.Background(), database.GetFeedByIDsParams{
		UserID: user_id,
		FeedID: feed_id,
	})
	if err != nil {
		_ , err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
			ID: uuid.NewString(),
			UserID: user_id,
			FeedID: feed_id,
		})
		if err != nil {
			return fmt.Errorf("error creating feedfollow table")
		}
		log.Println("Feed follow also initialised")
	}
	printfeedFollow(feedFollow)
	return nil
}

func followingHandler(s *state, c command, user database.User) error {

	currentUser, err := s.db.GetUser(context.Background(),s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("unexpected error with getting following")
	}

	feedFollowings, err := s.db.GetFeedFollowsForUser(context.Background(), sql.NullString{
		String: currentUser.ID,
		Valid: true,
	})
	if err != nil {
		return fmt.Errorf("unexpected error with getting followings")
	}

	for _, feed := range feedFollowings {
		printFeedFollowing(feed)
	}
	return nil
}


func unfollowHandler(s *state, cmd command, user database.User) error {

	if len(cmd.ArgSlice) != 1 {
		return fmt.Errorf("give two args")
	}

	url := cmd.ArgSlice[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	user_id := sql.NullString {
		String: user.ID,
		Valid: true,
	}

	feed_id := sql.NullString {
		String: feed.ID,
		Valid: true,
	}

	err = s.db.DeleteFeedFollows(context.Background(), database.DeleteFeedFollowsParams{
		UserID: user_id,
		FeedID: feed_id,
	})
	if err != nil {
		return err
	}

	log.Println("feeds unfollowed for user: ", user.Name)

	return nil
}


func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.ArgSlice) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.ArgSlice[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	user_id := sql.NullString {
		String: user.ID,
		Valid: true,
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user_id,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.Title)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}


func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("user not logged in")
		}

		return handler(s, c, user)
	}
}



func printUser(user database.User, current string) {
	if user.Name == current {
		fmt.Printf(" * %v (current)\n", user.Name)
	} else {
		fmt.Printf(" * %v\n", user.Name)
	}
}

func printFeed(feed database.Feed)  {
	fmt.Println("User Id: ", feed.UserID.String)
	fmt.Println("Feed name: ", feed.Name)
	fmt.Println("Feed URL: ", feed.Url)
}

func printFeedFollow(feedFollow database.CreateFeedFollowRow) {
	fmt.Println("Feed Name: ", feedFollow.FeedName)
	fmt.Println("User Name: ", feedFollow.UserName)
}

func printfeedFollow(feedFollow database.FeedFollow) {
	fmt.Println("Feed Name: ", feedFollow.FeedID)
	fmt.Println("User Name: ", feedFollow.UpdatedAt)
}


func printFeedFollowing(feedFollow database.GetFeedFollowsForUserRow) {
	fmt.Println("Feed Name: ", feedFollow.FeedName)
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("there was an err", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {

		Time, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			return err
		}

		desc := sql.NullString {
			String: item.Description,
			Valid: true,
		}

		pub_time := sql.NullTime {
			Time: Time,
			Valid: true,
		}
		
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.NewString(),
			Title: item.Title,
			Url: item.Link,
			Description: desc,
			PublishedAt: pub_time,
			FeedID: nextFeed.ID,
		})
		if err != nil {
			if isDuplicateError(err) {
				continue
			} else {
				log.Println("error inserting post: ", err)
			}
		}
	}

	return nil
}

func isDuplicateError(e error) bool {
	if e == nil {
        return false
    }
    if pqErr, ok := e.(*pq.Error); ok {
        // Check if the error code is "23505"
        if pqErr.Code == "23505" {
            return true
        }
    }
    return false
}