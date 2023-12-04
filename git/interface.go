package git

type GitI interface {
	SetDirPath(dirPath string)
	SetFromCommit(fromCommit string)
	SetFilterFiles(excludedFiles []string)
	GetDirPath() string
	GetChangedFiles() []string
	GetCommitMessages() []string
	GetFromCommit() string
	GetLastCommit() string
	Initialize(gitPath string) error
	RetrieveLastCommit() error
	RetrieveData() error
	ConfirmChanges(files []string, commitMessage string, tagMessage string) error
	GetOutput() []string
	CleanOutput()
	Add(file string) error
	Commit(commitMessage string) error
	Tag(tagMessage string) error
}
