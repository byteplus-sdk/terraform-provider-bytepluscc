resource "bytepluscc_tos_object" "Example" {
  bucket       = "example-bucket"
  key          = "example-key.xml"
  content      = "<test>content</test>"
  content_type = "application/xml"
  metadata = [{
    value = "meta-value1"
    key   = "meta-key1"
  }]
  storage_class          = "STANDARD"
  public_acl             = "private"
  server_side_encryption = "AES256"
  tags = [{
    value = "tag-value1"
    key   = "tag-key1"
  }]
  account_acl = [{
    acl_type   = "CanonicalUser"
    account_id = "replace-with-account-id"
    permission = "READ"
    }, {
    acl_type   = "CanonicalUser"
    account_id = "replace-with-account-id"
    permission = "WRITE"
  }]
}
