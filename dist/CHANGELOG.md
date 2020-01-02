

## Changelog

fd6c6b3 **BACKUP**   
3bc06a6 //WIP: Add file.go for Changelog-File handling   
3ae31de Add "template" version.go   
434bf39 Add 'CHANGELOG.md' file creation   
f6af901 Add 'Fixed' Change-Type to default Types   
8206718 Add 'GetCurrentBranchFromRepository' Function   
eb62cac Add 'Info' field to TplRelease struct   
9a88a81 Add 'changelog-go' to ignore List   
bb08dbc Add 'colapse' Functionality to default Template   
0f032fd Add 'customScheme' Option to example Config   
9b054ae Add 'debug' flag   
4a8e9d1 Add 'entry.author' Config to Example Config   
dcefefb Add 'gitconfig' Package to dependencies   
411f227 Add 'github.com/blang/semver' to dependencies   
2383171 Add 'github.com/speedata/gogit' Package as dependencie   
0934da4 Add 'go-git' Package as dependencie   
550dd56 Add 'go-git' Packages as dependencie   
1a31f68 Add 'gopkg.in/yaml.v2' to dependencies   
0103051 Add 'init' command   
2a7291b Add 'logrus' Package to glide.yaml   
e551abe Add 'pre-release' Flag to 'release' Command   
23e5138 Add 'release-date' value to struct   
01f9b2b Add 'script' installation documentation   
096be62 Add 'title' Argument to 'new' Subcommand   
67654e9 Add 'version' Flag to rootCmd   
458a975 Add (Un)Marshal Functions to tools Package   
b0b91e1 Add (first) changelog entry   
4388bb4 Add (first) generated CHANGELOG.md   
b3cf70a Add Changelog Entry Base Interface and Struct   
9a3827a Add Cobra generated Go Files   
f36fac0 Add CodeQuery artefacts to gitignore   
107d59f Add EditorConfig   
41c6b38 Add Error Handling when getting Branchname   
a5ed0fb Add Example & Description for 'release' Command   
ab52fe9 Add GenerateChangelog function   
0e22398 Add GetAuthorName() to gut Package   
62de2f8 Add IDEA Files   
b254ac1 Add IDEA Files   
8b95cc1 Add IDEA Workspace   
b4f83c7 Add ListAvailableTypes() for EntryTypes struct   
88c359e Add Makefile   
fc76e9a Add Marshal Functions to create Entry File   
698d73d Add Project changelog-go Config   
f128dc8 Add ReleaseInfo and Release type structs   
3815b4a Add Search of Paths to cmd/root.go   
d72c3bc Add Tags for Changelog Entry File   
7767859 Add Test Code to root Cmd   
6527952 Add TplEntries struct and extend TplRelease struct   
9a0eb06 Add assets generation to Makefile   
a882d78 Add basic Git Ignore   
fa48046 Add basic/ default Changelog Entries   
8888fdf Add changelog Package   
3c5db06 Add cmd/new.go Skeleton   
8d0dc9d Add cmd/release.go Skeleton   
2d53956 Add constant maxDepth in gut Package   
83343b7 Add debug Output to SearchEntryType()   
34efbfe Add debug message to GetListEntries() func   
679fb19 Add debug messages in tools pkg   
b5ab86d Add debug msg if is Pre-Release to 'new' cmd   
052800a Add debug msg to sortEntries() func   
22d0169 Add default Changelog Scheme Template   
94d049f Add documentation to README.md   
b86304d Add example Configuration for pre-release   
322fd27 Add field 'ReleaseDate' to ReleaseInfo struct   
ecc8a31 Add func CheckDirExist() to tools pkg   
168b559 Add generated 'assets.go' to gitignore   
4990257 Add git Library   
5399950 Add glide Config and Lock File   
79d037c Add goreleaser Config   
1255b4e Add internal Global Variables   
4e4d19a Add library 'github.com/kr/pretty'   
9ca1d01 Add log messages in changelog pkg   
2051afc Add more Detailed Comment to ChangeEntry interface   
f89069f Add more documentation to GitPath variable   
07f97b0 Add new Command   
39e1d37 Add new Package 'config'   
9c5ff29 Add new Package 'gut'   
e06d7a9 Add new TODO Comments   
423de17 Add new flag 'release-date' to release cmd   
5ccac7b Add prepareReleaseDir() function   
248ea57 Add processChangelogTmpl() function   
69da49a Add skeleton of GetReleasedEntries() func   
d64da8c Add tags to fields of struct ReleaseInfo   
aefdada Add two debug messages before exiting   
b1ce9aa Add version Entry to Example Config   
62a1190 Allocate memory for files slice in changelog pkg   
61af2d4 Call 'Unmarshal' directly without tools pkg   
9dec2ed Call RegisterAll() in initConfig()   
535d69b Change Entry Signature of entry Pkg   
6e6000c Change Format Type to decimal   
0b87c6f Change Function Signature out to Pointer   
f2ef416 Change Singature of Fnction AddEntry()   
467676d Change check for max Depth in Config Search   
b295b6c Change field name 'Collapse' in default Scheme   
da5de18 Change fieldname to 'Collapse'   
048df23 Change function signature of prepareReleaseDir()   
975071d Change search criterias of SearchEntryType()   
b613f1e Change signature of function GetReleasedEntries()   
3c7942f Check if path is Directory in GetEntries()   
5770d0b Clarify Fatal message in 'new' cmd   
58aa624 Comment unnecessary lines out changelog pkg   
24a2a29 Comment unused var out   
0150266 Correct 'package' misstake in tools Pkg   
2a5a7d5 Correct Type Registering in internal Pkg   
9fdd092 Correct spelling in README.md   
8a375a5 Correct spellingmisstake in interfaces.go   
18fe03a Correct spellingmisstake in interfaces.go   
382dede Correct var documentation in internal pkg   
b4294f4 Create skeleton tools.go in changelog Package   
7f26f37 Describe Fields of SEntryType struct   
d9b09c2 Document 'filepath.Walk' in changelog pkg   
6d8cfe8 Document case if both vars are set to true in example Conf   
1b1f5fa Fill out 'Info' field of current release   
9e14ea0 Fix 'YAMLUnmarshal' function in tools pkg   
ab8b63f Fix 'config not found' Bug   
c1a91ec Fix 'panic: runtime error: index out of range'   
09e943d Fix Mkdir Bug   
578aea7 Fix Version string problems in GenerateChangelog()   
b4e871f Fix goreleaser problems   
ce5866e Fix latest release overwritting in GenerateChangelog()   
2d1eb59 Fix problems with 'init' Command   
219afdf Fix spelling misstake in 'Flags().BoolP'   
16adefd Fix spelling misstake in example Config   
285e956 Fix variable name misspelling   
22e884c Fix wrong Path creation in AddEntry()   
729baf0 Implement 'AddEntry' Functionality   
4b25b55 Implement 'ReleaseInfo' file parsing to GetReleasedEntries()   
ae617fa Implement 'deletePreRelease' config var into GetReleasedEntries() func   
33faef7 Implement 'new' Command Functionality   
ebdb815 Implement 'release' target in Makefile   
d22b356 Implement CheckDir() Function into changelog Pkg   
14bc067 Implement Function to generate Random String   
ef35cc2 Implement GetEntries() Function in tools Pkg   
e0bda44 Implement GetReleasedEntries() function   
5c33db2 Implement MoveEntries() Function in changelog Pkg   
979ed15 Implement functionality of release Command   
7afd839 Implement incremental change-count   
5a8c3dd Implement moveToReleaseFolder()   
ec7512e Implement processChangeEntries() in changelog pkg   
803b4fa Implement reading&processing of GenerateChangelog   
5c6b4c0 Implement sortEntries() func in changelog pkg   
4e31e7a Improve log message in SearchEntryType()   
9361f18 Improve long help of rootCmd   
e1ae184 Make debug print more identifiable   
37ca460 Merge branch 'release/1.0.0'   
0a4a8b9 Migrate from glide to go modules   
f04822a Move Function Documentation out of struct   
2865939 Move default Entry Types into Sub-Pacakge entries   
03cc3a6 Omit some Fields if emtpy of Entry struct   
41aeace Print Output to give User better UX   
b8bdda2 Print available change types better   
a0427c7 Print only one new-line per Change-Type   
86dc02c Reformat Code in basicEntries.go File   
64c6aa9 Reformat Code in entry.go   
4db7556 Reformat Code in interfaces.go File   
34c29bf Reformat Files   
c49557a Reformat entry.go   
a0b84b2 Register all other Entry Types   
23dae82 Remove '.idea' Directory   
b8371b0 Remove 'Example Workflow' section from README.md   
499108c Remove 'go-git.v4' Package dependencie   
cc5e186 Remove TODO Comment from AddEntry()   
8321470 Remove additional Variable to improve Performance   
20550f8 Remove completed TODO comment from generator.go   
8099f5f Remove old glide Files   
18e83e9 Remove uncommented & deprecated Code from AddEntry()   
7e9e965 Remove unneded Code   
bbdbed6 Remove unneded variable decleration/ usage in processChangelogTmpl()   
5f2e550 Remove unused imports in tools pkg   
c70f918 Revert "Add 'github.com/speedata/gogit' Package as dependencie"   
41a9d94 Revert "Add 'go-git' Package as dependencie"   
b9f3e20 Revert "Call 'Unmarshal' directly without tools pkg"   
617b648 Set 'entry.author' DataType to Boolean   
f09b3b2 Sort imports in entry.go   
388e64a Sort releases so we have the versions listed correctly   
770747e Split GetEntries() function into two new ones   
787ace3 Uncomment non-working Code in tools pkg   
45cc61d Uncomment old & unstable Code from AddEntry()   
070f908 Update 'Info' field comment   
a61962b Update 'release' cmd Example   
4044f02 Update Function Signatures of basic Types   
ca7404b Update GenerateChangelog() algorithm   
e5f5e0d Update TODO comment in GetReleasedEntries()   
241158b Update and re-write (default) Changelog Tmpl   
38c239b Update dependencies   
c5599db Update describption of TplRelease and TplEntries   
36ddce0 Update go.mod and go.sum   
013eb96 Update regex checks in release cmd   
2a1cac8 Use MkdirAll in MoveEntries() to prevent errors   
e7d3858 Use Short() instead of String() in GetBranchName()

