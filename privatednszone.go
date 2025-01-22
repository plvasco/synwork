package networking

import (
	"fmt"

	"code.il2.gamewarden.io/gamewarden/platform/gravity/pkg/utils"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const privateDNSComponentName = "gravity:azure:privateDNS"

type PrivateDNS struct {
	pulumi.ResourceState
	PrivateDNSZoneID pulumi.IDOutput `pulumi:"privateDnsZoneID"`
}

type PrivateDNSArgs struct {
	ResourceGroupName pulumi.StringInput `pulumi:"resourceGroupName"`
	ZoneName          pulumi.StringInput `pulumi:"zoneName"`
}

func NewPrivateDNS(ctx *pulumi.Context, name string, args *PrivateDNSArgs, opts ...pulumi.ResourceOption) (*PrivateDNS, error) {
	// Validate args and create  NewPrivateDNS creates a new Private DNS component
	if err := args.validate(); err != nil {
		return nil, fmt.Errorf("validation failed for PrivateDNSArgs: %w", err)
	}

	component := &PrivateDNS{}
	if err := ctx.RegisterComponentResource(privateDNSComponentName, name, component, opts...); err != nil {
		return nil, fmt.Errorf("failed to register Private DNS component resource: %w", err)
	}

	if err := component.createPrivateDNSResources(ctx, name, args); err != nil {
		return nil, fmt.Errorf("failed to create Private DNS resources: %w", err)
	}

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{
		"zoneID": component.PrivateDNSZoneID,
	}); err != nil {
		return nil, fmt.Errorf("failed to register Private DNS outputs: %w", err)
	}

	return component, nil
}

// createPrivateDNSResources handles the creation of resources for the Private DNS component
func (c *PrivateDNS) createPrivateDNSResources(ctx *pulumi.Context, name string, args *PrivateDNSArgs) error {
	zone, err := network.NewPrivateZone(ctx, fmt.Sprintf("%s-zone", name), &network.PrivateZoneArgs{
		ResourceGroupName: args.ResourceGroupName,
		PrivateZoneName:   args.ZoneName,
	}, pulumi.Parent(c))
	if err != nil {
		return fmt.Errorf("failed to create Private Zone: %w", err)
	}

	c.PrivateDNSZoneID = zone.ID().ToIDOutput()
	return nil
}

func (args *PrivateDNSArgs) validate() error {
	if err := utils.ValidateStruct(args); err != nil {
		return fmt.Errorf("%T validation failed: %w", args, err)
	}

	return nil
}

func (args *PrivateDNSArgs) UnmarshalJSON(b []byte) error {
	if err := utils.UnmarshalPulumiArgs(b, args); err != nil {
		return fmt.Errorf("unable to unmarshal %T, %w", args, err)
	}

	return nil
}
