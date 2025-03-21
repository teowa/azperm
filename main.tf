terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
  skip_provider_registration = false
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

data "azapi_client_config" "current" {
}

data "azapi_resource" "subscription" {
  type                   = "Microsoft.Resources/subscriptions@2021-01-01"
  resource_id            = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = ["*"]
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      sku = {
        family = "A"
        name   = "standard"
      }
      accessPolicies   = []
      enableSoftDelete = true
      tenantId         = data.azapi_client_config.current.tenant_id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.accessPolicies]
  }
}

data "azapi_resource_id" "secret" {
  type      = "Microsoft.KeyVault/vaults/secrets@2023-02-01"
  parent_id = azapi_resource.vault.id
  name      = var.resource_name
}

resource "azapi_resource_action" "put_secret" {
  type        = "Microsoft.KeyVault/vaults/secrets@2023-02-01"
  resource_id = data.azapi_resource_id.secret.id
  method      = "PUT"
  body = {
    properties = {
      value = "szechuan"
    }
  }
  response_export_values = ["*"]
}
