---
subcategory: "Container"
page_title: "AzureSim: azuresim_container_app"
description: |-
  Manages a simulated Azure Container App.
---

# azuresim_container_app

Manages a simulated Azure Container App.

This resource mimics the behavior of the [`azurerm_container_app`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/container_app) resource. It manages all state within Terraform's state file and does not make any API calls to Azure.

## Example Usage

```terraform
resource "azuresim_container_app" "example" {
  name                         = "ca-example"
  resource_group_name          = azuresim_resource_group.example.name
  container_app_environment_id = azuresim_container_app_environment.example.id
  revision_mode                = "Single"

  template {
    min_replicas = 1
    max_replicas = 10

    container {
      name   = "app"
      image  = "nginx:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }

  ingress {
    external_enabled = true
    target_port      = 80
    transport        = "auto"

    traffic_weight {
      percentage      = 100
      latest_revision = true
    }
  }

  tags = {
    environment = "production"
  }
}
```

## Argument Reference

* `name` - (Required, ForceNew) App name.
* `resource_group_name` - (Required, ForceNew) Resource Group.
* `container_app_environment_id` - (Required, ForceNew) Parent Container App Environment ID.
* `revision_mode` - (Required) `Single` or `Multiple`.
* `workload_profile_name` - (Optional) Workload profile name from the parent environment.
* `template` - (Required) One `template` block as defined below.
* `ingress` - (Optional) Zero or one `ingress` block as defined below.
* `tags` - (Optional) Tags.

---

A `template` block supports:

* `min_replicas` - (Optional) Minimum replicas.
* `max_replicas` - (Optional) Maximum replicas.
* `revision_suffix` - (Optional) Suffix appended to revision names.
* `termination_grace_period_seconds` - (Optional) Pod termination grace period.
* `container` - (Required) One or more `container` blocks as defined below.

A `container` block (under `template`) supports:

* `name` - (Required) Container name.
* `image` - (Required) Container image (e.g. `nginx:latest`).
* `cpu` - (Required) vCPU (e.g. `0.25`).
* `memory` - (Required) Memory (e.g. `0.5Gi`).
* `args` - (Optional) Container args.
* `command` - (Optional) Container entrypoint command.

---

An `ingress` block supports:

* `external_enabled` - (Optional) Whether the app is exposed externally.
* `target_port` - (Required) Target container port.
* `transport` - (Optional) `auto`, `http`, `http2`, or `tcp`.
* `allow_insecure_connections` - (Optional) Allow HTTP.
* `traffic_weight` - (Optional) One or more `traffic_weight` blocks as defined below.

A `traffic_weight` block (under `ingress`) supports:

* `percentage` - (Required) Weight 0-100.
* `latest_revision` - (Optional) Apply weight to the latest revision.
* `revision_suffix` - (Optional) Apply weight to a specific revision suffix.
* `label` - (Optional) Traffic label.

## Attributes Reference

* `id` - The Container App ID:

  ```
  /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.App/containerApps/{name}
  ```

* `latest_revision_fqdn` - Simulated FQDN of the latest revision (`<name>.azurecontainerapps.io`).

## Timeouts

This resource does not support timeouts. All operations complete instantaneously as no external API calls are made.

## Import

This resource does not currently support import.
