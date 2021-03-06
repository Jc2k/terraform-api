package aws

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/xanzy/terraform-api/helper/resource"
	"github.com/xanzy/terraform-api/terraform"
)

func TestAccAWSRedshiftCluster_basic(t *testing.T) {
	var v redshift.Cluster

	ri := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	config := fmt.Sprintf(testAccAWSRedshiftClusterConfig_basic, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRedshiftClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSRedshiftClusterExists("aws_redshift_cluster.default", &v),
				),
			},
		},
	})
}

func testAccCheckAWSRedshiftClusterDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_redshift_cluster" {
			continue
		}

		// Try to find the Group
		conn := testAccProvider.Meta().(*AWSClient).redshiftconn
		var err error
		resp, err := conn.DescribeClusters(
			&redshift.DescribeClustersInput{
				ClusterIdentifier: aws.String(rs.Primary.ID),
			})

		if err == nil {
			if len(resp.Clusters) != 0 &&
				*resp.Clusters[0].ClusterIdentifier == rs.Primary.ID {
				return fmt.Errorf("Redshift Cluster %s still exists", rs.Primary.ID)
			}
		}

		// Return nil if the cluster is already destroyed
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ClusterNotFound" {
				return nil
			}
		}

		return err
	}

	return nil
}

func testAccCheckAWSRedshiftClusterExists(n string, v *redshift.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Redshift Cluster Instance ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).redshiftconn
		resp, err := conn.DescribeClusters(&redshift.DescribeClustersInput{
			ClusterIdentifier: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		for _, c := range resp.Clusters {
			if *c.ClusterIdentifier == rs.Primary.ID {
				*v = *c
				return nil
			}
		}

		return fmt.Errorf("Redshift Cluster (%s) not found", rs.Primary.ID)
	}
}

func TestResourceAWSRedshiftClusterIdentifierValidation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "tEsting",
			ErrCount: 1,
		},
		{
			Value:    "1testing",
			ErrCount: 1,
		},
		{
			Value:    "testing--123",
			ErrCount: 1,
		},
		{
			Value:    "testing!",
			ErrCount: 1,
		},
		{
			Value:    "testing-",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateRedshiftClusterIdentifier(tc.Value, "aws_redshift_cluster_identifier")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Redshift Cluster cluster_identifier to trigger a validation error")
		}
	}
}

func TestResourceAWSRedshiftClusterDbNameValidation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "tEsting",
			ErrCount: 1,
		},
		{
			Value:    "testing1",
			ErrCount: 1,
		},
		{
			Value:    "testing-",
			ErrCount: 1,
		},
		{
			Value:    "",
			ErrCount: 2,
		},
		{
			Value:    randomString(65),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateRedshiftClusterDbName(tc.Value, "aws_redshift_cluster_database_name")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Redshift Cluster database_name to trigger a validation error")
		}
	}
}

func TestResourceAWSRedshiftClusterFinalSnapshotIdentifierValidation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "testing--123",
			ErrCount: 1,
		},
		{
			Value:    "testing-",
			ErrCount: 1,
		},
		{
			Value:    "Testingq123!",
			ErrCount: 1,
		},
		{
			Value:    randomString(256),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateRedshiftClusterFinalSnapshotIdentifier(tc.Value, "aws_redshift_cluster_final_snapshot_identifier")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Redshift Cluster final_snapshot_identifier to trigger a validation error")
		}
	}
}

func TestResourceAWSRedshiftClusterMasterUsernameValidation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "1Testing",
			ErrCount: 1,
		},
		{
			Value:    "Testing!!",
			ErrCount: 1,
		},
		{
			Value:    randomString(129),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateRedshiftClusterMasterUsername(tc.Value, "aws_redshift_cluster_master_username")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Redshift Cluster master_username to trigger a validation error")
		}
	}
}

var testAccAWSRedshiftClusterConfig_basic = `
provider "aws" {
	region = "us-west-2"
}

resource "aws_redshift_cluster" "default" {
  cluster_identifier = "tf-redshift-cluster-%d"
  availability_zone = "us-west-2a"
  database_name = "mydb"
  master_username = "foo"
  master_password = "Mustbe8characters"
  node_type = "dc1.large"
  cluster_type = "single-node"
  automated_snapshot_retention_period = 7
  allow_version_upgrade = false
}`
