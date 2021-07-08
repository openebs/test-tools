variable "aws_region" {
    default = "ap-south-1"
}
variable "vpc_name" {
    default = "eks-vpc"
}
variable "vpc_cidr" {
    default = "10.0.0.0/16"
}
variable "private_subnets_cidr" {
    type = list
    default = [ "10.0.1.0/24" , "10.0.2.0/24", "10.0.3.0/24",]
}

variable "public_subnets_cidr" {
    type = list
    default = [ "10.0.101.0/24" , "10.0.102.0/24", "10.0.103.0/24",]
}

variable "azs" {
    default = [ "ap-south-1a" , "ap-south-1b", "ap-south-1c"]
}

variable "cluster_name" {
    default =  "eks-adarsh"
}

variable "k8s_version" {
    default = "1.17"
 }
 variable "node_instance_type"{

     default = "t3.xlarge"
 }
 variable "desired_nodes"{
     type = number
     default = 3
 }
 variable "max_nodes"{
     type = number
     default = 3
 }
 variable "spot_price"{
     default= "0.10"
 }

variable "configure_kafka"{
  type = bool
  description= "Enter true to setup kafka else enter false"
}
