#!/usr/bin/env python3

import json
import subprocess
import sys

COMMAND = "docker swarm join-token worker --quiet"

query = json.loads("".join(sys.stdin.readlines()))

#TODO meh.
ssh = subprocess.Popen(["ssh", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "ubuntu@%s" % query["swarm_ip"], COMMAND],
                       shell=False,
                       stdout=subprocess.PIPE,
                       stderr=subprocess.PIPE)

token = ssh.stdout.readline().strip().decode("utf-8")

result = {"token": token}

print(json.dumps(result))
