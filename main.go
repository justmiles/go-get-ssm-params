package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var (
	ssmPaths    arraySSMPaths
	key         string
	outputPtr   = flag.String("output", "json", "set the desired output")
	versionPtr  = flag.Bool("version", false, "show current version")
	keyPtr      = flag.String("key", "", "if specified, gets a single key and sends value to stdout")
	templatePtr = flag.String("template", "", "if specified, renders custom output from a template file")

	ssmOptionWithDecryption       = true
	ssmOptionMaxResults     int64 = 10
	ssmOptionNextToken      string
)

func main() {
	flag.Var(&ssmPaths, "path", "SSM Parameter path")
	flag.Parse()

	if *versionPtr {
		fmt.Println("v1.7.0")
		os.Exit(0)
	}

  sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

  svc := ssm.New(sess)

	if *keyPtr != "" {
		if len(ssmPaths) < 1 {
			fmt.Println("Please supply a path with -path")
			os.Exit(1)
		}
		o, err := svc.GetParameter(&ssm.GetParameterInput{
			Name:           aws.String(filepath.Join("/", ssmPaths[0], *keyPtr)),
			WithDecryption: aws.Bool(true),
		})
		check(err)
		fmt.Println(*o.Parameter.Value)
		os.Exit(0)
	}

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

	// write custom template
	if *templatePtr != "" {
		t, err := template.ParseFiles(*templatePtr)
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(os.Stdout, output)

		// write output for shell
		//    export key="value"
	} else if *outputPtr == "shell" {
		for k, v := range output {
			fmt.Printf(`export %s="%s"%s`, k, v, "\n")
		}

		// write output for json
		//    {"key":"value"}
	} else if *outputPtr == "json" {
		jsonString, err := json.MarshalIndent(output, "", "  ")
		check(err)
		fmt.Println(string(jsonString))

		// write output for text
		//    key=value
	} else if *outputPtr == "text" {
		for k, v := range output {
			fmt.Printf(`%s=%s%s`, k, v, "\n")
		}

		// write output for ecs
		//    {"name":"key", "value":"value"}
	} else if *outputPtr == "ecs" {
		res := convertToECS(output, true)
		fmt.Println(string(res))

		// write output for terraform-ecs
		//    { "JSONString": "[{\"name\":\"key\",\"value\":\"value\"}]"
	} else if *outputPtr == "terraform-ecs" {
		res := convertToECS(output, false)
		output = map[string]string{
			"JSONString": res,
		}

		b, err := json.MarshalIndent(output, "", "  ")
		check(err)
		fmt.Println(string(b))

		// fail unknown output
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

func convertToECS(output map[string]string, formatted bool) string {
	type ECS struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	ecsOutput := []ECS{}

	for k, v := range output {
		ecsOutputRecord := ECS{
			Name:  k,
			Value: v,
		}
		ecsOutput = append(ecsOutput, ecsOutputRecord)
	}

	sort.Slice(ecsOutput, func(i, j int) bool {
		return ecsOutput[i].Name < ecsOutput[j].Name
	})

	if formatted {
		res, err := json.MarshalIndent(ecsOutput, "", "  ")
		check(err)
		return string(res)
	}

	res, err := json.Marshal(ecsOutput)
	check(err)
	return string(res)
}
