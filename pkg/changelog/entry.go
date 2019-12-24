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
	"path/filepath"
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

// GetReleasedEntries returns a slice which contains all releaed entries
func GetReleasedEntries(r *Release) error {
	// we need to APPEND all new Release data!

	basePath := path.Join(internal.GitPath, viper.GetString("changelog.entryPath"))
	releasedPath := path.Join(basePath, "released")

	err := filepath.Walk(releasedPath, func(_path string, info os.FileInfo, err error) error {
		log.Debugf("===> %s\n", _path)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// sortEntries will sort all changelog entries into TplEntries struct
func sortEntries(entries []*entry.Entry) ([]TplEntries, error) {
	var ret = []TplEntries{}
	var typeMap = make(map[int]int) // typeMap contains the map of TypeID-->Index_in_'ret'

	for i, entry := range entries {
		// check if entry type does not exists in typeMap
		if iRet, ok := typeMap[entry.TypeID]; !ok {
			// add entry-type to 'ret'
			ret = append(ret, TplEntries{
				ShortTypeName: (*entry.Type).GetShortTypeName(),
				NumString:     "1 change",
			})

			ret[len(ret)-1].Changes = append(ret[len(ret)-1].Changes, entries[i])

			// add index of ret to typeMap
			typeMap[entry.TypeID] = len(ret)
		} else {
			// add entry to list
			ret[iRet].Changes = append(ret[iRet].Changes, entries[i])
		}
	}

	return ret, nil
}
