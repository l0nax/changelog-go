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
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
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

// GetReleasedEntries returns a slice which contains all released entries
func GetReleasedEntries(r *Release) error {
	// we need to APPEND all new Release data!

	basePath := path.Join(internal.GitPath, viper.GetString("changelog.entryPath"))
	releasedPath := path.Join(basePath, "released")

	err := filepath.Walk(releasedPath, func(_path string, info os.FileInfo, err error) error {
		// skip if entity is not directory or start path
		if !info.IsDir() || _path == releasedPath {
			return nil
		}

		// create skeleton struct
		release := TplRelease{
			Info: &ReleaseInfo{
				// TODO: IsPreRelease can be checked with the Regex
			},
			Version:  info.Name(),
			Collapse: false,
			Entries:  []TplEntries{},
		}

		// get and parse all files from this directory
		files, err := ReadEntryFiles(_path)
		if err != nil {
			return err
		}

		// get ReleaseInfo file and remove it from the map
		releaseInfo, ok := files[path.Join(_path, "ReleaseInfo")]
		if !ok {
			log.Fatalf("No 'ReleaseInfo' file was found...!\n")
		}

		delete(files, path.Join(_path, "ReleaseInfo"))

		_releaseInfo := &ReleaseInfo{}
		err = yaml.Unmarshal(releaseInfo, _releaseInfo)
		if err != nil {
			return err
		}

		release.Info.ReleaseDate = _releaseInfo.ReleaseDate
		_releaseInfo = nil

		// parse changelog-entry files
		entries, err := ParseFiles(files)
		if err != nil {
			return err
		}

		// get change-type
		for i, entry := range entries {
			if entry.Type != nil {
				// we can skip this entry because the Type is
				// already set
				continue
			}

			entries[i].Type, err = internal.EntryT.SearchEntryType(&internal.SEntryType{
				TypeID: entry.TypeID,
			})

			if err != nil {
				return err
			}
		}

		// "sort" all entries, this means that all entries will be added
		// to the TplEntries struct
		release.Entries, err = sortEntries(&entries)
		if err != nil {
			return err
		}

		// append "new" old release to *r
		r.Releases = append(r.Releases, release)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// sortEntries will sort all changelog entries into TplEntries struct
func sortEntries(entries *[]entry.Entry) ([]TplEntries, error) {
	var ret = []TplEntries{}
	var typeMap = make(map[int]int) // typeMap contains the map of TypeID-->Index_in_'ret'

	for i, entry := range *entries {
		// check if entry type does not exists in typeMap
		if iRet, ok := typeMap[entry.TypeID]; !ok {
			// add entry-type to 'ret'
			ret = append(ret, TplEntries{
				ShortTypeName: (*entry.Type).GetShortTypeName(),
				NumString:     "1 change",
			})

			ret[len(ret)-1].Changes = append(ret[len(ret)-1].Changes, &(*entries)[i])

			// add index of ret to typeMap
			typeMap[entry.TypeID] = len(ret)
		} else {
			// add entry to list
			ret[iRet].Changes = append(ret[iRet].Changes, &(*entries)[i])

			// increase NumString
			num, _ := strconv.Atoi(strings.Split(ret[iRet].NumString, " ")[0])
			num++

			ret[iRet].NumString = string(num) + " changes"
		}
	}

	return ret, nil
}
