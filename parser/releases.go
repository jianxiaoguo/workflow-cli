package parser

import (
	"fmt"
	"strconv"

	docopt "github.com/docopt/docopt-go"
	"github.com/drycc/workflow-cli/cmd"
)

// Releases routes releases commands to their specific function.
func Releases(argv []string, cmdr cmd.Commander) error {
	usage := `
Valid commands for releases:

releases:list        list an application's release history
releases:info        print information about a specific release
releases:deploy      deploy the latest release by process types
releases:rollback    return to a previous release

Use 'drycc help [command]' to learn more.
`

	switch argv[0] {
	case "releases:list":
		return releasesList(argv, cmdr)
	case "releases:info":
		return releasesInfo(argv, cmdr)
	case "releases:deploy":
		return releasesDeploy(argv, cmdr)
	case "releases:rollback":
		return releasesRollback(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "releases" {
			argv[0] = "releases:list"
			return releasesList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func releasesList(argv []string, cmdr cmd.Commander) error {
	usage := `
Lists release history for an application.

Usage: drycc releases:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -p --ptypes=<ptypes>
    the processes name as defined in your Procfile, comma separated.
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	results, err := responseLimit(safeGetString(args, "--limit"))
	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")
	ptypes := safeGetString(args, "--ptypes")

	return cmdr.ReleasesList(app, ptypes, results)
}

func releasesInfo(argv []string, cmdr cmd.Commander) error {
	usage := `
Prints info about a particular release.

Usage: drycc releases:info <version> [options]

Arguments:
  <version>
    the release of the application, such as 'v1'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	version, err := versionFromString(args["<version>"].(string))

	if err != nil {
		return err
	}

	app := safeGetString(args, "--app")

	return cmdr.ReleasesInfo(app, version)
}

func releasesDeploy(argv []string, cmdr cmd.Commander) error {
	usage := `
Deploy the latest release by process types.

Usage: drycc releases:deploy [<ptype>...] [options]

Arguments:
  <ptype>
    the process name as defined in your Procfile, such as 'web' or 'web worker'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -f --force
    force deploy.
  --confirm=yes
    to proceed, type "yes".
`

	args, err := docopt.ParseArgs(usage, argv, "")

	if err != nil {
		return err
	}

	apps := safeGetString(args, "--app")
	confirm := safeGetString(args, "--confirm")
	force := args["--force"].(bool)
	return cmdr.ReleasesDeploy(apps, args["<ptype>"].([]string), force, confirm)
}

func releasesRollback(argv []string, cmdr cmd.Commander) error {
	usage := `
Rolls back to a previous application release.

Usage: drycc releases:rollback [<ptype>...] [<version>] [options]

Arguments:
<ptype>
    the process name as defined in your Procfile, such as 'web'.
<version>
    the release of the application, such as 'v1'.

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.ParseArgs(usage, argv, "")
	if err != nil {
		return err
	}

	var version int

	if args["<version>"] == nil {
		version = -1
	} else {
		version, err = versionFromString(args["<version>"].(string))

		if err != nil {
			return err
		}
	}

	app := safeGetString(args, "--app")

	return cmdr.ReleasesRollback(app, args["<ptype>"].([]string), version)
}

func versionFromString(version string) (int, error) {
	if version[:1] == "v" {
		if len(version) < 2 {
			return -1, fmt.Errorf("%s is not in the form 'v#'", version)
		}

		return strconv.Atoi(version[1:])
	}

	return strconv.Atoi(version)
}
