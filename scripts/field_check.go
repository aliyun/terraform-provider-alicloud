package scripts

import (
	"flag"
	"os"
	"testing"

	old "github.com/aliyun/terraform-provider-alicloud-prev/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
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
			log.Errorf("Resource: %s, Field: %s Field incompatible has occurred in the current version,Please Check the Optional/Required type", *resourceName, fieldName)
			t.Fatal()
		}
		if !oldField.ForceNew && newField.ForceNew {
			log.Errorf("Resource: %s, Field: %s Field incompatible c`hanges have occurred in the current version,Please Check the Field type", *resourceName, fieldName)
			t.Fatal()
		}
		if oldField.Type != newField.Type {
			log.Errorf("Resource: %s, Field: %s Field incompatible c`hanges have occurred in the current version,Please Check the Field type", *resourceName, fieldName)
			t.Fatal()
		}
	}
}