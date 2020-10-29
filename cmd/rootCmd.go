// Copyright 2020 Kubestr Developers

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kastenhq/kubestr/pkg/kubestr"
	"github.com/spf13/cobra"
)

var (
	output  string
	rootCmd = &cobra.Command{
		Use:   "kubestr",
		Short: "A tool to validate kubernetes storage",
		Long: `kubestr is a tool that will scan your k8s cluster
		and validate that the storage systems in place as well as run
		performance tests.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			Baseline(ctx, output)
		},
	}

	fioCheckerStorageClass string
	fioCheckerConfigMap    string
	fioCheckerTestName     string
	fioCmd                 = &cobra.Command{
		Use:   "fio",
		Short: "Runs an fio test",
		Long:  `Run an fio test`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			Fio(ctx, output, fioCheckerStorageClass, fioCheckerConfigMap, fioCheckerTestName)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Options(json)")

	rootCmd.AddCommand(fioCmd)
	fioCmd.Flags().StringVarP(&fioCheckerStorageClass, "storageclass", "s", "", "The Name of a storageclass")
	fioCmd.Flags().StringVarP(&fioCheckerConfigMap, "configmap", "c", "", "The Name of a configmap in the current namespace")
	fioCmd.Flags().StringVarP(&fioCheckerTestName, "testname", "t", "", "The Name of a predefined kubestr fio test")
	// //rootCmd.AddCommand(provCmd)
}

// Execute executes the main command
func Execute() error {
	return rootCmd.Execute()
}

// Baseline executes the baseline check
func Baseline(ctx context.Context, output string) {
	p, err := kubestr.NewKubestr()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Print(kubestr.Logo)
	result := p.KubernetesChecks()
	if output == "json" {
		jsonRes, _ := json.MarshalIndent(result, "", "    ")
		fmt.Println(string(jsonRes))
		return
	}
	for _, retval := range result {
		retval.Print()
		fmt.Println()
		time.Sleep(500 * time.Millisecond)
	}

	provisionerList, err := p.ValidateProvisioners(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if output == "json" {
		jsonRes, _ := json.MarshalIndent(result, "", "    ")
		fmt.Println(string(jsonRes))
		return
	}
	fmt.Println("Available Storage Provisioners:")
	fmt.Println()
	time.Sleep(500 * time.Millisecond) // Added to introduce lag.
	for _, provisioner := range provisionerList {
		provisioner.Print()
		fmt.Println()
		time.Sleep(500 * time.Millisecond)
	}
}

// Fio executes the FIO test.
func Fio(ctx context.Context, output, storageclass, configMapName, testName string) {
	p, err := kubestr.NewKubestr()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	result := p.FIO(ctx, storageclass, configMapName, testName)
	if output == "json" {
		jsonRes, _ := json.MarshalIndent(result, "", "    ")
		fmt.Println(string(jsonRes))
		return
	}
	result.Print()
}
