// Package web is the main package of tournify which defines all routes and database connections and settings, the glue to the entire application
package web

import (
	"github.com/hako/durafmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"strings"
	"time"
)

func loadTemplates() (*template.Template, error) {
	var err4 error
	t := template.New("")
	t = t.Funcs(template.FuncMap{
		"timeToAgo": func(timeToCompare time.Time) string {
			dur := time.Now().Sub(timeToCompare)
			duration := durafmt.Parse(dur)
			return duration.LimitFirstN(1).String()
		},
	})
	err := fs.WalkDir(staticFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		f, err2 := staticFS.Open(path)
		if err2 != nil {
			return err2
		}
		h, err3 := ioutil.ReadAll(f)
		if err3 != nil {
			return err3
		}
		parts := strings.Split(path, "/")
		if len(parts) > 0 && strings.HasSuffix(parts[len(parts)-1], ".html") {
			t, err4 = t.New(parts[len(parts)-1]).Parse(string(h))
			if err4 != nil {
				return err4
			}
		}
		return nil
	})

	return t, err
}
