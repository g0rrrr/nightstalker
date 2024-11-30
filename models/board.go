package models

import (
	"fmt"
	"github.com/lib/pq"
	"math"
	"time"
)

const THREADS_PER_PAGE = 30

type Board struct {
	Id          int64  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Order       int    `db:"ordering"`
}

type BoardLatest struct {
	Op     *Post
	Latest *Post
}

type JoinBoardView struct {
	Board       *Board      `db:"-"`
	Id          int64       `db:"id"`
	Title       string      `db:"title"`
	Description string      `db:"description"`
	Order       int         `db:"ordering"`
	ViewedOn    pq.NullTime `db:"viewed_on"`
  ViewsCount  int         `db:"view_count"`
}

type JoinThreadView struct {
	Thread      *Post       `db:"-"`
	Id          int64       `db:"id"`
	BoardId     int64       `db:"board_id"`
  BoardTitle  string      `db:"board_title"`
  Views       int64       `db:"views"`
	Author      *User       `db:"-"`
	AuthorId    int64       `db:"author_id"`
	Title       string      `db:"title"`
  Reputations int64       `db:"reputations"`
	CreatedOn   time.Time   `db:"created_on"`
	LatestReply time.Time   `db:"latest_reply"`
	Sticky      bool        `db:"sticky"`
	Locked      bool        `db:"locked"`
	ViewedOn    pq.NullTime `db:"viewed_on"`
}

func NewBoard(title, desc string, order int) *Board {
	return &Board{
		Title:       title,
		Description: desc,
		Order:       order,
	}
}

func UpdateBoard(title, desc string, order int, id int64) *Board {
	return &Board{
		Title:       title,
		Description: desc,
		Order:       order,
		Id:          id,
	}
}

func GetBoard(id int) (*Board, error) {
	db := GetDbSession()
	obj, err := db.Get(&Board{}, id)
	if obj == nil {
		return nil, err
	}

	return obj.(*Board), err
}

func GetBoards() ([]*Board, error) {
	db := GetDbSession()

	var boards []*Board
	_, err := db.Select(&boards, "SELECT * FROM boards ORDER BY ordering ASC")

	return boards, err
}

func GetBoardsUnread(user *User) ([]*JoinBoardView, error) {
	db := GetDbSession()

	user_id := int64(-1)
	if user != nil {
		user_id = user.Id
	}

	var boards []*JoinBoardView
	_, err := db.Select(&boards, `
        SELECT
            boards.*,
            views.time AS viewed_on
        FROM boards
        LEFT OUTER JOIN views ON
            views.post_id=(SELECT id FROM posts WHERE board_id=boards.id AND parent_id IS NULL ORDER BY latest_reply DESC LIMIT 1) AND
            views.user_id=$1
        ORDER BY
            ordering ASC
    `, user_id)

	for i := range boards {
		if user_id == -1 {
			boards[i].ViewedOn = pq.NullTime{Time: time.Now(), Valid: true}
		}

		boards[i].Board = &Board{
			Id: boards[i].Id,
		}
	}
	return boards, err
}

func (board *Board) GetLatestPost() BoardLatest {
	db := GetDbSession()
	op := &Post{}
	latest := &Post{}

	err := db.SelectOne(op, "SELECT * FROM posts WHERE board_id=$1 AND parent_id IS NULL ORDER BY latest_reply DESC LIMIT 1", board.Id)

	if err != nil {
		fmt.Printf("[error] Could not get latest post in board %d (%s)\n", board.Id, err.Error())
	}

	err = db.SelectOne(latest, "SELECT * FROM posts WHERE board_id=$1 AND parent_id=$2 ORDER BY created_on DESC LIMIT 1", board.Id, op.Id)

	if latest.Author == nil {
		latest = nil
	}


	return BoardLatest{
		Op:     op,
		Latest: latest,
	}
}

func (board *Board) GetThreads(sort_by string, show_all bool, page int, user *User) ([]*JoinThreadView) {
	db := GetDbSession()

	db_query := `
    SELECT 
        posts.id, 
        posts.author_id,
        posts.title,
        posts.created_on,
        posts.latest_reply,
        posts.sticky,
        posts.locked,
        posts.board_id,
        posts.views,
        boards.title AS board_title,
        views.time AS viewed_on 
    FROM posts
    LEFT OUTER JOIN views ON 
        posts.id = views.post_id AND
        views.user_id = $1
    JOIN boards ON
        posts.board_id = boards.id
    `

    THREADS_PER_PAGE := 8
    THREADS_PER_PAGE_on_boards := 35

	var threads []*JoinThreadView
	i_begin := page * (THREADS_PER_PAGE - 1)

	user_id := int64(-1)
	if user != nil {
		user_id = user.Id
		_ = user_id
	}

  if sort_by == "following" {
		db_query += `JOIN followers ON 
        followers.followed_id=posts.author_id 
        AND followers.follower_id=$1 
        `
	}

	db_query += `WHERE posts.parent_id IS NULL `

	if !show_all {
		db_query += " AND posts.board_id=$2 "
	}

  db_query += `ORDER BY `
  if sort_by != "following" {
    db_query += sort_by
  } else {
    db_query += "created_on"
  }
  db_query += ` DESC `


	if !show_all {
    db_query += `LIMIT $3 OFFSET $4`
    _, _ = db.Select(&threads, db_query, user_id, board.Id, THREADS_PER_PAGE_on_boards-1, i_begin)
	} else {
    db_query += `LIMIT $2 OFFSET $3`
    _, _ = db.Select(&threads, db_query, user_id, THREADS_PER_PAGE-1, i_begin)
	}

	for i := range threads {
		if user_id == -1 {
			threads[i].ViewedOn = pq.NullTime{Time: time.Now(), Valid: true}
		}

		obj, _ := db.Get(&User{}, threads[i].AuthorId)
		user := obj.(*User)

		threads[i].Author = user
		threads[i].Thread = &Post{
			Id: threads[i].Id,
		}
	}

	return threads
}

func (board *Board) GetPagesInBoard() int {
	db := GetDbSession()
	count, err := db.SelectInt("SELECT COUNT(*) FROM posts WHERE board_id=$1 AND parent_id IS NULL", board.Id)

	if err != nil {
		fmt.Printf("[error] Could not get pages in the board %d (%s)\n", board.Id, err.Error())
	}


	return int(math.Floor(float64(count) / float64(THREADS_PER_PAGE)))
}

// Deletes a board and all of the posts it contains
func (board *Board) Delete() {
	db := GetDbSession()
    db.Exec("DELETE FROM views WHERE post_id IN (SELECT id FROM posts WHERE board_id=$1)", board.Id)
	db.Exec("DELETE FROM posts WHERE board_id=$1", board.Id)
	db.Delete(board)
}
