// nightly Buckets
resource "google_storage_bucket" "nightly" {
  name     = "knative-nightly"
  location = "US"
  storage_class = "MULTI_REGIONAL"
  project  = "knative-nightly"

  uniform_bucket_level_access = true
}

module "buckets_nightly" {
  source          = "terraform-google-modules/iam/google//modules/storage_buckets_iam"
  storage_buckets = [google_storage_bucket.nightly.name]
  mode            = "authoritative"

  bindings = {
    "roles/storage.legacyBucketOwner" = [
      "projectEditor:knative-nightly",
      "projectOwner:knative-nightly"
    ]
    "roles/storage.legacyBucketReader" = [
      "projectViewer:knative-nightly"
    ]
    "roles/storage.objectViewer" = [
      "allUsers"
    ]
  }
}
