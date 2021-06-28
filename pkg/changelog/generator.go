package changelog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"text/template"

	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"gitlab.com/l0nax/changelog-go/internal"
	"gitlab.com/l0nax/changelog-go/pkg/tools"
)

// GenerateChangelog will generate a new CHANGELOG.md
func GenerateChangelog(r *Release) {
	// load all unreleased Files
	err := GetEntries(r)
	if err != nil {
		log.Fatal(err)
	}

	// prepend our NEW release (data) to the list of releases at index 0
	r.Releases = append([]TplRelease{{
		Version:  r.Info.Version[0],
		Collapse: false,
		Entries:  []TplEntries{},
		Info:     r.Info,
	}}, r.Releases...)

	// get Entries by Change Type
	for _, cType := range internal.EntryT.ListAvailableTypes() {
		log.Debugf("Listing Entries of Type '%s': %# v\n", (*cType).GetShortTypeName(), (*cType).GetListEntries())

		changeEntries := (*cType).GetListEntries()
		if len(changeEntries) == 0 {
			// skip this Change entry type because it has no changes
			continue
		}

		// change exist
		newEntry := TplEntries{
			ShortTypeName: (*cType).GetShortTypeName(),
		}

		log.Debugf("len(changeEntries) before Itoa conversion: '%d'\n", len(changeEntries))

		// write 'n changes' instead of '1 change' if more than one
		if len(changeEntries) > 1 {
			newEntry.NumString = strconv.Itoa(len(changeEntries)) + " changes"
		} else {
			newEntry.NumString = "1 change"
		}

		// add change-entry to the list of changes
		newEntry.Changes = append(newEntry.Changes, changeEntries...)

		log.Debugf("Adding new Entries: %# v\n", pretty.Formatter(newEntry))

		// append the new Change entry to the list of changes
		r.Releases[0].Entries = append(r.Releases[0].Entries, newEntry)
	}

	// debug: print the Releases

	log.Debug("Processing changelog entries")
	rawEntries := ""

	// get the raw changelog output of the new entries
	for _, entries := range r.Releases[0].Entries {
		rawEntries += "\n" + processChangeEntries(&entries)
	}

	// fmt.Println("---------------------------------------------")
	// fmt.Println(rawEntries)
	// fmt.Println("---------------------------------------------")

	// read old releases into Release struct
	err = GetReleasedEntries(r)
	if err != nil {
		log.Fatal(err)
	}

	// sort change entries
	sortReleaseEntries(r)

	// create release dir
	err = createReleaseDir(r.Info.Version[0])
	if err != nil {
		log.Fatal(err)
	}

	// prepare release dir
	err = prepareReleaseDir(*r.Info)
	if err != nil {
		log.Fatal(err)
	}

	err = moveToReleaseFolder(r.Info.Version[0])
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("%# v\n", pretty.Formatter(r))

	data := processChangelogTmpl(r)
	file, err := os.Create(path.Join(internal.GitPath, "CHANGELOG.md"))
	if err != nil {
		log.Fatal(err)
	}

	file.Write([]byte(data))

	file.Close()

	// fmt.Println("+++---------------------------------------------+++")
	// fmt.Println(processChangelogTmpl(r))
	// fmt.Println("+++---------------------------------------------+++")
}

// sortReleaseEntries does sort the changelog entries from `r.[]Releases.Entries`
// by their `ShortTypeName` field.
// Additionally the "raw" change(log) entries (per-change-type) will be sorted.
func sortReleaseEntries(r *Release) {
	// Sort by Version->Change Type
	for i := range r.Releases {
		r.Releases[i].SortEntries()
	}

	// Sort by Version->Change Type->Change Entry
	for i := range r.Releases {
		for j := range r.Releases[i].Entries {
			sortRawChangeLogEntries(&r.Releases[i].Entries[j])
		}
	}
}

func sortRawChangeLogEntries(r *TplEntries) {
	sort.SliceStable(r.Changes, func(i, j int) bool {
		return r.Changes[i].ChangeTitle < r.Changes[j].ChangeTitle
	})
}

// prepareReleaseDir creates all needed files and directory for the release
func prepareReleaseDir(info ReleaseInfo) error {
	basePath := path.Join(internal.GitPath, viper.GetString("changelog.entryPath"))
	releasedPath := path.Join(basePath, "released", info.Version[0])

	// create 'ReleaseInfo' file
	infoFile, err := os.Create(path.Join(releasedPath, "ReleaseInfo"))
	if err != nil {
		return err
	}

	rawData, err := yaml.Marshal(info)
	if err != nil {
		return errors.Wrap(err, "Error while marshaling ReleaseInfo struct")
	}

	// write informations to file
	infoFile.Write(rawData)
	err = infoFile.Close()
	if err != nil {
		return err
	}

	return nil
}

// processChangeEntries generates the RAW output string of every change type
// and returns it as a string.
func processChangeEntries(entries *TplEntries) string {
	var raw string = ""

	// check if field `Changes` is not empty
	if len(entries.Changes) == 0 {
		panic("'entries.Changes' slice is empty!!")
	}

	// first we have to add our Title
	// since the CHANGELOG.md output must look in the "raw" view also pretty
	// so there must be a empty line between every change type header.
	raw += fmt.Sprintf("\n### %s (%s)\n", entries.ShortTypeName, entries.NumString)

	for _, entry := range entries.Changes {
		raw += "\n- " + entry.ChangeTitle
	}

	return raw
}

// processChangelogTmpl generates the CHANGELOG.md via the text template
func processChangelogTmpl(r *Release) string {
	var tmplStr string

	// use default changelog scheme/ template if `changelog.customScheme` is set to false
	if viper.GetBool("changelog.customScheme") {
		fp := viper.GetString("changelog.changelog")
		log.Infof("Custom CHANGELOG.md template has been enabled. File path: %s\n", fp)

		f, err := ioutil.ReadFile(fp)
		if err != nil {
			log.Fatalf("Unable to read custom CHANGELOG.md from '%s': %v\n", fp, err)
		}

		tmplStr = tools.ByteSlice2String(f)
	} else {
		tmplStr = defaultChangelogScheme
	}

	tmpl, err := template.New("changelog").Parse(tmplStr)
	if err != nil {
		log.Fatal(err)
	}

	var tmplOut bytes.Buffer

	err = tmpl.Execute(&tmplOut, r)
	if err != nil {
		log.Fatal(err)
	}

	return tmplOut.String()
}
