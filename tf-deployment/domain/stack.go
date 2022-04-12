package domain

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type IMyStack interface {
	CreateStack(scope constructs.Construct, id string) cdktf.TerraformStack
}