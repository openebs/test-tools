module "vpc" {
    source             = "terraform-aws-modules/vpc/aws"
    name               = var.vpc_name
    cidr               = var.vpc_cidr
    azs                = var.azs
    private_subnets    = var.private_subnets_cidr
    public_subnets     = var.public_subnets_cidr
    enable_nat_gateway = true
    enable_vpn_gateway = true
    tags = {
        Terraform      = "true"
        Environment    = "kafka"
    }
}
 data "aws_eks_cluster" "cluster" {
  name = module.eks.cluster_id
}

data "aws_eks_cluster_auth" "cluster" {
  name = module.eks.cluster_id
}

module "eks" {
  source               = "terraform-aws-modules/eks/aws"
  cluster_name         = var.cluster_name
    cluster_version    = var.k8s_version
    vpc_id             = module.vpc.vpc_id
    subnets            = module.vpc.public_subnets
    //enable_irsa      = true
  worker_groups = [
    
    {
      name                          = "worker-group"
      instance_type                 = var.node_instance_type
      //disk_size                   = var.bot_disk_size
      //ami_id                      = "ami-0c1a7f89451184c8b"
      asg_desired_capacity          = var.desired_nodes
      asg_max_size                  = var.max_nodes
      kubelet_extra_args            = "--node-labels=node.kubernetes.io/lifecycle=spot"
      suspended_processes           = ["AZRebalance"]
      spot_price                    = var.spot_price

      tags = [
        {
          "key"                     = "k8s.io/cluster-autoscaler/enabled"
          "propagate_at_launch"     = "false"
          "value"                   = "true"
        },
        {
          "key"                     = "k8s.io/cluster-autoscaler/${var.cluster_name}"
          "propagate_at_launch"     = "false"
          "value"                   = "true"
        }
      ]
    },
  ]

  workers_additional_policies = ["arn:aws:iam::aws:policy/AutoScalingFullAccess"]
}

resource "null_resource" "kafka_deployment" {
   depends_on = [module.eks]
   count      = var.configure_kafka ? 1 : 0
   provisioner "local-exec" {
    command   = "aws eks --region ${var.aws_region} update-kubeconfig --name ${var.cluster_name}"
   }
   provisioner "local-exec" {
   command    = "kubectl apply -f https://raw.githubusercontent.com/litmuschaos/litmus/master/litmus-portal/cluster-k8s-manifest.yml"
}
   provisioner "local-exec" {
     command  = "kubectl apply -f litmus-kafka-deployer/mysecret.yml"
   }
   provisioner "local-exec" {
     command  = "kubectl apply -f litmus-kafka-deployer/rbac.yml"
   }
   provisioner "local-exec" {
     command  = "kubectl create -f litmus-kafka-deployer/litmus-kafka-deployer-pod.yml"
   }
   provisioner "local-exec" {
     command  = "while [[ $(kubectl get pods litmus-kafka-deployer -o jsonpath='{.status.containerStatuses[0].state..reason}') != 'Completed' ]]; do echo 'waiting for pod to come in Complete state'; done"

   }
}
resource "null_resource" "Cleanup" {
  provisioner "local-exec" {
    when    = destroy
    command = "kubectl get svc  --all-namespaces -o json  | jq -r '.items[] |  select(.spec.type == \"LoadBalancer\") | .metadata.name ,.metadata.namespace' |while read -r svc_name ;read -r svc_ns;do  kubectl delete svc \"$svc_name\" -n \"$svc_ns\"; done"
  }
}
