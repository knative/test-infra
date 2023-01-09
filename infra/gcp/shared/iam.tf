module "iam" {
  source  = "terraform-google-modules/iam/google//modules/projects_iam"
  version = "~> 7"

  projects = [module.project.project_id]

  mode = "authoritative"

  bindings = {

    "roles/artifactregistry.reader" = [
      "serviceAccount:${google_service_account.gke_nodes.email}",
    ]

    "roles/logging.logWriter" = [
      "serviceAccount:${google_service_account.gke_nodes.email}",
    ]

    "roles/monitoring.metricWriter" = [
      "serviceAccount:${google_service_account.gke_nodes.email}",
    ]

    "roles/monitoring.viewer" = [ # Required for Managed Prometheus
      "serviceAccount:${google_service_account.gke_nodes.email}",
    ]

    "roles/viewer" = [ # Will change this in the near future
      "group:wg-leads@knative.team",
      "group:gke-security-groups@knative.dev",
      "group:k8s-infra-rbac-perf-tests@knative.dev",
    ]
    # "roles/container.viewer" = [
    #   "group:gke-security-groups@knative.dev"
    # ]
  }
}

resource "google_service_account" "gke_nodes_contribfest" {
  account_id   = "gke-nodes"
  display_name = "GKE Nodes"
  project      = module.contribfest.project_id
}


module "iam_contribfest" {
  source  = "terraform-google-modules/iam/google//modules/projects_iam"
  version = "~> 7"

  projects = [module.contribfest.project_id]

  mode = "authoritative"

  bindings = {

    "roles/artifactregistry.reader" = [
      "serviceAccount:${google_service_account.gke_nodes_contribfest.email}",
    ]

    "roles/logging.logWriter" = [
      "serviceAccount:${google_service_account.gke_nodes_contribfest.email}",
    ]

    "roles/monitoring.metricWriter" = [
      "serviceAccount:${google_service_account.gke_nodes_contribfest.email}",
    ]

    "roles/monitoring.viewer" = [ # Required for Managed Prometheus
      "serviceAccount:${google_service_account.gke_nodes_contribfest.email}",
    ]

    "roles/viewer" = [ # Will change this in the near future
      "group:wg-leads@knative.team",
      # "group:gke-security-groups@knative.dev"
    ]
    # "roles/container.viewer" = [
    #   "group:gke-security-groups@knative.dev"
    # ]
  }
}

resource "google_service_account" "gke_nodes" {
  account_id   = "gke-nodes"
  display_name = "GKE Nodes"
  project      = module.project.project_id
}
