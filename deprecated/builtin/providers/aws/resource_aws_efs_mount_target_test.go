package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/efs"

	"github.com/xanzy/terraform-api/helper/resource"
	"github.com/xanzy/terraform-api/terraform"
)

func TestAccAWSEFSMountTarget(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEfsMountTargetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSEFSMountTargetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEfsMountTarget(
						"aws_efs_mount_target.alpha",
					),
				),
			},
			resource.TestStep{
				Config: testAccAWSEFSMountTargetConfigModified,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEfsMountTarget(
						"aws_efs_mount_target.alpha",
					),
					testAccCheckEfsMountTarget(
						"aws_efs_mount_target.beta",
					),
				),
			},
		},
	})
}

func testAccCheckEfsMountTargetDestroy(s *terraform.State) error {
	if len(s.RootModule().Resources) > 0 {
		return fmt.Errorf("Expected all resources to be gone, but found: %#v", s.RootModule().Resources)
	}

	return nil
}

func testAccCheckEfsMountTarget(resourceID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceID]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceID)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		fs, ok := s.RootModule().Resources[resourceID]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceID)
		}

		conn := testAccProvider.Meta().(*AWSClient).efsconn
		mt, err := conn.DescribeMountTargets(&efs.DescribeMountTargetsInput{
			MountTargetId: aws.String(fs.Primary.ID),
		})
		if err != nil {
			return err
		}

		if *mt.MountTargets[0].MountTargetId != fs.Primary.ID {
			return fmt.Errorf("Mount target ID mismatch: %q != %q",
				*mt.MountTargets[0].MountTargetId, fs.Primary.ID)
		}

		return nil
	}
}

const testAccAWSEFSMountTargetConfig = `
resource "aws_efs_file_system" "foo" {
	reference_name = "radeksimko"
}

resource "aws_efs_mount_target" "alpha" {
	file_system_id = "${aws_efs_file_system.foo.id}"
	subnet_id = "${aws_subnet.alpha.id}"
}

resource "aws_vpc" "foo" {
	cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "alpha" {
	vpc_id = "${aws_vpc.foo.id}"
	availability_zone = "us-west-2a"
	cidr_block = "10.0.1.0/24"
}
`

const testAccAWSEFSMountTargetConfigModified = `
resource "aws_efs_file_system" "foo" {
	reference_name = "radeksimko"
}

resource "aws_efs_mount_target" "alpha" {
	file_system_id = "${aws_efs_file_system.foo.id}"
	subnet_id = "${aws_subnet.alpha.id}"
}

resource "aws_efs_mount_target" "beta" {
	file_system_id = "${aws_efs_file_system.foo.id}"
	subnet_id = "${aws_subnet.beta.id}"
}

resource "aws_vpc" "foo" {
	cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "alpha" {
	vpc_id = "${aws_vpc.foo.id}"
	availability_zone = "us-west-2a"
	cidr_block = "10.0.1.0/24"
}

resource "aws_subnet" "beta" {
	vpc_id = "${aws_vpc.foo.id}"
	availability_zone = "us-west-2b"
	cidr_block = "10.0.2.0/24"
}
`
