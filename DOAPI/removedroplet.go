package main

import (
    "context" 
    "golang.org/x/oauth2" 
    "github.com/digitalocean/godo" 
	"fmt"
)

var pat = "c2919e02dbcd8841abd96383356ea662207004c9c98a1cf9e2ece0d384caf4c7" 

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

	droplet := 45333165
	_, err := client.LoadBalancers.RemoveDroplets(ctx, "60904d58-0764-4b13-b52a-9d030fd2c96f", droplet)

	if err != nil {
		fmt.Printf("Something bad happened: %s", err.Error())
	} else {
		fmt.Printf("Removed droplet successfully")
	}
	
}