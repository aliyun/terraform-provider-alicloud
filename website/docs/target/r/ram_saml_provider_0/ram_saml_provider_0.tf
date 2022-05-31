resource "alicloud_ram_saml_provider" "example" {
  saml_provider_name            = "tf-testAcc"
  encodedsaml_metadata_document = "your encodedsaml metadata document"
  description                   = "For Terraform Test"
}

