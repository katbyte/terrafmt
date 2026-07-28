Example Document
================

An example resource:

.. code:: terraform
  
  resource "azurerm_resource_group" "example" {
    name     = "example"
    location = "West Europe"
  }
  
A correctly formatted one:

.. code:: terraform
  
  resource "azurerm_storage_account" "example" {
    name                     = "examplesa"
    resource_group_name      = azurerm_resource_group.example.name
    location                 = azurerm_resource_group.example.location
    account_tier             = "Standard"
    account_replication_type = "LRS"
  }
  
The end.
