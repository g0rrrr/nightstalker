package utils

import (
	"fmt"
    "os"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/russross/blackfriday"
	"nightstalker/models"
)

// Returns a list of all available themes
func ListTemplates() []string {
	names := []string{"default"}

    static_path := "./templates/"
	files, _ := os.ReadDir(static_path)

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		names = append(names, f.Name())
	}

	return names
}

func tplAdd(first, second int) int {
	return first + second
}

func tplParseMarkdown(input string) template.HTML {
	byte_slice := []byte(input)
	return template.HTML(string(blackfriday.MarkdownBasic(byte_slice)))
}

func tplGetCurrentUser(r *http.Request) func() *models.User {
	return func() *models.User {
		return GetCurrentUser(r)
	}
}

func tplGetStringSetting(key string) string {
	val, _ := models.GetStringSetting(key)
	return val
}

func tplIsValidTime(in time.Time) bool {
	return in.Year() > 1
}

var default_funcmap = template.FuncMap{
	"TimeRelativeToNow": TimeRelativeToNow,
	"Joined":            Joined,
	"Add":               tplAdd,
	"ParseMarkdown":     tplParseMarkdown,
	"IsValidTime":       tplIsValidTime,
	"GetStringSetting":  tplGetStringSetting,
}

func RenderTemplate(
    out http.ResponseWriter,
    r *http.Request,
    tpl_file string,
    context map[string]interface{},
    funcs template.FuncMap) {

    current_user := GetCurrentUser(r)

    send := map[string]interface{}{
        "current_user":   current_user,
        "request":        r,
        "IsIndexPage":    context["IsIndexPage"],
    }

    // Merge the global template variables with the local context
    for key, val := range context {
        send[key] = val
    }

    // Same with the function map
    func_map := default_funcmap
    func_map["GetCurrentUser"] = tplGetCurrentUser(r)
    for key, val := range funcs {
        func_map[key] = val
    }

    // Get the base template path
    var base_path string

    base_path = "./templates"

    base_tpl := filepath.Join(base_path, "base.html")
    rend_tpl := filepath.Join(base_path, tpl_file)

    tpl, err := template.New("tpl").Funcs(func_map).ParseFiles(base_tpl, rend_tpl)
    if err != nil {
        fmt.Printf("[error] Could not parse template (%s)\n", err.Error())
    }

    // Attempt to execute the template we're on
    err = tpl.ExecuteTemplate(out, tpl_file, send)
    if err != nil {
        fmt.Printf("[error] Could not parse template (%s)\n", err.Error())
    }

    // And now the base template
    err = tpl.ExecuteTemplate(out, "base.html", send)
    if err != nil {
        fmt.Printf("[error] Could not parse template (%s)\n", err.Error())
    }
}

