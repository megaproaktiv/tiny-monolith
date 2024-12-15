# Switch Security Group Rules

https://checkip.amazonaws.com

Open/close a security group for ssh for your external IP address.

Store the security group id in environment variable `WEB_SECURITY_GROUP_ID`.

## Description

There is an AWS security group that has an incoming rule that allows SSH access.
With the website "https://checkip.amazonaws.com" you get your public IP address.

Write a GO programm with AWS SDKV2 that has these functions:
- Called with the parameter `open` it opens the security group for an ssh port only for your IP address
- Called with the parameter `close` it deletes all ssh port entries in the security group
- Called with the parameter `list` it prints all entries of the security group

## Try it

Inside directory `toggle`:

- Have a default VPC in your AWS account
- Install https://taskfile.dev
- Install AWS CLI

1) Build the toggle program
  - `task build`
2) Create a security group with an incoming rule for SSH
  - `task create`
3) List the security group
  - `task list`
4) Open the security group for your IP address
  - `task open`
5) Close the security group
  - `task close`
6) List the changed security group
  - `task list`
7) Delete the Stack with the security group
  - `task delete`
