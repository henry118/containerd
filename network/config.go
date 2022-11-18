/*
   Copyright The containerd Authors.

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

package network

// PluginConfig contains toml config related Network service
type PluginConfig struct {
	// DefaultEnv is the default network env that serves client requests if not otherwise specified
	DefaultEnv string `toml:"default_env" json:"defaultEnv"`
	// Envs is a map of all configured envs
	Envs map[string]EnvConfig `toml:"envs" json:"envs"`
}

// EnvConfig contains toml config related network env
type EnvConfig struct {
	// PluginBinDir is the directory in which the binaries for the plugin are kept.
	PluginBinDir string `toml:"bin_dir" json:"binDir"`
	// PluginConfDir is the directory in which the admin places a network(CNI) conf.
	PluginConfDir string `toml:"conf_dir" json:"confDir"`
	// NetworkPluginMaxConfNum is the max number of plugin config files that will be loaded.
	// Set the value to 0 to load all config files.
	MaxConfNum int `toml:"max_conf_num" json:"maxConfNum"`
	// NetworkPluginMinConfNum is the min number of plugin config files that are loaded
	// for the network to be considered initialized successfully.
	MinConfNum int `toml:"min_conf_num" json:"minConfNum"`
	// IFPrefix is the prefix of all interfaces setup by the network.
	IFPrefix string `toml:"if_prefix" json:"ifPrefix"`
	// LoNetwork indicates whether a loopback network should be automatically setup.
	LoNetwork bool `toml:"lo_network" json:"loNetwork"`
}

func DefaultConfig() *PluginConfig {
	return &PluginConfig{
		DefaultEnv: "default",
		Envs: map[string]EnvConfig{
			"default": {
				PluginBinDir:  "/opt/cni/bin",
				PluginConfDir: "/etc/cni/net.d",
				MaxConfNum:    0,
				MinConfNum:    1,
				IFPrefix:      "eth",
				LoNetwork:     false,
			},
		},
	}
}
