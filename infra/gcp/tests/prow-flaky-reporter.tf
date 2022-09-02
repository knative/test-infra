resource "google_service_account_iam_binding" "flaky_test_reporter" {
  service_account_id = google_service_account.flaky_test_reporter.name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "serviceAccount:knative-tests.svc.id.goog[test-pods/flaky-test-reporter]",
  ]
}

resource "google_service_account" "flaky_test_reporter" {
  account_id   = "flaky-test-reporter"
  display_name = "Flaky Test Reporter"
  project      = "knative-tests"
}
