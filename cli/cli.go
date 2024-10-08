package cli

// Shortcuts is a map of all the shortcuts supported by the CLI
var Shortcuts = map[string]string{
	"create":         "apps:create",
	"destroy":        "apps:destroy",
	"info":           "apps:info",
	"login":          "auth:login",
	"logout":         "auth:logout",
	"logs":           "apps:logs",
	"open":           "apps:open",
	"pull":           "builds:create",
	"rollback":       "releases:rollback",
	"run":            "apps:run",
	"scale":          "pts:scale",
	"sharing":        "perms:list",
	"sharing:list":   "perms:list",
	"sharing:add":    "perms:add",
	"sharing:remove": "perms:remove",
	"whoami":         "auth:whoami",
}
