package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"nightstalker/controllers"
	"net/http"
)

const PORT = "8080"

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/", controllers.Index)
	r.HandleFunc("/register", controllers.Register)
	r.HandleFunc("/login", controllers.Login)
	r.HandleFunc("/logout", controllers.Logout)
	r.HandleFunc("/admin", controllers.Admin)
	r.HandleFunc("/admin/boards", controllers.AdminBoards)
	r.HandleFunc("/admin/users/{id:[0-9]+}", controllers.AdminUser)
	r.HandleFunc("/admin/users", controllers.AdminUsers)
    r.HandleFunc("/action/like", controllers.ActionLikeThread)
	r.HandleFunc("/action/stick", controllers.ActionStickThread)
	r.HandleFunc("/action/lock", controllers.ActionLockThread)
	r.HandleFunc("/action/delete", controllers.ActionDeleteThread)
	r.HandleFunc("/action/move", controllers.ActionMoveThread)
	r.HandleFunc("/action/mark_read", controllers.ActionMarkAllRead)
    r.HandleFunc("/action/mark_notifs_read", controllers.ActionMarkNotificationsAllRead)
	r.HandleFunc("/action/edit", controllers.PostEditor)
	r.HandleFunc("/board/{id:[0-9]+}", controllers.Board)
	r.HandleFunc("/board/{board_id:[0-9]+}/new", controllers.PostEditor)
	r.HandleFunc("/board/{board_id:[0-9]+}/{post_id:[0-9]+}", controllers.Thread)
	r.HandleFunc("/user/{id:[0-9]+}", controllers.User)
	r.HandleFunc("/user/{id:[0-9]+}/settings", controllers.UserSettings)
    r.HandleFunc("/static/images", controllers.Images)
    r.HandleFunc("/static/assets", controllers.Images)
    r.HandleFunc("/static/data", controllers.Images)

    static_path := "./templates/"
    r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(static_path)))

    r.PathPrefix("/assets/").Handler(http.FileServer(http.Dir(static_path)))
	http.Handle("/", r)

	fmt.Println("[NIGHSTALKER LOG] started server on port " + PORT)
	http.ListenAndServe(":"+PORT, nil)
}
