package upgrade012

import (
	"strings"
	"testing"

	"github.com/katbyte/terrafmt/lib/common"
)

func TestBlock(t *testing.T) {
	tests := []struct {
		name     string
		block    string
		expected string
		error    bool
	}{
		{
			name:     "noblocks",
			block:    "// nothing here",
			expected: "// nothing here",
		},
		{
			name: "oneline",
			block: `data "google_compute_lb_ip_ranges" "some" {}
`,
			expected: `data "google_compute_lb_ip_ranges" "some" {
}
`,
		},
		{
			name: "basic",
			block: `resource "google_dns_managed_zone" "foo" {
	name		= "qa-zone-%s"
	dns_name	= "qa.tf-test.club."
	description	= "QA DNS zone"
}

data "google_dns_managed_zone" "qa" {
	name	= "${google_dns_managed_zone.foo.name}"
}
`,
			expected: `resource "google_dns_managed_zone" "foo" {
  name        = "qa-zone-%s"
  dns_name    = "qa.tf-test.club."
  description = "QA DNS zone"
}

data "google_dns_managed_zone" "qa" {
  name = google_dns_managed_zone.foo.name
}
`,
		},
		{
			name: "complicated",
			block: `resource "google_compute_network" "container_network" {
	name = "container-net-%s"
	auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
	name					 = "${google_compute_network.container_network.name}"
	network					 = "${google_compute_network.container_network.name}"
	ip_cidr_range			 = "10.0.36.0/24"
	region				 = "us-central1"
	private_ip_google_access = true

	secondary_ip_range {
		range_name	  = "pod"
		ip_cidr_range = "10.0.0.0/19"
	}

	secondary_ip_range {
		range_name	  = "svc"
		ip_cidr_range = "10.0.32.0/22"
	}
}

resource "google_container_cluster" "with_private_cluster" {
	name     = "cluster-test-%s"
	location = "us-central1-a"
	initial_node_count = 1

	network = "${google_compute_network.container_network.name}"
	subnetwork = "${google_compute_subnetwork.container_subnetwork.name}"

	private_cluster_config {
		enable_private_endpoint = false
		enable_private_nodes = true
		master_ipv4_cidr_block = "10.42.0.0/28"
	}

	ip_allocation_policy {
		cluster_secondary_range_name  = "${google_compute_subnetwork.container_subnetwork.secondary_ip_range.0.range_name}"
		services_secondary_range_name = "${google_compute_subnetwork.container_subnetwork.secondary_ip_range.1.range_name}"
	}
}
`,
			expected: `resource "google_compute_network" "container_network" {
  name                    = "container-net-%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "container_subnetwork" {
  name                     = google_compute_network.container_network.name
  network                  = google_compute_network.container_network.name
  ip_cidr_range            = "10.0.36.0/24"
  region                   = "us-central1"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "pod"
    ip_cidr_range = "10.0.0.0/19"
  }

  secondary_ip_range {
    range_name    = "svc"
    ip_cidr_range = "10.0.32.0/22"
  }
}

resource "google_container_cluster" "with_private_cluster" {
  name               = "cluster-test-%s"
  location           = "us-central1-a"
  initial_node_count = 1

  network    = google_compute_network.container_network.name
  subnetwork = google_compute_subnetwork.container_subnetwork.name

  private_cluster_config {
    enable_private_endpoint = false
    enable_private_nodes    = true
    master_ipv4_cidr_block  = "10.42.0.0/28"
  }

  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.container_subnetwork.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.container_subnetwork.secondary_ip_range[1].range_name
  }
}
`,
		},
		{
			name: "invalid",
			block: `
Hi there i am going to fail... =C
`,
			expected: ``,
			error:    true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var errB strings.Builder
			log := common.CreateLogger(&errB)
			result, err := Block(log, test.block)
			if err != nil && !test.error {
				t.Fatalf("Got an error when none was expected: %v", err)
			}
			if err == nil && test.error {
				t.Errorf("Expected an error and none was generated")
			}
			if result != test.expected {
				t.Errorf("Got: \n%#v\nexpected:\n%#v\n", result, test.expected)
			}
		})
	}
}
