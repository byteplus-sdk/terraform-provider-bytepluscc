resource "bytepluscc_tos_bucket_inventory" "Example" {
  bucket_name = "example-bucket"
  destination = {
    tos_bucket_destination = {
      account_id = "replace-with-account-id"
      bucket     = "example-inventory-bucket"
      format     = "CSV"
      prefix     = "tos_bucket_inventory"
      role       = "replace-with-inventory-role-name"
    }
  }
  filter = {
    prefix = "inventory/"
  }
  included_object_versions = "Current"
  inventory_id             = "example-inventory"
  is_enabled               = true
  is_un_compressed         = false
  optional_fields = {
    field = ["StorageClass", "IntelligentTieringAccessTier"]
  }
  schedule = {
    frequency = "Weekly"
  }
}
