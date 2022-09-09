module "iam" {
  source  = "terraform-google-modules/iam/google//modules/projects_iam"
  version = "~> 7"

  projects = ["knative-releases"]

  mode = "authoritative"

  bindings = {

    "roles/storage.admin" = [
      "serviceAccount:${google_service_account.prow_job.email}",
    ]
    "projects/knative-releases/roles/ServiceAccountIAMEditor" = [
      "serviceAccount:${google_service_account.prow_job.email}",
    ]
  }
}

// Service Account used by Knative Releases
resource "google_service_account" "prow_job" {
  account_id   = "prow-job"
  display_name = "Prow Job Knative Release Creator"
  project      = module.project.project_id
}

resource "google_service_account_iam_binding" "prow_job" {
  service_account_id = google_service_account.prow_job.name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "serviceAccount:knative-tests.svc.id.goog[test-pods/release]",
  ]
}

resource "google_service_account" "signer" {
  account_id   = "signer"
  display_name = "Used for signing release images with cosign"
  project      = module.project.project_id
}

resource "google_service_account_iam_binding" "signer" {
  service_account_id = google_service_account.signer.name
  role               = "roles/iam.serviceAccountTokenCreator"

  members = [
    "serviceAccount:${google_service_account.prow_job.email}",
    // DO NOT ADD ANY OTHER PRINCIPAL HERE AS THIS ROLE ALLOWS OIDC TOKENS TO BE CREATED FOR THIS USER
  ]
}
