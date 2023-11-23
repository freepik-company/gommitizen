package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gommitizen/version"
	"os"
	"path/filepath"
)

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Genera el tag de versión",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		if len(args) == 0 {
			bumpVersion()
			return
		}

		if len(args) == 2 {
			bumpProjectVersion(args[1])
		}
	},
}

func init() {
	rootCmd.AddCommand(bumpCmd)
}

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

	config := version.VersionData{}

	errData := config.ReadData(filePath)
	if errData != nil {
		fmt.Println("Error al leer los datos de la versión:", errData)
		return errData
	}

	modified, err := config.IsSomeFileModified()
	if err != nil {
		fmt.Println("Error al verificar si algún archivo ha sido modificado en Git:", err)
		return err
	}

	// Si el archivo ha sido modificado, actualiza la versión
	if modified {
		currentVersion := config.GetVersion()
		newVersion, err := config.UpdateVersion()
		if err != nil {
			fmt.Println("Error al actualizar la versión:", err)
			return err
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
