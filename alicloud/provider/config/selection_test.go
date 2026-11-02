package config

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The static registration checker relies on this constructor being configuration
// only: resource declarations belong to product ServicePackages, never here.
func TestUnitSharedConfigurationHasNoRegistrations(t *testing.T) {
	p := Provider()
	if len(p.ResourcesMap) != 0 || len(p.DataSourcesMap) != 0 {
		t.Fatal("shared provider configuration contains product registrations")
	}
}

func TestUnitConfigurePreservesRegistrationMaps(t *testing.T) {
	resource := &schema.Resource{}
	dataSource := &schema.Resource{}
	p := &schema.Provider{ResourcesMap: map[string]*schema.Resource{"alicloud_legacy": resource}, DataSourcesMap: map[string]*schema.Resource{"alicloud_legacy": dataSource}}
	Configure(p)
	if p.ResourcesMap["alicloud_legacy"] != resource || p.DataSourcesMap["alicloud_legacy"] != dataSource {
		t.Fatal("shared configuration replaced caller registrations")
	}
	if p.Schema["region"] == nil || p.ConfigureContextFunc == nil {
		t.Fatal("shared configuration was not installed")
	}
}
