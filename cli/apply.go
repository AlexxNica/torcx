// Copyright 2017 CoreOS Inc.
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

package cli

import (
	"io/ioutil"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/coreos/torcx/pkg/torcx"
)

var (
	cmdApply = &cobra.Command{
		Use:   "apply",
		Short: "apply a profile",
		RunE:  runApply,
	}
)

func init() {
	TorcxCmd.AddCommand(cmdApply)
}

func runApply(cmd *cobra.Command, args []string) error {
	var err error

	commonCfg, err := fillCommonRuntime()
	if err != nil {
		return errors.Wrap(err, "common configuration failed")
	}
	applyCfg, err := fillApplyRuntime(commonCfg)
	if err != nil {
		return errors.Wrap(err, "apply configuration failed")
	}

	if torcx.IsFuseBlown() {
		return errors.New("fuse already blown")
	}

	err = torcx.ApplyProfile(applyCfg)
	if err != nil {
		return errors.Wrap(err, "apply failed")
	}

	err = torcx.BlowFuse(applyCfg)
	if err != nil {
		return errors.Wrapf(err, "blowing fuse")
	}
	return nil
}

// fillApplyRuntime generate runtime config for apply subcommand starting from
// system-wide configuration
func fillApplyRuntime(commonCfg *torcx.CommonConfig) (*torcx.ApplyConfig, error) {
	var (
		profile   string
		storePath []string
	)

	if commonCfg == nil {
		return nil, errors.New("missing common configuration")
	}

	fc, err := ioutil.ReadFile(filepath.Join(commonCfg.ConfDir, profile))
	if err != nil {
		profile = string(fc)
	}
	if profile == "" {
		profile = "vendor"
		logrus.Debug("no profile configured, using default")
	}

	storePath = []string{filepath.Join(commonCfg.BaseDir, "packs")}
	extraStorePath := viper.GetStringSlice("storepath")
	if extraStorePath != nil {
		storePath = append(storePath, extraStorePath...)
	}

	logrus.WithFields(logrus.Fields{
		"profile":   profile,
		"storepath": storePath,
	}).Debug("apply configuration parsed")

	return &torcx.ApplyConfig{
		CommonConfig: *commonCfg,
		Profile:      profile,
		StorePath:    storePath,
	}, nil
}