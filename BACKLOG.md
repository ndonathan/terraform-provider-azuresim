# Resource Backlog

Proposed resources to add to `terraform-provider-azuresim`, ranked by deployment frequency and ordered so dependencies ship before dependents.

> **See also `PARITY.md`** — the canonical, auto-generated coverage report against AzureRM (61 / 1128 resources as of the latest snapshot). The "Recommended next-up resources" table there is the curated short-list this backlog should pull from going forward. This file remains the source of truth for tiering / dependency-ordering across batches; PARITY.md feeds the next 3-5 picks.

## Already implemented
- `azuresim_resource_group`
- `azuresim_virtual_network`
- `azuresim_subnet`
- `azuresim_virtual_machine` (generic; see Tier 2 split)
- `azuresim_storage_account`
- `azuresim_network_security_group`
- `azuresim_public_ip`
- `azuresim_network_interface`
- `azuresim_managed_disk`
- `azuresim_user_assigned_identity`
- `azuresim_key_vault`
- `azuresim_log_analytics_workspace`
- `azuresim_storage_container`
- `azuresim_subnet_network_security_group_association`
- `azuresim_network_security_rule`
- `azuresim_key_vault_secret`
- `azuresim_key_vault_key`
- `azuresim_role_assignment`
- `azuresim_application_insights`
- `azuresim_service_plan`
- `azuresim_windows_virtual_machine`
- `azuresim_route_table`
- `azuresim_subnet_route_table_association`
- `azuresim_private_dns_zone`
- `azuresim_private_dns_zone_virtual_network_link`
- `azuresim_virtual_network_peering`
- `azuresim_linux_web_app`, `azuresim_windows_web_app`
- `azuresim_linux_function_app`, `azuresim_windows_function_app`
- `azuresim_mssql_server`, `azuresim_mssql_database`
- `azuresim_container_registry`
- `azuresim_kubernetes_cluster`
- `azuresim_cosmosdb_account`
- `azuresim_redis_cache`
- `azuresim_servicebus_namespace`
- `azuresim_eventhub_namespace`
- `azuresim_lb`
- `azuresim_recovery_services_vault`
- `azuresim_servicebus_queue`, `azuresim_servicebus_topic`, `azuresim_servicebus_subscription`
- `azuresim_eventhub`
- `azuresim_storage_blob`
- `azuresim_nat_gateway`
- `azuresim_application_gateway`
- `azuresim_firewall`
- `azuresim_postgresql_flexible_server`, `azuresim_mysql_flexible_server`
- `azuresim_container_app_environment`, `azuresim_container_app`
- `azuresim_linux_virtual_machine_scale_set`, `azuresim_windows_virtual_machine_scale_set`
- `azuresim_private_endpoint`
- `azuresim_monitor_diagnostic_setting`
- `azuresim_monitor_action_group`, `azuresim_monitor_metric_alert`
- `azuresim_bastion_host`
- `azuresim_api_management`
- `azuresim_cdn_frontdoor_profile`

## Tier 1 — Foundational (no missing deps; pre-reqs for almost everything)

| # | Resource | AzureRM Equivalent | Dependencies |
|---|----------|--------------------|--------------|
| 1 | ~~Network Security Group~~ ✅ | `azurerm_network_security_group` | Resource Group |
| 2 | ~~Public IP~~ ✅ | `azurerm_public_ip` | Resource Group |
| 3 | ~~Network Interface (NIC)~~ ✅ | `azurerm_network_interface` | Subnet, optional Public IP / NSG |
| 4 | ~~Managed Disk~~ ✅ | `azurerm_managed_disk` | Resource Group |
| 5 | ~~User-Assigned Managed Identity~~ ✅ | `azurerm_user_assigned_identity` | Resource Group |
| 6 | ~~Key Vault~~ ✅ | `azurerm_key_vault` | Resource Group |
| 7 | ~~Log Analytics Workspace~~ ✅ | `azurerm_log_analytics_workspace` | Resource Group |
| 8 | ~~Storage Container~~ ✅ | `azurerm_storage_container` | Storage Account |

## Tier 2 — Common, depend on Tier 1

| # | Resource | AzureRM Equivalent | Dependencies |
|---|----------|--------------------|--------------|
| 9 | ~~Subnet–NSG Association~~ ✅ | `azurerm_subnet_network_security_group_association` | Subnet, NSG |
| 10 | ~~NSG Rule (standalone)~~ ✅ | `azurerm_network_security_rule` | NSG |
| 11 | ~~Key Vault Secret~~ ✅ | `azurerm_key_vault_secret` | Key Vault |
| 12 | ~~Key Vault Key~~ ✅ | `azurerm_key_vault_key` | Key Vault |
| 13 | ~~Role Assignment~~ ✅ | `azurerm_role_assignment` | Any scope, principal (UAI) |
| 14 | ~~Application Insights~~ ✅ | `azurerm_application_insights` | Resource Group (+ LAW for workspace-based) |
| 15 | ~~App Service Plan~~ ✅ | `azurerm_service_plan` | Resource Group |
| 16 | ~~Windows Virtual Machine~~ ✅ | `azurerm_windows_virtual_machine` | NIC, Resource Group |
| 17 | ~~Route Table + Subnet Association~~ ✅ | `azurerm_route_table`, `azurerm_subnet_route_table_association` | Resource Group, Subnet |
| 18 | ~~Private DNS Zone + VNet Link~~ ✅ | `azurerm_private_dns_zone`, `azurerm_private_dns_zone_virtual_network_link` | Resource Group, VNet |
| 19 | ~~VNet Peering~~ ✅ | `azurerm_virtual_network_peering` | Two VNets |

## Tier 3 — Layered platform services

| # | Resource | AzureRM Equivalent | Dependencies |
|---|----------|--------------------|--------------|
| 20 | ~~Linux / Windows Web App~~ ✅ | `azurerm_linux_web_app`, `azurerm_windows_web_app` | App Service Plan |
| 21 | ~~Function App~~ ✅ | `azurerm_linux_function_app`, `azurerm_windows_function_app` | Plan, Storage Account, optional App Insights |
| 22 | ~~SQL Server~~ ✅ | `azurerm_mssql_server` | Resource Group |
| 23 | ~~SQL Database~~ ✅ | `azurerm_mssql_database` | SQL Server |
| 24 | ~~Container Registry (ACR)~~ ✅ | `azurerm_container_registry` | Resource Group |
| 25 | ~~AKS Cluster~~ ✅ | `azurerm_kubernetes_cluster` | Resource Group, Subnet, optional LAW + UAI |
| 26 | ~~Cosmos DB Account~~ ✅ | `azurerm_cosmosdb_account` | Resource Group |
| 27 | ~~Redis Cache~~ ✅ | `azurerm_redis_cache` | Resource Group |
| 28 | ~~Service Bus Namespace~~ ✅ | `azurerm_servicebus_namespace` | Resource Group |
| 29 | ~~Event Hub Namespace~~ ✅ | `azurerm_eventhub_namespace` | Resource Group |
| 30 | ~~Load Balancer~~ ✅ | `azurerm_lb` | Resource Group (+ Public IP or Subnet) |
| 31 | ~~Recovery Services Vault~~ ✅ | `azurerm_recovery_services_vault` | Resource Group |

## Tier 4 — Children of Tier 3 or specialized

| # | Resource | AzureRM Equivalent | Dependencies |
|---|----------|--------------------|--------------|
| 32 | ~~Service Bus Queue / Topic / Subscription~~ ✅ | `azurerm_servicebus_queue`, `_topic`, `_subscription` | SB Namespace |
| 33 | ~~Event Hub~~ ✅ | `azurerm_eventhub` | EH Namespace |
| 34 | ~~Storage Blob~~ ✅ | `azurerm_storage_blob` | Storage Container |
| 35 | ~~NAT Gateway~~ ✅ | `azurerm_nat_gateway` | Resource Group, Public IP, Subnet |
| 36 | ~~Application Gateway~~ ✅ | `azurerm_application_gateway` | Subnet, Public IP |
| 37 | ~~Azure Firewall~~ ✅ | `azurerm_firewall` | Subnet, Public IP |
| 38 | ~~PostgreSQL / MySQL Flexible Server~~ ✅ | `azurerm_postgresql_flexible_server`, `azurerm_mysql_flexible_server` | Resource Group (optional delegated Subnet) |
| 39 | ~~Container App Environment + Container App~~ ✅ | `azurerm_container_app_environment`, `azurerm_container_app` | Resource Group, optional LAW |
| 40 | ~~VM Scale Set (Linux / Windows)~~ ✅ | `azurerm_linux_virtual_machine_scale_set`, `azurerm_windows_virtual_machine_scale_set` | Subnet |
| 41 | ~~Private Endpoint + DNS Zone Group~~ ✅ | `azurerm_private_endpoint` | Subnet, target resource, Private DNS Zone |
| 42 | ~~Diagnostic Setting~~ ✅ | `azurerm_monitor_diagnostic_setting` | Any target resource + LAW / SA sink |
| 43 | ~~Monitor Action Group / Metric Alert~~ ✅ | `azurerm_monitor_action_group`, `azurerm_monitor_metric_alert` | Resource Group (+ targets) |
| 44 | ~~Bastion Host~~ ✅ | `azurerm_bastion_host` | Subnet (`AzureBastionSubnet`), Public IP |
| 45 | ~~API Management~~ ✅ | `azurerm_api_management` | Resource Group |
| 46 | ~~Front Door / CDN Profile~~ ✅ | `azurerm_cdn_frontdoor_profile`, `azurerm_cdn_profile` | Resource Group |

## Notes

- `azuresim_virtual_machine` is currently a single generic resource. Tier 2 splits it into Linux + Windows variants to mirror AzureRM (`azurerm_linux_virtual_machine` / `azurerm_windows_virtual_machine`). Decide whether to keep the generic resource as an alias or deprecate it.
- The NIC is implicit today — `virtual_machine.network_interface_ids` accepts strings but no NIC resource exists. Promoting it to Tier 1 unlocks realistic VM, LB, and App Gateway configurations.
- Tier 4 child resources (queue/topic/blob/etc.) can be deferred or batched with their parent so the parent ships usefully on its own.
