package controllers

import (
	"github.com/gorilla/mux"
	"nightstalker/models"
	"nightstalker/utils"
	"net/http"
	"strconv"
)

func Board(w http.ResponseWriter, r *http.Request) {
	db := models.GetDbSession()

	page_id_str := r.FormValue("page")
	page_id, err := strconv.Atoi(page_id_str)
	if err != nil {
		page_id = 0
	}

	board_id_str := mux.Vars(r)["id"]
	board_id, _ := strconv.Atoi(board_id_str)
	obj, err := db.Get(&models.Board{}, board_id)
	if err != nil || obj == nil {
		http.NotFound(w, r)
		return
	}
	board := obj.(*models.Board)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	current_user := utils.GetCurrentUser(r)
	threads:= board.GetThreads("latest_reply", false, page_id, current_user)

	num_pages := board.GetPagesInBoard()

	utils.RenderTemplate(w, r, "board.html", map[string]interface{}{
		"board":     board,
		"threads":   threads,
		"page_id":   page_id,
		"prev_page": (page_id != 0),
		"next_page": (page_id < num_pages),
	}, map[string]interface{}{
		"IsUnread": func(join *models.JoinThreadView) bool {
			if current_user != nil && !current_user.LastUnreadAll.Time.Before(join.LatestReply) {
				return false
			}
			return !join.ViewedOn.Valid || join.ViewedOn.Time.Before(join.LatestReply)
		},
	})
}
