module "cosign" {
  source  = "terraform-google-modules/kms/google"
  version = "~> 2.2"

  project_id          = module.project.project_id
  purpose             = "ASYMMETRIC_SIGN"
  key_algorithm       = "EC_SIGN_P384_SHA384"
  key_rotation_period = ""
  location            = "global"
  keyring             = "cosign"
  keys                = ["signing-key-v2"]
}

module "cosign_iam" {
  source        = "terraform-google-modules/iam/google//modules/kms_key_rings_iam"
  version       = "~> 7"
  kms_key_rings = [module.cosign.keyring_resource.id]
  mode          = "authoritative"

  bindings = {
    "roles/cloudkms.signerVerifier" = [
      "serviceAccount:${google_service_account.prow_job.email}"
    ]
  }
}
