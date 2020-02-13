/*
Copyright © 2019 Red Hat Inc.

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
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func defineFlags(cmd *cobra.Command) {
	cmd.Flags().String("address", "127.0.0.1", "Server address")
	cmd.Flags().String("port", "443", "Server port")
	cmd.Flags().String("volume", "", "Volume name ??")
	cmd.Flags().String("path", "/", "Content path")
	cmd.Flags().String("owner", "", "Object owner")
}

type config struct {
	Address string
	Port    string
	Path    string
}

func parseConfig(cmd *cobra.Command) *config {
	conf := &config{
		Address: getValidStringArg(cmd, "address"),
		Port:    getValidStringArg(cmd, "port"),
		Path:    getValidStringArg(cmd, "path"),
	}
	return conf
}

func getValidStringArg(cmd *cobra.Command, name string) string {
	val, _ := cmd.Flags().GetString(name)
	if val == "" {
		fmt.Fprintf(os.Stderr, "The command line argument '%s' is mandatory.\n", name)
		os.Exit(1)
	}
	return val
}

func main() {
	var srvCmd = &cobra.Command{
		Use:   "resultserver",
		Short: "A tool to serve raw SCAP scan results.",
		Long:  "A tool to serve raw SCAP scan results.",
		Run: func(cmd *cobra.Command, args []string) {
			server(parseConfig(cmd))
		},
	}

	defineFlags(srvCmd)

	if err := srvCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func server(c *config) {
	http.Handle("/", http.FileServer(http.Dir(c.Path)))
	log.Println("Listening...")
	http.ListenAndServe(c.Address+":"+c.Port, nil)
}
