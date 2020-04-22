data "terraform_remote_state" "match" {
  backend = "s3"
  
  config = {
    key    = "terraform.tfstate"
    bucket = "terraform-state"
    region = "us-east-1"
  }
}

data "terraform_remote_state" "wrong_type" {
  backend = "remote"
  
  config = {}
}

data "terraform_remote_state" "wrong_config" {
  backend = "s3"
  
  config = {
    key    = "a-different-terraform.tfstate"
    bucket = "terraform-state"
    region = "us-east-1"
  }
}
