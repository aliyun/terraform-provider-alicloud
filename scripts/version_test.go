package scripts

import (
	"flag"
	old "github.com/aliyun/terraform-provider-alicloud-prev/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
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

var resourceName = flag.String("resource", "", "the name of the terraform resource to diff")

func TestFieldCheck(t *testing.T) {
	flag.Parse()
	if len(*resourceName) == 0 {
		log.Errorf("The Resource Name is Empty")
		t.Fatal()
	}

	oldObj := old.Provider().(*schema.Provider).ResourcesMap[*resourceName].Schema
	n := alicloud.Provider().(*schema.Provider).ResourcesMap[*resourceName].Schema
	for fieldName, newField := range n {
		oldField := oldObj[fieldName]
		if oldField.Optional && newField.Required {
			log.Errorf("Resource: %s, Field: %s Field incompatible c`hanges have occurred in the current version,Please Check the Optional/Required type", *resourceName, fieldName)
			t.Fatal()
			//t.Errorf("Resource: %s, Field: %s Field incompatible c`hanges have occurred in the current version,Please Check the Optional/Required type",resourceName,fieldName)
		}
		if !oldField.ForceNew && newField.ForceNew {
			log.Errorf("Resource: %s, Field: %s Field incompatible c`hanges have occurred in the current version,Please Check the Field type", *resourceName, fieldName)
			t.Fatal()
			//t.Errorf("Resource: %s, Field: %s Field incompatible c`hanges have occurred in the current version,Please Check the Field type",resourceName,fieldName)
		}
		if oldField.Type != newField.Type {
			log.Errorf("Resource: %s, Field: %s Field incompatible c`hanges have occurred in the current version,Please Check the Field type", *resourceName, fieldName)
			t.Fatal()
		}
	}
}


//func TestCheck(t *testing.T) {
//	flag.Parse()
//	if len(*resourceName) == 0 {
//		log.Errorf("The Resource Name is Empty")
//		t.Fatal()
//	}
//	n := alicloud.Provider().(*schema.Provider).ResourcesMap[*resourceName].Schema
//	splitRes := strings.Split(*resourceName, "alicloud_")
//	if len(splitRes) < 2 {
//		return
//	}
//	basePath := "../website/docs/r/"
//	filename := strings.Join([]string{basePath, splitRes[1], ".html.markdown"}, "")
//	resourceMd, err := parseResourse(filename)
//	if err != nil {
//		log.Error(err)
//	}
//	if resourceMd == nil {
//		log.Errorf("Resource: %s Filename: %s The Markdown id nil ", *resourceName, filename)
//	}
//	argument, attr := 0, 0
//	schemaMd := make(map[string]Field,0)
//	if resourceMd != nil && resourceMd.Arguments != nil {
//		argument = len(resourceMd.Arguments)
//		schemaMd = mergeMaps(schemaMd,resourceMd.Arguments)
//	}
//	if resourceMd != nil && resourceMd.Attributes != nil {
//		attr = len(resourceMd.Attributes)
//		schemaMd = mergeMaps(schemaMd,resourceMd.Attributes)
//	}
//	log.Infof("The Resource: %v, the Arguments: %v, the Attribute: %v", resourceName, argument, attr)
//
//	if resourceMd.Arguments == nil || resourceMd.Attributes == nil {
//		t.Logf("Resource: %s, the Arguments or Attributes is empty with the doc %s", *resourceName, filename)
//	}
//
//	// checkout the consistency with the schema, "+1" stands for "id" in argumen
//	if len(n)+1 != len(schemaMd){
//		t.Errorf("Resource: %s, the number of the field defined in schema is not consistent with the doc %s",*resourceName,filename)
//	}
//}

//func mergeMaps(Dst map[string]interface{},maps ...map[string]interface{}) map[string]interface{} {
//	for _, m := range maps {
//		for k, v := range m {
//			Dst[k] = v
//		}
//	}
//	return Dst
//}