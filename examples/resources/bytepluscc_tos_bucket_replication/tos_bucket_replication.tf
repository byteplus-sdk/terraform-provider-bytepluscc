resource "bytepluscc_tos_bucket_replication" "Example" {
  role   = "ServiceRoleforReplicationAccessTOS"
  bucket = "example-source-bucket"
  rules = [{
    id         = "example-rule"
    prefix_set = ["prefix_1", "prefix_2"]
    status     = "Enabled"
    destination = {
      bucket                          = "example-destination-bucket"
      location                        = "ap-southeast-1"
      storage_class                   = "STANDARD"
      storage_class_inherit_directive = "SOURCE_OBJECT"
    }
    historical_object_replication = "Enabled"
    access_control_translation = {
      owner = "BucketOwnerEntrusted"
    }
  }]
}
