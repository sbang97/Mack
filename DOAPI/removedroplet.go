package main

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// enter your access token into this variable
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

	// enter the droplet id into this variable
	droplet := 123456
	// enter the load balancer id here for the second argument
	_, err := client.LoadBalancers.RemoveDroplets(ctx, "", droplet)

	if err != nil {
		fmt.Printf("Something bad happened: %s", err.Error())
	} else {
		fmt.Printf("Removed droplet successfully")
	}

}
