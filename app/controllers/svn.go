package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"os"
	"os/exec"
	"path/filepath"
	"svnadmin/app/routes"
	"regexp"
)

type Svn struct {
	App
}

func (c Svn) getSvnParentPath() string {
	return revel.Config.StringDefault("svn.parent_path", "")
}

func (c Svn) getSvnAdminBin() string {
	return revel.Config.StringDefault("svn.admin_bin", "")
}

func (c Svn) Index() revel.Result {
	parent_path := c.getSvnParentPath()
	repos, _ := filepath.Glob(parent_path + "/*")

	for i, path := range repos {
		repos[i] = filepath.Base(path)
	}

	svn_url_base := revel.Config.StringDefault("svn.url", "http://xxxxxxxxx/")
	
	return c.Render(repos, svn_url_base)
}

func (c Svn) Create(Name string) revel.Result {
	repo := c.getSvnParentPath() + "/" + Name
		
	_, e := os.Stat(repo)
	is_exists := false
	if e != nil {
		is_exists = true
	}
	
	re := regexp.MustCompile("^[a-z0-9_.-]+$")
	
	c.Validation.Required(Name).Message("リポジトリ名は必須です。")
	c.Validation.Required(Name != "websvnadmin").Message("websvnadminは予約語のた指定できません。")
	c.Validation.Required(Name != "svn").Message("svnは予約語のた指定できません。")
	c.Validation.Required(is_exists).Message("入力のリポジトリ名はすでに存在しています。")
	c.Validation.Match(Name, re).Message("リポジトリ名は半角英数記号で指定してください。")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Svn.Index())
	}

	err := exec.Command("sudo", c.getSvnAdminBin(), "create", repo).Run()
	if err != nil {
		c.FlashParams()
		fmt.Println(err)
		
		c.Flash.Error(fmt.Sprintf("%sの作成に失敗しました。", Name))
		return c.Redirect(routes.Svn.Index())
	}

	owner := revel.Config.StringDefault("svn.owner", "apache")
	group := revel.Config.StringDefault("svn.group", "apache")
	permit := revel.Config.StringDefault("svn.permit", "775")
	
	exec.Command("sudo", "chown", "-R", fmt.Sprintf("%s.%s", owner, group), repo).Run()
	exec.Command("sudo", "chmod", "-R", permit, repo).Run()

	c.Flash.Success(fmt.Sprintf("%sを作成しました。", Name))
	return c.Redirect(routes.Svn.Index())
}
