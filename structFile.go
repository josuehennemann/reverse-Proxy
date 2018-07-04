package main

import (
	"encoding/json"
	"github.com/josuehennemann/logger"
	"io/ioutil"
	"strings"
	"sync"
)

type RuleList struct {
	sync.RWMutex
	list map[string]*Rule
	path string //file path
}

type Rule struct {
	Origin       string //handle
	Destiny      string //url to redirect or make request
	RemoveOrigin bool   // flag to remove Origin from URL
	Redirect     bool   // redirect or make request
}

func (ru *Rule) Validate() bool {
	ru.Origin = strings.TrimSpace(ru.Origin)
	ru.Destiny = strings.TrimSpace(ru.Destiny)
	if isEmpty(ru.Origin) {
		return false
	}

	if isEmpty(ru.Destiny) {
		return false
	}
	return true
}

func initRuleList() (*RuleList, error) {

	rl := RuleList{}
	rl.list = map[string]*Rule{}
	rl.path = config.GetFilePath()

	if rl.path != "" {
		if err := rl.load(); err != nil {
			return nil, err
		}
	}

	return &rl, nil
}

//private functions
/*
	Method to add new or update rule

	@handle is url in server http
	@rule is rule
*/
func (rl *RuleList) add(handle string, rule *Rule) {
	rl.list[handle] = rule
}

func (rl *RuleList) load() error {
	content, err := ioutil.ReadFile(rl.path)
	if err != nil {
		return err
	}
	tmp := []Rule{}
	err = json.Unmarshal(content, &tmp)
	if err != nil {
		return err
	}

	for _, r := range tmp {
		internalRule := r
		if !internalRule.Validate() {
			Logger.Printf(logger.WARN, "Invalid rule [%+v]", internalRule)
			continue
		}
		rl.add(internalRule.Origin, &internalRule)
	}
	return nil
}

//public functions

//Reload rule files
func (rl *RuleList) Reload() error {
	rl.Lock()
	defer rl.Unlock()
	rl.list = map[string]*Rule{}
	return rl.load()
}

func (rl *RuleList) GetRule(handle string) *Rule {
	rl.RLock()
	defer rl.RUnlock()
	
	for  {
		item, exists := rl.list[handle]
		if exists {
			return item
		}
		count := strings.LastIndex(handle, "/")
		if count <= 0 {
			break
		}
		handle = handle[:count]
	}

	//if not exist, item is nil
	return nil
}

func (rl *RuleList) List() map[string]*Rule {
	newlist := map[string]*Rule{}
	for k, v := range rl.list {
		tmp := *v
		newlist[k] = &tmp
	}
	return newlist
}
