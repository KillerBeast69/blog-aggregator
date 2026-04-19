package main

import (
	"fmt"
	"log"
	"os"

	"github.com/KillerBeast69/blog-aggregator/internal/config"
)

type state struct {
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

	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}

	fmt.Printf("user has been set to: %s\n", username)

	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
	fmt.Printf("initial config: %+v\n", cfg)

	program_state := &state{
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

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
