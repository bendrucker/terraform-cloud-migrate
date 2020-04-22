output "attribute" {
	value = terraform.workspace
}

output "interpolated" {
	value = "The workspace is ${terraform.workspace}"
}

output "function" {
	value = lookup({}, terraform.workspace, false)
}
