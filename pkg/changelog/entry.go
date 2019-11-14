package changelog

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/l0nax/changelog-go/internal"
	"gitlab.com/l0nax/changelog-go/pkg/entry"
	"gitlab.com/l0nax/changelog-go/pkg/gut"
	"gitlab.com/l0nax/changelog-go/pkg/tools"
	git "gopkg.in/src-d/go-git.v4"
	"os"
	"path"
)

// AddEntry() creates a new Changelog Entry by creating the Entry File
// under (default)
//
// Entryname Scheme: $ROOT/.changelogs/unrelease/$BRANCH_NAME-$RAND_STRING
func AddEntry(entry entry.Entry) {
	// open the Git Repository
	repo, err := git.PlainOpen(internal.GitPath)
	if err != nil {
		log.Fatal(err)
	}

	// get Branchname
	branchName, err := gut.GetCurrentBranchFromRepository(repo)
	if err != nil {
		log.Fatal(err)
	}

	_path := path.Join(internal.GitPath, viper.GetString("changelog.entryPath"))
	_path = path.Join(_path, "unreleased")

	// get random String
	fileName := path.Join(_path, fmt.Sprintf("%s-%s", branchName, tools.RandomString(8)))

	// crete File and write Data down
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Marshal Changelog Type into byte slice
	data, err := tools.YAMLMarshal(entry)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}
