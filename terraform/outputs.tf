output "public_ip" {
  description = "Public IP address of the EC2 instance"
  value       = aws_instance.main.public_ip
}

output "instance_id" {
  description = "EC2 instance ID"
  value       = aws_instance.main.id
}

output "ssh_private_key" {
  description = "SSH private key for the instance (save to a .pem file)"
  value       = tls_private_key.main.private_key_pem
  sensitive   = true
}

output "ssh_command" {
  description = "SSH command to connect to the instance"
  value       = "ssh -i key.pem admin@${aws_instance.main.public_ip}"
}
