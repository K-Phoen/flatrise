provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region     = "${var.aws_region}"
}

resource "aws_instance" "manager" {
  count = 1

  ami           = "ami-08182c55a1c188dee"
  instance_type = "t2.micro"

  tags = {
    Name = "swarm-manager"
  }

  connection {
    user = "ubuntu"
  }

  # The name of our SSH keypair we created above.
  key_name = "${aws_key_pair.auth.id}"

  # Our Security group to allow HTTP and SSH access
  vpc_security_group_ids = ["${aws_security_group.swarm.id}"]

  provisioner "remote-exec" {
    inline = [
      "sudo apt-get update",
      "sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common",
      "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
      "sudo add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable\"",
      "sudo apt-get update",
      "sudo apt-get install -y docker-ce",
      "sudo usermod -aG docker ubuntu",
      "sudo docker swarm init --advertise-addr ${aws_instance.manager.private_ip}",
    ]
  }
}

data "external" "swarm_token" {
  program = ["python3", "${path.module}/get_swarm_token.py"]
  query = {
    swarm_ip = "${aws_instance.manager.public_ip}"
  }
}

output "manager.ip" {
  value = "${aws_instance.manager.public_ip}"
}

resource "aws_instance" "small_worker" {
  count = 1

  ami           = "ami-08182c55a1c188dee"
  instance_type = "t2.small"

  tags = {
    Name = "swarm-small-worker-${count.index}"
  }

  connection {
    user = "ubuntu"
  }

  # The name of our SSH keypair we created above.
  key_name = "${aws_key_pair.auth.id}"

  # Our Security group to allow HTTP and SSH access
  vpc_security_group_ids = ["${aws_security_group.swarm.id}"]

  provisioner "remote-exec" {
    inline = [
      "sudo apt-get update",
      "sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common",
      "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
      "sudo add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable\"",
      "sudo apt-get update",
      "sudo apt-get install -y docker-ce",
      "sudo usermod -aG docker ubuntu",
      "echo 'Swarm token: ${data.external.swarm_token.result.token}'",
      "echo 'Swarm manager private IP: ${aws_instance.manager.private_ip}'",
      "sudo docker swarm join --token ${data.external.swarm_token.result.token} ${aws_instance.manager.private_ip}:2377",
    ]
  }
}

resource "aws_instance" "micro_worker" {
  count = 3

  ami           = "ami-08182c55a1c188dee"
  instance_type = "t2.micro"

  tags = {
    Name = "swarm-micro-worker-${count.index}"
  }

  connection {
    user = "ubuntu"
  }

  # The name of our SSH keypair we created above.
  key_name = "${aws_key_pair.auth.id}"

  # Our Security group to allow HTTP and SSH access
  vpc_security_group_ids = ["${aws_security_group.swarm.id}"]

  provisioner "remote-exec" {
    inline = [
      "sudo apt-get update",
      "sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common",
      "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
      "sudo add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable\"",
      "sudo apt-get update",
      "sudo apt-get install -y docker-ce",
      "sudo usermod -aG docker ubuntu",
      "echo 'Swarm token: ${data.external.swarm_token.result.token}'",
      "echo 'Swarm manager private IP: ${aws_instance.manager.private_ip}'",
      "sudo docker swarm join --token ${data.external.swarm_token.result.token} ${aws_instance.manager.private_ip}:2377",
    ]
  }
}
