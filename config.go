package main

import (
	"errors"
	"fmt"
	"github.com/josuehennemann/conf"
	"reflect"
	"strings"
)

var pathConfig string
var config *StConfig

type StConfig struct {
	Httplisten       string
	Httpslisten      string
	Httpscertificate string
	Filepath         string
	Environment      string
	Autosave         bool
	Logfile          string
}

func NewConfig() {
	cf, err := conf.ReadConfigFile(pathConfig)
	if err != nil {
		CheckErrorAndKillMe(errors.New(err.String()))
	}

	hasFieldError := false
	fields := map[string][]string{}
	tmp := cf.GetSections()
	for _, section := range tmp {
		tmp2, _ := cf.GetOptions(section)
		fields[section] = tmp2
	}
	config = &StConfig{}
	t := reflect.ValueOf(config).Elem()
	for section, entries := range fields {
		for _, entry := range entries {
			fieldStruct := t.FieldByName(strings.Title(entry))
			if !fieldStruct.CanSet() {
				fmt.Printf("Entry [%s] not found in struct config [%+v]\n", entry, config)
				hasFieldError = true
				continue
			}

			switch fieldStruct.Kind() {
			case reflect.String:
				v, _ := cf.GetString(section, entry)
				fieldStruct.SetString(v)
			case reflect.Bool:
				v, _ := cf.GetBool(section, entry)
				fieldStruct.SetBool(v)
			}
		}
	}
	if hasFieldError {
		CheckErrorAndKillMe(fmt.Errorf("Invalid file conf"))
	}
}

func (this *StConfig) GetHttpListen() string {
	return this.Httplisten
}
func (this *StConfig) GetHttpsListen() string {
	return this.Httpslisten
}
func (this *StConfig) GetHttpsCertificate() string {
	return this.Httpscertificate
}
func (this *StConfig) GetFilePath() string {
	return this.Filepath
}

func (this *StConfig) GetLogFile() string {
	return this.Logfile
}
