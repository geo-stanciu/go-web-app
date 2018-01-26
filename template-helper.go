package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func parseTemplate(basePath string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// don't process folders themselves
		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".html") {
			return nil
		}

		templateName := strings.Replace(path[len(basePath):], "\\", "/", -1)

		if templates == nil {
			templates = template.New(templateName)
			templates.Delims(templateDelims[0], templateDelims[1])
			_, err = parseTemplateFile(templates, templateName, path)
		} else {
			_, err = parseTemplateFile(templates, templateName, path)
		}

		return err
	}
}

func parseTemplateFile(t *template.Template, templateName string, filename string) (*template.Template, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	s := string(b)

	var tmpl *template.Template

	if templateName == t.Name() {
		tmpl = t
	} else {
		tmpl = t.New(templateName)
	}

	_, err = tmpl.Parse(s)

	if err != nil {
		return nil, err
	}
	return t, nil
}
