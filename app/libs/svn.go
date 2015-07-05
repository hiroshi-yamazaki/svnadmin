package libs

import (
	"github.com/revel/revel"
)

func GetSvnParentPath() string {
	return revel.Config.StringDefault("svn.parent_path", "/home/svn/repos")
}

func GetSvnAdminBin() string {
	return revel.Config.StringDefault("svn.svnadmin", "svnadmin")
}

func GetSvnOwner() string {
	return revel.Config.StringDefault("svn.owner", "apache")
}

func GetSvnGroup() string {
	return revel.Config.StringDefault("svn.group", "apache")
}

func GetSvnPermission() string {
	return revel.Config.StringDefault("svn.permit", "775")
}

func GetSvnLookBin() string {
	return revel.Config.StringDefault("svn.svnlook", "svnlook")
}

func GetSvnUrlBase() string {
	return revel.Config.StringDefault("svn.url", "http://xxxxxxxxx/")
}
