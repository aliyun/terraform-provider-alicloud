package scripts

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type Resource struct {
	Name       string
	Arguments  map[string]Field
	Attributes map[string]Field
}

type Field struct {
	Name        string
	Optional    bool
	Required    bool
	ForceNew    bool
	Description string
}

func parseResourse(filePath string) (*Resource,error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("resource parsing: %v", err)
	}
	name := filepath.Base(filePath)
	re := regexp.MustCompile("[a-zA-Z_]*")
	resourceName := "alicloud_" + re.FindString(name)
	result := &Resource{Name: resourceName, Arguments: nil, Attributes: nil}

	argsRegex := regexp.MustCompile("## Argument Reference")
	secondLevelRegex := regexp.MustCompile("##")
	attribRegex := regexp.MustCompile("## Attributes Reference")

	argsRange := argsRegex.FindIndex(bytes)
	argumentsStart,argumentsEnd := 0,0
	if argsRange != nil {
		argumentsStart = argsRange[1]
	}
	attributesStart := 0
	attrRange := attribRegex.FindIndex(bytes[argumentsStart:])
	if attrRange == nil{
		logrus.Warningf("the File: %s does not have the Attributes Reference", filePath)
		return nil,nil
	}
	ot := secondLevelRegex.FindIndex(bytes[argumentsStart:])
	if attrRange != nil {
		argumentsEnd = attrRange[0] + argumentsStart
		attributesStart = attrRange[1] + argumentsStart
	}
	if ot != nil{
		argumentsEnd = ot[0] + argumentsStart
	}
	var argumentsBytes []byte
	if argumentsStart <= argumentsEnd {
		argumentsBytes = bytes[argumentsStart:argumentsEnd]
	} else {
		argumentsBytes = make([]byte, 0)
	}
	otherRange := secondLevelRegex.FindIndex(bytes[attributesStart:])
	attributesBytes := bytes[attributesStart:]
	if otherRange != nil{
		attributesBytes = bytes[attributesStart:otherRange[0]+attributesStart]
	}
	argumentsFieldRegex := regexp.MustCompile("\\* `([a-zA-Z_0-9]*)`[ ]*-? ?(\\(.*\\)) ?(.*)")
	attributeFieldRegex := regexp.MustCompile("\\* `([a-zA-Z_0-9]*)`[ ]*-?(.*)")
	attributesMatched := attributeFieldRegex.FindAllSubmatch(attributesBytes, -1)
	result.Attributes = make(map[string]Field,0)
	for _, attributeParsed := range attributesMatched {
		Field := parseMatchLine(attributeParsed,false)
		result.Attributes[(*Field).Name] = *Field
	}

	argumentsMatched := argumentsFieldRegex.FindAllSubmatch(argumentsBytes, -1)
	result.Arguments = make(map[string]Field, 0)
	for _, argumentMatched := range argumentsMatched {
		Field := parseMatchLine(argumentMatched,true)
		result.Arguments[(*Field).Name] = *Field
	}
	return result,nil
}

func parseMatchLine(words [][]byte, argumentFlag bool) *Field {
	result := &Field{Name: "", Optional: true, Description: ""}
	if argumentFlag && len(words) >= 4{
		result.Name = string(words[1])
		result.Description = string(words[3])
		if strings.Contains(string(words[2]),"Optional"){
			result.Optional = true
		}
		if strings.Contains(string(words[2]),"Required"){
			result.Required = true
		}
		if strings.Contains(string(words[2]),"ForceNew"){
			result.ForceNew = true
		}
		return result
	}
	if !argumentFlag && len(words) >= 3{
		result.Name = string(words[1])
		result.Description = string(words[2])
		return result
	}
	return nil
}
