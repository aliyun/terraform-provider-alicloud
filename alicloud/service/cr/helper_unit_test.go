package cr

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// unitTestProvider mirrors the root Provider().ResourcesMap for the migrated
// cr types. The split-package internal test cannot reach the root Provider:
// importing acctest pulls in the root alicloud package, which imports the cr
// service via the registry, re-introducing an import cycle. The unit tests
// therefore look the cr resource schemas up through this local provider built
// from the cr ServicePackage factories. Only the cr factories are wired here,
// which is the only set the unit tests inspect.
func unitTestProvider() *schema.Provider {
	p := &schema.Provider{
		ResourcesMap:   map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{},
	}
	sp := ServicePackage()
	for _, r := range sp.SDKResources {
		p.ResourcesMap[r.TypeName] = r.Factory()
	}
	for _, d := range sp.SDKDataSources {
		p.DataSourcesMap[d.TypeName] = d.Factory()
	}
	return p
}
