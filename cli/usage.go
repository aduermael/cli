package cli

// Command is the struct containing the command name and description
type Command struct {
	Name        string
	Description string
}

// DockerCommandUsage lists the top level docker commands and their short usage
var DockerCommandUsage = []Command{
	{"cp", "Copy files/folders between a container and the local filesystem"},
	{"exec", "Run a command in a running container"},
	{"info", "Display system-wide information"},
	{"inspect", "Return low-level information on a container or image"},
	{"load", "Load an image from a tar archive or STDIN"},
	{"login", "Log in to a Docker registry"},
	{"logout", "Log out from a Docker registry"},
	{"ps", "List containers"},
	{"pull", "Pull an image or a repository from a registry"},
	{"push", "Push an image or a repository to a registry"},
	{"save", "Save one or more images to a tar archive"},
	{"update", "Update configuration of one or more containers"},
}

// DockerCommands stores all the docker command
var DockerCommands = make(map[string]Command)

func init() {
	for _, cmd := range DockerCommandUsage {
		DockerCommands[cmd.Name] = cmd
	}
}
