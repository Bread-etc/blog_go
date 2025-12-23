package config

import (
	"go-blog/pkg/logger"
	"strings"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// DatabaseConfig 存储数据库连接信息
type DatabaseConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Name      string `mapstructure:"name"`
	Charset   string `mapstructure:"charset"`
	ParseTime bool   `mapstructure:"parseTime"`
	Loc       string `mapstructure:"loc"`
}

type JWTConfig struct {
	Algorithm      string `mapstructure:"algorithm"`
	Secret         string `mapstructure:"secret"`
	PrivateKeyPath string `mapstructure:"private_key_path"`
	PublicKeyPath  string `mapstructure:"public_key_path"`
	ExpireHours    int    `mapstructure:"expire_hours"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

var AppConfig Config

func InitConfig() {
	viper.SetConfigName("config") // 文件名不带扩展名
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		logger.Log.Errorf("❌ Failed to read the configuration file: %v", err)
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		logger.Log.Errorf("❌ Failed to unmarshal config: %v", err)
	}

	// env overridesd
	if s := viper.GetString("JWT_SECRET"); s != "" {
		AppConfig.JWT.Secret = s
	}
	if privateKey := viper.GetString("JWT_PRIVATE_KEY_PATH"); privateKey != "" {
		AppConfig.JWT.PrivateKeyPath = privateKey
	}
	if publicKey := viper.GetString("JWT_PUBLIC_KEY_PATH"); publicKey != "" {
		AppConfig.JWT.PublicKeyPath = publicKey
	}

	logger.Log.Infof("✅ Configuration file loaded successfully!")
}
