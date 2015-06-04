package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"svnadmin/app/routes"
)

type Svn struct {
	App
}

type Svninfo struct {
	Name     string
	Url      string
	LastDate string
	LastRev  string
}

func (c Svn) getSvnParentPath() string {
	return revel.Config.StringDefault("svn.parent_path", "")
}

func (c Svn) getSvnAdminBin() string {
	return revel.Config.StringDefault("svn.svnadmin", "svnadmin")
}

func (c Svn) getSvnLookBin() string {
	return revel.Config.StringDefault("svn.svnlook", "svnlook")
}

func (c Svn) Index() revel.Result {
	parent_path := c.getSvnParentPath()
	repos, _ := filepath.Glob(parent_path + "/*")

	svn_url_base := revel.Config.StringDefault("svn.url", "http://xxxxxxxxx/")

	svnlook := c.getSvnLookBin()

	var svninfos []Svninfo
	for _, path := range repos {
		date, _ := exec.Command(svnlook, "date", path).Output()
		rev, _ := exec.Command(svnlook, "youngest", path).Output()

		name := filepath.Base(path)
		info := Svninfo{
			Name:     name,
			Url:      svn_url_base + name,
			LastDate: string(date),
			LastRev:  "r" + string(rev),
		}

		svninfos = append(svninfos, info)
	}

	return c.Render(svninfos)
}

func (c Svn) Create(Name string) revel.Result {
	repo := c.getSvnParentPath() + "/" + Name

	_, e := os.Stat(repo)
	is_exists := false
	if e != nil {
		is_exists = true
	}

	c.Validation.Required(Name).Message("リポジトリ名は必須です。")
	c.Validation.Required(Name != "websvnadmin").Message(Name + "は予約語のた指定できません。")
	c.Validation.Required(Name != "svn").Message(Name + "は予約語のた指定できません。")
	c.Validation.Required(is_exists).Message(Name + "はすでに存在しています。")
	c.Validation.Match(Name, regexp.MustCompile("^[a-z0-9_.-]+$")).Message("リポジトリ名は半角英数記号で指定してください。")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Svn.Index())
	}

	err := exec.Command("sudo", c.getSvnAdminBin(), "create", repo).Run()
	if err != nil {
		c.FlashParams()
		revel.TRACE.Println(err)

		c.Flash.Error(fmt.Sprintf("%sの作成に失敗しました。", Name))
		return c.Redirect(routes.Svn.Index())
	}

	owner := revel.Config.StringDefault("svn.owner", "apache")
	group := revel.Config.StringDefault("svn.group", "apache")
	permit := revel.Config.StringDefault("svn.permit", "775")

	exec.Command("sudo", "chown", "-R", fmt.Sprintf("%s.%s", owner, group), repo).Run()
	exec.Command("sudo", "chmod", "-R", permit, repo).Run()

	revel.TRACE.Println(Name)

	c.Flash.Success(fmt.Sprintf("%sを作成しました。", Name))
	return c.Redirect(routes.Svn.Index())
}
