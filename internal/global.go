package internal

var (
	// Path to the top git Directory. This variable does NOT contain the path
	// to the `.git` directory, it contains the path to the Root directory
	// of the project – which contains the `.git` directory.
	GitPath string
	Cwd     string // Path to the current working directory
)
