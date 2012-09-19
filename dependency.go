// Dependency inejction library
//
// Workflow:
// Tell where to read the config file with SetConfig()
// Introduce constructor with Add()
// Refresh() the dependency information
// Get() your struct
// and cast it to intended type.
//
// NOTE: constructor returns an empty interface which should be castable to whatever type you intend it to be.
package dependency

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Constructor func() (interface{}, error)

var (
	EmptyCategoryErr       = errors.New("No constructors in this category")
	EmptySelectedErr       = errors.New("Constructor not set in config file or Refresh() was not run after Add()")
	ConfigNotSetErr        = errors.New("Config file not set")
	IncorrectConfigTypeErr = errors.New("Incorrect configType")
)

// map[category][name] = constructor, error
var deps map[string](map[string]Constructor)

var configSource, configType, configData string
var selected map[string]Constructor

func init() {
	deps = make(map[string](map[string]Constructor))
	selected = make(map[string]Constructor)
}

// TODO: Lisätessä voisi teoriasa tarkistaa reflektiolla että täyttää jonkun tietyn interfacen tj.
//       reflect.Type.Implements(reflect.Type)
//          otetaan mapissa ensimmäisenä olevan konstruktorin paluuarvo, katsotaan onko jotain interfacea
//          jos on niin testaan uudet sitä vasten? tuntuu huonolta idealta.

func Add(category, name string, constructor Constructor) {
	if deps[category] == nil {
		deps[category] = make(map[string]Constructor)
	}
	deps[category][name] = constructor
}

// Paluuarvo pitää castata kutsujan puolella oikeaksi tyypiksi
func Get(category string) (interface{}, error) {
	if deps[category] == nil {
		return nil, EmptyCategoryErr
	}

	if len(deps[category]) == 1 {
		for _, constructor := range deps[category] {
			constructor()
		}
	}

	if selected[category] == nil {
		return nil, EmptySelectedErr
	}
	return selected[category]()
}

// TODO: Errors?
func Refresh() (err error) {
	var line []byte

	if configSource == "" || configType == "" {
		return ConfigNotSetErr
	}
	loadConfig()

	reader := bufio.NewReader(strings.NewReader(configData))

	for {
		line, _, err = reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}

		if !strings.HasPrefix(string(line), "@dep") {
			continue
		}

		parts := strings.SplitAfterN(string(line), " ", 4)
		if parts[0] != "@dep " || parts[2] != "= " {
			// incorrect @dep line
			continue
		}

		var constructor Constructor
		if constructors, ok := deps[strings.TrimSpace(parts[1])]; !ok {
			continue
		} else if constructor, ok = constructors[strings.TrimSpace(parts[3])]; !ok {
			continue
		}
		selected[strings.TrimSpace(parts[1])] = constructor
	}
	return
}

// Set config source
// configtype:      text|file
// configsource:    newline separated config information if configtype == text
//                  name of config file if configtype == file
func SetConfig(configsource, configtype string) {
	configType = configtype
	configSource = configsource
}

func loadConfig() (err error) {
	var temp []byte

	switch configType {
	case "text":
		configData = configSource
	case "file":
		if temp, err = ioutil.ReadFile(configSource); err != nil {
			return
		}
		configData = string(temp)
	default:
		return IncorrectConfigTypeErr
	}
	return
}
