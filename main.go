package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"strings"
)

var (
	ssmPaths         arraySSMPaths
	outputAsShellPtr = flag.Bool("shell", false, "optionally output shell command to export variables. otherwise, output as json")
	regionPtr        = flag.String("region", "us-east-1", "aws region")
)

type Parameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Parameters []Parameter

func main() {
	flag.Var(&ssmPaths, "path", "SSM Parameters path to source")
	flag.Parse()

	sess := session.Must(session.NewSession())

	svc := ssm.New(sess, &aws.Config{
		Region: regionPtr,
	})

	var (
		ssmOptionWithDecryption       = true
		ssmOptionMaxResults     int64 = 10
		ssmOptionNextToken      string
	)

	var output Parameters

	for _, ssmPath := range ssmPaths {
		for {
			ssmOpts := ssm.GetParametersByPathInput{
				Path:           &ssmPath,
				WithDecryption: &ssmOptionWithDecryption,
				MaxResults:     &ssmOptionMaxResults,
			}
			if ssmOptionNextToken != "" {
				ssmOpts.NextToken = &ssmOptionNextToken
			}
			ssmResponse, err := svc.GetParametersByPath(&ssmOpts)
			if err != nil {
				fmt.Println(err)
			}
			for _, parameter := range ssmResponse.Parameters {
				s := strings.Split(*parameter.Name, "/")
				output = append(output, Parameter{
					Key:   s[len(s)-1],
					Value: *parameter.Value,
				})
			}
			if ssmResponse.NextToken != nil {
				ssmOptionNextToken = *ssmResponse.NextToken
			} else {
				ssmOptionNextToken = ""
				break
			}
		}
	}

	// write output
	if *outputAsShellPtr {
		for _, p := range output {
			fmt.Printf(`export %s="%s"`, p.Key, p.Value)
			fmt.Println("")
		}
	} else {
		jsonString, _ := json.MarshalIndent(struct{ Parameters }{Parameters: output}, "", "  ")
		fmt.Println(string(jsonString))
	}
}

type arraySSMPaths []string

func (i *arraySSMPaths) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *arraySSMPaths) String() string {
	return ""
}
