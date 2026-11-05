package conns_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

var consumers = map[string]string{
	"Name":               "registry.ServicePackages, for the duplicate-name error messages",
	"SDKResources":       "alicloud.assembleProvider, which rewrites Provider.ResourcesMap",
	"SDKDataSources":     "alicloud.assembleProvider, which rewrites Provider.DataSourcesMap",
	"Resources":          "framework.(*alicloudProvider).Resources",
	"DataSources":        "framework.(*alicloudProvider).DataSources",
	"Functions":          "framework.(*alicloudProvider).Functions",
	"EphemeralResources": "framework.(*alicloudProvider).EphemeralResources",
	"ListResources":      "framework.(*alicloudProvider).ListResources",
	"Actions":            "framework.(*alicloudProvider).Actions",
}

func TestUnitEveryServicePackageFieldHasAConsumer(t *testing.T) {
	rt := reflect.TypeOf(conns.ServicePackage{})

	declared := make([]string, 0, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		declared = append(declared, rt.Field(i).Name)
	}

	for _, name := range declared {
		if _, ok := consumers[name]; !ok {
			t.Errorf("ServicePackage.%s has no registered consumer: add one and record it in this table, "+
				"or the category is declared by product packages and never registered", name)
		}
	}
	for name := range consumers {
		if !slices.Contains(declared, name) {
			t.Errorf("this table names ServicePackage.%s, which no longer exists: drop the entry", name)
		}
	}
}
