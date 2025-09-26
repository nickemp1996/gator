package main

import (
	"fmt"
	"time"
	"context"
	"strconv"
	"github.com/google/uuid"
	"github.com/nickemp1996/gator/internal/database"
)

type command struct {
	name	string
	args	[]string
}

type commands struct {
	handlers	map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("No handler for command name ", cmd.name)
	}

	err := handler(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.queries.GetUser(context.Background(), s.cfg.CurrentUser)
		if err != nil {
			return fmt.Errorf("Error getting user: %v", err)
		}

		err = handler(s, cmd, user)
		if err != nil {
			return err
		}

		return nil
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Login command requires a username!")
	}

	user, err := s.queries.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error getting user: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Current user set to %s\n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Register command requires a username!")
	}

	userParams := database.CreateUserParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		cmd.args[0],
	}

	user, err := s.queries.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("Error creating user: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User %s successfully created!\nID: %v\nCreated At: %v\nUpdated At: %v\nName: %s\n", 
		user.Name, user.ID, user.CreatedAt, user.UpdatedAt, user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.queries.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error deleting users: %v", err)
	}

	fmt.Println("Successfully deleted all users!")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.queries.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting users: %v", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUser {
			fmt.Println("*", user.Name, "(current)")
		} else {
			fmt.Println("*", user.Name)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Agg command requires a duration string!")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("addfeed command requires a name and a url!")
	}

	feedParams := database.CreateFeedParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		cmd.args[0],
		Url:		cmd.args[1],
		UserID:		user.ID,
	}

	feed, err := s.queries.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("Error adding feed to database: %v", err)
	}

	fmt.Printf("Feed %s successfully created!\nID: %v\nCreated At: %v\nUpdated At: %v\nName: %s\nURL: %s\nUser ID: %v\n", 
		feed.Name, feed.ID, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, feed.UserID)

	feedFollowParams := database.CreateFeedFollowParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		UserID:		user.ID,
		FeedID:		feed.ID,
	}

	feedFollowInfo, err := s.queries.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("Error creating feed follow: %v", err)
	}

	fmt.Printf("Feed name: %s, User name: %s\n", feedFollowInfo.FeedName, feedFollowInfo.UserName)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feedRows, err := s.queries.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting feeds: %v", err)
	}

	for _, feedRow := range feedRows {
		fmt.Printf("Feed Name: %s\nFeed URL: %s\nUser Name: %s\n\n",
			feedRow.FeedName, feedRow.Url, feedRow.UserName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Follow command requires a feed URL!")
	}

	feedID, err := s.queries.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error getting feed: %v", err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		UserID:		user.ID,
		FeedID:		feedID,
	}

	feedFollowInfo, err := s.queries.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("Error creating feed follow: %v", err)
	}

	fmt.Printf("Feed name: %s, User name: %s\n", feedFollowInfo.FeedName, feedFollowInfo.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.queries.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("Error getting users feeds: %v", err)
	}

	for _, feedFollow := range feedFollows {
		fmt.Println(feedFollow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Unfollow command requires a feed URL!")
	}

	feedID, err := s.queries.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error getting feed: %v", err)
	}

	deleteFeedFollowParams := database.DeleteFeedFollowForUserParams{
		UserID:		user.ID,
		FeedID:		feedID,
	}

	err = s.queries.DeleteFeedFollowForUser(context.Background(), deleteFeedFollowParams)
	if err != nil {
		return fmt.Errorf("Error deleting user's feed: %v", err)
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	if len(cmd.args) == 0 {
		limit = 2
	} else {
		num, err := strconv.ParseInt(cmd.args[0], 10, 32) 
		if err != nil {
			return fmt.Errorf("Error converting decimal string: %v\n", err)
		}
		limit = int32(num)
	}

	params := database.GetPostsForUserParams{
		UserID:		user.ID,
		Limit:		limit,
	}

	posts, err := s.queries.GetPostsForUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error getting user's posts: %v", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Printf("Description: %s\n", post.Description.String)
		fmt.Printf("Published at: %v\n\n", post.PublishedAt.Time)
	}

	return nil
}