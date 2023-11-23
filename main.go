package main

import (
	"fmt"
	"os"
	"flag"
	"path/filepath"

	"gommitizen/version"
)

func main() {
	// Definir los flags para los parámetros de línea de comandos
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	// Flags específicos para cada comando
	initPath := initCmd.String("path", "", "Ruta del repositorio para la inicialización")

	// Verificar si se proporcionó al menos un comando
	if len(os.Args) < 2 {
		fmt.Printf("Uso: %s <comando> [opciones]\n", os.Args[0])
		os.Exit(1)
	}

	// Analizar los parámetros de línea de comandos
	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
		initialize(*initPath)
	case "bump":
		if len(os.Args) == 2 {
			bumpVersion()
		} else if len(os.Args) == 3 {
			bumpProjectVersion(os.Args[2])
		}
	case "help":
		fmt.Println("Uso:")
		fmt.Printf("  %s init -path <ruta>\n", os.Args[0])
		fmt.Printf("  %s bump\n", os.Args[0])
	default:
		fmt.Printf("Comando no reconocido. Ejecuta '%s help' para obtener ayuda.\n", os.Args[0])
		os.Exit(1)
	}
}

func initialize(path string) {
	// Lógica para inicializar el repositorio
	fmt.Printf("Inicializando en la ruta: %s\n", path)
	// ...
}

// Ejecuta el comando bump para un proyecto en concreto
func bumpProjectVersion(project string) {
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio actual:", err)
		os.Exit(1)
	}

	filePath := filepath.Join(rootDir, project, ".version.json")

	if bumpRun(rootDir, filePath) != nil {
		fmt.Println("Error al ejecutar bump:", err)
		os.Exit(1)
	}
}

// Ejecuta el comando bump para todos los archivos .version.json en el directorio actual y sus subdirectorios
func bumpVersion() {
	// Obtén el directorio actual
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio actual:", err)
		os.Exit(1)
	}

	// Encuentra todos los archivos .version.json en el directorio actual y sus subdirectorios
	fileList, err := version.FindFCVersionFiles(rootDir)
	if err != nil {
		fmt.Println("Error al encontrar archivos .version.json:", err)
		os.Exit(1)
	}

	// Itera sobre los archivos encontrados
	for _, filePath := range fileList {
		err := bumpRun(rootDir, filePath)
		if err != nil {
			fmt.Println("Error al ejecutar bump:", err)
			continue
		}
	}
}

// Ejecuta el comando bump para un archivo .version.json
func bumpRun(rootDir string, filePath string) error {
	// Obtiene la ruta relativa al directorio actual
	relativePath, err := filepath.Rel(rootDir, filePath)
	if err != nil {
		fmt.Println("Error al obtener la ruta relativa:", err)
		return err
	}

	// Imprime el mensaje de inicio
	fmt.Printf("\n* Ejecutando bump en proyecto %s\n\n", filepath.Dir(relativePath))

	version := version.VersionData{}

	errData := version.ReadData(filePath)
	if errData != nil {
		fmt.Println("Error al leer los datos de la versión:", errData)
		return errData
	}

	modified, err := version.IsSomeFileModified()
	if err != nil {
		fmt.Println("Error al verificar si algún archivo ha sido modificado en Git:", err)
		return err
	}

	// Si el archivo ha sido modificado, actualiza la versión
	if modified {
		currentVersion := version.GetVersion()
		newVersion, err := version.UpdateVersion()
		if err != nil {
			fmt.Println("Error al actualizar la versión:", err)
			return err
		}

		if newVersion == currentVersion {
			fmt.Printf("No hay actualización de version en %s\n", relativePath)
		} else {
			fmt.Printf("Versión actualizada en %s\n", relativePath)
		}
	} else {
		fmt.Printf("No se realizaron cambios en %s\n", relativePath)
	}
	fmt.Printf("\n")

	return nil
}

