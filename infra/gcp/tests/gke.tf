// WARNING, MAKE SURE YOU DON"T DESTROY THESE CLUSTERS ACCIDENTALLY


module "prow_trusted" {
  source                     = "terraform-google-modules/kubernetes-engine/google"
  project_id                 = module.project.project_id
  name                       = "prow-trusted"
  regional                   = false
  zones                      = ["us-central1-a"]
  release_channel            = "STABLE"
  network                    = "default"
  subnetwork                 = "default"
  ip_range_pods              = "gke-prow-trusted-pods-f579223d"
  ip_range_services          = "gke-prow-trusted-services-f579223d"
  http_load_balancing        = true
  network_policy             = false
  horizontal_pod_autoscaling = true
  filestore_csi_driver       = false
  create_service_account     = false
  cluster_resource_labels = {
    cluster     = "prow-trusted"
    role        = "prow"
    environment = "production"
  }

  node_pools = [
    {
      name               = "prod-v1"
      machine_type       = "e2-medium"
      node_locations     = "us-central1-f"
      min_count          = 1
      max_count          = 3
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
}

module "prow" {
  source                     = "terraform-google-modules/kubernetes-engine/google"
  project_id                 = module.project.project_id
  name                       = "prow"
  region                     = "us-central1"
  release_channel            = "RAPID"
  network                    = "default"
  subnetwork                 = "default"
  ip_range_pods              = ""
  ip_range_services          = ""
  http_load_balancing        = true
  network_policy             = false
  horizontal_pod_autoscaling = true
  filestore_csi_driver       = false
  create_service_account     = false
  cluster_resource_labels = {
    cluster     = "prow"
    role        = "prow"
    environment = "production"
  }

  node_pools = [
    {
      name               = "prod-v1"
      machine_type       = "e2-standard-4"
      node_locations     = "us-central1-a,us-central1-b"
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

module "prow_build" {
  source                     = "terraform-google-modules/kubernetes-engine/google"
  project_id                 = module.project.project_id
  name                       = "prow-build"
  region                     = "us-central1"
  release_channel            = "RAPID"
  network                    = "default"
  subnetwork                 = "default"
  ip_range_pods              = ""
  ip_range_services          = ""
  http_load_balancing        = true
  network_policy             = false
  horizontal_pod_autoscaling = true
  filestore_csi_driver       = false
  create_service_account     = false
  cluster_resource_labels = {
    cluster     = "prow-build"
    role        = "prow"
    environment = "production"
  }

  node_pools = [
    {
      name               = "system-pool-v1"
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
    {
      name               = "testing-pool-v1"
      machine_type       = "e2-standard-16"
      node_locations     = "us-central1-c,us-central1-f"
      local_ssd_count    = 1
      min_count          = 0
      max_count          = 4
      disk_size_gb       = 100
      disk_type          = "pd-ssd"
      image_type         = "COS_CONTAINERD"
      auto_repair        = true
      auto_upgrade       = true
      service_account    = google_service_account.gke_nodes.email
      enable_secure_boot = true
      initial_node_count = 0
    },
  ]

  node_pools_labels = {
    all = {
      environment = "production"
    }

    testing-pool-v1 = {
      type = "testing"
    }

    system-pool-v1 = {
      type = "system"
    }
  }
}
