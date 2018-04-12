package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var (
	ssmPaths   arraySSMPaths
	outputPtr  = flag.String("output", "json", "set the desired output")
	versionPtr = flag.Bool("version", false, "show current version")
	regionPtr  = flag.String("region", "us-east-1", "aws region")

	ssmOptionWithDecryption       = true
	ssmOptionMaxResults     int64 = 10
	ssmOptionNextToken      string
)

func main() {
	flag.Var(&ssmPaths, "path", "SSM Parameter path")
	flag.Parse()

	if *versionPtr {
		fmt.Println("v1.2.0")
		os.Exit(0)
	}

	sess := session.Must(session.NewSession())

	svc := ssm.New(sess, &aws.Config{
		Region: regionPtr,
	})

	output := make(map[string]string)

	// Range over the provided -path arguments
	for _, ssmPath := range ssmPaths {
		ssmOpts := ssm.GetParametersByPathInput{
			Path:           &ssmPath,
			WithDecryption: &ssmOptionWithDecryption,
			MaxResults:     &ssmOptionMaxResults,
		}

		// Loop through the SSM GetParametersByPathInput call until Pagination is complete
		for {
			// consume pagination NextToken if exists
			if ssmOptionNextToken != "" {
				ssmOpts.NextToken = &ssmOptionNextToken
			}

			// perform the request
			ssmResponse, err := svc.GetParametersByPath(&ssmOpts)
			check(err)

			// range over response and store results in memory
			for _, parameter := range ssmResponse.Parameters {
				s := strings.Split(*parameter.Name, "/")
				key := s[len(s)-1]
				output[key] = *parameter.Value
			}

			// if pagination NextToken exists, set it and continue loop. otherwise break loop
			if ssmResponse.NextToken != nil {
				ssmOptionNextToken = *ssmResponse.NextToken
			} else {
				ssmOptionNextToken = ""
				break
			}
		}
	}

	// write output
	if *outputPtr == "shell" {
		for k, v := range output {
			fmt.Printf(`export %s="%s"%s`, k, v, "\n")
		}
	} else if *outputPtr == "json" {
		jsonString, err := json.MarshalIndent(output, "", "  ")
		check(err)
		fmt.Println(string(jsonString))
	} else if *outputPtr == "text" {
		for k, v := range output {
			fmt.Printf(`%s=%s%s`, k, v, "\n")
		}
	} else {
		log.Fatalf(`unknown output "%s"`, *outputPtr)
	}
}

type arraySSMPaths []string

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (i *arraySSMPaths) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *arraySSMPaths) String() string {
	return ""
}
