// WARNING, MAKE SURE YOU DON"T DESTROY THIS CLUSTER ACCIDENTALLY

module "shared" {
  source                               = "terraform-google-modules/kubernetes-engine/google//modules/beta-public-cluster"
  version                              = "~> 21"
  project_id                           = module.project.project_id
  name                                 = "shared"
  region                               = "us-central1"
  release_channel                      = "RAPID"
  network                              = module.vpc.network_name
  subnetwork                           = module.vpc.subnets["us-central1/gke-subnet-01"].name
  ip_range_pods                        = "gke-pods-01"
  ip_range_services                    = "gke-services-01"
  http_load_balancing                  = true
  network_policy                       = false
  horizontal_pod_autoscaling           = true
  filestore_csi_driver                 = false
  create_service_account               = false
  remove_default_node_pool             = true
  gce_pd_csi_driver                    = true
  # monitoring_enable_managed_prometheus = true
  authenticator_security_group         = "gke-security-groups@knative.dev"
  cluster_autoscaling = {
    enabled             = false
    autoscaling_profile = "OPTIMIZE_UTILIZATION"
    gpu_resources       = []
    max_cpu_cores       = null
    max_memory_gb       = null
    min_cpu_cores       = null
    min_memory_gb       = null
  }

  cluster_resource_labels = {
    cluster     = "shared"
    role        = "shared"
    environment = "production"
  }

  node_pools = [
    {
      name               = "main-pool-v1"
      machine_type       = "e2-standard-2"
      node_locations     = "us-central1-c,us-central1-f"
      min_count          = 1
      max_count          = 2
      disk_size_gb       = 100
      disk_type          = "pd-standard"
      image_type         = "COS_CONTAINERD"
      auto_repair        = true
      auto_upgrade       = true
      service_account    = google_service_account.gke_nodes.email
      enable_secure_boot = true
      initial_node_count = 1
    },
  ]

  node_pools_labels = {
    all = {
      environment = "production"
    }
  }
}
