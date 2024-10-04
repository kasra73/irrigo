package config

import "go.uber.org/zap/zapcore"

type Server struct {
	BindIP      string `yaml:"bind_ip"`
	BindPort    int    `yaml:"bind_port"`
	KeepIdle    int    `yaml:"keep_idle"`
	BindTimeout int    `yaml:"bind_timeout"`
	Backlog     int    `yaml:"backlog"`
	SwiftDir    string `yaml:"swift_dir"`
	User        string `yaml:"user"`
}

// Logging struct represents the logging key in yaml
type Logging struct {
	Output        string        `yaml:"output"`
	FilePath      string        `yaml:"file_path"`
	Facility      string        `yaml:"facility"`
	Level         zapcore.Level `yaml:"level"`
	Name          string        `yaml:"name"`
	MaxLineLength int           `yaml:"max_line_length"`
	UDPHost       string        `yaml:"udp_host"`
	UDPPort       int           `yaml:"udp_port"`
	Address       string        `yaml:"address"`
}

type Metrics struct {
	StatsdHost             string  `yaml:"statsd_host"`
	StatsdPort             int     `yaml:"statsd_port"`
	StatsdDefaultRate      float64 `yaml:"statsd_default_sample_rate"`
	StatsdSampleRateFactor float64 `yaml:"statsd_sample_rate_factor"`
	StatsdMetricPrefix     string  `yaml:"statsd_metric_prefix"`
}

type GeneralSettings struct {
	Server  Server  `yaml:"server"`
	Logging Logging `yaml:"logging"`
	Metrics Metrics `yaml:"metrics"`
	Custom  any     `yaml:"custom"`
}

type Filter struct {
	Use        string          `yaml:"use"`
	Type       string          `yaml:"type"`
	PluginPath string          `yaml:"plugin_path"`
	Options    map[string]any  `yaml:"options"`
	Overrides  GeneralSettings `yaml:"overrides"`
}

type Global struct {
	Default  GeneralSettings   `yaml:"default"`
	Pipeline []string          `yaml:"pipeline"`
	Filters  map[string]Filter `yaml:"filters"`
}
