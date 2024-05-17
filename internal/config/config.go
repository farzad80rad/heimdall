package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type HeimdallConfig struct {
	HeimdallPort int         `yaml:"heimdall_port"`
	ApisConfig   []ApiConfig `yaml:"apis_config"`
}

type ApiConfig struct {
	Match                MatchPolicy          `yaml:"match_policy"`
	LoadBalancePolicy    HostLoadPolicy       `yaml:"load_balance"`
	HealthCheckConfig    *HealthCheckConfig   `yaml:"health_check_config"`
	CircuitBreakerConfig CircuitBreakerConfig `yaml:"circuit_breaker_config"`
}

type RequestBodyCheckConfig struct {
	MandatoryFields []RequestValidationUnit `yaml:"mandatory_fields"`
}

type RequestValidationUnit struct {
	Name string `yaml:"field_name"`
	Type string `yaml:"type"`
}

type ConnectionType string

const (
	ConnectionType_HTTP1 ConnectionType = "http"
	ConnectionType_GPRC  ConnectionType = "grpc"
)

type HostMatchInfo struct {
	RequestBodyCheckConfig *RequestBodyCheckConfig `yaml:"request_body_check"`
	SupportedType          string                  `yaml:"type"`
}

type MatchPolicy struct {
	ConnectionType     ConnectionType  `yaml:"connection_type"`
	Name               string          `yaml:"name"`
	Path               string          `yaml:"path"`
	SupportedRestTypes []HostMatchInfo `yaml:"per_method"`
}

type LoadBalanceType string

const (
	LoadBalanceType_ROUNDROBIN          LoadBalanceType = "round_robin"
	LoadBalanceType_WEIGHTED_ROUNDROBIN LoadBalanceType = "weighted_round_robin"
)

type HealthCheckConfig struct {
	Path              string        `yaml:"path"`
	FailureThreshHold int           `yaml:"failure_threshold"`
	Interval          time.Duration `yaml:"interval"`
}

type HostUnit struct {
	Host   string `yaml:"host"`
	Weight int    `yaml:"load_balance_weight"`
}

type CircuitBreakerConfig struct {
	QuarantineDuration    time.Duration `yaml:"quarantine_duration"`
	FailureToleranceCount uint32        `yaml:"failure_tolerance_count"`
}

type HostLoadPolicy struct {
	LoadBalanceType LoadBalanceType `yaml:"type"`
	HostUnits       []HostUnit      `yaml:"host_units"`
}

func ReadConfig(filename string) (*HeimdallConfig, error) {
	// Read the YAML file contents
	data, err := os.ReadFile(filename) // Or use os.ReadFile(filename)
	if err != nil {
		log.Println("error reading YAML file: ", err)
		return nil, err
	}

	// Unmarshal the YAML data into the config struct
	config := &HeimdallConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Println("error parsing YAML data: ", err)
		return nil, err
	}

	return config, nil
}
