package cloudmonitorservice

import (
	"testing"

	"github.com/alibabacloud-go/tea/tea"
)

func TestUnitContactDeleteError(t *testing.T) {
	for _, absent := range []string{"ContactNotExists", "ContactGroupNotExists"} {
		for _, code := range []string{absent, "ResourceNotFound", "400", "403", "404", "Forbidden", "InvalidParameter"} {
			err := tea.NewSDKError(map[string]interface{}{"code": code, "message": "403 " + absent})
			if got, want := cmsContactAbsent(err, absent), code == absent || code == "ResourceNotFound"; got != want {
				t.Fatalf("%s / %s: got %v", absent, code, got)
			}
		}
	}
	// DeleteContactGroup reports ContactNotExists for a missing group, so the
	// group deletion accepts both not-found codes.
	for _, code := range []string{"ContactNotExists", "ContactGroupNotExists", "ResourceNotFound"} {
		err := tea.NewSDKError(map[string]interface{}{"code": code, "message": "403 ContactNotExists"})
		if got := cmsContactAbsent(err, "ContactNotExists", "ContactGroupNotExists"); !got {
			t.Fatalf("group deletion / %s: got %v, want true", code, got)
		}
	}
	if got := cmsContactAbsent(tea.NewSDKError(map[string]interface{}{"code": "404", "message": "403"}), "ContactNotExists", "ContactGroupNotExists"); got {
		t.Fatal("group deletion / 404: got true, want false")
	}
}
