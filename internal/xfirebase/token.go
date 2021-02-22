package xfirebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"log"
)

type AuthToken struct {
	*auth.Token
}

type Client interface {
	VerifyIDToken(ctx context.Context, idToken string) (*AuthToken, error)
}

func NewClient(app *firebase.App) Client {
	return &client{
		app: app,
	}
}

func NewDefaultClient() Client {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	return NewClient(app)
}

type client struct {
	app *firebase.App
}

func (c *client) VerifyIDToken(ctx context.Context, idToken string) (*AuthToken, error) {
	return verifyIDToken(ctx, c.app, idToken)
}

func verifyIDToken(ctx context.Context, app *firebase.App, idToken string) (*AuthToken, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Verified ID token: %v\n", token)
	return &AuthToken{Token: token}, nil
}
