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

	//
	createRequest := &godo.DropletMultiCreateRequest{
		Names:   []string{"apiserver","webclient", "deleteme"},
		Region: "sfo1",
		Size:   "512mb",
		Image: godo.DropletCreateImage{
		Slug: "ubuntu-16-04-x64",
		},
		IPv6: true,
		Tags: []string{"web"},
	}

	droplets, _, err := client.Droplets.CreateMultiple(ctx, createRequest)
	for num := range droplets {
		fmt.Printf("success with droplet# %v", num "\n")
	}
	if err != nil {
		fmt.Printf("Something bad happened: %s", err.Error())
	}
}

