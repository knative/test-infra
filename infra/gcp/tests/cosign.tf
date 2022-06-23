module "cosign" {
  source  = "terraform-google-modules/kms/google"
  version = "~> 2.1"

  project_id          = module.project.project_id
  purpose             = "ASYMMETRIC_SIGN"
  key_rotation_period = "7776000s" # 90 days
  location            = "global"
  keyring             = "cosign"
  keys                = ["signing-key-v2"]
}

module "cosign_iam" {
  source        = "terraform-google-modules/iam/google//modules/kms_key_rings_iam"
  version       = "~> 7"
  kms_key_rings = [module.cosign.keyring_name]
  mode          = "authoritative"

  bindings = {
    "roles/cloudkms.signerVerifier" = [
      "serviceAccount:${google_service_account.prow_job.email}"
    ]
  }
}
