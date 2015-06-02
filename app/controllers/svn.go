package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"os"
	"os/exec"
	"path/filepath"
	"svnadmin/app/routes"
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

	return c.Render(repos)
}

func (c Svn) Create(Name string) revel.Result {
	repo := c.getSvnParentPath() + "/" + Name

	_, e := os.Stat(repo)
	is_exists := false
	if e != nil {
		is_exists = true
	}

	c.Validation.Required(Name).Message("Name must input")
	c.Validation.Required(is_exists).Message("Name already exist")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Svn.Index())
	}

	err := exec.Command("sudo", c.getSvnAdminBin(), "create", repo).Run()
	if err != nil {
		c.FlashParams()
		fmt.Println(err)
		return c.Redirect(routes.Svn.Index())
	}

	owner := revel.Config.StringDefault("svn.owner", "apache")
	group := revel.Config.StringDefault("svn.group", "apache")

	exec.Command("sudo", "chown", "-R", fmt.Sprintf("%s.%s", owner, group), repo).Run()
	exec.Command("sudo", "chmod", "-R", "775", repo).Run()

	return c.Redirect(routes.Svn.Index())
}
