// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redshift_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfredshift "github.com/hashicorp/terraform-provider-aws/internal/service/redshift"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccRedshiftSnapshotSchedule_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v redshift.SnapshotSchedule
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_redshift_snapshot_schedule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSnapshotScheduleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotScheduleConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSnapshotScheduleExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "definitions.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "definitions.*", "rate(12 hours)"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "identifier", rName),
					resource.TestCheckResourceAttr(resourceName, "identifier_prefix", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy",
				},
			},
		},
	})
}

func TestAccRedshiftSnapshotSchedule_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v redshift.SnapshotSchedule
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_redshift_snapshot_schedule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSnapshotScheduleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotScheduleConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotScheduleExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfredshift.ResourceSnapshotSchedule(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRedshiftSnapshotSchedule_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var v redshift.SnapshotSchedule
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_redshift_snapshot_schedule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSnapshotScheduleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotScheduleConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotScheduleExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy",
				},
			},
			{
				Config: testAccSnapshotScheduleConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotScheduleExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccSnapshotScheduleConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotScheduleExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccRedshiftSnapshotSchedule_identifierGenerated(t *testing.T) {
	ctx := acctest.Context(t)
	var v redshift.SnapshotSchedule
	resourceName := "aws_redshift_snapshot_schedule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSnapshotScheduleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotScheduleConfig_identifierGenerated(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotScheduleExists(ctx, resourceName, &v),
					acctest.CheckResourceAttrNameGenerated(resourceName, "identifier"),
					resource.TestCheckResourceAttr(resourceName, "identifier_prefix", id.UniqueIdPrefix),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy",
				},
			},
		},
	})
}

func TestAccRedshiftSnapshotSchedule_identifierPrefix(t *testing.T) {
	ctx := acctest.Context(t)
	var v redshift.SnapshotSchedule
	resourceName := "aws_redshift_snapshot_schedule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSnapshotScheduleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotScheduleConfig_identifierPrefix("tf-acc-test-prefix-"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotScheduleExists(ctx, resourceName, &v),
					acctest.CheckResourceAttrNameFromPrefix(resourceName, "identifier", "tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "identifier_prefix", "tf-acc-test-prefix-"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy",
				},
			},
		},
	})
}

func TestAccRedshiftSnapshotSchedule_update(t *testing.T) {
	ctx := acctest.Context(t)
	var v redshift.SnapshotSchedule
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_redshift_snapshot_schedule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSnapshotScheduleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotScheduleConfig_multipleDefinitions(rName, "cron(30 12 *)", "cron(15 6 *)"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotScheduleExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "definitions.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "definitions.*", "cron(30 12 *)"),
					resource.TestCheckTypeSetElemAttr(resourceName, "definitions.*", "cron(15 6 *)"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy",
				},
			},
			{
				Config: testAccSnapshotScheduleConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotScheduleExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "definitions.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "definitions.*", "rate(12 hours)"),
				),
			},
			{
				Config: testAccSnapshotScheduleConfig_multipleDefinitions(rName, "cron(30 8 *)", "cron(15 10 *)"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotScheduleExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "definitions.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "definitions.*", "cron(30 8 *)"),
					resource.TestCheckTypeSetElemAttr(resourceName, "definitions.*", "cron(15 10 *)"),
				),
			},
		},
	})
}

func TestAccRedshiftSnapshotSchedule_withDescription(t *testing.T) {
	ctx := acctest.Context(t)
	var v redshift.SnapshotSchedule
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_redshift_snapshot_schedule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSnapshotScheduleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotScheduleConfig_description(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotScheduleExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "description", "Test Schedule"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy",
				},
			},
		},
	})
}

func TestAccRedshiftSnapshotSchedule_withForceDestroy(t *testing.T) {
	ctx := acctest.Context(t)
	var snapshotSchedule redshift.SnapshotSchedule
	var cluster redshift.Cluster
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_redshift_snapshot_schedule.test"
	clusterResourceName := "aws_redshift_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSnapshotScheduleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotScheduleConfig_forceDestroy(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotScheduleExists(ctx, resourceName, &snapshotSchedule),
					testAccCheckClusterExists(ctx, clusterResourceName, &cluster),
					testAccCheckSnapshotScheduleCreateSnapshotScheduleAssociation(ctx, &cluster, &snapshotSchedule),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy",
				},
			},
		},
	})
}

func testAccCheckSnapshotScheduleDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).RedshiftConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_redshift_snapshot_schedule" {
				continue
			}

			_, err := tfredshift.FindSnapshotScheduleByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Redshift Snapshot Schedule %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckSnapshotScheduleExists(ctx context.Context, n string, v *redshift.SnapshotSchedule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Redshift Cluster Snapshot Schedule ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).RedshiftConn(ctx)

		output, err := tfredshift.FindSnapshotScheduleByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckSnapshotScheduleCreateSnapshotScheduleAssociation(ctx context.Context, cluster *redshift.Cluster, snapshotSchedule *redshift.SnapshotSchedule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).RedshiftConn(ctx)

		if _, err := conn.ModifyClusterSnapshotScheduleWithContext(ctx, &redshift.ModifyClusterSnapshotScheduleInput{
			ClusterIdentifier:    cluster.ClusterIdentifier,
			ScheduleIdentifier:   snapshotSchedule.ScheduleIdentifier,
			DisassociateSchedule: aws.Bool(false),
		}); err != nil {
			return fmt.Errorf("Error associate Redshift Cluster and Snapshot Schedule: %s", err)
		}

		id := fmt.Sprintf("%s/%s", aws.StringValue(cluster.ClusterIdentifier), aws.StringValue(snapshotSchedule.ScheduleIdentifier))
		if _, err := tfredshift.WaitScheduleAssociationActive(ctx, conn, id); err != nil {
			return err
		}

		return nil
	}
}

func testAccSnapshotScheduleConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_redshift_snapshot_schedule" "test" {
  identifier = %[1]q
  definitions = [
    "rate(12 hours)",
  ]
}
`, rName)
}

func testAccSnapshotScheduleConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_redshift_snapshot_schedule" "test" {
  identifier = %[1]q
  definitions = [
    "rate(12 hours)",
  ]

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccSnapshotScheduleConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_redshift_snapshot_schedule" "test" {
  identifier = %[1]q
  definitions = [
    "rate(12 hours)",
  ]

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccSnapshotScheduleConfig_identifierGenerated() string {
	return `
resource "aws_redshift_snapshot_schedule" "test" {
  definitions = [
    "rate(12 hours)",
  ]
}
`
}

func testAccSnapshotScheduleConfig_identifierPrefix(prefix string) string {
	return fmt.Sprintf(`
resource "aws_redshift_snapshot_schedule" "test" {
  identifier_prefix = %[1]q
  definitions = [
    "rate(12 hours)",
  ]
}
`, prefix)
}

func testAccSnapshotScheduleConfig_multipleDefinitions(rName, definition1, definition2 string) string {
	return fmt.Sprintf(`
resource "aws_redshift_snapshot_schedule" "test" {
  identifier = %[1]q
  definitions = [
    %[2]q,
    %[3]q,
  ]
}
`, rName, definition1, definition2)
}

func testAccSnapshotScheduleConfig_description(rName string) string {
	return fmt.Sprintf(`
resource "aws_redshift_snapshot_schedule" "test" {
  identifier  = %[1]q
  description = "Test Schedule"
  definitions = [
    "rate(12 hours)",
  ]
}
`, rName)
}

func testAccSnapshotScheduleConfig_forceDestroy(rName string) string {
	return acctest.ConfigCompose(testAccClusterConfig_basic(rName), fmt.Sprintf(`
resource "aws_redshift_snapshot_schedule" "test" {
  identifier = %[1]q
  definitions = [
    "rate(12 hours)",
  ]
  force_destroy = true
}
`, rName))
}
