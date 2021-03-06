package backendutil

import (
	"context"
	"io"

	"github.com/emersion/go-smtp"
)

// TransformBackend is a backend that transforms messages.
type TransformBackend struct {
	Backend smtp.Backend

	TransformMail func(from string) (string, error)
	TransformRcpt func(to string) (string, error)
	TransformData func(r io.Reader) (io.Reader, error)
}

// Login implements the smtp.Backend interface.
func (be *TransformBackend) Login(ctx context.Context, state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	s, err := be.Backend.Login(ctx, state, username, password)
	if err != nil {
		return nil, err
	}
	return &transformSession{s, be}, nil
}

// AnonymousLogin implements the smtp.Backend interface.
func (be *TransformBackend) AnonymousLogin(ctx context.Context, state *smtp.ConnectionState) (smtp.Session, error) {
	s, err := be.Backend.AnonymousLogin(ctx, state)
	if err != nil {
		return nil, err
	}
	return &transformSession{s, be}, nil
}

type transformSession struct {
	Session smtp.Session

	be *TransformBackend
}

func (s *transformSession) Reset(ctx context.Context) {
	s.Session.Reset(ctx)
}

func (s *transformSession) Mail(ctx context.Context, from string, opts smtp.MailOptions) error {
	if s.be.TransformMail != nil {
		var err error
		from, err = s.be.TransformMail(from)
		if err != nil {
			return err
		}
	}
	return s.Session.Mail(ctx, from, opts)
}

func (s *transformSession) Rcpt(ctx context.Context, to string) error {
	if s.be.TransformRcpt != nil {
		var err error
		to, err = s.be.TransformRcpt(to)
		if err != nil {
			return err
		}
	}
	return s.Session.Rcpt(ctx, to)
}

func (s *transformSession) Data(ctx context.Context, r io.Reader) error {
	if s.be.TransformData != nil {
		var err error
		r, err = s.be.TransformData(r)
		if err != nil {
			return err
		}
	}
	return s.Session.Data(ctx, r)
}

func (s *transformSession) Logout(ctx context.Context) error {
	return s.Session.Logout(ctx)
}
