package controllers

import (
  "fmt"
	"github.com/gorilla/mux"
	"nightstalker/models"
	"nightstalker/utils"
	"net/http"
	"strconv"
)

func User(w http.ResponseWriter, r *http.Request) {
	db := models.GetDbSession()

	user_id_str := mux.Vars(r)["id"]
	user_id, err := strconv.Atoi(user_id_str)
  current_user := utils.GetCurrentUser(r)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	user, err := db.Get(&models.User{}, user_id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

  var error error
  if r.Method == "POST" {
    err = current_user.FollowUser(int64(user_id))
    fmt.Println(err)
    if err == nil {
      if int64(user_id) != current_user.Id {
        message := fmt.Sprintf("started to follow you")
        this, err := db.Exec("INSERT INTO notifications (user_id, notif_user_id, message) VALUES ($1, $2, $3)", user_id, current_user.Id, message)
        _ = this 
        if err != nil {
          return
        }
      }
    } else {
      error = err
    }
  }
  progress, needed := models.GetLevelProgress(int64(user_id))

	utils.RenderTemplate(w, r, "user.html", map[string]interface{}{
    "already_following": current_user.AlreadyFollowing(int64(user_id)),
    "progress": progress,
    "needed": needed,
    "error": error,
		"user": user,
	}, map[string]interface{}{
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
	})
}
