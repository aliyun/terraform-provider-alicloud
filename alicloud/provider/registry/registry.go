// Package registry aggregates the service package declarations consumed by the
// provider assembly loop, and validates uniqueness, this being the only place both
// provider servers agree on.
package registry

import (
	"fmt"
	"strings"
	"sync"

	alicloudfunction "github.com/aliyun/terraform-provider-alicloud/alicloud/function"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ims"
)

// A managed resource and a data source may share a name; two managed resources may
// not, whichever server implements them.
const (
	nsResource          = "resource"
	nsDataSource        = "data source"
	nsFunction          = "function"
	nsEphemeralResource = "ephemeral resource"
	nsListResource      = "list resource"
	nsAction            = "action"
)

var servicePackages = []conns.ServicePackage{
	ims.ServicePackage(),
	alicloudfunction.ServicePackage(),
}

var validateOnce sync.Once

// Validates the set once on first use, and panics on a violation: failing at
// startup beats a nil dereference on a later RPC.
func ServicePackages() []conns.ServicePackage {
	validateOnce.Do(func() {
		if err := validate(servicePackages); err != nil {
			panic("provider registration error: " + err.Error())
		}
	})
	return servicePackages
}

// Reports every violation at once. A declaration colliding with a provider.go map
// literal is assembleProvider's to find, not visible here.
func validate(sps []conns.ServicePackage) error {
	v := &validator{seen: map[string]map[string]string{}}
	for _, sp := range sps {
		for _, d := range sp.SDKResources {
			v.claim(nsResource, d.TypeName, sp.Name, "SDKResources", d.Factory != nil)
		}
		for _, d := range sp.Resources {
			v.claim(nsResource, d.TypeName, sp.Name, "Resources", d.Factory != nil)
		}
		for _, d := range sp.SDKDataSources {
			v.claim(nsDataSource, d.TypeName, sp.Name, "SDKDataSources", d.Factory != nil)
		}
		for _, d := range sp.DataSources {
			v.claim(nsDataSource, d.TypeName, sp.Name, "DataSources", d.Factory != nil)
		}
		for _, d := range sp.Functions {
			v.claim(nsFunction, d.TypeName, sp.Name, "Functions", d.Factory != nil)
		}
		for _, d := range sp.EphemeralResources {
			v.claim(nsEphemeralResource, d.TypeName, sp.Name, "EphemeralResources", d.Factory != nil)
		}
		for _, d := range sp.ListResources {
			v.claim(nsListResource, d.TypeName, sp.Name, "ListResources", d.Factory != nil)
		}
		for _, d := range sp.Actions {
			v.claim(nsAction, d.TypeName, sp.Name, "Actions", d.Factory != nil)
		}
	}
	return v.err()
}

type validator struct {
	seen     map[string]map[string]string // namespace → type name → "Service.Category"
	problems []string
}

func (v *validator) claim(namespace, typeName, service, category string, hasFactory bool) {
	origin := service + "." + category
	if typeName == "" {
		v.problems = append(v.problems, fmt.Sprintf("%s declared in %s has an empty TypeName", namespace, origin))
		return
	}
	if !hasFactory {
		v.problems = append(v.problems, fmt.Sprintf("%s %q declared in %s has a nil Factory", namespace, typeName, origin))
	}
	names := v.seen[namespace]
	if names == nil {
		names = map[string]string{}
		v.seen[namespace] = names
	}
	if prev, dup := names[typeName]; dup {
		v.problems = append(v.problems, fmt.Sprintf("%s %q is declared twice: %s and %s", namespace, typeName, prev, origin))
		return
	}
	names[typeName] = origin
}

func (v *validator) err() error {
	if len(v.problems) == 0 {
		return nil
	}
	return fmt.Errorf("%d invalid service package declaration(s):\n  %s", len(v.problems), strings.Join(v.problems, "\n  "))
}
