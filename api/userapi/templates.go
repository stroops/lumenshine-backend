package main

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/Soneso/lumenshine-backend/constants"

	mw "github.com/Soneso/lumenshine-backend/api/middleware"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/gin"
)

//Templates holds all the html templates that are used for sending e.g. emails
//the first key is the language, the second the name for the template
var Templates map[string]map[string]*template.Template

func loadTemplates(log *logrus.Entry) {
	var cnf rice.Config
	//we will search for the file localy on FS first, then embedded
	//this is easyer for development but not so good for production, as one can forget to rebuild the resources
	//we assume that every template is in a lang directory, and that every template is present in every lang
	cnf.LocateOrder = []rice.LocateMethod{
		rice.LocateFS, rice.LocateEmbedded, rice.LocateAppended, rice.LocateWorkingDirectory,
	}
	Templates = make(map[string]map[string]*template.Template)
	box := cnf.MustFindBox("templates/mail")

	//add all files from the box. name is the filename without extention
	err := box.Walk("", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".html" {
			data, err := box.String(path)
			if err != nil {
				log.WithError(err).WithField("file", path).Fatalf("error unboxing file")
			}

			dir, file := filepath.Split(path)
			lang := dir[:len(dir)-1]
			templateName := strings.TrimSuffix(file, filepath.Ext(file))
			if _, ok := Templates[lang]; !ok {
				Templates[lang] = make(map[string]*template.Template)
			}
			Templates[lang][templateName] = template.Must(template.New(templateName).Parse(data))
			log.WithFields(logrus.Fields{"file": file, "lang": lang}).Debug("Added mail-template")
		}
		return nil
	})
	if err != nil {
		log.WithError(err).Fatalf("Error loading Templates")
		os.Exit(1)
	}
}

func notUsedAddTemplate(log *logrus.Entry, box *rice.Box, templateName string, fileNames ...string) {
	for _, language := range constants.ServerLanguages {
		lang := language.String()
		if _, ok := Templates[lang]; !ok {
			Templates[lang] = make(map[string]*template.Template)
		}

		templateString := ""
		for _, file := range fileNames {
			tmpString, err := box.String(lang + "/" + file)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{
					"template": file,
					"language": lang,
				}).Fatalf("Error parsing template")
			}
			templateString += tmpString
		}

		Templates[lang][templateName] = template.Must(template.New(templateName).Parse(templateString))
	}
}

//RenderTemplateToString renders the given template to a string
func RenderTemplateToString(uc *mw.IcopContext, c *gin.Context, templateName string, data interface{}) string {
	langCode := uc.Language
	if langCode == "" {
		langCode = "en"
	}

	t, ok := Templates[langCode][templateName]
	if !ok {
		uc.Log.WithField("template", templateName).Error("Template not found!")
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		uc.Log.WithError(err).WithField("template", templateName).Error("Error executing template")
		return "Error rendering template"
	}

	return tpl.String()
}
