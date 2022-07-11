module "vpc" {
  source  = "terraform-google-modules/network/google"
  version = "~> 5"

  project_id   = module.project.project_id
  network_name = "gke"
  routing_mode = "GLOBAL"

  subnets = [
    {
      subnet_name           = "gke-subnet-01"
      subnet_ip             = "10.251.0.0/20"
      subnet_region         = "us-central1"
      subnet_private_access = "true"
    },
  ]

  secondary_ranges = {
    gke-subnet-01 = [
      {
        range_name    = "gke-services-01"
        ip_cidr_range = "10.251.32.0/20"
      },
      {
        range_name    = "gke-pods-01"
        ip_cidr_range = "10.251.64.0/20"
      },
    ]
  }
}
