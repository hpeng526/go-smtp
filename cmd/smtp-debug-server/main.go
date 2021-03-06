package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	"github.com/emersion/go-smtp"
)

var addr = "127.0.0.1:1025"

func init() {
	flag.StringVar(&addr, "l", addr, "Listen address")
}

type backend struct{}

func (bkd *backend) Login(ctx context.Context, state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return &session{}, nil
}

func (bkd *backend) AnonymousLogin(ctx context.Context, state *smtp.ConnectionState) (smtp.Session, error) {
	return &session{}, nil
}

type session struct{}

func (s *session) Mail(ctx context.Context, from string, opts smtp.MailOptions) error {
	return nil
}

func (s *session) Rcpt(ctx context.Context, to string) error {
	return nil
}

func (s *session) Data(ctx context.Context, r io.Reader) error {
	return nil
}

func (s *session) Reset(ctx context.Context) {}

func (s *session) Logout(ctx context.Context) error {
	return nil
}

func main() {
	flag.Parse()

	s := smtp.NewServer(&backend{})

	s.Addr = addr
	s.Domain = "localhost"
	s.AllowInsecureAuth = true
	s.Debug = os.Stdout

	log.Println("Starting SMTP server at", addr)
	log.Fatal(s.ListenAndServe())
}
