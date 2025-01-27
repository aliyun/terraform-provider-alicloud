//nolint:all
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/waigani/diffparser"
)

func init() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.DisableTimestamp = false
	customFormatter.DisableColors = false
	customFormatter.ForceColors = true
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

var (
	fileNames = flag.String("fileNames", "", "the files to check diff")
)

func main() {
	exitCode := 0
	flag.Parse()
	if fileNames != nil && len(*fileNames) == 0 {
		log.Warningf("the diff file is empty")
		return
	}
	byt, _ := ioutil.ReadFile(*fileNames)
	diff, _ := diffparser.Parse(string(byt))
	fileRegex := regexp.MustCompile("alicloud/resource[0-9a-zA-Z_]*.go")
	fileTestRegex := regexp.MustCompile("alicloud/resource[0-9a-zA-Z_]*_test.go")
	for _, file := range diff.Files {
		fmt.Println()
		if fileRegex.MatchString(file.NewName) {
			if fileTestRegex.MatchString(file.NewName) {
				continue
			}
			resourceName := strings.TrimPrefix(strings.TrimSuffix(strings.Split(file.NewName, "/")[1], ".go"), "resource_")
			log.Infof("==> Checking resource %s breaking change...", resourceName)
			oldAttrs := make(map[string]map[string]interface{})
			newAttrs := make(map[string]map[string]interface{})
			for _, hunk := range file.Hunks {
				if hunk != nil {
					ParseResourceSchema(hunk.OrigRange, hunk.OrigRange.Length, oldAttrs)
					ParseResourceSchema(hunk.NewRange, hunk.NewRange.Length, newAttrs)
				}
			}
			if !IsBreakingChange(oldAttrs, newAttrs) {
				log.Infof("--- PASS")
			} else {
				log.Errorf("--- FAIL")
			}
		}
	}

	os.Exit(exitCode)
}

func IsBreakingChange(oldAttrs, newAttrs map[string]map[string]interface{}) (res bool) {
	for filedName, oldAttr := range oldAttrs {
		// Optional -> Required
		_, exist1 := oldAttr["Optional"]
		_, exist2 := newAttrs[filedName]["Required"]
		if exist1 && exist2 {
			res = true
			log.Errorf("[Breaking Change]: '%v' should not been changed to required!", filedName)
		}
		// Type changed
		typPrev, exist1 := oldAttr["Type"]
		typCurr, exist2 := newAttrs[filedName]["Type"]

		if exist1 && exist2 && typPrev != typCurr {
			res = true
			log.Errorf("[Breaking Change]: '%v' type should not been changed!", filedName)
		}

		// Non-ForceNwe -> ForceNew
		_, exist1 = oldAttr["ForceNew"]
		_, exist2 = newAttrs[filedName]["ForceNew"]
		if !exist1 && exist2 {
			res = true
			log.Errorf("[Breaking Change]: '%v' should not been changed to ForceNew!", filedName)
		}

		// Type string/int: valid values
		validateValuesOld, exist1 := oldAttr["ValidateFuncValues"]
		validateValuesNew, exist2 := newAttrs[filedName]["ValidateFuncValues"]
		if exist1 {
			if !exist2 {
				log.Warningf("[Warning]: '%v' ValidateFunc should not been removed!", filedName)
			} else {
				for key, _ := range validateValuesOld.(map[string]struct{}) {
					if _, ok := validateValuesNew.(map[string]struct{})[key]; !ok {
						res = true
						log.Errorf("[Breaking Change]: '%v' valid value %s should not been removed!", filedName, key)
					}
				}
			}
		}
	}
	return
}

func ParseResourceSchema(hunk diffparser.DiffRange, length int, attributeMap map[string]map[string]interface{}) {
	schemaRegex := regexp.MustCompile("^\\t*\"([a-zA-Z_]*)\"")
	typeRegex := regexp.MustCompile("^\\t*Type:\\s+schema.([a-zA-Z]*)")
	optionRegex := regexp.MustCompile("^\\t*Optional:\\s+([a-z]*),")
	forceNewRegex := regexp.MustCompile("^\\t*ForceNew:\\s+([a-z]*),")
	requiredRegex := regexp.MustCompile("^\\t*Required:\\s+([a-z]*),")
	validateStringRegex := regexp.MustCompile("^\\t*ValidateFunc: +(validation\\.)?StringInSlice\\(\\[\\]string\\{([a-z\\-A-Z_,\"\\s]*)")
	validateIntRegex := regexp.MustCompile("^\\t*ValidateFunc: +(validation\\.)?IntInSlice\\(\\[\\]int\\{([0-9,\\s]*)")

	temp := map[string]interface{}{}
	schemaName := ""
	for i := 0; i < length; i++ {
		currentLine := hunk.Lines[i]
		content := currentLine.Content
		fieldNameMatched := schemaRegex.FindAllStringSubmatch(content, -1)
		if fieldNameMatched != nil && fieldNameMatched[0] != nil {
			if len(schemaName) != 0 && schemaName != fieldNameMatched[0][1] {
				temp["Name"] = schemaName
				attributeMap[schemaName] = temp
				temp = map[string]interface{}{}
			}
			schemaName = fieldNameMatched[0][1]
		}

		if !schemaRegex.MatchString(currentLine.Content) && currentLine.Mode == diffparser.UNCHANGED {
			continue
		}

		typeMatched := typeRegex.FindAllStringSubmatch(content, -1)
		typeValue := ""
		if typeMatched != nil && typeMatched[0] != nil {
			typeValue = typeMatched[0][1]
			temp["Type"] = typeValue
		}

		optionalMatched := optionRegex.FindAllStringSubmatch(content, -1)
		optionValue := ""
		if optionalMatched != nil && optionalMatched[0] != nil {
			optionValue = optionalMatched[0][1]
			op, _ := strconv.ParseBool(optionValue)
			temp["Optional"] = op
		}

		forceNewMatched := forceNewRegex.FindAllStringSubmatch(content, -1)
		forceNewValue := ""
		if forceNewMatched != nil && forceNewMatched[0] != nil {
			forceNewValue = forceNewMatched[0][1]
			fc, _ := strconv.ParseBool(forceNewValue)
			temp["ForceNew"] = fc
		}

		requiredMatched := requiredRegex.FindAllStringSubmatch(content, -1)
		requiredValue := ""
		if requiredMatched != nil && requiredMatched[0] != nil {
			requiredValue = requiredMatched[0][1]
			rq, _ := strconv.ParseBool(requiredValue)
			temp["Required"] = rq
		}

		for _, validateRegex := range []*regexp.Regexp{validateStringRegex, validateIntRegex} {
			validateMatched := validateRegex.FindAllStringSubmatch(content, -1)
			if validateMatched != nil && validateMatched[0] != nil {
				validateValue := validateMatched[0][2]
				stringMap := make(map[string]struct{})
				for _, v := range strings.Split(validateValue, ",") {
					stringMap[strings.TrimSpace(v)] = struct{}{}
				}
				temp["ValidateFuncValues"] = stringMap
			}
		}

	}
	if _, exist := attributeMap[schemaName]; !exist && len(temp) >= 1 {
		temp["Name"] = schemaName
		attributeMap[schemaName] = temp
	}
	return
}
