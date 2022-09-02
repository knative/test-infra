module "iam" {
  source  = "terraform-google-modules/iam/google//modules/projects_iam"
  version = "~> 7"

  projects = [module.project.project_id]

  mode = "authoritative"

  bindings = {

    "roles/storage.admin" = [
      "serviceAccount:${google_service_account.prow_job.email}",
    ]
    "projects/knative-nightly/roles/ServiceAccountIAMEditor" = [
      "serviceAccount:${google_service_account.prow_job.email}",
    ]
  }
}

// Service Account used by Knative Nightly Jobs
resource "google_service_account" "prow_job" {
  account_id   = "prow-job"
  display_name = "Prow Job Knative Nightly Release Creator"
  description  = "Service account for Prow Jobs that create Knative nightly releases"
  project      = module.project.project_id
}

resource "google_service_account_iam_binding" "prow_job" {
  service_account_id = google_service_account.prow_job.name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "serviceAccount:knative-tests.svc.id.goog[test-pods/nightly]",
  ]
}
