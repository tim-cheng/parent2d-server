package models

import (
	"fmt"
	"github.com/coopernurse/gorp"
	"time"
)

type User struct {
	Id          int64     `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	Type        string    `db:"type"`
	FbId        string    `db:"fb_id"`
	Admin       bool      `db:"admin"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	NumDegree1  int64     `db:"num_degree1"`
	NumDegree2  int64     `db:"num_degree2"`
	Description string    `db:"description"`
	Picture     []byte    `db:"picture"`
	Location    string    `db:"location"`
	Zip         string    `db:"zip"`
	Interests   string    `db:"interests"`
}

// Validation Hooks
func (u *User) PreInsert(s gorp.SqlExecutor) error {
	u.CreatedAt = time.Now()
	return nil
}

// helper
func (db *MyDb) genSaltedHash(password string) string {
	hashQuery := fmt.Sprintf("select crypt('%s', gen_salt('md5'))", password)
	var hash string
	db.SelectOne(&hash, hashQuery)
	return hash
}

func (db *MyDb) NewUser(typ, email, password, firstName, lastName, fb_id, location, zip string) (*User, error) {

	hash := db.genSaltedHash(password)

	return &User{
		Type:        typ,
		Email:       email,
		Password:    hash,
		Admin:       false,
		FirstName:   firstName,
		LastName:    lastName,
		Description: "",
		FbId:        fb_id,
		Location:    location,
		Zip:         zip,
	}, nil
}

func (db *MyDb) GetUser(id interface{}) (*User, error) {
	var user User
	err := db.SelectOne(&user, "select id, created_at, email, type, first_name, last_name, num_degree1, num_degree2, description, location, zip, interests from users where id=$1", id)
	return &user, err
}

func (db *MyDb) GetUserName(id interface{}) (*User, error) {
	var user User
	err := db.SelectOne(&user, "select id, first_name, last_name from users where id=$1", id)
	return &user, err
}

func (db *MyDb) GetUserNameByEmail(email string) (*User, error) {
	var user User
	err := db.SelectOne(&user, "select id, first_name, last_name from users where lower(email)=$1", email)
	return &user, err
}

func (db *MyDb) UpdateUserDesc(userId, desc string) error {
	_, err := db.Exec("update users set description=$2 where id=$1", userId, desc)
	return err
}

func (db *MyDb) Update1dConnection(userId int64, nConn int) error {
	_, err := db.Exec("update users set num_degree1=$2 where id=$1", userId, nConn)
	return err
}

func (db *MyDb) Update2dConnection(userId int64, nConn int) error {
	_, err := db.Exec("update users set num_degree2=$2 where id=$1", userId, nConn)
	return err
}

func (db *MyDb) UpdateFirstName(userId int64, name string) error {
	_, err := db.Exec("update users set first_name=$2 where id=$1", userId, name)
	return err
}

func (db *MyDb) UpdateLastName(userId int64, name string) error {
	_, err := db.Exec("update users set last_name=$2 where id=$1", userId, name)
	return err
}

func (db *MyDb) UpdateInterests(userId int64, val string) error {
	_, err := db.Exec("update users set interests=$2 where id=$1", userId, val)
	return err
}

func (db *MyDb) UpdateZip(userId int64, val string) error {
	_, err := db.Exec("update users set zip=$2 where id=$1", userId, val)
	return err
}

func (db *MyDb) UpdateLocation(userId int64, val string) error {
	_, err := db.Exec("update users set location=$2 where id=$1", userId, val)
	return err
}

func (db *MyDb) PostUser(user *User) error {
	err := db.Insert(user)
	return err
}

func (db *MyDb) UpdatePassword(user *User, password string) error {
	hash := db.genSaltedHash(password)
	_, err := db.Exec("update users set password=$1 where id=$2", hash, user.Id)
	return err
}

func (db *MyDb) PostUserPicture(userId int64, image []byte) error {
	u := new(User)
	err := db.SelectOne(u, "select * from users where id=$1", userId)
	u.Picture = image
	_, err = db.Update(u)
	return err
}

func (db *MyDb) GetUserPicture(userId int64) ([]byte, error) {
	u := new(User)
	err := db.SelectOne(u, "select picture from users where id=$1", userId)
	return u.Picture, err
}

func (db *MyDb) GetPassword(email string) (string, error) {
	p := ""
	err := db.SelectOne(&p, "select password from users where lower(email)=$1", email)
	return p, err
}

func (db *MyDb) GetUsersByFirstName(name string) ([]User, error) {
	var users []User
	_, err := db.Select(&users, "select id, first_name, last_name from users where first_name ilike $1", name)
	return users, err
}

func (db *MyDb) GetUsersByLastName(name string) ([]User, error) {
	var users []User
	_, err := db.Select(&users, "select id, first_name, last_name from users where last_name ilike $1", name)
	return users, err
}

func (db *MyDb) GetUsersByFullName(firstName, lastName string) ([]User, error) {
	var users []User
	_, err := db.Select(&users, "select id, first_name, last_name, location from users where first_name ilike $1 and last_name ilike $2", firstName, lastName)
	return users, err
}

func (db *MyDb) GetUsersByFbIdList(fbIds []string) ([]User, error) {
	var users []User
	cond := "fb_id = '1'"
	for _, v := range fbIds {
		cond += " or fb_id = '" + v + "'"
	}
	_, err := db.Select(&users, "select id, first_name, last_name, location from users where "+cond)
	return users, err
}
