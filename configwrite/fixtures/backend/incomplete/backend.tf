terraform {
  backend "s3" {
    key    = "terraform.tfstate"
    bucket = "terraform-state"
    region = "us-east-1"
  }
}
