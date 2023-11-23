package git

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// Gestiona los errores personalizados
type GitError struct {
	Message string
}

func (e *GitError) Error() string {
	return e.Message
}

// Controla la información de Git para nuestro proyecto
type Git struct {
	DirPath        string
	FromCommit     string
	LastCommit     string
	ChangedFiles   []string
	CommitMessages []string
	ExcludedFiles  []string
}

// funciones setter
func (git *Git) SetDirPath(dirPath string) {
	git.DirPath = dirPath
}

func (git *Git) SetFromCommit(fromCommit string) {
	git.FromCommit = fromCommit
}

func (git *Git) setExcludedFiles(excludedFiles []string) {
	for i := 0; i < len(excludedFiles); i++ {
		git.ExcludedFiles = append(git.ExcludedFiles, filepath.Join(git.DirPath, excludedFiles[i]))
	}
}

func (git *Git) setLastCommit() error {
	lastCommit, err := getLastCommitFromGit()
	if err != nil {
		return err
	}

	git.LastCommit = lastCommit

	return nil
}

// Funciones getter
func (git *Git) GetChangedFiles() []string {
	return git.ChangedFiles
}

func (git *Git) GetCommitMessages() []string {
	return git.CommitMessages
}

func (git *Git) GetFromCommit() string {
	return git.FromCommit
}

func (git *Git) GetLastCommit() string {
	return git.LastCommit
}

func (git *Git) GetExcludedFiles() []string {
	return git.ExcludedFiles
}

// Métodos públicos

// Obtiene la lista de archivos modificados en Git desde un commit dado en un directorio dado y los almacena en el atributo ChangedFiles
// También recupera los mensajes de commit y los almacena en el atributo CommitMessages.
func (git *Git) UpdateData() error {
	if git.DirPath == "" {
		return &GitError{
			Message: "Error: no se ha especificado el directorio de trabajo",
		}
	}

	if git.FromCommit == "" {
		return &GitError{
			Message: "Error: no se ha especificado el valor del commit",
		}
	}

	changedFiles, err := git.getListOfModifiedFilesInGitFromAGivenCommitInDirExcludingFiles()
	if err != nil {
		return err
	}
	git.ChangedFiles = changedFiles

	commitMesages, err := git.getCommitMessages()
	if err != nil {
		return err
	}
	git.CommitMessages = commitMesages

	lastCommit, err := getLastCommitFromGit()
	if err != nil {
		return err
	}
	git.LastCommit = lastCommit

	return nil
}

// Actualiza los datos de Git
func (git *Git) UpdateGit(files []string, commitMessage string, tagMessage string) ([]string, error) {
	output := []string{}

	for _, file := range files {
		outputAdd, errAdd := add(file)
		if errAdd != nil {
			return nil, errAdd
		}
		output = append(output, outputAdd)
	}

	outputCommit, errCommit := commit(commitMessage)
	if errCommit != nil {
		return nil, errCommit
	}
	output = append(output, outputCommit)

	outputTag, errTag := tag(tagMessage)
	if errTag != nil {
		return nil, errTag
	}
	output = append(output, outputTag)

	return output, nil
}

// Métodos privados

// Obtiene la lista de archivos modificados en Git desde un commit dado en un directorio dado
// Permite excluir archivos de la lista de archivos modificados
func (git *Git) getListOfModifiedFilesInGitFromAGivenCommitInDirExcludingFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", git.FromCommit, "HEAD", git.DirPath)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Divide la salida en líneas
	lines := strings.Split(string(output), "\n")

	// Elimina la última línea vacía
	lines = lines[:len(lines)-1]

	// Elimina los archivos excluidos de la lista de archivos modificados
	for _, excludeFile := range git.ExcludedFiles {
		lines = removeStringFromSlice(lines, excludeFile)
	}

	return lines, nil
}

// Obtiene los mensajes de commit para los archivos modificados en Git desde un commit dado en un directorio dado
func (git *Git) getCommitMessages() ([]string, error) {
	// Construir el comando git log con opciones para obtener mensajes y archivos modificados
	args := append([]string{"log", "--pretty=%s", "--name-only", git.FromCommit + "..", "--"}, git.ChangedFiles...)
	cmd := exec.Command("git", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Dividir la salida en líneas y eliminar líneas vacías
	lines := strings.Split(string(output), "\n")
	var CommitMessages []string

	// Iterar sobre las líneas para construir la lista de mensajes de commit
	for i := 0; i < len(lines); i++ {
		message := lines[i]
		if strings.TrimSpace(message) != "" {
			// Agregar el mensaje a la lista
			CommitMessages = append(CommitMessages, message)

			// Avanzar al próximo mensaje
			for i < len(lines) && strings.TrimSpace(lines[i]) != "" {
				i++
			}
		}
	}

	return CommitMessages, nil
}

// Funciones privadas

// Añade un archivo al repositorio Git
func add(filePath string) (string, error) {
	cmd := exec.Command("git", "add", filePath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// Crea un nuevo commit en Git
func commit(message string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", message)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// Crea un nuevo tag en Git
func tag(tag string) (string, error) {
	cmd := exec.Command("git", "tag", tag)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// Obtiene el commit actual en Git
func getLastCommitFromGit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// Funciones auxiliares

// Elimina una cadena de texto de un slice de cadenas
func removeStringFromSlice(slice []string, s string) []string {
	var result []string

	for _, str := range slice {
		if str != s {
			result = append(result, str)
		}
	}

	return result
}
