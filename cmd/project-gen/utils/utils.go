package utils

import (
	"encoding/json"
	"fmt"
	"github.com/flosch/pongo2/v4"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func cloneProjectTemplate(from, to string) (err error) {
	_, err = git.PlainClone(to, false, &git.CloneOptions{
		URL:      from,
		Progress: os.Stdout,
	})
	return
}

func parseTemplateJSON(path string) (rv map[string]interface{}, err error) {
	rv = make(map[string]interface{})
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &rv)
	if err != nil {
		return nil, err
	}
	return
}

func generateFromTemplateProject(tplPath string, destPath string) (err error) {
	// get template vars
	conf, err := parseTemplateJSON(path.Join(tplPath, "kratos.json"))
	if err != nil {
		return err
	}

	// find template root
	files, err := ioutil.ReadDir(tplPath)
	if err != nil {
		return err
	}
	var tplRootName string
	for _, f := range files {
		if f.IsDir() && strings.HasPrefix(f.Name(), "{{") && strings.HasSuffix(f.Name(), "}}") {
			tplRootName = f.Name()
			break
		}
	}

	tplRoot := path.Join(tplPath, tplRootName)
	tpl, err := pongo2.FromString(tplRootName)
	if err != nil {
		return err
	}

	// create project root
	projectName, err := tpl.Execute(pongo2.Context(conf))
	projectRoot := path.Join(destPath, projectName)
	err = os.MkdirAll(projectRoot, 0755)
	if err != nil {
		return err
	}

	// walk template directory and generate to project directory
	err = filepath.Walk(tplRoot, func(path string, info os.FileInfo, err error) error {
		fmt.Printf("path:%v name:%v\n", path, info.Name())
		return nil
	})
	return
}
