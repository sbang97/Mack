package main

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// enter the access token into this variable
var pat = ""

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func main() {
	tokenSource := &TokenSource{
		AccessToken: pat,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	ctx := context.TODO()

	_, err := client.Droplets.Delete(ctx, 45333165)

	if err != nil {
		fmt.Printf("Something bad happened: %s", err.Error())
	} else {
		fmt.Printf("Deleted droplet successfully")
	}
}
