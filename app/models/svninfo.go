package models

import (
	"github.com/revel/revel"
	"os/exec"
	"path/filepath"
)

type Svninfo struct {
	Name       string
	Url        string
	LastDate   string
	LastRev    string
	LastAuthor string
}

func GetSvnLookBin() string {
	return revel.Config.StringDefault("svn.svnlook", "svnlook")
}

func GetSvnUrlBase() string {
	return revel.Config.StringDefault("svn.url", "http://xxxxxxxxx/")
}

func GetSvninfoList(list []string) []Svninfo {
	svn_url_base := GetSvnUrlBase()
	svnlook := GetSvnLookBin()

	var svninfos []Svninfo
	for _, path := range list {
		date, _ := exec.Command(svnlook, "date", path).Output()
		rev, _ := exec.Command(svnlook, "youngest", path).Output()
		author, _ := exec.Command(svnlook, "author", path).Output()

		name := filepath.Base(path)
		info := Svninfo{
			Name:       name,
			Url:        svn_url_base + name,
			LastDate:   string(date),
			LastRev:    "r" + string(rev),
			LastAuthor: string(author),
		}

		svninfos = append(svninfos, info)
	}

	return svninfos
}
