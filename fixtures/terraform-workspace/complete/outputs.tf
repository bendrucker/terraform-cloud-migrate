output "attribute" {
  value = var.environment
}

output "interpolated" {
  value = "The workspace is ${var.environment}"
}

output "function" {
  value = lookup({}, var.environment, false)
}
