package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceBusNamespaceAndChildren(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_servicebus_namespace" "ns" {
  name                = "sb-ns-001"
  resource_group_name = "rg"
  location            = "eastus"
  sku                 = "Standard"
}

resource "azuresim_servicebus_queue" "queue" {
  name         = "orders"
  namespace_id = azuresim_servicebus_namespace.ns.id
  max_delivery_count = 10
}

resource "azuresim_servicebus_topic" "topic" {
  name         = "events"
  namespace_id = azuresim_servicebus_namespace.ns.id
}

resource "azuresim_servicebus_subscription" "sub" {
  name               = "sub1"
  topic_id           = azuresim_servicebus_topic.topic.id
  max_delivery_count = 10
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azuresim_servicebus_namespace.ns", "endpoint",
						"sb://sb-ns-001.servicebus.windows.net/"),
					resource.TestMatchResourceAttr("azuresim_servicebus_namespace.ns", "default_primary_connection_string",
						regexp.MustCompile(`Endpoint=sb://sb-ns-001\.servicebus\.windows\.net/;SharedAccessKeyName=RootManageSharedAccessKey`)),
					resource.TestCheckResourceAttr("azuresim_servicebus_queue.queue", "id",
						arm("/resourceGroups/rg/providers/Microsoft.ServiceBus/namespaces/sb-ns-001/queues/orders")),
					resource.TestCheckResourceAttr("azuresim_servicebus_topic.topic", "id",
						arm("/resourceGroups/rg/providers/Microsoft.ServiceBus/namespaces/sb-ns-001/topics/events")),
					resource.TestCheckResourceAttr("azuresim_servicebus_subscription.sub", "id",
						arm("/resourceGroups/rg/providers/Microsoft.ServiceBus/namespaces/sb-ns-001/topics/events/subscriptions/sub1")),
				),
			},
		},
	})
}

func TestAccEventHubNamespaceAndHub(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "azuresim_eventhub_namespace" "ns" {
  name                = "eh-ns-001"
  resource_group_name = "rg"
  location            = "eastus"
  sku                 = "Standard"
  capacity            = 1
}

resource "azuresim_eventhub" "hub" {
  name              = "telemetry"
  namespace_id      = azuresim_eventhub_namespace.ns.id
  partition_count   = 4
  message_retention = 1
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuresim_eventhub_namespace.ns", "default_primary_key"),
					resource.TestCheckResourceAttr("azuresim_eventhub.hub", "id",
						arm("/resourceGroups/rg/providers/Microsoft.EventHub/namespaces/eh-ns-001/eventhubs/telemetry")),
				),
			},
		},
	})
}
