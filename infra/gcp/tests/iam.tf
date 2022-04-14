module "iam" {
  source  = "terraform-google-modules/iam/google//modules/projects_iam"
  version = "~> 7"

  projects = ["knative-tests"]

  mode = "authoritative"

  bindings = {
    "roles/storage.objectViewer" = [
      "serviceAccount:${google_service_account.gke_nodes.email}",
    ]

    "roles/logging.logWriter" = [
      "serviceAccount:${google_service_account.gke_nodes.email}",
    ]

    "roles/monitoring.metricWriter" = [
      "serviceAccount:${google_service_account.gke_nodes.email}",
    ]

    "roles/stackdriver.resourceMetadata.writer" = [
      "serviceAccount:${google_service_account.gke_nodes.email}",
    ]

    "roles/artifactregistry.reader" = [
      "serviceAccount:${google_service_account.gke_nodes.email}",
    ]

    "roles/cloudbuild.builds.editor" = [
      "serviceAccount:prow-job@knative-tests.iam.gserviceaccount.com",
      "serviceAccount:pwg-admins@knative-tests.iam.gserviceaccount.com"
    ]

    "roles/container.admin" : [
      "serviceAccount:prow-deployer@knative-tests.iam.gserviceaccount.com",
      "serviceAccount:pwg-admins@knative-tests.iam.gserviceaccount.com"
    ]
  }
}

resource "google_service_account" "gke_nodes" {
  account_id   = "gke-nodes"
  display_name = "GKE Nodes"
  project      = module.project.project_id
}
