resource "bytepluscc_tos_bucket" "Example" {
  name                  = "example-bucket"
  storage_class         = "STANDARD"
  enable_version_status = "Enabled"
  bucket_type           = "fns"
  tags = [
    {
      key   = "env"
      value = "test"
    }
  ]
  # Replace the principal and bucket resource identifiers with your own values.
  policy = "{\"Version\":\"1.0\",\"Statement\":[{\"Sid\":\"example\",\"Effect\":\"Allow\",\"Principal\":[\"replace-with-account-id\"],\"Action\":[\"tos:Get*\",\"tos:List*\",\"tos:HeadBucket\"],\"Resource\":[\"replace-with-bucket-resource-id\",\"replace-with-bucket-resource-id/*\"]}]}"
  lifecycle_config = [
    {
      expiration = {
        date = "2027-01-18T00:00:00Z"
      }
      filter = {
        greater_than_include_equal = "Disabled"
        less_than_include_equal    = "Disabled"
        object_size_greater_than   = 123
        object_size_less_than      = 789
      }
      lifecycle_rule_id = "example-rule"
      prefix            = "prefix"
      status            = "Enabled"
    }
  ]
  object_lock_configuration = {
    object_lock_enabled = "Enabled"
    rule = {
      default_retention = {
        mode = "COMPLIANCE"
        days = 21
      }
    }
  }
}
