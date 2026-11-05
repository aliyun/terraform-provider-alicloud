package expressconnect

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Router-interface helpers relocated from the root alicloud package
// (extension_vpc.go, diff_suppress_funcs.go). The leaf-package rule forbids
// importing root alicloud, so the minimal set consumed by the
// alicloud_router_interface resource and its tests is re-declared here with
// identical values and bodies. RouterType / VRouter / VBR live in
// common_expressconnect.go.

type Role string
type Spec string

const (
	InitiatingSide = Role("InitiatingSide")
	AcceptingSide  = Role("AcceptingSide")

	Mini2   = Spec("Mini.2")
	Mini5   = Spec("Mini.5")
	Small1  = Spec("Small.1")
	Small2  = Spec("Small.2")
	Small5  = Spec("Small.5")
	Middle1 = Spec("Middle.1")
	Middle2 = Spec("Middle.2")
	Middle5 = Spec("Middle.5")
	Large1  = Spec("Large.1")
	Large2  = Spec("Large.2")
	Large5  = Spec("Large.5")
	Xlarge1 = Spec("Xlarge.1")

	Negative = Spec(("Negative"))
)

func GetAllRouterInterfaceSpec() (specifications []string) {
	specifications = append(specifications, string(Mini2), string(Mini5),
		string(Small1), string(Small2), string(Small5),
		string(Middle1), string(Middle2), string(Middle5),
		string(Large1), string(Large2), string(Large5), string(Xlarge1),
		string(Negative))
	return
}

func PostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	// payment_type is the instance_charge_type's replacement.
	// If both instance_charge_type and payment_type are "", it means hiding a default "PostPaid"
	if v, ok := d.GetOk("instance_charge_type"); ok && strings.ToLower(v.(string)) == "prepaid" {
		return false
	}
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
		return false
	}
	return true
}

func routerInterfaceAcceptsideDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("role").(string) == string(AcceptingSide)
}

func routerInterfaceVBRTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("role").(string) == string(AcceptingSide) {
		return true
	}
	if d.Get("router_type").(string) == string(VRouter) {
		return true
	}
	return false
}
