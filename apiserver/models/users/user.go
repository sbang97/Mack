package users

import (
	"crypto/md5"
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar profile photos
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

const defaultCost = 13

//UserID defines the type for user IDs
type UserID string

//User represents a user account in the database
type User struct {
	ID        UserID `json:"id" bson:"_id"`
	Email     string `json:"email"`
	PassHash  []byte `json:"-" bson:"passHash"` //stored in mongo, but never encoded to clients
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//UserUpdates represents updates one can make to a user
type UserUpdates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user
func (nu *NewUser) Validate() error {
	//ensure Email field is a valid Email
	//HINT: use mail.ParseAddress()
	//https://golang.org/pkg/net/mail/#ParseAddress
	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return fmt.Errorf("enter a valid email address")
	}
	//ensure Password is at least 6 chars
	if len(nu.Password) < 6 {
		return fmt.Errorf("password must be atleast 6 characters long")
	}
	//ensure Password and PasswordConf match
	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("passwords do not match")
	}
	//ensure UserName has non-zero length
	if len(nu.UserName) == 0 {
		return fmt.Errorf("username cannot be empty")
	}
	//if you made here, it's valid, so return nil
	return nil
}

//ToUser converts the NewUser to a User
func (nu *NewUser) ToUser() (*User, error) {
	//build the Gravatar photo URL by creating an MD5
	//hash of the new user's email address, converting
	//that to a hex string, and appending it to their base URL:
	//https://www.gravatar.com/avatar/ + hex-encoded md5 has of email
	hash := md5.New()
	hash.Write([]byte(nu.Email))
	photoURL := gravatarBasePhotoURL + string(hash.Sum(nil))
	//construct a new User setting the various fields
	//but don't assign a new ID here--do that in your
	//concrete Store.Insert() method
	user := &User{
		Email:     nu.Email,
		PhotoURL:  photoURL,
		UserName:  nu.UserName,
		FirstName: nu.FirstName,
		LastName:  nu.LastName,
	}
	//call the User's SetPassword() method to set the password,
	//which will hash the plaintext password
	err := user.SetPassword(nu.Password)
	if err != nil {
		return nil, fmt.Errorf("error setting password")
	}
	//return the User and nil
	return user, nil
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	//hash the plaintext password using an adaptive
	//crytographic hashing algorithm like bcrypt
	//https://godoc.org/golang.org/x/crypto/bcrypt
	passhash, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}
	//set the User's PassHash field to the resulting hash
	u.PassHash = passhash
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	//compare the plaintext password with the PassHash field
	//using the same hashing algorithm you used in SetPassword
	err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
	if err != nil {
		return fmt.Errorf("Invalid Password")
	}
	return nil
}
