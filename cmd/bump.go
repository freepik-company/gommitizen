package cmd

import (
	"fmt"
	"gommitizen/version"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var projectDir string
var changelog bool
var incrementType string

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Genera el tag de versión",
	//Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if projectDir == "" {
			fmt.Printf("\n* Ejecutando bump en todos los proyectos\n")
			bumpVersion()
			return
		} else {
			fmt.Printf("\n* Ejecutando bump en el proyecto %s\n", projectDir)
			bumpProjectVersion(projectDir)
			return
		}
	},
}

func init() {
	bumpCmd.Flags().StringVarP(&projectDir, "directory", "d", "", "Select a project directory to bump")
	bumpCmd.Flags().BoolVarP(&changelog, "changelog", "c", false, "Create CHANGELOG.md")
	bumpCmd.Flags().StringVar(&incrementType, "increment", "", "Version increment type (MINOR, MAJOR, PATCH)")

	rootCmd.AddCommand(bumpCmd)
}

func bumpProjectVersion(project string) {
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio actual:", err)
		os.Exit(1)
	}

	filePath := filepath.Join(rootDir, project, ".version.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("No se encontró el archivo %s\n", filePath)
		os.Exit(1)
	}

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

	if len(fileList) == 0 {
		fmt.Println("No se encontraron archivos .version.json")
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
		return fmt.Errorf("Error al obtener la ruta relativa: %s", err)
	}

	// Imprime el mensaje de inicio
	fmt.Printf("\n* Ejecutando bump en proyecto %s\n\n", filepath.Dir(relativePath))

	config := version.VersionData{}

	errData := config.ReadData(filePath)
	if errData != nil {
		return fmt.Errorf("Error al leer los datos de la versión: %s", errData)
	}

	modified, err := config.IsSomeFileModified()
	if err != nil {
		return fmt.Errorf("Error al verificar si algún archivo ha sido modificado en Git: %s", err)
	}

	// Si el archivo ha sido modificado, actualiza la versión
	if modified {
		currentVersion := config.GetVersion()
		newVersion, err := config.UpdateVersion()
		if err != nil {
			return fmt.Errorf("Error al actualizar la versión: %s", err)
		}

		if newVersion == currentVersion {
			fmt.Printf("No hay actualización de config en %s\n", relativePath)
		} else {
			fmt.Printf("Versión actualizada en %s\n", relativePath)
		}
	} else {
		fmt.Printf("No se realizaron cambios en %s\n", relativePath)
	}
	fmt.Printf("\n")

	return nil
}
