package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const (
	sshPort = 22
)

var SecurityGroupID string

func getPublicIP() (string, error) {
	resp, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(ip)), nil
}

func openSSHAccess(cfg aws.Config, ip string) error {
	client := ec2.NewFromConfig(cfg)

	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(SecurityGroupID),
		IpPermissions: []types.IpPermission{
			{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(sshPort),
				ToPort:     aws.Int32(sshPort),
				IpRanges: []types.IpRange{
					{
						CidrIp: aws.String(fmt.Sprintf("%s/32", ip)),
					},
				},
			},
		},
	}

	_, err := client.AuthorizeSecurityGroupIngress(context.TODO(), input)
	return err
}

func closeSSHAccess(cfg aws.Config) error {
	client := ec2.NewFromConfig(cfg)

	// Describe the security group to get its current ingress rules
	describeInput := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{SecurityGroupID},
	}

	describeOutput, err := client.DescribeSecurityGroups(context.TODO(), describeInput)
	if err != nil {
		return fmt.Errorf("error describing security group: %w", err)
	}

	if len(describeOutput.SecurityGroups) == 0 {
		return fmt.Errorf("security group %s not found", SecurityGroupID)
	}

	// Iterate over the ingress rules and find all that match the SSH port
	var sshPermissions []types.IpPermission
	for _, permission := range describeOutput.SecurityGroups[0].IpPermissions {
		if permission.FromPort != nil && permission.ToPort != nil &&
			*permission.FromPort == sshPort && *permission.ToPort == sshPort &&
			*permission.IpProtocol == "tcp" {
			sshPermissions = append(sshPermissions, permission)
		}
	}

	if len(sshPermissions) == 0 {
		fmt.Println("No SSH ingress rules found to revoke")
		return nil
	}

	// Revoke the SSH ingress rules
	revokeInput := &ec2.RevokeSecurityGroupIngressInput{
		GroupId:       aws.String(SecurityGroupID),
		IpPermissions: sshPermissions,
	}

	_, err = client.RevokeSecurityGroupIngress(context.TODO(), revokeInput)
	if err != nil {
		return fmt.Errorf("error revoking SSH ingress rules: %w", err)
	}

	return nil
}

func listSecurityGroupEntries(cfg aws.Config, securityGroupID string) error {
	client := ec2.NewFromConfig(cfg)

	// Describe the security group to get its current ingress rules
	describeInput := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{SecurityGroupID},
	}

	describeOutput, err := client.DescribeSecurityGroups(context.TODO(), describeInput)
	if err != nil {
		return fmt.Errorf("error describing security group: %w", err)
	}

	if len(describeOutput.SecurityGroups) == 0 {
		return fmt.Errorf("security group %s not found", securityGroupID)
	}

	// Print the ingress rules
	fmt.Println("Ingress rules for security group:", securityGroupID)
	for _, permission := range describeOutput.SecurityGroups[0].IpPermissions {
		fmt.Printf("Protocol: %s, FromPort: %d, ToPort: %d\n", *permission.IpProtocol, *permission.FromPort, *permission.ToPort)
		for _, ipRange := range permission.IpRanges {
			fmt.Printf("  CIDR: %s\n", *ipRange.CidrIp)
		}
		for _, userIdGroupPair := range permission.UserIdGroupPairs {
			fmt.Printf("  UserIdGroupPair: %s\n", *userIdGroupPair.GroupId)
		}
		for _, ipv6Range := range permission.Ipv6Ranges {
			fmt.Printf("  IPv6 CIDR: %s\n", *ipv6Range.CidrIpv6)
		}
	}

	return nil
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <open|close>")
		return
	}

	action := os.Args[1]

	SecurityGroupID = os.Getenv("WEB_SECURITY_GROUP_ID")
	if SecurityGroupID == "" {
		fmt.Println("Environment variable WEB_SECURITY_GROUP_ID is not set")
		return
	}

	fmt.Printf("Security Group ID: >%s<\n", SecurityGroupID)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS configuration:", err)
		return
	}

	switch action {
	case "open":
		ip, err := getPublicIP()
		if err != nil {
			fmt.Println("Error getting public IP:", err)
			return
		}
		fmt.Printf("Public IP: %s\n", ip)

		err = openSSHAccess(cfg, ip)
		if err != nil {
			fmt.Println("Error opening SSH access:", err)
			return
		}

		fmt.Println("SSH access opened for IP:", ip)

	case "close":
		err = closeSSHAccess(cfg)
		if err != nil {
			fmt.Println("Error closing SSH access:", err)
			return
		}

		fmt.Println("SSH access closed")
	case "list":
		err = listSecurityGroupEntries(cfg, SecurityGroupID)
		if err != nil {
			fmt.Println("Error listing security group entries:", err)
			return
		}

	default:
		fmt.Println("Invalid action. Use 'open' or 'close'.")
	}
}
