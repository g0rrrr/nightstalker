package controllers

import (
    "fmt"
    "nightstalker/models"
    "nightstalker/utils"
    "net/http"
)

type PageData struct {
    Title       string
    IsIndexPage bool
}

func Index(w http.ResponseWriter, r *http.Request) {
    db := models.GetDbSession()

    obj, err := db.Get(&models.Board{}, 1)
    if err != nil || obj == nil {
        http.NotFound(w, r)
        fmt.Printf("[error] Something went wrong in index (%s)\n", err.Error())
        return
    }

    board := obj.(*models.Board)
  page_id := 0
  session, _ := utils.GetCookieStore(r).Get(r, "sirsid")
  sort_by := r.URL.Query().Get("sort")

  if sort_by == "" {
    if session.Values["sort_by"] == nil {
      sort_by = "latest_reply"
    } else {
      sort_by = session.Values["sort_by"].(string)
    }
  } else {
    session.Values["sort_by"] = sort_by
  }

  /*if sort_by == "" && session.Values["sort_by"] == nil {
    sort_by = "latest_reply"
    session.Values["sort_by"] = sort_by
  } else if sort_by == "" && session.Values["sort_by"] != nil {
    sort_by = session.Values["sort_by"].(string)
  } else if sort_by != "" {
    session.Values["sort_by"] = sort_by
  }*/

  err = session.Save(r, w)
  if err != nil {
    fmt.Printf("[error] Could not save session (%s)\n", err.Error())
  }

  current_user := utils.GetCurrentUser(r)
  boards, err := models.GetBoardsUnread(current_user)
  threads:= board.GetThreads(sort_by, true, page_id, current_user)

  user_count, _ := models.GetUserCount()
  latest_user, _ := models.GetLatestUser()
  total_posts, _ := models.GetPostCount()

  latest_post := board.GetLatestPost()
  fmt.Println(latest_post.Op.LatestReply)

  utils.RenderTemplate(w, r, "index.html", map[string]interface{}{
      "board":        board,
      "boards":       boards,
      "threads":      threads,
      "user_count":   user_count,
      "online_users": models.GetOnlineUsers(),
      "latest_user":  latest_user,
      "total_posts":  total_posts,
      "SortBy":       sort_by,
      "IsIndexPage":  true,
  }, map[string]interface{}{
      "IsUnread": func(join *models.JoinBoardView) bool {
          latest_post := join.Board.GetLatestPost()

          if current_user != nil && !current_user.LastUnreadAll.Time.Before(latest_post.Op.LatestReply) {
              return false
          }

          return !join.ViewedOn.Valid || join.ViewedOn.Time.Before(latest_post.Op.LatestReply)
      },
      "IsThreadUnread": func(join *models.JoinThreadView) bool {
        if current_user != nil && !current_user.LastUnreadAll.Time.Before(join.LatestReply) {
          return false
        }
        return !join.ViewedOn.Valid || join.ViewedOn.Time.Before(join.LatestReply)
      },
  })
}
