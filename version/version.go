package version

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"

	"gommitizen/git"
)

// Tipo de error personalizado
type VersionError struct {
	Message string
}

func (e *VersionError) Error() string {
	return e.Message
}

// Gestiona la información de la versión para nuestro proyecto
type VersionData struct {
	Version string `json:"version"`
	Commit string `json:"commit"`
	filePath string
}

// Métodos getter
func (version *VersionData) GetVersion() string {
	return version.Version
}

func (version *VersionData) GetCommit() string {
	return version.Commit
}

func (version *VersionData) GetFilePath() string {
	return version.filePath
}

// Funciones públicas

// Busca archivos .version.json en un directorio dado y sus subdirectorios
func FindFCVersionFiles(rootDir string) ([]string, error) {
	var fileList []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".version.json") {
			fileList = append(fileList, path)
		}

		return nil
	})

	return fileList, err
}

// Métodos públicos

// Obtene los valores de la versión y el commit del archivo .version.json
func (version *VersionData) ReadData(filePath string) error {
	version.filePath = filePath

	ver, errVersion := version.getCurrentVersionFromJsonFile()
	if errVersion != nil {
		return errVersion
	}
	version.Version = ver

	commit, errCommit := version.getCommitValueFromJsonFile()
	if errCommit != nil {
		return errCommit
	}
	version.Commit = commit

	return nil
}

// Devuelve si algún archivo ha sido modificado en Git desde un commit dado en un directorio dado
func (version *VersionData) IsSomeFileModified() (bool, error) {
	if version.Commit == "" || version.filePath == "" {
		return false, &VersionError{
			Message: "No se ha especificado un commit o un archivo .version.json",
		}
	}

	// Obtén el directorio actual
	currentDir, err := os.Getwd()
	if err != nil {
		return false, &VersionError{
			Message: "Error al obtener el directorio actual",
		}
	}

	// Obtiene la ruta relativa al directorio actual
	relativePath, err := filepath.Rel(currentDir, version.filePath)
	if err != nil {
		return false, &VersionError{
			Message: "Error al obtener la ruta relativa",
		}
	}

	// Obtiene la ruta base del archivo
	dirPath := filepath.Dir(relativePath)

	// Obtiene la lista de archivos modificados en Git desde un commit dado en un directorio dado
	git := git.Git{
		DirPath: dirPath,
		FromCommit: version.Commit,
	}
	errUpdate := git.UpdateData()
	if errUpdate != nil {
		return false, &VersionError{
			Message: "Error al actualizar los datos de Git",
		}
	}
	changedFiles := git.GetChangedFiles()

	// Verifica si la lista de archivos modificados está vacía
	return len(changedFiles) > 0, nil
}


// Actualiza el valor de la versión en el archivo .version.json en función de los cambios en Git
func (version *VersionData) UpdateVersion() (string, error) {
	if version.Version == "" || version.Commit == "" {
		return "", &VersionError{
			Message: "Error: no se han especificado los valores de la versión y el commit",
		}
	}

	// Obtén el directorio actual
	currentDir, err := os.Getwd()
	if err != nil {
		return "", &VersionError{
			Message: "Error al obtener el directorio actual",
		}
	}

	// Obtiene la ruta relativa al directorio actual
	relativePath, err := filepath.Rel(currentDir, version.filePath)
	if err != nil {
		return "", &VersionError{
			Message: "Error al obtener la ruta relativa",
		}
	}

	// Obtiene la ruta base del archivo
	dirPath := filepath.Dir(relativePath)

	// Crea una instancia de Git
	git := git.Git{
		DirPath: dirPath,
		FromCommit: version.Commit,
	}

	// Actualiza los datos de Git
	errUpdate := git.UpdateData()
	if errUpdate != nil {
		fmt.Println("Error al actualizar los datos de Git:", errUpdate)
		return "", errUpdate
	}

	// Determina el tipo de incremento de versión en función de los mensajes de confirmación
	incType := determineVersionBump(git.CommitMessages)

	// Incrementa la versión actual
	currentVersion, newVersion, err := incrementVersion(version.Version, incType)
	if err != nil {
		fmt.Println("Error al incrementar la versión actual:", err)
		return "", err
	}

	if incType != "none" {
		// Informamos del incremento de versión, actualizamos el valor de la versión y el commit y actualizamos Git
		fmt.Println("Incrementando la versión de " + currentVersion + " a " + newVersion)
	  version.Commit = git.LastCommit
		version.Version = newVersion

		// Serializa la estructura actualizada de nuevo en JSON
		updatedContent, err := json.MarshalIndent(version, "", "  ")
		if err != nil {
			fmt.Println("Error al serializar la estructura actualizada:", err)
			return "", err
		}

		// Escribe el contenido actualizado en el archivo
		err = ioutil.WriteFile(version.filePath, updatedContent, os.ModePerm)
		if err != nil {
			fmt.Println("Error al escribir el contenido actualizado en el archivo:", err)
			return "", err
		}

		addFiles := []string{relativePath}
		commitMessage := "Versión actualizada (" + version.Version + ") en " + getBaseDirFromFilePath(git.DirPath)
		tagMessage := version.Version + "_" + getBaseDirFromFilePath(git.DirPath)
		output, err := git.UpdateGit(addFiles, commitMessage, tagMessage)
		if err != nil {
			fmt.Println("Error al actualizar Git:", err)
			return "", err
		}

		for _, file := range output {
			fmt.Println(file)
		}

	} else {
		fmt.Println("No se incrementa la versión: " + currentVersion)
		version.Version = currentVersion
	}

	return version.Version, nil
}

// Métodos privados

// Obtiene el commit almacenado en el archivo .version.json
func (version *VersionData) getCommitValueFromJsonFile() (string, error) {
	// Lee el contenido del archivo .version.json
	content, err := ioutil.ReadFile(version.filePath)
	if err != nil {
		fmt.Println("Error al leer el contenido del archivo:", err)
		return "", err
	}

	// Deserializa el contenido en una estructura Version
	err = json.Unmarshal(content, version)
	if err != nil {
		fmt.Println("Error al deserializar el contenido del archivo:", err)
		return "", err
	}

	// Devuelve el valor del commit
	return version.Commit, nil
}

// Obtiene el valir de commit almacenado en el archivo .version.json
func (version *VersionData) getCurrentVersionFromJsonFile() (string, error) {
	// Lee el contenido del archivo .version.json
	content, err := ioutil.ReadFile(version.filePath)
	if err != nil {
		fmt.Println("Error al leer el contenido del archivo:", err)
		return "", err
	}

	// Deserializa el contenido en una estructura Version
	err = json.Unmarshal(content, version)
	if err != nil {
		fmt.Println("Error al deserializar el contenido del archivo:", err)
		return "", err
	}

	// Devuelve el valor de la versión
	return version.Version, nil
}

// Funciones auxiliares privadas

// Obtiene el directorio base de un archivo dado
func getBaseDirFromFilePath(filePath string) string {
	return filepath.Base(filePath)
}

// Determina el tipo de incremento de versión en función de los mensajes de confirmación
func determineVersionBump(commitMessages []string) string {
	for _, message := range commitMessages {
		// Un mensaje contiene al inicio de la cadena dada el siguiente prefijo "feat:", "fix:" o "BREAKING CHANGE:"
		if strings.Contains(message, "BREAKING CHANGE:") {
			return "major"
		} else if strings.Contains(message, "feat:") {
			return "minor"
		} else if strings.Contains(message, "fix:") {
			return "patch"
		}
	}

	return "none"
}

// Incrementa la versión actual en función del tipo de incremento dado y devuelve la nueva versión 
func incrementVersion(version string, incType string) (string, string , error) {
	currentVersion, err := semver.NewVersion(version)
	if err != nil {
		return "", "", err
	}
	
	var newVersion semver.Version
	if incType == "major" {
		newVersion = currentVersion.IncMajor() // Incrementa el mayor (por ejemplo, de 1.2.3 a 2.0.0)
	} else if incType == "minor" {
		newVersion = currentVersion.IncMinor() // Incrementa el menor (por ejemplo, de 1.2.3 a 1.3.0)
	} else if incType == "patch" {
		newVersion = currentVersion.IncPatch() // Incrementa el parche (por ejemplo, de 1.2.3 a 1.2.4)
	} else {
		newVersion = *currentVersion
	}

	return currentVersion.String(), newVersion.String(), nil
}
