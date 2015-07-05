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

func (c Svn) Index() revel.Result {
	parent_path := libs.GetSvnParentPath()
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
	repo := libs.GetSvnParentPath() + "/" + Name

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

	err := exec.Command("sudo", libs.GetSvnAdminBin(), "create", repo).Run()
	if err != nil {
		c.FlashParams()
		revel.INFO.Println(err)

		c.Flash.Error(fmt.Sprintf("%sの作成に失敗しました。", Name))
		return c.Redirect(routes.Svn.Index())
	}

	err = exec.Command("sudo", "chown", "-R", fmt.Sprintf("%s.%s", libs.GetSvnOwner(), libs.GetSvnOwner()), repo).Run()
	if err != nil {
		revel.INFO.Println(err)
	}

	err = exec.Command("sudo", "chmod", "-R", libs.GetSvnPermission(), repo).Run()
	if err != nil {
		revel.INFO.Println(err)
	}

	revel.INFO.Println(Name)

	c.Flash.Success(fmt.Sprintf("%sを作成しました。", Name))
	return c.Redirect(routes.Svn.Index())
}
