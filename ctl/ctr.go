/*
Copyright (C) 2022-2024 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package ctl

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	kzap "sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/apecloud/lorry/dcs"
	"github.com/apecloud/lorry/engines/register"
)

const cliVersionTemplateString = "CLI version: %s \nRuntime version: %s\n"

var configDir string
var disableDNSChecker bool
var opts = kzap.Options{
	Development: true,
	Level:       zap.NewAtomicLevelAt(zap.DPanicLevel),
}

var RootCmd = &cobra.Command{
	Use:   "lorry",
	Short: "Lorry command line interface",
	Long: `
 _      ____  _____  _______     __
 | |    / __ \|  __ \|  __ \ \   / /
 | |   | |  | | |__) | |__) \ \_/ / 
 | |   | |  | |  _  /|  _  / \   /  
 | |___| |__| | | \ \| | \ \  | |   
 |______\____/|_|  \_\_|  \_\ |_|  
===============================
Lorry command line interface`,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		err := viper.BindPFlags(pflag.CommandLine)
		if err != nil {
			return errors.Wrap(err, "fatal error viper bindPFlags")
		}

		// Initialize logger
		kopts := []kzap.Opts{kzap.UseFlagOptions(&opts)}
		if strings.EqualFold("debug", viper.GetString("zap-log-level")) {
			kopts = append(kopts, kzap.RawZapOpts(zap.AddCaller()))
		}
		ctrl.SetLogger(kzap.New(kopts...))

		// Initialize DCS (Distributed Control System)
		err = dcs.InitStore()
		if err != nil {
			return errors.Wrap(err, "DCS initialize failed")
		}

		// Initialize DB Manager
		err = register.InitDBManager(configDir)
		if err != nil {
			return errors.Wrap(err, "DB manager initialize failed")
		}
		return nil
	},

	Run: func(cmd *cobra.Command, _ []string) {
		if versionFlag {
			printVersion()
		} else {
			_ = cmd.Help()
		}
	},
}

type lorryVersion struct {
	CliVersion     string `json:"Cli version"`
	RuntimeVersion string `json:"Runtime version"`
}

var (
	cliVersion       string
	versionFlag      bool
	lorryVer         lorryVersion
	lorryRuntimePath string
)

// Execute adds all child commands to the root command.
func Execute(cliVersion, apiVersion string) {
	lorryVer = lorryVersion{
		CliVersion:     cliVersion,
		RuntimeVersion: apiVersion,
	}

	cobra.OnInitialize(initConfig)

	setVersion()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func setVersion() {
	template := fmt.Sprintf(cliVersionTemplateString, lorryVer.CliVersion, lorryVer.RuntimeVersion)
	RootCmd.SetVersionTemplate(template)
}

func printVersion() {
	fmt.Printf(cliVersionTemplateString, lorryVer.CliVersion, lorryVer.RuntimeVersion)
}

func initConfig() {
	// err intentionally ignored since lorry may not yet be installed.
	runtimeVer := GetRuntimeVersion()

	lorryVer = lorryVersion{
		// Set in Execute() method in this file before initConfig() is called by cmd.Execute().
		CliVersion:     cliVersion,
		RuntimeVersion: strings.ReplaceAll(runtimeVer, "\n", ""),
	}
}

func init() {
	klog.InitFlags(flag.CommandLine)
	opts.BindFlags(flag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	RootCmd.PersistentFlags().StringVar(&configDir, "config-path", "/config/lorry/components/", "Lorry default config directory for builtin type")
	RootCmd.PersistentFlags().BoolVar(&disableDNSChecker, "disable-dns-checker", false, "disable dns checker, for test&dev")
	RootCmd.PersistentFlags().StringVarP(&lorryRuntimePath, "kb-runtime-dir", "", "/kubeblocks/", "The directory of kubeblocks binaries")
	RootCmd.PersistentFlags().AddFlagSet(pflag.CommandLine)
}

// GetRuntimeVersion returns the version for the local lorry runtime.
func GetRuntimeVersion() string {
	lorryCMD := filepath.Join(lorryRuntimePath, "lorry")

	out, err := exec.Command(lorryCMD, "--version").Output()
	if err != nil {
		return "n/a\n"
	}
	return string(out)
}
