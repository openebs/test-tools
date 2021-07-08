# Setup EKS And Configure Kafka
A terraform script to create a managed kubernetes cluster on AWS EKS with the spot Instance.This script also gives you an option to setup kafka along with prometheus and kafka-provisioned grafana and also litmus portal.
## Prerequisites
Make sure You have installed all of the following prerequisites on your system:
* Kubectl - [Download & Install kubectl](https://kubernetes.io/docs/tasks/tools/) 
* aws-cli - [Download & Install aws-cli](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2-windows.html) and also configure the awscli.
* Terraform - [Download & Install terraform](https://www.terraform.io/downloads.html)

## Step 1: 
  Clone the repo
  ```
  git clone https://github.com/Adarshkumar14/kafka-eks-infra-setup.git
  ```
## Step 2:
### Variables
You can change the value of variable in the variable.tf
 <table>
    <tr>
      <th> Name </th>
      <th> Description</th>
      <th> Example </th>
    </tr>
    <tr>
      <td> <code>vpc_name</code> </td>
      <td> The name of VPC</td>
        <td><code>eks-vpc</code> </td>
    </tr>
    <tr>
      <td><code>vpc_cidr</code></td>
      <td> Subnet CIDR</td>
      <td><code>10.0.0.0/16</code></td>
    </tr>
    <tr>
      <td><code>private_subnets_cidr</code></td>
      <td>Private Subnet CIDR</td>
      <td><code>[ "10.0.1.0/24" , "10.0.2.0/24", "10.0.3.0/24",]</code></td>
  </tr>
  <tr>
    <td><code>public_subnets_cidr</code></td>
    <td>Public Subnet CIDR</td>
    <td><code>[ "10.0.101.0/24" , "10.0.102.0/24", "10.0.103.0/24",]</code></td>
  </tr>
  <tr>
    <td><code>azs</code></td>
    <td>AWS Availability Zones</td>
    <td> <code>[ "ap-south-1a" , "ap-south-1b", "ap-south-1c"]</code></td>
  </tr>
  <tr>
    <td><code>cluster_name</code></td>
    <td>The name of your EKS Cluster</td>
    <td><code>eks-adarsh</code></td>
  </tr>
  <tr>
    <td><code>k8s_version</code></td>
    <td>The desired K8s version to launch</td>
    <td><code>1.17</code></td>
  </tr>
  <tr>
    <td><code>node_instance_type</code></td>
    <td>Worker Node instance type</td>
    <td><code>t3.xlarge</code></td>
  </tr>
  <tr>
    <td><code>desired_nodes</code></td>
    <td> Autoscaling Desired node capacity</td>
    <td><code>3</code></td>
  </tr>
  <tr>
    <td><code>max_nodes</code></td>
    <td>Autoscaling Maximum node capacity</td>
    <td><code>3</code></td>
  </tr>
  <tr>
    <td><code>spot_price</code></td>
    <td> Max spot Price</td>
    <td><code>0.10</code></td>
  </tr>
  <tr>
    <td><code>configure_kafka</code></td>
    <td>To configure kafka enter true else false </td>
    <td> <code>true</code> or <code>false</code> </td>
  </tr>
  </table>
  
##  Step 3:
   Provide the provider in providers.tf
   ```
   provider "aws" {
    region = "ap-south-1"
    profile = "default"
}
```
## Step 4:
   Initialize the  working directory
   ```
    terraform init
  ```
## Step 5:
  **Note:** If you are also configuring  kafka , then also you need to give the aws credentials in the secrets.For better understanding [refer this](https://github.com/litmuschaos/test-tools/tree/master/custom/app-setup/kafka).
  
 After Terraform has been successfully initialized, run <code>terraform apply</code>,It will ask for the setup of kafka ,enter <code>true</code> to setup else enter <code>false</code>.
 
```
terraform apply
```
It takes approximately 15 minutes to complete.\
 **Note**: If you have not configure the kafka , then to interact with your cluster .You have to run the following command in your terminal:
  ```
  aws eks --region aws_region_name update-kubeconfig --name eks_cluster_name
 ```
  Now Your eks-cluster is ready and configured.
  
## CleanUp
  To delete the eks-cluster, run the following command:
  ```
   terraform destroy
 ```
It will delete all the resources that are created by Terraform.

