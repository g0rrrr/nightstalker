package models

import (
    crypto_rand "crypto/rand"
    "math"
    "math/rand"
    _ "strings"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"io"
	"strconv"
	"time"
)

type User struct {
    Id            int64          `db:"id"`
	GroupId       int64          `db:"group_id"`
	CreatedOn     time.Time      `db:"created_on"`
	Username      string         `db:"username"`
	Password      string         `db:"password"`
    Experience    int64          `db:"experience"`
    Level         int64          `db:"level"`
	Avatar        string         `db:"avatar"`
    Reputations   int64          `db:"reputations"`
	Signature     sql.NullString `db:"signature"`
	Salt          string         `db:"salt"`
	StylesheetUrl sql.NullString `db:"stylesheet_url"`
	UserTitle     string         `db:"user_title"`
	LastSeen      time.Time      `db:"last_seen"`
	HideOnline    bool           `db:"hide_online"`
	LastUnreadAll pq.NullTime    `db:"last_unread_all"`
}

type Notifications struct {
    Id            int64         `db:"id"`
    UserId        int64         `db:"user_id"`
    NotifUserId   int64         `db:"notif_user_id"`
    Read          bool          `db:"read"`
	Author        *User         `db:"-"`
    Message       string        `db:"message"`
	CreatedOn     time.Time     `db:"created_on"`
}

func NewUser(username, password string) *User {
	user := &User{
		CreatedOn: time.Now(),
		Username:  username,
		LastSeen:  time.Now(),
	}

	user.SetPassword(password)
	return user
}

func AuthenticateUser(username, password string) (error, *User) {
	db := GetDbSession()
	user := &User{}
	err := db.SelectOne(user, "SELECT * FROM users WHERE username=$1", username)
	if err != nil {
		fmt.Printf("[error] Cannot select user (%s)\n", err.Error())
		return err, nil
	}

	if user.Id == 0 {
		return errors.New("Invalid username/password"), nil
	}

	hasher := sha1.New()
	io.WriteString(hasher, password)
	io.WriteString(hasher, user.Salt)
	password = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if password != user.Password {
		return errors.New("Invalid username/password"), nil
	}

	// Update the user's last seen
	user.LastSeen = time.Now()
	db.Update(user)

	return nil, user
}

/*func AuthenticateUser(username, password string) (error, *User) {
	db := GetDbSession()
    if len(username) > 10 || len(username) < 4 {
        return errors.New("Username should have 4 to 10 characters"), nil
    }
    username = strings.ToLower(username)
	user := &User{}
	err := db.SelectOne(user, "SELECT * FROM users WHERE username=$1", username)
	if err != nil {
		fmt.Printf("[error] Cannot select user (%s)\n", err.Error())
		return err, nil
	}

	if user.Id == 0   {
		return errors.New("Invalid username/password"), nil
	}

	hasher := sha1.New()
	io.WriteString(hasher, password)
	io.WriteString(hasher, user.Salt)
	password = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	if password != user.Password {
		return errors.New("Invalid username/password"), nil
	}

	// Update the user's last seen
	user.LastSeen = time.Now()
	db.Update(user)

	return nil, user
} */

func GetUserCount() (int64, error) {
	db := GetDbSession()

	count, err := db.SelectInt("SELECT COUNT(*) FROM users")
	if err != nil {
		fmt.Printf("[error] Error selecting user count (%s)\n", err.Error())
		return 0, errors.New("Database error: " + err.Error())
	}

	return count, nil
}

func GetLatestUser() (*User, error) {
	db := GetDbSession()

	user := &User{}
	err := db.SelectOne(user, "SELECT * FROM users ORDER BY created_on DESC LIMIT 1")

	if err != nil {
		fmt.Printf("[error] Error selecting latest user (%s)\n", err.Error())
		return nil, errors.New("Database error: " + err.Error())
	}

	if user.Username == "" {
		return nil, nil
	}

	return user, nil
}

func GetOnlineUsers() (users []*User) {
	db := GetDbSession()

	db.Select(&users, "SELECT * FROM users WHERE last_seen > current_timestamp - interval '5 seconds' AND hide_online != true")
	return users
}

func GetUser(id int) (*User, error) {
	db := GetDbSession()
	obj, err := db.Get(&User{}, id)
	if obj == nil {
		return nil, err
	}

	return obj.(*User), err
}

// Converts the given string into an appropriate hash, resets the salt,
// and sets the Password attribute. Does *not* commit to the database.
func (user *User) SetPassword(password string) {
	var int_salt int32
	binary.Read(crypto_rand.Reader, binary.LittleEndian, &int_salt)
	salt := strconv.Itoa(int(int_salt))

	hasher := sha1.New()
	io.WriteString(hasher, password)
	io.WriteString(hasher, salt)
	user.Password = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	user.Salt = salt
}

func (user *User) IsAdmin() bool {
	if user.GroupId == 2 {
		return true
	}

	return false
}

func (user *User) CanModerate() bool {
	if user.GroupId > 0 {
		return true
	}

	return false
}

/*
first level color : 87ff00
*/

func ConvertXpToLevel(xp int64) int64 {
	return int64(0.07 * math.Sqrt(float64(xp)))
}

func ConvertLevelToXp(level int64) int64 {
  return int64(math.Pow(float64(level) / 0.07, 2))
}

func GetLevelProgress(user_id int64) (int64, string) {
	db := GetDbSession()
	level, _ := db.SelectInt("SELECT level FROM users WHERE id=$1", user_id)

	current_xp, _ := db.SelectInt("SELECT experience FROM users WHERE id=$1", user_id)

	currentLevelXP := ConvertLevelToXp(level)
	nextLevelXP := ConvertLevelToXp(level + 1)

	neededXP := nextLevelXP - currentLevelXP
	earnedXP := current_xp - currentLevelXP

	fmt.Println("neededXp: ", neededXP, "earnedXp: ", earnedXP, "next level xp: ", nextLevelXP, "current xp: ", current_xp)

	levelProgress := int64(100 - int(math.Ceil(float64(neededXP-earnedXP) / float64(neededXP) * 100)))

	// Return both level progress and neededXP
	return levelProgress, fmt.Sprintf("%dXP/%dXP", earnedXP, neededXP)
}

func Experience(user_id int64, experience_type int64) {
    db := GetDbSession()

    experience_amount := 0
    switch experience_type { 
    case 1: // replying to a thread
        experience_amount = 50
    case 2: // getting reply to your thread
        experience_amount = 100
    case 3: // post a thread 
        experience_amount = 100
    case 4: // getting a like 
        experience_amount = 150
    case 5: // someone followed you 
        experience_amount = 200
    case 6: // someone unfollowed you 
        experience_amount = -200
    case 7: // random xp after getting new level
      experience_amount = rand.Intn(100)
    }

  /*Level	        XP	                Difference
    0	            0	                  0
    1	            204.0816327	        204
    2	            816.3265306	        612
    3	            1836.734694	        1020
    4	            3265.306122	        1429
    5	            5102.040816	        1837
    6	            7346.938776	        2245
    7	            10000	              2653
    8	            13061.22449	        3061
    9	            16530.61224	        3469
    10	          20408.16327         3878
  */

  /* 100 - (current_level + 1).difference
     x - current_level
     x = current_level * 100 / level different 
  */


  db.Exec("UPDATE users SET experience = experience + $1 WHERE id=$2", experience_amount, user_id)
  current_xp, _ := db.SelectInt("SELECT experience FROM users WHERE id=$1", user_id)

  db.Exec("UPDATE users SET level=$1 WHERE id=$2", ConvertXpToLevel(current_xp), user_id)
}

func (user *User) GetPostCount() int64 {
	db := GetDbSession()
	count, err := db.SelectInt("SELECT COUNT(*) FROM posts WHERE author_id=$1", user.Id)

	if err != nil {
		return 0
	}

	return count
}

func (user *User) GetRepCount() int64 {
  db := GetDbSession()
  count, err := db.SelectInt("SELECT reputations FROM users WHERE id=$1", user.Id)

  if err != nil {
    return 0
  }

  return count
}

func (user *User) AlreadyFollowing(user_id int64) bool {
  db := GetDbSession()
  count, _ := db.SelectInt("SELECT COUNT(1) FROM followers WHERE follower_id=$1 AND followed_id=$2", user.Id, user_id)

  if count > 0 {
    return true
  }
  return false
}

func (user *User) FollowUser(followed_user int64) error {
    db := GetDbSession()
    count, err := db.SelectInt("SELECT COUNT(1) FROM followers WHERE follower_id=$1 AND followed_id=$2", user.Id, followed_user)

    if err != nil {
      return errors.New("Cant get the count of user followers")
    }

    fmt.Println(user.Id, " : ", followed_user)
    if user.Id == followed_user {
      return nil
    }

    if count > 0 {
      db.Exec("DELETE FROM followers WHERE follower_id=$1 AND followed_id=$2", user.Id, followed_user)
    Experience(followed_user, 50)
    return nil
  }

  Experience(followed_user, 5)
  //current_xp, _ := db.SelectInt("SELECT experience FROM users WHERE id=$1", followed_user)
  //user.Progress = GetLevelProgress(ConvertXpToLevel(current_xp), current_xp)
  _, err = db.Exec("INSERT INTO followers (follower_id, followed_id) VALUES ($1, $2)", user.Id, followed_user)
  return nil
}

func (user *User) GetFollowersCount() int64 {
    db := GetDbSession() 
    count, err := db.SelectInt("SELECT COUNT(*) FROM followers WHERE followed_id=$1", user.Id)

    if err != nil {
      return 0
    }

    return count
}

func (user *User) GetFollowingsCount() int64 {
    db := GetDbSession() 
    count, err := db.SelectInt("SELECT COUNT(*) FROM followers WHERE follower_id=$1", user.Id)

    if err != nil {
      return 0
    }

    return count
}

func (user *User) GetMessagesCount() int64 {
    db := GetDbSession()
    count, err := db.SelectInt("SELECT COUNT(*) FROM messages WHERE user_id=$1 AND read=FALSE", user.Id)

    if err !=nil {
        return 0;
    }

    return count
}

func (user *User) GetNotificationsCount() int64 {
    db := GetDbSession() 
    count, err := db.SelectInt("SELECT COUNT(*) FROM notifications WHERE user_id=$1 AND read=FALSE", user.Id)

    if err != nil {
      return 0
    }

    return count
}

func (user *User) GetNotifications() ([]*Notifications, error) { 
  db := GetDbSession()
  var notifications []*Notifications

  _, err := db.Select(&notifications, "SELECT * FROM notifications WHERE user_id=$1 ORDER BY created_on DESC",  user.Id)

  if err != nil {
		fmt.Printf("[error] Could not get user's followers (%s)", err.Error())
  }

	for i := range notifications {
		obj, _ := db.Get(&User{}, notifications[i].NotifUserId)
		user := obj.(*User)

		notifications[i].Author = user
	}

  return notifications, err
}

func (user *User) GetPosts(page int) []*Post {
	db := GetDbSession()
	var posts []*Post

	offset := POSTS_PER_PAGE * int64(page)

	_, err := db.Select(&posts, "SELECT * FROM posts WHERE author_id=$1 ORDER BY created_on DESC LIMIT $2 OFFSET $3", user.Id, POSTS_PER_PAGE, offset)

	if err != nil {
		fmt.Printf("[error] Could not get user's posts (%s)", err.Error())
	}

	return posts
}

