package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	MaxUniverse      int64
	MaxAgents        int64
	MaxHP            int64
	MaxIdeology      int64
	RecoveryRate     float64
	InitMaxFollowing int64
	MaxSteps         int64
	Pride            float64
	IsNorm           bool
}

var config Config

func Parse(confPath string) error {
	if confPath == "" {
		confPath = "./config.json"
	}

	bytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Printf("failed to open confPath. %s", confPath)
		return err
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return err
	}

	return nil
}

func MaxUniverse() int64 {
	return config.MaxUniverse
}

func MaxAgents() int64 {
	return config.MaxAgents
}

func MaxHP() int64 {
	return config.MaxHP
}

func MaxIdeology() int64 {
	return config.MaxIdeology
}

func RecoveryRate() float64 {
	return config.RecoveryRate
}

func InitMaxFollowing() int64 {
	return config.InitMaxFollowing
}

func MaxSteps() int64 {
	return config.MaxSteps
}

func Pride() float64 {
	return config.Pride
}

func IsNorm() bool {
	return config.IsNorm
}
