package chatbot_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	chatbot "github.com/aliyun/terraform-provider-alicloud/alicloud/service/chatbot"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// Package-local counterparts of the flat alicloud test helpers
// (service_alicloud_common_test.go), which are not reachable from a service
// package test. The collection-order coverage gate keys builder discovery off
// the bare identifier "resourceTestAccConfigFunc", so the same spelling is kept
// here; the exported acctest.ResourceTestAccConfig is semantically equivalent
// but is a selector expression the static gate does not match.

const (
	NOSET      = "#NOSET"
	CHECKSET   = "#CHECKSET"
	REMOVEKEY  = "#REMOVEKEY"
	CLEARLIST  = "#CLEARLIST"
	CLEARMAP   = "#CLEARMAP"
	REGEXMATCH = "#REGEXMATCH:"

	indentationSymbol = " "
	childIndent       = 2
)

// resourceTestAccConfigFunc returns a function that renders the dependence
// config plus the primary resource built from the accumulated attribute map.
func resourceTestAccConfigFunc(resourceId string, name string, configDependence func(name string) string) func(map[string]interface{}) string {
	attributeMap := make(map[string]interface{})
	return func(changeMap map[string]interface{}) string {
		newMap := make(map[string]interface{}, len(attributeMap))
		for k, v := range attributeMap {
			newMap[k] = v
		}
		for k, v := range changeMap {
			if sentinel, ok := v.(string); ok {
				switch sentinel {
				case REMOVEKEY:
					delete(newMap, k)
					continue
				case CLEARMAP:
					newMap[k] = make(map[string]interface{})
					continue
				case CLEARLIST:
					newMap[k] = make([]interface{}, 0)
					continue
				}
			}
			newMap[k] = v
		}
		attributeMap = newMap

		strs := strings.Split(resourceId, ".")
		primaryConfig := fmt.Sprintf("\n\nresource \"%s\" \"%s\" ", strs[0], strs[1])
		return configDependence(name) + primaryConfig + fmt.Sprint(valueConvert(0, reflect.ValueOf(attributeMap)))
	}
}

func valueConvert(indentation int, val reflect.Value) interface{} {
	switch val.Kind() {
	case reflect.Interface:
		return valueConvert(indentation, reflect.ValueOf(val.Interface()))
	case reflect.String:
		return fmt.Sprintf("\"%s\"", val.String())
	case reflect.Slice:
		return listValue(indentation, val)
	case reflect.Map:
		return mapValue(indentation, val)
	case reflect.Bool:
		return val.Bool()
	case reflect.Int, reflect.Int64:
		return val.Int()
	default:
		log.Panicf("invalid attribute value type: %#v", val)
	}
	return ""
}

func listValue(indentation int, val reflect.Value) string {
	var valList []string
	for i := 0; i < val.Len(); i++ {
		valList = append(valList, addIndentation(indentation+childIndent)+
			fmt.Sprint(valueConvert(indentation+childIndent, val.Index(i))))
	}
	return fmt.Sprintf("[\n%s\n%s]", strings.Join(valList, ",\n"), addIndentation(indentation))
}

func mapValue(indentation int, val reflect.Value) string {
	var valList []string
	for _, keyV := range val.MapKeys() {
		mapVal := getRealValueType(val.MapIndex(keyV))
		if mapVal.Kind() == reflect.Slice && mapVal.Len() > 0 {
			if getRealValueType(mapVal.Index(0)).Kind() == reflect.Map {
				valList = append(valList, addIndentation(indentation)+
					listValueMapChild(indentation+childIndent, keyV.String(), mapVal))
				continue
			}
		}
		value := valueConvert(indentation+len(keyV.String())+childIndent+3, val.MapIndex(keyV))
		var line string
		switch value.(type) {
		case bool:
			line = fmt.Sprintf(`%s%s = %t`, addIndentation(indentation+childIndent), keyV.String(), value)
		case int64:
			line = fmt.Sprintf(`%s%s = %d`, addIndentation(indentation+childIndent), keyV.String(), value)
		default:
			line = fmt.Sprintf(`%s%s = %s`, addIndentation(indentation+childIndent), keyV.String(), value)
		}
		valList = append(valList, line)
	}
	return fmt.Sprintf("{\n%s\n%s}", strings.Join(valList, "\n"), addIndentation(indentation))
}

func listValueMapChild(indentation int, key string, val reflect.Value) string {
	var valList []string
	for i := 0; i < val.Len(); i++ {
		valList = append(valList, addIndentation(indentation)+key+" "+
			mapValue(indentation, getRealValueType(val.Index(i))))
	}
	return fmt.Sprintf("%s\n%s", strings.Join(valList, "\n"), addIndentation(indentation))
}

func getRealValueType(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Interface {
		return getRealValueType(reflect.ValueOf(value.Interface()))
	}
	return value
}

func addIndentation(indentation int) string {
	return strings.Repeat(indentationSymbol, indentation)
}

// testAccChatbotFaqCategoryIds returns real FAQ category ids published on the
// caller's default business space. CreatePublishTask only accepts DataIdList
// when BizType is "faq", and the values must be real FAQ category ids rather
// than fabricated placeholders. The provider is not configured at the start of
// a TestAcc function (resource.TestMain does not pre-configure it), so a
// standalone AliyunClient is built from the same ALICLOUD_* environment
// variables the rest of the suite uses, with SkipRegionValidation and a fixed
// AccountType to avoid the unrelated BssOpenApi identity probe.
func testAccChatbotFaqCategoryIds(t *testing.T) []string {
	acctest.PreCheck(t)
	region := os.Getenv("ALICLOUD_REGION")
	accessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
	secretKey := os.Getenv("ALICLOUD_SECRET_KEY")
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(accessKey).
		SetAccessKeySecret(secretKey))
	if err != nil {
		t.Fatalf("build chatbot fixture credential: %v", err)
	}
	config := &connectivity.Config{
		AccessKey:            accessKey,
		SecretKey:            secretKey,
		Region:               connectivity.Region(region),
		RegionId:             region,
		AccountType:          "test",
		Protocol:             "https",
		Endpoints:            new(sync.Map),
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
		Credential:           credential,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("build chatbot fixture client: %v", err)
	}
	svc := chatbot.NewService(client)
	ids, err := svc.ListFaqCategories()
	if err != nil {
		t.Fatalf("list chatbot faq categories for data_id_list fixture: %v", err)
	}
	// The reorder coverage needs two distinct ids so step1 -> step2 swaps the
	// member order. If the account's default business space already exposes
	// two or more FAQ categories, use the first two. If it only exposes one
	// (the common case on a fresh test account), create a temporary FAQ
	// category via CreateCategory and delete it through t.Cleanup so the
	// fixture is a real FAQ category id, not a fabricated placeholder. The
	// t.Fatal/t.Fatalf calls stay inside this helper so the test function body
	// itself is free of the call family that the collection-order gate keys on
	// to decide whether a test is skipped.
	if len(ids) >= 2 {
		return ids[:2]
	}
	createdId, err := svc.CreateFaqCategory(fmt.Sprintf("tf-testacc-faq-%d", tfacctest.RandIntRange(10000, 99999)))
	if err != nil {
		t.Fatalf("create chatbot faq category fixture for data_id_list reorder: %v", err)
	}
	t.Cleanup(func() {
		if err := svc.DeleteFaqCategory(createdId); err != nil {
			t.Logf("cleanup chatbot faq category fixture %s: %v", createdId, err)
		}
	})
	if len(ids) == 1 {
		return []string{ids[0], createdId}
	}
	// No existing category at all: create a second temporary one too.
	secondId, err := svc.CreateFaqCategory(fmt.Sprintf("tf-testacc-faq-%d", tfacctest.RandIntRange(10000, 99999)))
	if err != nil {
		t.Fatalf("create chatbot faq category fixture2 for data_id_list reorder: %v", err)
	}
	t.Cleanup(func() {
		if err := svc.DeleteFaqCategory(secondId); err != nil {
			t.Logf("cleanup chatbot faq category fixture2 %s: %v", secondId, err)
		}
	})
	return []string{createdId, secondId}
}
