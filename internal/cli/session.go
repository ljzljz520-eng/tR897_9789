package cli

import (
	"fmt"
	"strings"
)

type Session struct {
	Actor string
	Store string
}

func NewSession(actor, store string) Session {
	if strings.TrimSpace(store) == "" {
		store = "troubleshooter.db"
	}
	return Session{Actor: strings.ToLower(strings.TrimSpace(actor)), Store: store}
}

func (s Session) Valid() bool {
	return s.Actor != "" && s.Store != ""
}

func (s Session) Header() string {
	return fmt.Sprintf("actor=%s store=%s", s.Actor, s.Store)
}

func (s Session) Can(command Command) bool {
	if !s.Valid() {
		return false
	}
	if command.Name == "diagnose" {
		return command.Employee != ""
	}
	return command.Name == "history" || command.Name == "health"
}

func ParseSession(args []string) (Session, Command, error) {
	command, err := Parse(args)
	if err != nil {
		return Session{}, Command{}, err
	}
	session := NewSession(command.Actor, command.Store)
	if command.Name != "diagnose" {
		session = NewSession("operator", command.Store)
	}
	if !session.Can(command) {
		return Session{}, Command{}, fmt.Errorf("session cannot execute command")
	}
	return session, command, nil
}

func NormalizeCommandName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
