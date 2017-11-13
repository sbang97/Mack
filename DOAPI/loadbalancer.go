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

	createRequest := &godo.LoadBalancerRequest{
		Name:      "lb",
		Algorithm: "round_robin",
		Region:    "sfo1",
		ForwardingRules: []godo.ForwardingRule{
			{
				EntryProtocol:  "http",
				EntryPort:      80,
				TargetProtocol: "http",
				TargetPort:     80,
			},
			{
				EntryProtocol:  "https",
				EntryPort:      443,
				TargetProtocol: "https",
				TargetPort:     443,
				TlsPassthrough: true,
			},
		},
		HealthCheck: &godo.HealthCheck{
			Protocol:               "http",
			Port:                   80,
			Path:                   "/",
			CheckIntervalSeconds:   10,
			ResponseTimeoutSeconds: 5,
			HealthyThreshold:       5,
			UnhealthyThreshold:     3,
		},
		StickySessions: &godo.StickySessions{
			Type: "none",
		},
		// Enter the droplet ids youd like to add to the load balancer here
		DropletIDs:          []int{},
		RedirectHttpToHttps: false,
	}

	lb, _, err := client.LoadBalancers.Create(ctx, createRequest)

	if lb != nil {
		fmt.Print("successfully created load balancer")
	}

	if err != nil {
		fmt.Printf("Something bad happened: %s", err.Error())
	}

}
