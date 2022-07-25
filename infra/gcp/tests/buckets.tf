
resource "google_storage_bucket" "prow" {
  name                        = "knative-prow"
  location                    = "US"
  storage_class               = "STANDARD"
  uniform_bucket_level_access = true
  lifecycle_rule {
    action {
      storage_class = "NEARLINE"
      type          = "SetStorageClass"
    }
    condition {
      age                        = 180
    }
  }
  lifecycle_rule {
    action {
      type = "Delete"
    }

    condition {
      age                        = 210
    }
  }
}

module "iam_prow_bucket" {
  source          = "terraform-google-modules/iam/google//modules/storage_buckets_iam"
  storage_buckets = [google_storage_bucket.prow.name]
  version         = "~> 7"

  mode = "authoritative"

  bindings = {
    "roles/storage.objectAdmin" = [
      "serviceAccount:prow-job@knative-nightly.iam.gserviceaccount.com",
      "serviceAccount:prow-job@knative-releases.iam.gserviceaccount.com",
      "serviceAccount:${google_service_account.prow_job.email}",
      "serviceAccount:${google_service_account.prow_pod_utils.email}",
      "serviceAccount:${google_service_account.gsuite.email}",
    ]
    "roles/storage.objectViewer" = [
      "allUsers",
    ]
  }
}
