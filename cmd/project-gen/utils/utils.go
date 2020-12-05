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

	tplRoot := path.Join(tplPath, tplRootName, "/")
	fmt.Printf("tplRoot :%v\n", tplRoot)

	// walk template directory and generate to project directory
	err = filepath.Walk(tplRoot, func(p string, info os.FileInfo, err error) error {
		rp := strings.TrimPrefix(p, tplRoot)
		fmt.Printf("\npath:%v\n", rp)

		if info.IsDir() {
			//fmt.Printf("dir path:%v name:%v\n", p, info.Name())
			dnTpl, err := pongo2.FromString(rp)
			if err != nil {
				return err
			}
			dirName, err := dnTpl.Execute(pongo2.Context(conf))
			if err != nil {
				return err
			}
			destDirPath := path.Join(destPath, dirName)
			fmt.Printf("destDirPath:%v \n", destDirPath)
			err = os.MkdirAll(destDirPath, 0755)
		} else {
			fmt.Printf("file path:%v \n", p)
			fnTpl, err := pongo2.FromString(rp)
			if err != nil {
				return err
			}
			fcTpl, err := pongo2.FromFile(p)
			if err != nil {
				return err
			}
			fn, err := fnTpl.Execute(pongo2.Context(conf))
			if err != nil {
				return err
			}
			fc, err := fcTpl.Execute(pongo2.Context(conf))
			if err != nil {
				return err
			}
			destFilePath := path.Join(destPath, fn)
			fmt.Printf("write to:%v \n", destFilePath)
			if err = ioutil.WriteFile(destFilePath, []byte(fc), info.Mode()); err != nil {
				return err
			}
		}
		return err
	})
	return
}
