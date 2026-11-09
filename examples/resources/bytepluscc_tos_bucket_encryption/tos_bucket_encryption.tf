resource "bytepluscc_tos_bucket_encryption" "Example" {
  # Use a KMS key in the same region as the bucket.
  kms_master_key_id = "replace-with-kms-key-resource-id"
  name              = "example-bucket"
  sse_algorithm     = "kms"
}
