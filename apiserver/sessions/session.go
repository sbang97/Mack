package sessions

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("scheme used in Authorization header is not supported")

//BeginSession creates a new session ID, saves the state to the store, adds a
//header to the response with the session ID, and returns the new session ID
func BeginSession(signingKey string, store Store, state interface{}, w http.ResponseWriter) (SessionID, error) {
	//create a new SessionID
	//if you get an error, return InvalidSessionID and the error
	sessID, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	//save the state to the store
	//if you get an error, return InvalidSessionID and the error
	err = store.Save(sessID, state)
	if err != nil {
		return InvalidSessionID, err
	}
	//Add a response header like this:
	//  Authorization: Bearer <sid>
	//where <sid> is the new SessionID
	w.Header().Add(headerAuthorization, schemeBearer+string(sessID))
	//return the new SessionID and nil
	return sessID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	//get the value of the Authorization header
	auth := r.Header.Get("Authorization")
	//if it's zero-length, return InvalidSessionID and ErrNoSessionID
	if len(auth) == 0 {
		log.Println("Auth Header missing falling back to Form Paramater")
		auth = r.FormValue("auth")
	}
	if len(auth) == 0 {
		log.Println("No Bearer token found")
		return InvalidSessionID, ErrNoSessionID
	}
	//if it doesn't start with "Bearer ",
	//return InvalidSessionID and ErrInvalidScheme
	if !strings.HasPrefix(auth, schemeBearer) {
		return InvalidSessionID, ErrInvalidScheme
	}
	//trim off the "Bearer " prefix and validate the remaining id
	//if you get an error return InvalidSessionID and the error
	auth = strings.TrimPrefix(auth, schemeBearer)
	validatedID, err := ValidateID(auth, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	//return the validated SessionID and nil
	return validatedID, nil
}

//GetState extracts the SessionID from the request,
//and gets the associated state from the provided store
func GetState(r *http.Request, signingKey string, store Store, state interface{}) (SessionID, error) {
	//get the SessionID from the request
	//if you get an error, return the SessionID and error
	sessID, err := GetSessionID(r, signingKey)
	if err != nil {
		return sessID, err
	}
	//get the associated state data from the provided store
	//if you get an error return the SessionID and the error
	err = store.Get(sessID, state)
	if err != nil {
		return sessID, err
	}
	//return the SessionID and nil
	return sessID, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	//get the SessionID from the request
	//if you get an error return the SessionID and error
	sessID, err := GetSessionID(r, signingKey)
	if err != nil {
		return sessID, err
	}
	//delete the associated data in the provided store
	err = store.Delete(sessID)
	if err != nil {
		return sessID, err
	}
	//return the SessionID and nil
	return sessID, nil
}
