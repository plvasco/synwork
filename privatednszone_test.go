package networking_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"code.il2.gamewarden.io/gamewarden/platform/gravity/pkg/components/cloud/azure/networking"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestNewPrivateDNS(t *testing.T) {
	t.Parallel()

	type want struct {
		name string
	}

	type args struct {
		name string
		args *networking.PrivateDNSArgs
	}

	tests := []struct {
		name    string
		input   args
		want    want
		wantErr bool
	}{
		{
			name: "should create a new private DNS",
			input: args{
				name: "test-private-dns",
				args: &networking.PrivateDNSArgs{
					ResourceGroupName: pulumi.String("test-resource-group"),
					ZoneName:          pulumi.String("test-zone"),
				},
			},
			want: want{
				name: "test-private-dns",
			},
			wantErr: false,
		},
		{
			name: "should return error for nil resource group name",
			input: args{
				name: "test-private-dns",
				args: &networking.PrivateDNSArgs{
					ResourceGroupName: nil,
					ZoneName:          pulumi.String("test-zone"),
				},
			},
			want: want{
				name: "test-private-dns",
			},
			wantErr: true,
		},
		{
			name: "should return error for nil zone name",
			input: args{
				name: "test-private-dns",
				args: &networking.PrivateDNSArgs{
					ResourceGroupName: pulumi.String("test-resource-group"),
					ZoneName:          nil,
				},
			},
			want: want{
				name: "test-private-dns",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				got, err := networking.NewPrivateDNS(ctx, test.input.name, test.input.args)
				if err != nil {
					return err
				}

				require.NotNil(t, got)

				return nil
			}, nil) // Add testutils mocks if needed.
			if (err != nil) != test.wantErr {
				t.Errorf("error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestPrivateDNSArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    *networking.PrivateDNSArgs
		wantErr bool
	}{
		{
			name: "valid input",
			input: `{
				"resourceGroupName":"test-resource-group",
				"zoneName":"test-zone"
			}`,
			want: &networking.PrivateDNSArgs{
				ResourceGroupName: pulumi.String("test-resource-group"),
				ZoneName:          pulumi.String("test-zone"),
			},
			wantErr: false,
		},
		{
			name: "invalid input",
			input: `{
				"resourceGroupName":42,
				"zoneName":"test-zone"
			}`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing required field",
			input: `{
				"zoneName":"test-zone"
			}`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var args networking.PrivateDNSArgs

			b := []byte(test.input)
			err := json.Unmarshal(b, &args)

			if test.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)

			if !reflect.DeepEqual(test.want, &args) {
				t.Errorf("unexpected result: got %+v, want %+v", args, *test.want)
			}
		})
	}
}
