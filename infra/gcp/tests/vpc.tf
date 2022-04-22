module "vpc" {
  source  = "terraform-google-modules/network/google"
  version = "~> 5"

  project_id   = module.project.project_id
  network_name = "prow"
  routing_mode = "GLOBAL"

  subnets = [
    {
      subnet_name           = "prow-subnet-01"
      subnet_ip             = "10.250.0.0/20"
      subnet_region         = "us-central1"
      subnet_private_access = "true"

    },
  ]

  secondary_ranges = {
    prow-subnet-01 = [
      {
        range_name    = "prow-build-services"
        ip_cidr_range = "10.250.32.0/20"
      },
      {
        range_name    = "prow-build-pods"
        ip_cidr_range = "10.250.64.0/20"
      },
      {
        range_name    = "prow-services"
        ip_cidr_range = "10.250.128.0/20"
      },
      {
        range_name    = "prow-pods"
        ip_cidr_range = "10.250.160.0/20"
      },
      {
        range_name    = "prow-trusted-services"
        ip_cidr_range = "10.250.192.0/20"
      },
      {
        range_name    = "prow-trusted-pods"
        ip_cidr_range = "10.250.224.0/20"
      },
    ]
  }
}
