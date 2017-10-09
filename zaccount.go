package main

import (
	"log"

	"github.com/zhengxiaoyao0716/zcli/client"

	"github.com/kardianos/service"
	"github.com/zhengxiaoyao0716/zmodule"
)

func main() {
	zmodule.Main("zaccount",
		&service.Config{
			Name:        "ZhengAccountService",
			DisplayName: "Account Service",
			Description: "Daemon service for ZAccount.",
		}, func() {
			log.Println("Running.")
		})
}

// In this way that override those values,
// you can use `main` as the module name, instead of `github.com/zhengxiaoyao0716/zmodule`.
var (
	Version   string // `git describe --tags`
	Built     string // `date +%FT%T%z`
	GitCommit string // `git rev-parse --short HEAD`
	GoVersion string // `go version`
)

func initInfo() {
	zmodule.Author = "zhengxiaoyao0716"
	zmodule.Homepage = "https://zhengxiaoyao0716.github.io/zaccound"
	zmodule.Repository = "https://github.com/zhengxiaoyao0716/zaccound"
	zmodule.License = "${Repository}/blob/master/LICENSE"

	zmodule.Version = Version
	zmodule.Built = Built
	zmodule.GitCommit = GitCommit
	zmodule.GoVersion = GoVersion
}

func initArgs() {
	zmodule.Args["db_path"] = zmodule.Argument{
		Type:    "string",
		Default: ".zaccount." + Version + ".db",
		Usage:   "Path of the database file.",
	}

	// // Dangerous! All Argument would be storage into a temp file in order to deliver into daemon service!
	// // So you should never push secret arguments into `zmodule.Args`.
	// zmodule.Args["pepper"] = zmodule.Argument{
	// 	Type:    "string",
	// 	Default: "",
	// 	Usage:   "-pepper YOUR_PEPPER",
	// }
	// // Instead, set the pepper by the `cli` command.

	// Other args used only in some modules was declared inner those module.
}

func initCmds() {
	zmodule.Cmds["cli"] = zmodule.Command{
		Usage:   "Enter the command line",
		Handler: func(parsed string, args []string) { client.Start(args) },
	}
}

func init() {
	initInfo()
	initArgs()
	initCmds()
}
