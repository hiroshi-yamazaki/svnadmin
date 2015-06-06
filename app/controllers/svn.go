package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"svnadmin/app/libs"
	"svnadmin/app/models"
	"svnadmin/app/routes"
)

type Svn struct {
	App
}

func (c Svn) getSvnParentPath() string {
	return revel.Config.StringDefault("svn.parent_path", "/home/svn/repos")
}

func (c Svn) getSvnAdminBin() string {
	return revel.Config.StringDefault("svn.svnadmin", "svnadmin")
}

func (c Svn) getSvnOwner() string {
	return revel.Config.StringDefault("svn.owner", "apache")
}

func (c Svn) getSvnGroup() string {
	return revel.Config.StringDefault("svn.group", "apache")
}

func (c Svn) getSvnPermission() string {
	return revel.Config.StringDefault("svn.permit", "775")
}

func (c Svn) Index() revel.Result {
	parent_path := c.getSvnParentPath()
	repos, _ := filepath.Glob(parent_path + "/*")

	pager := &libs.Pager{
		Params: c.Params,
		Limit:  15,
		Left:   5,
		Right:  5,
		Items:  repos,
		Total:  len(repos),
	}

	page := pager.Result()
	svninfos := models.GetSvninfoList(page.Items)

	return c.Render(svninfos, page)
}

func (c Svn) IsExists(Repo string) bool {
	_, e := os.Stat(Repo)
	is_exists := false
	if e != nil {
		is_exists = true
	}

	return is_exists
}

func (c Svn) Create(Name string) revel.Result {
	repo := c.getSvnParentPath() + "/" + Name

	c.Validation.Required(Name).Message("リポジトリ名は必須です。")
	c.Validation.Required(Name != "websvnadmin").Message(Name + "は予約語のた指定できません。")
	c.Validation.Required(Name != "svn").Message(Name + "は予約語のた指定できません。")
	c.Validation.Required(c.IsExists(repo)).Message(Name + "はすでに存在しています。")
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

	err = exec.Command("sudo", "chown", "-R", fmt.Sprintf("%s.%s", c.getSvnOwner(), c.getSvnOwner()), repo).Run()
	if err != nil {
		revel.TRACE.Println(err)
	}

	err = exec.Command("sudo", "chmod", "-R", c.getSvnPermission(), repo).Run()
	if err != nil {
		revel.TRACE.Println(err)
	}

	revel.TRACE.Println(Name)

	c.Flash.Success(fmt.Sprintf("%sを作成しました。", Name))
	return c.Redirect(routes.Svn.Index())
}
