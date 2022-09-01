# We grant all the boskos permissions on the folder and they get inherited
module "folder_iam" {
  source  = "terraform-google-modules/iam/google//modules/folders_iam"
  version = "~> 7.4"
  folders = [google_folder.boskos.id]
  mode    = "authoritative"

  bindings = {
    "roles/editor" = [
      "serviceAccount:prow-job@knative-tests.iam.gserviceaccount.com",
      "serviceAccount:prow-job@knative-nightly.iam.gserviceaccount.com",
      "serviceAccount:prow-job@knative-releases.iam.gserviceaccount.com",
    ]

    "roles/viewer" = [
      "group:knative-dev@googlegroups.com",
    ]

    "roles/container.admin" = [
      "serviceAccount:prow-job@knative-tests.iam.gserviceaccount.com",
      "serviceAccount:prow-job@knative-nightly.iam.gserviceaccount.com",
      "serviceAccount:prow-job@knative-releases.iam.gserviceaccount.com",
    ]

    "roles/logging.configWriter" = [
      "serviceAccount:prow-job@knative-tests.iam.gserviceaccount.com",
      "serviceAccount:prow-job@knative-nightly.iam.gserviceaccount.com",
      "serviceAccount:prow-job@knative-releases.iam.gserviceaccount.com",
    ]
    "roles/pubsub.admin" = [
      "serviceAccount:prow-job@knative-tests.iam.gserviceaccount.com",
      "serviceAccount:prow-job@knative-nightly.iam.gserviceaccount.com",
      "serviceAccount:prow-job@knative-releases.iam.gserviceaccount.com",
    ]

    "roles/storage.admin" = [
      "serviceAccount:prow-job@knative-tests.iam.gserviceaccount.com",
      "serviceAccount:prow-job@knative-nightly.iam.gserviceaccount.com",
      "serviceAccount:prow-job@knative-releases.iam.gserviceaccount.com",
    ]
  }
}
