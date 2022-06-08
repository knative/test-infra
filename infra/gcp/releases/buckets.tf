// Releases Buckets
resource "google_storage_bucket" "releases" {
  name     = "knative-releases"
  location = "US"
  project  = "knative-releases"

  uniform_bucket_level_access = true

  logging {
    log_bucket        = "knative-releases-analysis"
    log_object_prefix = "knative-releases"
  }

  versioning {
    enabled = true
  }
}

module "buckets_releases" {
  source          = "terraform-google-modules/iam/google//modules/storage_buckets_iam"
  storage_buckets = [google_storage_bucket.releases.name]
  mode            = "authoritative"

  bindings = {
    "roles/storage.legacyBucketOwner" = [
      "projectEditor:knative-releases",
      "projectOwner:knative-releases"
    ]
    "roles/storage.legacyBucketReader" = [
      "projectViewer:knative-releases"
    ]

    "roles/storage.objectAdmin" = [
      "user:mark@chmarny.com"
    ]

    "roles/storage.objectViewer" = [
      "allUsers"
    ]
  }
}
