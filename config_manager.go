/********
 * config_manager
 *
 * This file provides all the implementation necessary to read and write INI style
 * configurationManager files and to intelligently accept all configurationManager options as command
 * line flags. Getters are provided to retrieve configurationManager data.
 *
 * Reading configurationManager Options:
 *  When configurationManager data is read in from a file, any file paths present in the
 *  configurationManager values will be expanded relative to the location of said configurationManager
 *  file. When configurationManager data containing paths is passed in via the command line
 *  paths will be expanded relative to the current working directory. Command
 *  line flags take precedence and override any values loaded from a configurationManager file.
 *
 * Writing configurationManager Options:
 *  For convenience a configurationManager wizard is implemented which will allow the user to
 *  create a configurationManager file interactively. Any values containing a file path will
 *  expand that path relative to the current working directory. As a result the generated
 *  configurationManager file will contain only absolute paths.
 *
 * Searching For a Config File:
 *  When this module is loaded it attempts to load a configurationManager file, if the -config
 *  flag was passed in and the file exists it will be loaded, otherwise the program will
 *  begin searching up through the current directory hierarchy until it finds one. If it
 *  is unable to find one it explicitly looks in ~/.config/pact/. If it still does not
 *  find a configurationManager file, a warning will be output suggesting the user generate one.
 *  It is possible to pass all configurationManager options with command line flags and avoid
 *  using a configurationManager file.
 */

package main

import (
	"flag"
	"github.com/mitchellh/cli"
	"github.com/rakyll/globalconf"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// In the event the config flag isn't passed this is the filename that will be searched for
const CONFIG_FILE_NAME = ".pact"
const FILE_SEPERATOR = string(filepath.Separator)

var currentWorkingDirectory string
var userHomeDir string
var defaultGlobalConfigPath string

var flagBasePath string
var configFilePath string

// configurationManager Values
var naclPublicKey string
var naclPrivateKey string

// Globalconf is used to intelligently merge flags and INI config values as well as
// persist changes to disk
var configurationManager *globalconf.GlobalConf
var ui = &cli.BasicUi{
	Writer: os.Stdout,
	Reader: os.Stdin,
}

//////////////// PUBLIC GETTERS /////////////////////

func GetConfigFilePath() string {
	return configFilePath
}

func GetNACLPrivateKeyPath() string {
	return naclPrivateKey
}

func GetNACLPublicKeyPath() string {
	return naclPublicKey
}

/**
 * init
 * Initialize flags and set helper variables like currentWorkingDirectory and userHomeDir.
 * Search for a configurationManager file if necessary and correctly parse all values found.
 */
func init() {
	currentUser, _ := user.Current()
	userHomeDir = currentUser.HomeDir
	defaultGlobalConfigPath = userHomeDir + FILE_SEPERATOR + ".config" + FILE_SEPERATOR + "pact"

	flag.StringVar(&configFilePath, "config", "", "What is the path to the configurationManager file?")
	flag.StringVar(&naclPublicKey, "public-key", defaultGlobalConfigPath+FILE_SEPERATOR+"naclPub.key", "Where should the NaCl public key be saved to or loaded from?")
	flag.StringVar(&naclPrivateKey, "private-key", defaultGlobalConfigPath+FILE_SEPERATOR+"naclPriv.key", "Where should the NaCl private key be saved to or loaded from?")

	var err error
	currentWorkingDirectory, err = os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// If any flags containing a file path were passed in
	// on the command line we want to resolve them to absolute paths
	// relative to the current working directory. We'll make this call
	// again when we've parsed the config file, but those will be resolved
	// relative to the directory containing the config file.
	flag.Parse()
	flagBasePath = currentWorkingDirectory + FILE_SEPERATOR
	flag.Visit(resolveAbsoluteFlagPaths)

	// If a config file wasn't passed in on the command line we go looking for one
	// starting at the current working directory and traveling up the hierarchy
	// then checking in ~/.config/pact/
	_, err = os.Stat(configFilePath)
	if configFilePath == "" || os.IsNotExist(err) {
		configFilePath = findConfigFile(currentWorkingDirectory)
	}

	if configFilePath == "" {
		configFilePath = checkDirForConfigFile(defaultGlobalConfigPath)
	}

	if configFilePath != "" {
		loadConfFile()
	} else {
		ui.Info("WARNING: No config file found, run the config command to generate one.")
	}
}

/**
 * loadGlobalConf
 * Helper method to load a configurationManager file and parse values, used by init and during
 * interactive configurationManager file generation as the globalconf package is able to handle
 * persisting flag values to disk.
 */
func loadConfFile() *globalconf.GlobalConf {
	var err error
	configurationManager, err = globalconf.NewWithOptions(&globalconf.Options{
		Filename: configFilePath,
	})

	if err != nil {
		log.Println(err)
	}

	absConfigFilePath, _ := filepath.Abs(configFilePath)
	configFileBasePath, _ := filepath.Split(absConfigFilePath)
	// Reads configurationManager data as provided in the config file
	// Path data provided will be expanded relative to the config file
	// any flags passed in via CLI will already be absolute and unaffected by
	// this repeated call
	configurationManager.ParseAll()
	flagBasePath = configFileBasePath
	flag.Visit(resolveAbsoluteFlagPaths)
	return configurationManager
}

/**
 * findConfigFile
 * Recursive method to walk backwards through a path looking for a configurationManager file
 */
func findConfigFile(directory string) string {
	filePath := checkDirForConfigFile(directory)
	if filePath == "" {
		climbIndex := strings.LastIndex(directory, FILE_SEPERATOR)
		if climbIndex != -1 {
			return findConfigFile(directory[0:climbIndex])
		}
	} else {
		return filePath
	}

	return ""
}

/**
 * checkDirForConfigFile
 * Helper which either returns the full path of the configurationManager file if one is found
 * or returns the empty string if one is not.
 */
func checkDirForConfigFile(directory string) string {
	if len(directory) > 0 && directory[:len(directory)-1] != FILE_SEPERATOR {
		directory += FILE_SEPERATOR
	}
	filePath := directory + CONFIG_FILE_NAME
	if _, err := os.Stat(filePath); err == nil {
		return filePath
	} else {
		return ""
	}
}

/**
 * getAbsPath
 * Returns the absolute path given a string representation of a file path.In contrast to the
 * built in filepath.Abs method which always evaluates a path relative to the current
 * working directory, this method allows you to set a starting point for file path resolution.
 * This method also expands ~ to the home directory, which Abs does not support.
 *
 * To avoid confusion this method is always used to expand paths even when evaluating
 * paths relative to the current working directory.
 */
func getAbsPath(path string, base string) string {
	if base != currentWorkingDirectory {
		err := os.Chdir(base)
		// change working directory only while necessary
		defer os.Chdir(currentWorkingDirectory)
		if err != nil {
			panic(err)
		}
	}

	if "~" == path[:1] {
		path = userHomeDir + path[1:]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return absPath
}

/**
 * resolveAbsoluteFlagPaths
 * Helper method which is passed to flag.Visit and ensures the file path strings
 * are properly expanded.
 */
func resolveAbsoluteFlagPaths(f *flag.Flag) {
	f.Value.Set(getAbsPath(f.Value.String(), flagBasePath))
}

/**
 * GenerateConfigFile
 * Helper method that guides the user through interactive prompts and determines whether the
 * intention is to create a new configuration file or update the existing one.
 */
func GenerateConfigFile() {
	if configFilePath == "" {
		ok := getOrCreateConfigFile()
		if !ok {
			return
		}
	} else {
		ui.Info("A config file is already loaded from: " + configFilePath)
		resp, err := ui.Ask("U to Update the existing file, C to Create a new file somewhere else [U/C]:")
		if err != nil {
			log.Fatal(err)
		}
		if len(resp) > 1 {
			resp = resp[:1]
		}
		resp = strings.ToLower(resp)
		switch {
		default:
			ui.Info("No input detected, exiting...")
			return
		case "c" == resp:
			// this method will create an empty file
			ok := getOrCreateConfigFile()
			if !ok {
				return
			}
		case "u" == resp:
			ui.Info("Updating existing file...\n")
		}
	}

	ui.Info("Config file will be written to: " + configFilePath + "\n")
	loadConfFile()

	ui.Info("All paths entered on these prompts are relative to the current working directory")
	ui.Info("Any of the following options can be skipped by hitting return.")
	ui.Info("Skipped responses do not overwrite existing settings.\n")

	flag.VisitAll(promptForAndPersistFlagValue)
	ui.Info("All provided configuration information has been persisted to disk.\n")

	// The user will need an NACL keypair for encryption, let's see if one exists and
	// remind the user to generate one in the case where it does not.
	if !DoesNACLKeypairExist() {
		ui.Info("Your public/private keypair could not be found or appears corrupt!")
		ui.Info("Run the key-gen command to create a new keypair.")
	}

}

/**
 * promptForAndSetConfigFilePath
 * Helper method to handle the logic of creating a new configurationManager file.
 */
func getOrCreateConfigFile() bool {
	cp, err := ui.Ask("Where should we write a new config file? [" + defaultGlobalConfigPath + "]:")
	if err != nil {
		log.Fatal(err)
	}
	if cp == "" {
		ui.Warn("Default location chosen...")
		os.Mkdir(defaultGlobalConfigPath, 0755)
		cp = defaultGlobalConfigPath
	}
	fullPath := getAbsPath(cp, currentWorkingDirectory)
	filePath := checkDirForConfigFile(fullPath)
	if filePath == "" {
		filePath = fullPath + "/" + CONFIG_FILE_NAME
		_, err := os.Create(filePath)
		if err != nil {
			ui.Warn("Path invalid, path directories must already exist, exiting...")
			return false
		}
	} else {
		ui.Info("Config file " + filePath + " exists, updating in place...\n")
	}
	configFilePath = filePath
	return true
}

/**
 * promptForAndPersistFlagValues
 * Helper method which uses the flag's usage string to prompt the user to enter a value
 * for the flag. That value is then assigned to the flag and persisted to disk using
 * globalconf's Set method.
 */
func promptForAndPersistFlagValue(f *flag.Flag) {
	// If we're in this method we don't need to
	// ask where the config file should be saved.
	if f.Name == "config" {
		return
	}

	message := f.Usage
	fValue := f.Value.String()
	if fValue != "" {
		message += " [" + fValue + "]:"
	}

	var response, err = ui.Ask(message)
	if err != nil {
		log.Fatal(err)
	}

	// If we're running config command and flagValues are set, persist them
	if response == "" && fValue != "" {
		response = fValue
	}

	if response != "" {
		response = getAbsPath(response, currentWorkingDirectory)
		f.Value.Set(response)
		configurationManager.Set("", f)
	}
}

/**
 * DoesNACLKeypairExist
 * Helper method which checks for the existence of public and private keys based on
 * the currently loaded configurationManager values.
 */
func DoesNACLKeypairExist() bool {
	_, pubKeyError := os.Stat(GetNACLPublicKeyPath())
	_, privKeyError := os.Stat(GetNACLPrivateKeyPath())
	return pubKeyError == nil && privKeyError == nil
}
