package main

import (
	"encoding/json"
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var ConfigFileType string
var ConfigFileName string
var ConfigDirectory string

type config struct {
	PublicKeyPath  string
	PrivateKeyPath string
	Pacts          map[string][]string
}

var Configuration config

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Generates a new configuration file",
	Long:  `Generates a new configuration file and will refuse to overwrite an existing one.`,
	Run: func(cmd *cobra.Command, args []string) {
		if configFileExists() {
			fmt.Println("Configuration file already exists, refusing to overwrite.")
			os.Exit(400)
		}
		SetDefaultPact()
		HydrateConfigurationModel()
		PersistConfiguration()
	},
}

var VersionCmd = &cobra.Command{
	Use: "version",
	Short: "Outputs the version number to STDOUT",
	Long: `Outputs the version number to STDOUT`,
	Run: func(cmd *cobra.Command, args []string){
		fmt.Println("Pact CLI v" + VERSION)
	},
}
func init() {
	userHomeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	fileSeperator := string(filepath.Separator)

	ConfigDirectory = userHomeDir + fileSeperator + ".config" + fileSeperator + "pact" + fileSeperator
	ConfigFileName = "pact"
	ConfigFileType = "json"

	os.MkdirAll(ConfigDirectory, 0755)

	viper.SetConfigType(ConfigFileType)
	viper.SetConfigName(ConfigFileName)

	viper.AddConfigPath(ConfigDirectory)

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("No config file found. The config command will generate one.")
	} else {
		HydrateConfigurationModel()
	}
}

func HydrateConfigurationModel() {
	viper.Marshal(&Configuration)
}

/**
 * PersistConfigurtaion
 * Method which writes the current configuration model to disk.
 */
func PersistConfiguration() {

	//viper.Marshal(&Configuration)

	configurationString, err := json.MarshalIndent(Configuration, "", "  ")
	if err != nil {
		panic(err)
	}

	f, err := os.Create(GetConfigFilePath())
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(string(configurationString))
}

/**
 * configFileExists
 * @returns bool
 * Heloer method to determine if a configuration file exists.
 */
func configFileExists() bool {
	_, configExistsError := os.Stat(GetConfigFilePath())
	return !os.IsNotExist(configExistsError)
}

/**
 * GetConfigFilePath
 * @returns string Absolute path to the configuration file
 * Helper method which returns the path to the configuration file.
 */
func GetConfigFilePath() string {
	return ConfigDirectory + ConfigFileName + "." + ConfigFileType
}
