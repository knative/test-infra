resource "google_artifact_registry_repository" "images" {
  project       = module.project.project_id
  location      = "us"
  repository_id = "images"
  description   = "Nightly Knative Images"
  format        = "DOCKER"
}

resource "google_artifact_registry_repository_iam_member" "images" {
  project    = google_artifact_registry_repository.images.project
  location   = google_artifact_registry_repository.images.location
  repository = google_artifact_registry_repository.images.name
  role       = "roles/artifactregistry.reader"
  member     = "allUsers"
}
