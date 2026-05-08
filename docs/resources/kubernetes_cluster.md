---
subcategory: "Container"
page_title: "AzureSim: azuresim_kubernetes_cluster"
description: |-
  Manages a simulated Azure Kubernetes Service (AKS) cluster.
---

# azuresim_kubernetes_cluster

Manages a simulated Azure Kubernetes Service (AKS) cluster.

This resource mimics the behavior of the [`azurerm_kubernetes_cluster`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/kubernetes_cluster) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"
}

resource "azuresim_kubernetes_cluster" "example" {
  name                = "aks-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  dns_prefix          = "aksexample"
  kubernetes_version  = "1.30.0"
  sku_tier            = "Standard"

  default_node_pool {
    name       = "system"
    vm_size    = "Standard_D2s_v5"
    node_count = 3
    os_disk_size_gb = 128
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    network_policy    = "calico"
    service_cidr      = "10.100.0.0/16"
    dns_service_ip    = "10.100.0.10"
    load_balancer_sku = "standard"
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) Cluster name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `location` - (Required, ForceNew) Azure region.
* `dns_prefix` - (Required, ForceNew) DNS prefix used in the API server FQDN.
* `kubernetes_version` - (Optional) Kubernetes version (e.g. `1.30.0`).
* `sku_tier` - (Optional) `Free`, `Standard`, or `Premium`.
* `private_cluster_enabled` - (Optional) Enable private cluster.
* `default_node_pool` - (Required) One `default_node_pool` block as defined below.
* `identity` - (Optional) One `identity` block as defined below.
* `network_profile` - (Optional) One `network_profile` block as defined below.
* `tags` - (Optional) Tags.

---

A `default_node_pool` block supports:

* `name` - (Required) Node pool name.
* `vm_size` - (Required) VM size (e.g. `Standard_D2s_v5`).
* `node_count` - (Optional) Initial node count.
* `min_count` - (Optional) Min nodes (autoscale).
* `max_count` - (Optional) Max nodes (autoscale).
* `vnet_subnet_id` - (Optional) Subnet ID for the node pool.
* `os_disk_size_gb` - (Optional) OS disk size in GB.
* `os_disk_type` - (Optional) `Managed` or `Ephemeral`.
* `zones` - (Optional) List of availability zones.

---

An `identity` block supports:

* `type` - (Required) `SystemAssigned` or `UserAssigned`.
* `identity_ids` - (Optional) User-assigned identity IDs (when `type = "UserAssigned"`).

---

A `network_profile` block supports:

* `network_plugin` - (Required) `azure`, `kubenet`, or `none`.
* `network_policy` - (Optional) `calico`, `azure`, or `cilium`.
* `service_cidr` - (Optional) Service CIDR.
* `dns_service_ip` - (Optional) DNS service IP.
* `load_balancer_sku` - (Optional) `basic` or `standard`.

## Attributes Reference

* `id` - The AKS Cluster ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.ContainerService/managedClusters/{name}
  ```

* `node_resource_group` - Auto-generated node resource group name (`MC_<rg>_<name>_<location>`).
* `fqdn` - Simulated public API server FQDN (`<dns_prefix>-deadbeef.hcp.<location>.azmk8s.io`).
* `private_fqdn` - Simulated private API server FQDN (when `private_cluster_enabled = true`): `<dns_prefix>-deadbeef.privatelink.<location>.azmk8s.io`.

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
