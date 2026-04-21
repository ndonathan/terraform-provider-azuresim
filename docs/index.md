---
page_title: "Provider: AzureSim"
description: |-
  The AzureSim provider simulates Azure Resource Manager resources for use in Terraform without provisioning real infrastructure or requiring Azure credentials.
---

# AzureSim Provider

The AzureSim provider is a simulation layer that mimics the behavior of the [AzureRM provider](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs). It allows practitioners to write, plan, and apply Terraform configurations that follow Azure resource conventions without provisioning real cloud resources or authenticating to an Azure environment.

All resource lifecycle operations (create, read, update, delete) are handled entirely within Terraform's state file. No external API calls are made. Resources are assigned realistic Azure-style resource IDs that mirror the format used by the Azure Resource Manager API.

## Use Cases

* **Learning and training** - Practice writing Terraform configurations with Azure-style resources in an offline or sandboxed environment.
* **Module development** - Develop and test Terraform modules that target Azure without incurring costs or requiring credentials.
* **CI/CD pipeline validation** - Validate Terraform plan output and configuration structure in automated pipelines without Azure access.
* **Demos and workshops** - Run live demonstrations of Terraform workflows without depending on cloud connectivity.

## Example Usage

```terraform
terraform {
  required_providers {
    azuresim = {
      source = "ndonathan/azuresim"
    }
  }
}

provider "azuresim" {
  subscription_id = "12345678-1234-1234-1234-123456789012"
  tenant_id       = "87654321-4321-4321-4321-210987654321"
}

resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"

  tags = {
    environment = "dev"
  }
}
```

## Authentication

The AzureSim provider does not perform any authentication. The `subscription_id` and `tenant_id` arguments are accepted for compatibility with Azure-style configurations but are used only to generate realistic resource IDs. No credentials, service principals, or managed identities are required.

## Argument Reference

The following arguments are supported:

* `subscription_id` - (Optional) A simulated Azure Subscription ID. This value is embedded in all generated resource IDs. If not specified, defaults to `00000000-0000-0000-0000-000000000000`.

* `tenant_id` - (Optional) A simulated Azure Tenant ID. This value is accepted for configuration compatibility but is not used in resource ID generation.
