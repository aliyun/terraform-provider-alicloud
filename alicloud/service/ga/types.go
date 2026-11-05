package ga

// RenewalStatus and its values mirror the shared renewal-status declarations
// in the root package (alicloud/extension_ecs.go). The Ga accelerator and
// bandwidth package reference them in their renewal_status schema validators;
// the root package declarations are not reachable from this product
// subpackage, so the typed constants are replicated here verbatim.
type RenewalStatus string

const (
	RenewAutoRenewal = RenewalStatus("AutoRenewal")
	RenewNormal      = RenewalStatus("Normal")
	RenewNotRenewal  = RenewalStatus("NotRenewal")
)
