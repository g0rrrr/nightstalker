package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"nightstalker/models"
	"nightstalker/utils"
)

func Thread(w http.ResponseWriter, r *http.Request) {
	page_id_str := r.FormValue("page")
	page_id, err := strconv.Atoi(page_id_str)
	if err != nil {
		page_id = 0
	}

	board_id_str := mux.Vars(r)["board_id"]
	board_id, _ := strconv.Atoi(board_id_str)
	board, err := models.GetBoard(board_id)

	post_id_str := mux.Vars(r)["post_id"]
	post_id, _ := strconv.Atoi(post_id_str)
	err, op, posts := models.GetThread(post_id, page_id)
  fmt.Println(err)

	var posting_error error

	current_user := utils.GetCurrentUser(r)
	if r.Method == "POST" {
		db := models.GetDbSession()
		title := r.FormValue("title")
		content := r.FormValue("content")

		if current_user == nil {
			http.NotFound(w, r)
			return
		}

		if op.Locked && !current_user.CanModerate() {
			http.NotFound(w, r)
			return
		}

		post := models.NewPost(current_user, board, title, content)
		post.ParentId = sql.NullInt64{int64(post_id), true}
		op.LatestReply = time.Now()

		posting_error = post.Validate()

		if posting_error == nil {
			db.Insert(post)
			db.Update(op)

      if current_user.Id != op.AuthorId {
        models.Experience(current_user.Id, 1)
        models.Experience(op.AuthorId, 2)
        message := fmt.Sprintf("replied to your post")
        _, err := db.Exec("INSERT INTO notifications (user_id, notif_user_id, message) VALUES ($1, $2, $3)", op.AuthorId, current_user.Id, message)
        if err != nil {
            return
        }
      }

			if page := post.GetPageInThread(); page != page_id {
				http.Redirect(w, r, fmt.Sprintf("/board/%d/%d?page=%d#post_%d", post.BoardId, op.Id, page, post.Id), http.StatusFound)
			}

			err, op, posts = models.GetThread(post_id, page_id)
		}
	}

	if err != nil {
		http.NotFound(w, r)
		fmt.Printf("[error] Something went wrong in posts (%s)\n", err.Error())
		return
	}

	num_pages := op.GetPagesInThread()

	if page_id > num_pages {
		http.NotFound(w, r)
		return
	}

	var previous_text string
	if posting_error != nil {
		previous_text = r.FormValue("content")
	}

	// Mark the thread as read
	if current_user != nil {
		models.AddView(current_user, op)
	}

	utils.RenderTemplate(w, r, "thread.html", map[string]interface{}{
		"board":         board,
		"op":            op,
		"posts":         posts,
		"first_page":    (page_id > 0),
		"prev_page":     (page_id > 1),
		"next_page":     (page_id < num_pages-1),
		"last_page":     (page_id < num_pages),
		"page_id":       page_id,
		"posting_error": posting_error,
		"previous_text": previous_text,
	}, map[string]interface{}{

    "IfOp" : func(thread *models.Post) bool { 
      current_user := utils.GetCurrentUser(r)
			if current_user == nil {
				return false
			}
      return current_user.Id == op.Id 
    },

		"CurrentUserCanModerateThread": func(thread *models.Post) bool {
			current_user := utils.GetCurrentUser(r)
			if current_user == nil {
				return false
			}

			return (current_user.CanModerate() && thread.ParentId.Valid == false)
		},

		"CurrentUserCanDeletePost": func(thread *models.Post) bool {
			current_user := utils.GetCurrentUser(r)
			if current_user == nil {
				return false
			}

			return (current_user.Id == thread.AuthorId) || current_user.CanModerate()
		},

		"CurrentUserCanEditPost": func(post *models.Post) bool {
			current_user := utils.GetCurrentUser(r)
			if current_user == nil {
				return false
			}

			return (current_user.Id == post.AuthorId || current_user.CanModerate())
		},

		"CurrentUserCanModerate": func() bool {
			current_user := utils.GetCurrentUser(r)
			if current_user == nil {
				return false
			}

			return current_user.CanModerate()
		},

		"SignaturesEnabled": func() bool {
            enable_signatures := true
			return enable_signatures
		},

		"CurrentUserCanReply": func(post *models.Post) bool {
			current_user := utils.GetCurrentUser(r)
			if current_user != nil && (!post.Locked || current_user.CanModerate()) {
				return true
			}
			return false
		},

		"CurrentUserCanLike": func(post *models.Post) bool {
			current_user := utils.GetCurrentUser(r)
			if current_user != nil && current_user.Id != op.Id {
				return true
			}
			return false
		},
	})
}
