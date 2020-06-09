/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package kubectl

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/multi-tenancy/benchmarks/kubectl-mtb/pkg/benchmark"
	"sigs.k8s.io/multi-tenancy/benchmarks/kubectl-mtb/test"
)

var k8sClient *kubernetes.Clientset
var rootCmd *cobra.Command
var maxProfileLevel = 3
var benchmarks []*benchmark.Benchmark

// singular-plural
var benchmarkResourceNames = []string{"benchmark", "benchmarks"}

func init() {

	kubecfgFlags := genericclioptions.NewConfigFlags(false)

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "kubectl-mtb",
		Short: "Multi-Tenancy Benchmarking",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			config, err := kubecfgFlags.ToRESTConfig()
			if err != nil {
				return err
			}

			// create the K8s clientset
			k8sClient, err = kubernetes.NewForConfig(config)
			if err != nil {
				return err
			}

			validateResource(args)

			// Initiate new suite instance
			bs := test.NewBenchmarkSuite()

			profileLevel, _ := cmd.Flags().GetInt("profile-level")
			benchmarks = bs.GetBenchmarksOfProfileLevel(profileLevel)

			return nil
		},
	}

	rootCmd.PersistentFlags().StringP("category", "c", "", "Category of the benchmarks.")
	rootCmd.PersistentFlags().IntP("profile-level", "p", maxProfileLevel, "ProfileLevel of the benchmarks.")

	// Commands
	rootCmd.AddCommand(newGetCmd())
	rootCmd.AddCommand(newTestCmd())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func validateResource(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: Please specify any resource to get.")
		os.Exit(1)
	}
	// Check the resourceName
	resourceName := args[0]
	for _, nm := range benchmarkResourceNames {
		if strings.EqualFold(nm, resourceName) {
			resource = nm
			break
		}
	}

	switch resource {
	case "benchmark":
	case "benchmarks":
	default:
		fmt.Println("Error: Please specify any valid resource to get.")
		os.Exit(1)
	}
}
