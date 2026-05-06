package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/KillerBeast69/blog-aggregator/internal/config"
	"github.com/KillerBeast69/blog-aggregator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handlerFunc, exists := c.registeredCommands[cmd.name]
	if !exists {
		return fmt.Errorf("command not found: %s", cmd.name)
	}

	return handlerFunc(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}

	username := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("could not find user: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}

	fmt.Printf("user has been set to: %s\n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the register handler expects a single argument, the username")
	}
	username := cmd.args[0]
	// how do I create a new user in the database??

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to set user in config: %v", err)
	}

	fmt.Println("user has been created successfully:")
	fmt.Printf("%+v\n", user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset: %v", err)
	}

	fmt.Println("Database reset successfull")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	// how do I find out if the user is logged in?
	if err != nil {
		return fmt.Errorf("failed to list all users: %v", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("the addfeed handler expects two arguments, <the name of the feed> and <the url of the feed>")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	current_usr_name := s.cfg.CurrentUserName

	current_urs, err := s.db.GetUser(context.Background(), current_usr_name)
	if err != nil {
		return fmt.Errorf("failed to get user_id: %v", err)
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    current_urs.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed getting feed records: %v", err)
	}

	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed getting feed: %v", err)
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("User Name: %s\n", feed.UserName)
		fmt.Println("-------------------------------------")
	}
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
	//fmt.Printf("initial config: %+v\n", cfg)

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	dbQueries := database.New(db)

	program_state := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)

	args := os.Args

	if len(args) < 2 {
		log.Fatal("usage: gator <command> [args...]")
	}

	cmd_name := args[1]
	cmd_args := args[2:]

	cmd := command{
		name: cmd_name,
		args: cmd_args,
	}

	err = cmds.run(program_state, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
