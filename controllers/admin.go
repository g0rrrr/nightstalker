package controllers

import (
	"nightstalker/models"
	"nightstalker/utils"
	"net/http"
)

func Images(w http.ResponseWriter, r *http.Request) {
  current_user := utils.GetCurrentUser(r)
  if !current_user.IsAdmin() {
    http.NotFound(w, r)
    return
  }
}

func Admin(w http.ResponseWriter, r *http.Request) {
	current_user := utils.GetCurrentUser(r)
	if current_user == nil || !current_user.IsAdmin() {
		http.NotFound(w, r)
		return
	}

	var err error
	success := false
	current_template, _ := models.GetStringSetting("template")

	if r.Method == "POST" {
		current_template = r.FormValue("template")
		models.SetStringSetting("template", current_template)
		success = true
	}

	utils.RenderTemplate(w, r, "admin.html", map[string]interface{}{
		"error":            err,
		"success":          success,
		"current_template": current_template,
		"templates":        utils.ListTemplates(),
	}, map[string]interface{}{
		"IsCurrentTemplate": func(name string) bool {
			return name == current_template
		},
	})
}
