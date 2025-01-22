package azure

import (
	"errors"
	"fmt"

	"code.il2.gamewarden.io/gamewarden/platform/gravity/pkg/components/cloud/azure/networking"
	"code.il2.gamewarden.io/gamewarden/platform/gravity/pkg/components/cloud/azure/operationalinsights"
	"code.il2.gamewarden.io/gamewarden/platform/gravity/pkg/components/cloud/azure/policy"
	"code.il2.gamewarden.io/gamewarden/platform/gravity/pkg/components/cloud/azure/securityinsights"
	"code.il2.gamewarden.io/gamewarden/platform/gravity/pkg/utils"
	"github.com/pulumi/pulumi-azure-native-sdk/authorization/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/resources/v2"

	// "github.com/pulumi/pulumi-azure-native/sdk/go/azure/securityinsights"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const hubNetworkComponentName = "gravity:azure:hub"

var (
	ErrRequiredArgumentLocation        = errors.New("required argument Location is missing")
	ErrRequiredArgumentBillingAccount  = errors.New("required argument BillingAccount is missing")
	ErrRequiredArgumentCustomer        = errors.New("required argument Customer is missing")
	ErrRequiredArgumentEnvironment     = errors.New("required argument Environment is missing")
	ErrRequiredArgumentComplianceLevel = errors.New("required argument ComplianceLevel is missing")
	ErrProjectTypeConflict             = errors.New("project must be either a Host Project or a Service Project, not both and not neither, ya heard?")
)

type Hub struct {
	pulumi.ResourceState
	ResourceGroupName         pulumi.StringOutput  `pulumi:"resourceGroupName"`
	NatGatewayID              pulumi.IDOutput      `pulumi:"natGatewayID"`
	VirtualNetworkID          pulumi.IDOutput      `pulumi:"virtualNetworkID"`
	VirtualNetworkName        pulumi.StringOutput  `pulumi:"virtualNetworkName"`
	FirewallID                pulumi.IDOutput      `pulumi:"firewallID"`
	BastionID                 pulumi.IDOutput      `pulumi:"bastionID"`
	PrivateDNSZoneID          pulumi.IDOutput      `pulumi:"privateDnsZoneID"`
	LogAnalyticsWorkspaceID   pulumi.IDOutput      `pulmi:"logAnalyticsWorkspaceID"`
	CustomPolicyInitiativeIDs pulumi.IDArrayOutput `pulumi:"customPolicyInitiativeIDs"`
}

type HubArgs struct {
	Customer                    pulumi.StringInput                             `pulumi:"customer"                    validate:"required"`
	Environment                 pulumi.StringInput                             `pulumi:"environment"                 validate:"required"`
	ComplianceLevel             pulumi.StringInput                             `pulumi:"complianceLevel"`
	Location                    pulumi.StringInput                             `pulumi:"location"`
	Tags                        pulumi.StringMap                               `pulumi:"tags"`
	CustomPolicyDefinitionScope pulumi.String                                  `pulumi:"customPolicyDefinitionScope"`
	ResourceGroupName           pulumi.StringInput                             `pulumi:"resourceGroupName"`
	Networking                  *networking.VirtualNetworkArgs                 `pulumi:"networking"`
	Firewall                    *networking.FirewallArgs                       `pulumi:"firewall"`
	Bastion                     *networking.BastionArgs                        `pulumi:"bastion"`
	NatGateway                  *networking.NatGatewayArgs                     `pulumi:"natGateway"`
	PrivateDnsZones             *networking.PrivateDNSArgs                     `pulumi:"privatenszone"`
	LogAnalytics                *operationalinsights.LogAnalyticsWorkspaceArgs `pulumi:"logAnalytics"`
	Sentinel                    *securityinsights.SentinelArgs                 `pulumi:"sentinel"`
	PolicyAssignments           *policy.PolicyAssignmentsArgs                  `pulumi:"policyAssignments"`
}

func NewHub(ctx *pulumi.Context, name string, args *HubArgs, opts ...pulumi.ResourceOption) (*Hub, error) {
	component := &Hub{}

	if err := args.validate(); err != nil {
		return nil, err
	}

	if err := ctx.RegisterComponentResource(hubNetworkComponentName, name, component, opts...); err != nil {
		return nil, fmt.Errorf("unable to register component resource %s, %w", hubNetworkComponentName, err)
	}

	_, err := component.createCustomPolicyDefinitions(ctx, name, args)
	if err != nil {
		return nil, err
	}

	if err := component.createResourceGroup(ctx, name, args); err != nil {
		return nil, err
	}

	if err := component.createVirtualNetwork(ctx, name+"-vnet", args); err != nil {
		return nil, err
	}

	if err := component.createNatGateway(ctx, name+"-nat", args); err != nil {
		return nil, err
	}

	if err := component.createPrivateDNSResources(ctx, name+"-pvtdnszone", args); err != nil {
		return nil, err
	}

	if err := component.createFirewall(ctx, name+"-azfw", args); err != nil {
		return nil, err
	}

	if err := component.createBastion(ctx, name+"-bastion", args); err != nil {
		return nil, err
	}

	if _, err := component.createLogAnalytics(ctx, name+"-loga", args); err != nil {
		return nil, err
	}

	if err := component.assignPolicyInitiatives(ctx, name+"-policy", args); err != nil {
		return nil, err
	}

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{
		"resourceGroupName":         component.ResourceGroupName,
		"virtualNetworkId":          component.VirtualNetworkID,
		"virtualNetworkName":        component.VirtualNetworkName,
		"firewallID":                component.FirewallID,
		"bastionID":                 component.BastionID,
		"logAnalyticsWorkspaceID":   component.LogAnalyticsWorkspaceID,
		"natGatewayID":              component.NatGatewayID,
		"customPolicyInitiativeIDs": component.CustomPolicyInitiativeIDs,
		"privateDnsZoneID":          component.PrivateDNSZoneID,
	}); err != nil {
		return nil, fmt.Errorf("unable to register component resource outputs %s, %w", hubNetworkComponentName, err)
	}

	return component, nil
}

func (a *HubArgs) validate() error {
	if a.NatGateway == nil {
		a.NatGateway = &networking.NatGatewayArgs{}
	}

	if a.LogAnalytics == nil {
		a.LogAnalytics = &operationalinsights.LogAnalyticsWorkspaceArgs{}
	}

	if a.Sentinel == nil {
		a.Sentinel = &securityinsights.SentinelArgs{}
	}

	if a.PolicyAssignments == nil {
		a.PolicyAssignments = &policy.PolicyAssignmentsArgs{}
	}

	switch {
	case a.ComplianceLevel == nil:
		return ErrRequiredArgumentComplianceLevel
	case a.Customer == nil:
		return ErrRequiredArgumentCustomer
	case a.Environment == nil:
		return ErrRequiredArgumentEnvironment
	}

	a.setDefaultTags()

	if err := utils.ValidateStruct(a); err != nil {
		return fmt.Errorf("%T validation failed, %w", a, err)
	}

	return nil
}

func (a *HubArgs) setDefaultTags() {
	if a.Tags == nil {
		a.Tags = pulumi.StringMap{}
	}

	if _, ok := a.Tags["complianceLevel"]; !ok {
		a.Tags["complianceLevel"] = a.ComplianceLevel
	}

	if _, ok := a.Tags["environment"]; !ok {
		a.Tags["environment"] = a.Environment
	}

	if _, ok := a.Tags["customer"]; !ok {
		a.Tags["customer"] = a.Customer
	}
}

func (c *Hub) createCustomPolicyDefinitions(ctx *pulumi.Context, name string, args *HubArgs) (*policy.PolicyDefinitionsState, error) {
	policyDefs, err := policy.NewCustomPolicyDefinitions(ctx, name, args.CustomPolicyDefinitionScope, pulumi.Parent(c))
	if err != nil {
		return nil, fmt.Errorf("unable to create custom policy definitions %s, %w", name, err)
	}

	c.CustomPolicyInitiativeIDs = policyDefs.InitiativeIDs

	return policyDefs, nil
}

func (c *Hub) createResourceGroup(ctx *pulumi.Context, name string, args *HubArgs) error {
	rg, err := resources.NewResourceGroup(ctx, name, &resources.ResourceGroupArgs{
		ResourceGroupName: args.ResourceGroupName,
		Location:          args.Location,
		Tags:              args.Tags,
	}, pulumi.Parent(c))
	if err != nil {
		return fmt.Errorf("unable to create resource group %s, %w", name, err)
	}

	c.ResourceGroupName = rg.Name

	return nil
}

func (c *Hub) createVirtualNetwork(ctx *pulumi.Context, name string, args *HubArgs) error {
	args.Networking.ResourceGroupName = c.ResourceGroupName

	vnet, err := networking.NewVirtualNetwork(ctx, name, args.Networking, pulumi.Parent(c))
	if err != nil {
		return fmt.Errorf("error building networking %w", err)
	}
	c.VirtualNetworkID = vnet.VNetID
	c.VirtualNetworkName = vnet.VNetName

	return nil
}

func (c *Hub) createNatGateway(ctx *pulumi.Context, name string, args *HubArgs) error {
	args.NatGateway.ResourceGroupName = c.ResourceGroupName

	natGateway, err := networking.NewNatGateway(ctx, name, args.NatGateway, pulumi.Parent(c))
	if err != nil {
		return fmt.Errorf("error building networking %w", err)
	}
	c.NatGatewayID = natGateway.NatGatewayID

	return nil
}

func (c *Hub) createPrivateDNSResources(ctx *pulumi.Context, name string, args *HubArgs) error {
	args.NatGateway.ResourceGroupName = c.ResourceGroupName

	privateDnsZone, err := networking.NewPrivateDNS(ctx, name, args.PrivateDnsZones, pulumi.Parent(c))
	if err != nil {
		return fmt.Errorf("error building networking %w", err)
	}
	c.PrivateDNSZoneID = privateDnsZone.PrivateDNSZoneID

	return nil
}

func (c *Hub) createFirewall(ctx *pulumi.Context, name string, args *HubArgs) error {
	args.Firewall.NatGatewayID = c.NatGatewayID
	args.Firewall.VirtualNetworkName = c.VirtualNetworkName
	args.Firewall.ResourceGroupName = c.ResourceGroupName

	azfw, err := networking.NewFirewall(ctx, name, args.Firewall, pulumi.Parent(c))
	if err != nil {
		return fmt.Errorf("error building networking %w", err)
	}
	c.FirewallID = azfw.AzureFirewallID

	return nil
}

func (c *Hub) createBastion(ctx *pulumi.Context, name string, args *HubArgs) error {
	args.Bastion.VirtualNetworkName = c.VirtualNetworkName
	args.Bastion.ResourceGroupName = c.ResourceGroupName

	bastion, err := networking.NewBastion(ctx, name, args.Bastion, pulumi.Parent(c))
	if err != nil {
		return fmt.Errorf("error building networking %w", err)
	}
	c.BastionID = bastion.BastionID

	return nil
}

func (c *Hub) createLogAnalytics(ctx *pulumi.Context, name string, args *HubArgs) (*operationalinsights.LogAnalyticsWorkspace, error) {
	args.LogAnalytics.ResourceGroupName = c.ResourceGroupName

	laws, err := operationalinsights.NewLogAnalyticsWorkspace(ctx, name, args.LogAnalytics, pulumi.Parent(c))
	if err != nil {
		return nil, fmt.Errorf("error building laws %w", err)
	}
	c.LogAnalyticsWorkspaceID = laws.WorkspaceID

	args.Sentinel.ResourceGroupName = c.ResourceGroupName
	args.Sentinel.WorkspaceName = laws.WorkspaceName
	_, err = securityinsights.EnableSentinel(ctx, name+"-sentinel", args.Sentinel, pulumi.Parent(laws))
	if err != nil {
		return nil, fmt.Errorf("error enabling Sentinel %w", err)
	}

	return laws, nil
}

func (c *Hub) assignPolicyInitiatives(ctx *pulumi.Context, name string, args *HubArgs) error {
	client, err := authorization.GetClientConfig(ctx, pulumi.Parent(c))
	if err != nil {
		return fmt.Errorf("error creating authorization client, %w", err)
	}
	args.PolicyAssignments.SubscriptionID = pulumi.String(client.SubscriptionId)
	args.PolicyAssignments.LogAnalyticsWorkspaceID = c.LogAnalyticsWorkspaceID
	args.PolicyAssignments.Tags = args.Tags

	_, err = policy.NewPolicyAssignmentsAtSubscription(ctx, name, *args.PolicyAssignments, pulumi.Parent(c))
	if err != nil {
		return fmt.Errorf("unable to assign policy initiatives %s, %w", name, err)
	}
	return nil
}

func (a *HubArgs) UnmarshalJSON(b []byte) error {
	if err := utils.UnmarshalPulumiArgs(b, a); err != nil {
		return fmt.Errorf("unable to unmarshal infra args, %w", err)
	}

	return nil
}
