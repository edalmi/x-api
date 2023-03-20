// Copyright 2023 Edson Michaque
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	profile    string
)

const (
	flagConfig              = "config-file"
	envPrefix               = "X"
	configFileName          = "config"
	defaultConfigFileFormat = "yml"
)

var v = viper.New()

func Execute() error {
	return New().Execute()
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "x-api",
		SilenceUsage: true,
	}

	cmd.AddCommand(NewCmdStart())
	cmd.AddCommand(NewCmdMigrate())
	cobra.OnInitialize(configLoad(v))

	cmd.PersistentFlags().StringVarP(&configFile, flagConfig, "c", "", "Configuration file")
	cmd.PersistentFlags().StringP("mode", "m", "", "Configuration file")

	v.SetEnvPrefix(envPrefix)
	if err := v.BindPFlags(cmd.PersistentFlags()); err != nil {
		panic(err)
	}

	return cmd
}

func configLoad(v *viper.Viper) func() {
	return func() {
		if configFile != "" {
			v.SetConfigFile(configFile)
		} else if configFile := os.Getenv("X_CONFIG_FILE"); configFile != "" {
			v.SetConfigFile(configFile)
		} else {
			v.AddConfigPath("/etc/x-api")
			v.SetConfigType(defaultConfigFileFormat)

			v.SetConfigName(configFileName)
		}

		v.AutomaticEnv()

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				fmt.Println("Found error: ", err.Error())
			}
		}
	}
}
