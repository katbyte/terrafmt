package blocks

import (
	"bytes"
	"testing"

	"github.com/katbyte/terrafmt/lib/common"
	"github.com/kylelemons/godebug/diff"
	"github.com/spf13/afero"
)

func TestBlockDetection(t *testing.T) {
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
					text: `resource "aws_s3_bucket" "simple" {
  bucket = "tf-test-bucket-simple"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "aws_s3_bucket" "simple2" {
  bucket = "tf-test-bucket-simple2"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "aws_s3_bucket" "with-parameters" {
  bucket = "tf-test-bucket-with-parameters-%d"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "aws_s3_bucket" "with-parameters-and-append" {
  bucket = "tf-test-bucket-parameters-and-append-%d"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "aws_s3_bucket" "const" {
  bucket = "tf-test-bucket-const"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "aws_s3_bucket" "composed" {
  bucket = "tf-test-bucket-composed-%d"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `data "aws_s3_bucket" "simple" {
  bucket = "tf-test-bucket-simple"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `    resource "aws_s3_bucket" "leading-space" {
  bucket = "tf-test-bucket-leading-space-%d"
}
`,
				},
				{
					leadingPadding:  "\n    \n",
					trailingPadding: "\n",
					text: `    
    resource "aws_s3_bucket" "leading-space-and-line" {
  bucket = "tf-test-bucket-leading-space-and-line-%d"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "aws_s3_bucket" "%s" {
  bucket = "tf-test-bucket-with-quotedname"
}
`,
				},
				{
					leadingPadding:  "\n",
					trailingPadding: "\n",
					text: `resource "aws_s3_bucket" "UpperCase" {
  bucket = "tf-test-bucket-with-uppercase"
}
`,
				},
			},
		},
		{
			sourcefile: "testdata/test2.markdown",
			expectedBlocks: []block{
				{text: `resource "aws_s3_bucket" "hcl" {
  bucket = "tf-test-bucket-hcl"
}
`},
				{text: `resource "aws_s3_bucket" "tf" {
  bucket = "tf-test-bucket-tf"
}
`},
				{
					text: `    resource "aws_s3_bucket" "leading-space" {
  bucket = "tf-test-bucket-leading-space"
}
`,
				},
				{
					text: `    
    resource "aws_s3_bucket" "leading-space-and-line" {
  bucket = "tf-test-bucket-leading-space-and-line"
}
`,
				},
				{
					text: `resource "aws_s3_bucket" "UpperCase" {
  bucket = "tf-test-bucket-with-uppercase"
}
`,
				},
			},
		},
		{
			sourcefile: "testdata/test3.rst",
			expectedBlocks: []block{
				{
					text: `  resource "aws_s3_bucket" "terraform" {
    bucket = "tf-test-bucket-terraform"
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
			BlockRead: func(br *Reader, i int, b string, preserveIndent bool) error {
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
	testcases := []struct {
		text     string
		expected bool
	}{
		{
			text: `
resource "aws_s3_bucket" "simple-resource" {
  bucket = "tf-test-bucket-simple"
}`,
			expected: true,
		},
		{
			text: `
data "aws_s3_bucket" "simple-data" {
  bucket = "tf-test-bucket-simple"
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
  value = aws_s3_bucket.simple-resource.arn
}`,
			expected: true,
		},
		{
			text: `
resource "aws_s3_bucket" "%s" {
  bucket = "tf-test-bucket-simple"
}`,
			expected: true,
		},
		// 		{
		// 			text: `
		// resource "aws_s3_bucket" "%[1]s" {
		//   bucket = "tf-test-bucket-simple"
		// }`,
		// 			expected: true,
		// 		},
		// 		{
		// 			text: `
		// resource "aws_s3_bucket" %q {
		//   bucket = "tf-test-bucket-simple"
		// }`,
		// 			expected: true,
		// 		},
		// 		{
		// 			text: `
		// resource "aws_s3_bucket" %[1]q {
		//   bucket = "tf-test-bucket-simple"
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
