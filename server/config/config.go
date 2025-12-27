package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

var AppConfig *Config

// Init 从 YAML 文件加载配置
func Init() (err error) {
	// 读取配置文件
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return
	}

	// 解析 YAML
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return
	}

	AppConfig = config
	log.Println("Configuration loaded from config.yaml")
	return
}
