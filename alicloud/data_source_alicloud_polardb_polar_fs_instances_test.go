package alicloud

import (
	"reflect"
	"testing"
)

func TestUnitPolarDBPolarFsListHelpers(t *testing.T) {
	wantItems := []interface{}{map[string]interface{}{"PolarFsInstanceId": "pfs-1"}}
	if got := polarDBPolarFsListItems(map[string]interface{}{"PolarFsPaths": wantItems}); !reflect.DeepEqual(got, wantItems) {
		t.Fatalf("unexpected list items: %#v", got)
	}
	wantTags := []interface{}{map[string]interface{}{"Key": "env", "Value": "test"}}
	if got := polarDBPolarFsTags(map[string]interface{}{"Tag": wantTags}); !reflect.DeepEqual(got, wantTags) {
		t.Fatalf("unexpected tags: %#v", got)
	}
	expanded := expandPolarDBPolarFsTags(map[string]interface{}{"env": "test"})
	if len(expanded) != 1 || expanded[0]["Key"] != "env" || expanded[0]["Value"] != "test" {
		t.Fatalf("unexpected expanded tags: %#v", expanded)
	}
}
