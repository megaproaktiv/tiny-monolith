# Switch Security Group Rules

https://checkip.amazonaws.com

Open for ssh or close for ssh

flag: open
flag: close

## Description

There is an AWS security group that has an incoming rule that allows SSH access.
With the website "https://checkip.amazonaws.com" you get your public IP address.

Write a GO programm with AWS SDKV2 that has these functions:
- Called with the parameter `open` it opens the security group for an ssh port only for your IP address
- Called with the parameter `close` it deletes all ssh port entries in the security group
- Called with the parameter `list` it prints all entries of the security group
