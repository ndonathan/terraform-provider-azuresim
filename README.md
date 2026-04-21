# Terraform Provider AzureSim

A simulated Azure provider for Terraform that mimics the [AzureRM provider](https://registry.terraform.io/providers/hashicorp/azurerm/latest) interface. It lets you plan and apply Azure-style infrastructure configurations without provisioning real resources or authenticating to Azure.

All resource state is managed entirely within Terraform's state file. No API calls are made.

## Use Cases

- **Learning Terraform** with Azure-style resources in an offline environment
- **Testing modules and CI pipelines** without Azure credentials or costs
- **Validating Terraform configuration structure** before targeting a real environment
- **Demos and workshops** where live infrastructure is impractical

## Supported Resources

| Resource | Type Name | Azure Equivalent |
|----------|-----------|------------------|
| Resource Group | `azuresim_resource_group` | `azurerm_resource_group` |
| Virtual Network | `azuresim_virtual_network` | `azurerm_virtual_network` |
| Subnet | `azuresim_subnet` | `azurerm_subnet` |
| Virtual Machine | `azuresim_virtual_machine` | `azurerm_linux_virtual_machine` |
| Storage Account | `azuresim_storage_account` | `azurerm_storage_account` |

Each resource generates a realistic Azure-style resource ID (e.g. `/subscriptions/<id>/resourceGroups/<name>/providers/Microsoft.Compute/virtualMachines/<name>`).

## Provider Configuration

```hcl
provider "azuresim" {
  subscription_id = "12345678-1234-1234-1234-123456789012"  # optional
  tenant_id       = "87654321-4321-4321-4321-210987654321"  # optional
}
```

Both attributes are optional. If `subscription_id` is omitted, a zeroed UUID is used in generated resource IDs.

## Quick Start

### 1. Build the provider

```sh
go build -o terraform-provider-azuresim .
```

### 2. Configure a dev override

Add the following to `~/.terraformrc` (update the path to where you built the binary):

```hcl
provider_installation {
  dev_overrides {
    "ndonathan/azuresim" = "/path/to/terraform-provider-azuresim"
  }
  direct {}
}
```

### 3. Write a Terraform configuration

```hcl
terraform {
  required_providers {
    azuresim = {
      source = "ndonathan/azuresim"
    }
  }
}

provider "azuresim" {
  subscription_id = "12345678-1234-1234-1234-123456789012"
}

resource "azuresim_resource_group" "example" {
  name     = "rg-example"
  location = "eastus"

  tags = {
    environment = "dev"
  }
}

resource "azuresim_virtual_network" "example" {
  name                = "vnet-example"
  resource_group_name = azuresim_resource_group.example.name
  location            = azuresim_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}
```

### 4. Apply

With dev overrides, `terraform init` is not required:

```sh
terraform apply
```

A full working example is available in the [`examples/`](examples/main.tf) directory.

## Resource Reference

### azuresim_resource_group

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | yes | Resource group name (forces replacement) |
| `location` | string | yes | Azure region (forces replacement) |
| `tags` | map(string) | no | Tags |
| `id` | string | computed | Azure-style resource ID |

### azuresim_virtual_network

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | yes | VNet name (forces replacement) |
| `resource_group_name` | string | yes | Parent resource group (forces replacement) |
| `location` | string | yes | Azure region (forces replacement) |
| `address_space` | list(string) | yes | CIDR blocks (e.g. `["10.0.0.0/16"]`) |
| `dns_servers` | list(string) | no | DNS server IPs |
| `tags` | map(string) | no | Tags |
| `id` | string | computed | Azure-style resource ID |

### azuresim_subnet

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | yes | Subnet name (forces replacement) |
| `resource_group_name` | string | yes | Parent resource group (forces replacement) |
| `virtual_network_name` | string | yes | Parent VNet (forces replacement) |
| `address_prefixes` | list(string) | yes | CIDR blocks (e.g. `["10.0.1.0/24"]`) |
| `id` | string | computed | Azure-style resource ID |

### azuresim_virtual_machine

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | yes | VM name (forces replacement) |
| `resource_group_name` | string | yes | Parent resource group (forces replacement) |
| `location` | string | yes | Azure region (forces replacement) |
| `vm_size` | string | yes | VM size (e.g. `Standard_DS1_v2`) |
| `admin_username` | string | no | Admin user |
| `admin_password` | string | no | Admin password (sensitive) |
| `network_interface_ids` | list(string) | no | NIC IDs to attach |
| `tags` | map(string) | no | Tags |
| `id` | string | computed | Azure-style resource ID |

**Blocks:**

- `os_disk` — `caching` (required), `storage_account_type` (required), `disk_size_gb` (optional)
- `source_image_reference` — `publisher`, `offer`, `sku`, `version` (all required)

### azuresim_storage_account

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | yes | Account name (forces replacement) |
| `resource_group_name` | string | yes | Parent resource group (forces replacement) |
| `location` | string | yes | Azure region (forces replacement) |
| `account_tier` | string | yes | `Standard` or `Premium` |
| `account_replication_type` | string | yes | `LRS`, `GRS`, `RAGRS`, `ZRS` |
| `account_kind` | string | no | `StorageV2`, `BlobStorage`, etc. |
| `tags` | map(string) | no | Tags |
| `id` | string | computed | Azure-style resource ID |
| `primary_blob_endpoint` | string | computed | Simulated blob URL |
| `primary_access_key` | string | computed | Simulated access key (sensitive) |

## Project Structure

```
terraform-provider-azuresim/
├── main.go                                        # Entry point
├── go.mod
├── internal/provider/
│   ├── provider.go                                # Provider schema & config
│   ├── resource_resource_group.go                 # azuresim_resource_group
│   ├── resource_virtual_network.go                # azuresim_virtual_network
│   ├── resource_subnet.go                         # azuresim_subnet
│   ├── resource_virtual_machine.go                # azuresim_virtual_machine
│   └── resource_storage_account.go                # azuresim_storage_account
├── docs/                                          # Provider documentation
└── examples/
    └── main.tf                                    # Full working example
```

## Built With

- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) (v1.19)
- Go 1.26+
