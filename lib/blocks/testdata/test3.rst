# Test 3

Test fenced code block with `terraform`

.. code:: terraform
  resource "aws_s3_bucket" "terraform" {
    bucket = "tf-test-bucket-terraform"
  }

Stuff that is not to be formatted.

A bit more complicated example

.. code:: terraform
  resource "azurerm_resource_group" "example" {
    name     = "testaccbatch"
    location = "West Europe"
  }

  resource "azurerm_storage_account" "example" {
    name                     = "testaccsa"
    resource_group_name      = azurerm_resource_group.example.name
    location                 = azurerm_resource_group.example.location
    account_tier             = "Standard"
    account_replication_type = "LRS"
  }

  resource "azurerm_batch_account" "example" {
    name                 = "testaccbatch"
    resource_group_name  = azurerm_resource_group.example.name
    location             = azurerm_resource_group.example.location
    pool_allocation_mode = "BatchService"
    storage_account_id   = azurerm_storage_account.example.id

    tags = {
      env = "test"
    }
  }

  resource "azurerm_batch_certificate" "example" {
    resource_group_name  = azurerm_resource_group.example.name
    account_name         = azurerm_batch_account.example.name
    certificate          = filebase64("certificate.cer")
    format               = "Cer"
    thumbprint           = "312d31a79fa0cef49c00f769afc2b73e9f4edf34"
    thumbprint_algorithm = "SHA1"
  }

  resource "azurerm_batch_pool" "example" {
    name                = "testaccpool"
    resource_group_name = azurerm_resource_group.example.name
    account_name        = azurerm_batch_account.example.name
    display_name        = "Test Acc Pool Auto"
    vm_size             = "Standard_A1"
    node_agent_sku_id   = "batch.node.ubuntu 20.04"

    auto_scale {
      evaluation_interval = "PT15M"

      formula = <<EOF
        startingNumberOfVMs = 1;
        maxNumberofVMs = 25;
        pendingTaskSamplePercent = $PendingTasks.GetSamplePercent(180 * TimeInterval_Second);
        pendingTaskSamples = pendingTaskSamplePercent < 70 ? startingNumberOfVMs : avg($PendingTasks.GetSample(180 *   TimeInterval_Second));
        $TargetDedicatedNodes=min(maxNumberofVMs, pendingTaskSamples);
  EOF

    }

    storage_image_reference {
      publisher = "microsoft-azure-batch"
      offer     = "ubuntu-server-container"
      sku       = "20-04-lts"
      version   = "latest"
    }

    container_configuration {
      type = "DockerCompatible"
      container_registries {
        registry_server = "docker.io"
        user_name       = "login"
        password        = "apassword"
      }
    }

    start_task {
      command_line       = "echo 'Hello World from $env'"
      task_retry_maximum = 1
      wait_for_success   = true

      common_environment_properties = {
        env = "TEST"
      }

      user_identity {
        auto_user {
          elevation_level = "NonAdmin"
          scope           = "Task"
        }
      }
    }

    certificate {
      id             = azurerm_batch_certificate.example.id
      store_location = "CurrentUser"
      visibility     = ["StartTask"]
    }
  }

Stuff that is not to be parsed as block
