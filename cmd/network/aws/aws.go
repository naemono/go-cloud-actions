package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/naemono/go-cloud-actions/cmd/shared"
	shared_aws "github.com/naemono/go-cloud-actions/cmd/shared/aws"
	auth_aws "github.com/naemono/go-cloud-actions/pkg/auth/aws"
	"github.com/naemono/go-cloud-actions/pkg/logging"
	aws_network "github.com/naemono/go-cloud-actions/pkg/network/aws"
	"github.com/naemono/go-cloud-actions/pkg/validate"
)

var (
	// AWSCmd is the base aws network command
	AWSCmd = &cobra.Command{
		Use:              "aws",
		Short:            "Control networks in AWS's public clouds",
		Long:             `A cli to interact with networks in AWS's public cloud.`,
		PersistentPreRun: shared_aws.PersistentPreRun,
	}
	vpcCmd = &cobra.Command{
		Use:   "vpc",
		Short: "control VPCs in AWS's public clouds",
		Long:  `A cli to control VPCs in AWS's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			shared.RunParentsPersistentPreRun(cmd, args)
		},
	}
	regionCmd = &cobra.Command{
		Use:   "regions",
		Short: "control regions/azs in AWS's public clouds",
		Long:  `A cli to control Regions and AZs in AWS's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			shared.RunParentsPersistentPreRun(cmd, args)
		},
	}
	vpcCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create VPC in AWS's public clouds",
		Long:  `A cli to create VPC in AWS's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			shared.RunParentsPersistentPreRun(cmd, args)
			viper.BindPFlag("name", cmd.Flags().Lookup("name"))
			viper.BindPFlag("cidr", cmd.Flags().Lookup("cidr"))
			viper.BindPFlag("dry-run", cmd.Flags().Lookup("dry-run"))
			viper.BindPFlag("additional-tags", cmd.Flags().Lookup("additional-tags"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validate.NotEmpty(
				viper.GetViper(),
				[]string{"profile", "name", "region", "cidr"}); err != nil {
				return err
			}
			return createVPC()
		},
	}
	vpcDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete VPC in AWS's public clouds",
		Long:  `A cli to delete VPC in AWS's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			shared.RunParentsPersistentPreRun(cmd, args)
			viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validate.NotEmpty(viper.GetViper(), []string{"region", "profile", "id"}); err != nil {
				return err
			}
			return deleteVPC()
		},
	}
	vpcListCmd = &cobra.Command{
		Use:   "list",
		Short: "list VPCs in AWS's public clouds",
		Long:  `A cli to list VPCs in AWS's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			shared.RunParentsPersistentPreRun(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validate.NotEmpty(viper.GetViper(), []string{"region", "profile"}); err != nil {
				return err
			}
			return listVPCs()
		},
	}
	vpcCreateSubnetCmd = &cobra.Command{
		Use:   "create-subnet",
		Short: "create subnet in VPC in AWS's public clouds",
		Long:  `A cli to create subnets in vpcs in AWS's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			shared.RunParentsPersistentPreRun(cmd, args)
			viper.BindPFlag("id", cmd.Flags().Lookup("id"))
			viper.BindPFlag("cidr", cmd.Flags().Lookup("cidr"))
			viper.BindPFlag("az", cmd.Flags().Lookup("az"))
			viper.BindPFlag("additional-tags", cmd.Flags().Lookup("additional-tags"))
			viper.BindPFlag("dry-run", cmd.Flags().Lookup("dry-run"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validate.NotEmpty(
				viper.GetViper(),
				[]string{"region", "profile", "id", "cidr", "az"}); err != nil {
				return err
			}
			return createSubnetInVPC()
		},
	}
	vpcListSubnetsCmd = &cobra.Command{
		Use:   "list-subnets",
		Short: "list subnets within a vpc in AWS's public clouds",
		Long:  `A cli to list subnets in a given vpc in AWS's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			shared.RunParentsPersistentPreRun(cmd, args)
			viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validate.NotEmpty(viper.GetViper(), []string{"region", "profile", "id"}); err != nil {
				return err
			}
			return listSubnetsInVPC()
		},
	}
	azListCmd = &cobra.Command{
		Use:   "az-list",
		Short: "list AZs in AWS's public clouds",
		Long:  `A cli to list Availability Zones in AWS's public cloud.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			shared.RunParentsPersistentPreRun(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validate.NotEmpty(viper.GetViper(), []string{"region", "profile"}); err != nil {
				return err
			}
			return listAZs()
		},
	}
)

func init() {
	shared_aws.AddAuthFlagsToCommand(AWSCmd)

	vpcCreateCmd.Flags().StringP("name", "n", "", "name of vpc")
	vpcCreateCmd.Flags().StringP("cidr", "c", "10.4.240.0/21", "virtual network cidr to use")
	vpcCreateCmd.Flags().BoolP("dry-run", "d", false, "dry-run the vpc creation")
	vpcCreateCmd.Flags().StringSliceP("additional-tags", "t", []string{"environment", "development"}, "tags to apply to vpc")

	vpcDeleteCmd.Flags().StringP("id", "i", "", "vpc id to delete")

	vpcCreateSubnetCmd.Flags().StringP("id", "i", "", "vpc id to create subnet within")
	vpcCreateSubnetCmd.Flags().StringP("cidr", "c", "10.4.240.0/21", "virtual network cidr to use")
	vpcCreateSubnetCmd.Flags().StringP("az", "a", "us-east-1a", "availability zone to create cidr within")
	vpcCreateSubnetCmd.Flags().StringSliceP("additional-tags", "t", []string{"environment", "development"}, "tags to apply to vpc subnet")
	vpcCreateSubnetCmd.Flags().BoolP("dry-run", "d", false, "dry-run the vpc subnet creation")

	vpcListSubnetsCmd.Flags().StringP("id", "i", "", "vpc id to list subnets within")

	AWSCmd.AddCommand(vpcCmd)
	AWSCmd.AddCommand(regionCmd)
	vpcCmd.AddCommand(vpcCreateCmd)
	vpcCmd.AddCommand(vpcDeleteCmd)
	vpcCmd.AddCommand(vpcListCmd)
	vpcCmd.AddCommand(vpcCreateSubnetCmd)
	vpcCmd.AddCommand(vpcListSubnetsCmd)
	regionCmd.AddCommand(azListCmd)
}

func getLoggerAndNetworkClient() (*logrus.Entry, *aws_network.Client, error) {
	logger := logging.GetLogger(viper.GetString("loglevel"))
	client, err := aws_network.New(aws_network.Config{
		AuthConfig: auth_aws.AuthConfig{
			Profile: viper.GetString("profile"),
			Region:  viper.GetString("region"),
		},
	})
	return logger, client, err
}

func createVPC() error {
	logger, client, err := getLoggerAndNetworkClient()
	if err != nil {
		return err
	}
	logger.Infof("creating vpc")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	request := aws_network.CreateVpcRequest{
		CidrBlock:      viper.GetString("cidr"),
		InstanceTenacy: types.Tenancy(types.VpcTenancyDefault),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeVpc,
				Tags: append(tagsFromSlice(viper.GetStringSlice("additional-tags")), types.Tag{
					Key:   to.StringPtr("Name"),
					Value: to.StringPtr(viper.GetString("name")),
				}),
			},
		},
		DryRun: viper.GetBool("dry-run"),
	}
	err = client.CreateVPC(ctx, request)
	if err != nil {
		return err
	}
	logger.Infof("aws vpc '%s' created", viper.GetString("name"))
	return nil
}

func listVPCs() error {
	logger, client, err := getLoggerAndNetworkClient()
	if err != nil {
		return err
	}
	logger.Infof("listing vpcs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var vpcs []types.Vpc
	vpcs, err = client.ListVPCs(ctx)
	if err != nil {
		return err
	}
	for _, vpc := range vpcs {
		if vpc.IsDefault {
			logger.WithField("cidr", *vpc.CidrBlock).Infof("default vpc id %s", *vpc.VpcId)
			continue
		}
		if len(vpc.Tags) == 0 {
			logger.WithField("cidr", *vpc.CidrBlock).Infof("unnamed vpc id: %s", *vpc.VpcId)
			continue
		}
		tags := []string{}
		for _, tag := range vpc.Tags {
			tags = append(tags, fmt.Sprintf("%s: %s", *tag.Key, *tag.Value))
		}
		logger.WithField("cidr", *vpc.CidrBlock).Infof("vpc id %s with tags: %s", *vpc.VpcId, strings.Join(tags, ", "))
	}
	return nil
}

func deleteVPC() error {
	logger, client, err := getLoggerAndNetworkClient()
	if err != nil {
		return err
	}
	logger.Infof("deleting vpc")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return client.DeleteVPC(ctx, viper.GetString("id"))
}

func createSubnetInVPC() error {
	logger, client, err := getLoggerAndNetworkClient()
	if err != nil {
		return err
	}
	logger.Infof("creating subnet in vpc")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return client.CreateSubnetInVPC(ctx, aws_network.CreateVpcSubnetRequest{
		CidrBlock:        viper.GetString("cidr"),
		VPCId:            viper.GetString("id"),
		AvailabilityZone: viper.GetString("az"),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeSubnet,
				Tags: append(tagsFromSlice(viper.GetStringSlice("additional-tags")), types.Tag{
					Key:   to.StringPtr("Availability-Zone"),
					Value: to.StringPtr(viper.GetString("az")),
				}),
			},
		},
		DryRun: viper.GetBool("dry-run"),
	})
}

func listSubnetsInVPC() error {
	logger, client, err := getLoggerAndNetworkClient()
	if err != nil {
		return err
	}
	logger.Infof("listing subnets in vpc")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var subnets []types.Subnet
	subnets, err = client.ListSubnetsInVPC(ctx, viper.GetString("id"))
	if err != nil {
		return err
	}
	for _, subnet := range subnets {
		logger.Infof("subnet az: %s, cidr: %s", *subnet.AvailabilityZone, *subnet.CidrBlock)
	}
	return nil
}

func tagsFromSlice(s []string) (tags []types.Tag) {
	for i := 0; i < len(s); i += 2 {
		if i+1 > len(s) {
			return tags
		}
		tags = append(tags, types.Tag{Key: to.StringPtr(s[i]), Value: to.StringPtr(s[i+1])})
	}
	return tags
}

func listAZs() error {
	logger, client, err := getLoggerAndNetworkClient()
	if err != nil {
		return err
	}
	logger.Infof("listing azs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var azs []types.AvailabilityZone
	azs, err = client.ListAvailabilityZones(ctx)
	if err != nil {
		return err
	}
	for _, az := range azs {
		logger.Infof("az: %s", *az.ZoneName)
	}
	return nil
}
