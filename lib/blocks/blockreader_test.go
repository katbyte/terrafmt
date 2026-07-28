package blocks

import (
	"bytes"
	"testing"

	"github.com/katbyte/terrafmt/lib/common"
	"github.com/kylelemons/godebug/diff"
	"github.com/spf13/afero"
)

func TestBlockDetection(t *testing.T) {
	t.Parallel()
	type block struct {
		leadingPadding  string
		trailingPadding string
		text            string
	}

	testcases := []struct {
		sourcefile     string
		expectedBlocks []block
	}{
		{
			sourcefile: "testdata/test1.go",
			expectedBlocks: []block{
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "azurerm_storage_container" "simple" {
  name = "tf-test-container-simple"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "azurerm_storage_container" "simple2" {
  name = "tf-test-container-simple2"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "azurerm_storage_container" "with-parameters" {
  name = "tf-test-container-with-parameters-%d"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "azurerm_storage_container" "with-parameters-and-append" {
  name = "tf-test-container-parameters-and-append-%d"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "azurerm_storage_container" "const" {
  name = "tf-test-container-const"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "azurerm_storage_container" "composed" {
  name = "tf-test-container-composed-%d"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `data "azurerm_storage_container" "simple" {
  name = "tf-test-container-simple"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `    resource "azurerm_storage_container" "leading-space" {
  name = "tf-test-container-leading-space-%d"
}
`,
				},
				{
					leadingPadding:  "\n    \n",
					trailingPadding: "\n",
					text: `    
    resource "azurerm_storage_container" "leading-space-and-line" {
  name = "tf-test-container-leading-space-and-line-%d"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "azurerm_storage_container" "%s" {
  name = "tf-test-container-with-quotedname"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "azurerm_storage_container" "UpperCase" {
  name = "tf-test-container-with-uppercase"
}
`,
				},
			},
		},
		{
			sourcefile: "testdata/test2.markdown",
			expectedBlocks: []block{
				{text: `resource "azurerm_storage_container" "hcl" {
  name = "tf-test-container-hcl"
}
`},
				{text: `resource "azurerm_storage_container" "tf" {
  name = "tf-test-container-tf"
}
`},
				{
					text: `    resource "azurerm_storage_container" "leading-space" {
  name = "tf-test-container-leading-space"
}
`,
				},
				{
					text: `    
    resource "azurerm_storage_container" "leading-space-and-line" {
  name = "tf-test-container-leading-space-and-line"
}
`,
				},
				{
					text: `resource "azurerm_storage_container" "UpperCase" {
  name = "tf-test-container-with-uppercase"
}
`,
				},
				{
					text: `    resource "aws_s3_bucket" "indented" {
      bucket = "tf-test-bucket-indented"
    }
`,
				},
			},
		},
		{
			sourcefile: "testdata/test3.rst",
			expectedBlocks: []block{
				{
					text: `  resource "azurerm_storage_container" "terraform" {
    name = "tf-test-container-terraform"
  }

`,
				},
				{
					text: `  resource "azurerm_resource_group" "example" {
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

`,
				},
			},
		},
	}

	fs := afero.NewReadOnlyFs(afero.NewOsFs())

	errB := bytes.NewBufferString("")
	log := common.CreateLogger(errB)

	for _, testcase := range testcases {
		var actualBlocks []block
		// also test leading and trailing padding
		br := Reader{
			Log:      log,
			ReadOnly: true,
			LineRead: ReaderIgnore,
			BlockRead: func(br *Reader, _ int, b string, _ bool) error {
				actualBlocks = append(actualBlocks, block{
					leadingPadding:  br.CurrentNodeLeadingPadding,
					text:            b,
					trailingPadding: br.CurrentNodeTrailingPadding,
				})

				return nil
			},
		}
		err := br.DoTheThing(fs, testcase.sourcefile, nil, nil)
		if err != nil {
			t.Errorf("Case %q: Got an error when none was expected: %v", testcase.sourcefile, err)
			continue
		}

		if len(testcase.expectedBlocks) != len(actualBlocks) {
			t.Errorf("Case %q: expected %d blocks, got %d", testcase.sourcefile, len(testcase.expectedBlocks), len(actualBlocks))
			continue
		}

		for i, actual := range actualBlocks {
			expected := testcase.expectedBlocks[i]
			if actual.text != expected.text {
				t.Errorf("Case %q, block %d text: ('-' actual, '+' expected)\n%s", testcase.sourcefile, i+1, diff.Diff(actual.text, expected.text))
			}
			if actual.leadingPadding != expected.leadingPadding {
				t.Errorf("Case %q, block %d leading padding: expected %q, got %q", testcase.sourcefile, i+1, expected.leadingPadding, actual.leadingPadding)
			}
			if actual.trailingPadding != expected.trailingPadding {
				t.Errorf("Case %q, block %d trailing padding: expected %q, got %q", testcase.sourcefile, i+1, expected.trailingPadding, actual.trailingPadding)
			}
		}

		actualErr := errB.String()
		if actualErr != "" {
			t.Errorf("Case %q: Got error output:\n%s", testcase.sourcefile, actualErr)
		}
	}
}

func TestLooksLikeTerraform(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		text     string
		expected bool
	}{
		{
			text: `
resource "azurerm_storage_container" "simple-resource" {
  name = "tf-test-container-simple"
}`,
			expected: true,
		},
		{
			text: `
data "azurerm_storage_container" "simple-data" {
  name = "tf-test-container-simple"
}`,
			expected: true,
		},
		{
			text: `
list "azurerm_resource_group" "example" {
  provider = azurerm

  config {
    filter = "tagName eq 'query' and tagValue eq 'example'"
  }
}
`,
			expected: true,
		},
		{
			text: `
ephemeral "azurerm_key_vault_secret" "example" {
  secret_id = data.azurerm_key_vault_secret.example.id
}`,
			expected: true,
		},
		{
			text: `
action "azurerm_virtual_machine_run_command" "example" {
  config {
    virtual_machine_id = azurerm_virtual_machine.example.id
  }
}`,
			expected: true,
		},
		{
			text: `
variable "name" {
  type = string
}`,
			expected: true,
		},
		{
			text: `
output "arn" {
  value = azurerm_storage_container.simple-resource.id
}`,
			expected: true,
		},
		{
			text: `
resource "azurerm_storage_container" "%s" {
  name = "tf-test-container-simple"
}`,
			expected: true,
		},
		// 		{
		// 			text: `
		// resource "azurerm_storage_container" "%[1]s" {
		//   name = "tf-test-container-simple"
		// }`,
		// 			expected: true,
		// 		},
		// 		{
		// 			text: `
		// resource "azurerm_storage_container" %q {
		//   name = "tf-test-container-simple"
		// }`,
		// 			expected: true,
		// 		},
		// 		{
		// 			text: `
		// resource "azurerm_storage_container" %[1]q {
		//   name = "tf-test-container-simple"
		// }`,
		// 			expected: true,
		// 		},
		{
			text:     "%d: bad create: \n%#v\n%#v",
			expected: false,
		},
		{
			text: `<DescribeAccountAttributesResponse xmlns="http://ec2.amazonaws.com/doc/2016-11-15/">
  <requestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</requestId>
  <accountAttributeSet>
	  <item>
	    <attributeName>supported-platforms</attributeName>
	    <attributeValueSet>
	      <item>
	        <attributeValue>VPC</attributeValue>
	      </item>
	      <item>
	        <attributeValue>EC2</attributeValue>
	      </item>
	    </attributeValueSet>
	  </item>
  </accountAttributeSet>
</DescribeAccountAttributesResponse>`,
			expected: false,
		},
	}

	for _, testcase := range testcases {
		actual := looksLikeTerraform(testcase.text)

		if testcase.expected && !actual {
			t.Errorf("Expected match, but not identified as Terraform:\n%s", testcase.text)
		} else if !testcase.expected && actual {
			t.Errorf("Expected no match, but was identified as Terraform:\n%s", testcase.text)
		}
	}
}
