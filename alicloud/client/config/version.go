package config

const DefaultClientRetryCountSmall = 5

const DefaultClientRetryCountMedium = 10

const DefaultClientRetryCountLarge = 15

const Terraform = "HashiCorp-Terraform"

const Provider = "Terraform-Provider"

const Module = "Terraform-Module"

const TerraformTraceId = "TerraformTraceId"

// The main version number that is being run at the moment.
var providerVersion = "2.3.0"

func ProviderVersion() string {
	return providerVersion
}
