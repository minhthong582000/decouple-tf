package ec2

import (
	"cdk.tf/go/stack/domain"
	"cdk.tf/go/stack/generated/hashicorp/aws"
	"cdk.tf/go/stack/generated/hashicorp/aws/ec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type MyEC2Stack struct {
    name string // Instance name
}

func NewEC2Stack(name string) domain.IMyStack {
	return &MyEC2Stack{
        name: name,
    }
}

func (m MyEC2Stack) CreateStack(scope constructs.Construct, id string) cdktf.TerraformStack {
    stack := cdktf.NewTerraformStack(scope, &id)

    aws.NewAwsProvider(stack, jsii.String("AWS"), &aws.AwsProviderConfig{
        Region: jsii.String("us-east-1"),
    })

    instance := ec2.NewInstance(stack, jsii.String("compute"), &ec2.InstanceConfig{
        Ami:          jsii.String("ami-0c02fb55956c7d316"),
        InstanceType: jsii.String("t2.micro"),
        Tags: &map[string]*string{
            "Name": jsii.String(m.name),
        },
    })

    cdktf.NewTerraformOutput(stack, jsii.String("public_ip"), &cdktf.TerraformOutputConfig{
        Value: instance.PublicIp(),
    })

    return stack
}
