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

    // Existing bindings that should be removed I think?
    "serviceAccount:knative-boskos-01.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-01.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-01.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-01.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-02.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-02.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-02.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-02.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-03.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-03.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-04.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-04.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-04.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-04.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-05.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-05.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-05.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-05.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-06.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-06.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-07.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-07.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-08.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-08.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-08.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-08.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-09.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-09.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-09.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-09.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-10.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-10.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-10.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-10.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-100.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-100.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-100.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-100.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-11.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-11.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-11.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-11.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-12.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-12.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-12.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-12.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-13.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-13.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-14.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-14.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-14.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-14.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-15.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-15.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-15.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-15.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-16.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-16.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-17.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-17.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-17.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-17.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-18.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-18.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-18.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-18.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-19.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-19.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-19.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-19.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-20.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-20.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-21.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-21.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-21.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-21.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-22.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-22.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-23.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-23.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-23.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-23.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-24.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-24.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-25.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-25.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-26.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-26.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-26.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-26.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-27.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-27.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-28.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-28.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-29.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-29.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-29.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-29.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-30.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-30.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-31.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-31.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-31.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-31.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-32.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-32.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-33.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-33.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-34.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-34.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-35.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-35.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-35.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-35.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-36.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-36.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-36.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-36.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-37.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-37.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-37.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-37.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-38.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-38.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-38.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-38.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-39.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-39.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-39.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-39.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-40.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-40.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-41.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-41.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-42.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-42.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-43.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-43.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-44.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-44.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-44.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-44.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-44.svc.id.goog[test-pull-subscription-with-target-rgk4w/ksa-name]",
    "serviceAccount:knative-boskos-44.svc.id.goog[test-smoke-cloud-scheduler-source-ljx54/ksa-name]",
    "serviceAccount:knative-boskos-44.svc.id.goog[test-smoke-cloud-storage-source-4srz6/ksa-name]",
    "serviceAccount:knative-boskos-45.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-45.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-45.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-45.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-46.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-46.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-47.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-47.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-47.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-47.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-48.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-48.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-48.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-48.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-49.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-49.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-49.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-49.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-50.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-50.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-50.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-50.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-51.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-51.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-51.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-51.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-52.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-52.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-53.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-53.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-54.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-54.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-54.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-54.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-55.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-55.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-55.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-55.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-56.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-56.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-57.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-57.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-58.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-58.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-58.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-58.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-59.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-59.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-59.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-59.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-60.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-60.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-61.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-61.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-61.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-61.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-62.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-62.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-62.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-62.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-63.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-63.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-63.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-63.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-64.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-64.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-64.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-64.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-65.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-65.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-65.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-65.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-66.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-66.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-67.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-67.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-67.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-67.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-68.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-68.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-68.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-68.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-69.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-69.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-69.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-69.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-70.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-70.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-70.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-70.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-71.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-71.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-71.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-71.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-72.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-72.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-72.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-72.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-73.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-73.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-73.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-73.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-74.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-74.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-75.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-75.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-75.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-75.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-76.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-76.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-76.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-76.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-77.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-77.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-78.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-78.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-79.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-79.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-80.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-80.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-81.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-81.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-82.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-82.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-82.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-82.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-83.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-83.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-83.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-83.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-84.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-84.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-85.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-85.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-85.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-85.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-86.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-86.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-86.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-86.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-87.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-87.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-87.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-87.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-88.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-88.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-88.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-88.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-89.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-89.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-89.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-89.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-90.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-90.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-90.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-90.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-91.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-91.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-91.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-91.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-92.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-92.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-92.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-92.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-93.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-93.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-94.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-94.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-94.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-94.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-95.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-95.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-96.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-96.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-97.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-97.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-97.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-97.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-98.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-98.svc.id.goog[cloud-run-events/controller]",
    "serviceAccount:knative-boskos-98.svc.id.goog[events-system/broker]",
    "serviceAccount:knative-boskos-98.svc.id.goog[events-system/controller]",
    "serviceAccount:knative-boskos-99.svc.id.goog[cloud-run-events/broker]",
    "serviceAccount:knative-boskos-99.svc.id.goog[cloud-run-events/controller]"

  ]
}
