data "terraform_remote_state" "foo" {
  backend = "s3"
  
  config = {
    key    = "terraform.tfstate"
    bucket = "terraform-state"
    region = "us-east-1"
  }
}
