package handlers

import (
	"fmt"
	"net/http"
	"golang.org/x/net/html"
	"encoding/json"
	"strings"
	"net/url"
)

//openGraphPrefix is the prefix used for Open Graph meta properties
const openGraphPrefix = "og:"

//openGraphProps represents a map of open graph property names and values
type openGraphProps map[string]string

func getPageSummary(u string) (openGraphProps, error) {
	//Get the URL
	//If there was an error, return it
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("error getting url + %v\n" + err.Error())
	}
	//ensure that the response body stream is closed eventually
	//HINTS: https://gobyexample.com/defer
	//https://golang.org/pkg/net/http/#Response
	defer resp.Body.Close()
	//if the response StatusCode is >= 400
	//return an error, using the response's .Status
	//property as the error message
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf(resp.Status +"\n")
	}

	//if the response's Content-Type header does not
	//start with "text/html", return an error noting
	//what the content type was and that you were
	//expecting HTML
	/*ctype := resp.Header.Get("Content-Type");
	expectedContentType := "text/html; charset=utf-8"
	if ctype != expectedContentType {
		return nil, fmt.Errorf("Content-Type must be %s, was %s", expectedContentType, ctype)
	}
	*/

    ctype := resp.Header.Get("Content-Type");
	expectedContentType := "text/html; charset=utf-8"
	if strings.ToLower(ctype) != expectedContentType {
		return nil, fmt.Errorf("Content-Type must be %s, was %s", expectedContentType, ctype)
	}
	//create a new openGraphProps map instance to hold
	//the Open Graph properties you find
	//(see type definition above)
	m := make(openGraphProps)

	//tokenize the response body's HTML and extract
	//any Open Graph properties you find into the map,
	//using the Open Graph property name as the key, and the
	//corresponding content as the value.
	//strip the openGraphPrefix from the property name before
	//you add it as a new key, so that the key is just `title`
	//and not `og:title` (for example).
	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tokenType := tokenizer.Next()
        if tokenType == html.ErrorToken {
			link, _ := url.Parse(m["image"])
			link2, _ := url.Parse(u)
			if !link.IsAbs() {
				if (m["image"] != "") {
					m["image"] = link2.Scheme + "://" + link2.Hostname() + m["image"]
				} else {
					m["image"] = ""
				}
			}
			return m, fmt.Errorf("reached end of file")
		}
		token := tokenizer.Token()
			if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
				if token.Data == "meta"  {
					parseAttributes(m, token.Attr)
				}
		}
	}
}

//SummaryHandler fetches the URL in the `url` query string parameter, extracts
//summary information about the returned page and sends those summary properties
//to the client as a JSON-encoded object.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	//Add the following header to the response
	//Access-Control-Allow-Origin: *
	//this will allow JavaScript served from other origins
	//to call this API
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//get the `url` query string parameter
	//if you use r.FormValue() it will also handle cases where
	//the client did POST with `url` as a form field
	//HINT: https://golang.org/pkg/net/http/#Request.FormValue
	url := r.FormValue(`url`)
	//if no `url` parameter was provided, respond with
	//an http.StatusBadRequest error and return
	//HINT: https://golang.org/pkg/net/http/#Error
	if url == "" {
		http.Error(w, "the url is required", http.StatusBadRequest)
		return
	}

	//call getPageSummary() passing the requested URL
	//and holding on to the returned openGraphProps map
	//(see type definition above)
	m, err := getPageSummary(url)
	//if you get back an error, respond to the client
	//with that error and an http.StatusBadRequest code
	if err != nil && err.Error() != "reached end of file" {
		http.Error(w, "error tokenizing HTML", http.StatusBadRequest)
		return
	}
	//otherwise, respond by writing the openGrahProps
	//map as a JSON-encoded object
	//add the following headers to the response before
	//you write the JSON-encoded object:
	//   Content-Type: application/json; charset=utf-8
	//this tells the client that you are sending it JSON
	encoder := json.NewEncoder(w)
 	if err := encoder.Encode(m); err != nil {
		http.Error(w, "error encoding json: " + err.Error(),
		http.StatusInternalServerError)
	}
}

//given a map and a slice of attributes, parses the attributes
//and adds any open graph properties to the map.
func parseAttributes(m openGraphProps, attr []html.Attribute) {
	for _, b := range attr {
		if b.Key == "property" {
			if strings.HasPrefix(b.Val, openGraphPrefix) {
				//strips the "og:" off of the property before adding it as a key to the map
				split := strings.Split(b.Val, openGraphPrefix)
				key := split[1]
				//finds the corrresponding "content" attribute to add as a value for the key
				for i := 0; i < len(attr); i++ {
					attribute := attr[i]
					if attribute.Key == "content" {
						m[key] = attribute.Val
					}
				}
			}
		}
	}
}