package main

import (
	"ffiiitc/internal/classifier"
	"ffiiitc/internal/config"
	"github.com/go-pkgz/lgr"
	"github.com/navossoc/bayesian"
)

func main() {
	l := lgr.New(lgr.Debug, lgr.CallerFunc)
	l.Logf("INFO Inspecting model")
	cls, err := classifier.NewTrnClassifierFromFile(config.ModelFile, l)
	if err != nil {
		l.Logf("ERROR %v", err)
	}

	catList := cls.Classifier.Classes
	l.Logf("DEBUG learned classes: %v", catList)
	for _, cat := range catList {
		l.Logf("DEBUG class: %v", cat)
		words := cls.Classifier.WordsByClass(bayesian.Class(cat))
		l.Logf("DEBUG words: %v", words)
	}
}
