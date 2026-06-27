package config

import (
	"log"
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
		log.Fatalf("❌ Failed to read the configuration file: %v", err)
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("❌ Failed to unmarshal config: %v", err)
	}

	// env overrides
	if port := firstEnvInt("SERVER_PORT", "APP_PORT"); port != 0 {
		AppConfig.Server.Port = port
	}

	if host := firstEnvString("DATABASE_HOST", "DB_HOST"); host != "" {
		AppConfig.Database.Host = host
	}
	if port := firstEnvInt("DATABASE_PORT", "DB_PORT"); port != 0 {
		AppConfig.Database.Port = port
	}
	if user := firstEnvString("DATABASE_USER", "DB_USER"); user != "" {
		AppConfig.Database.User = user
	}
	if password := firstEnvString("DATABASE_PASSWORD", "DB_PASSWORD"); password != "" {
		AppConfig.Database.Password = password
	}
	if name := firstEnvString("DATABASE_NAME", "DB_NAME"); name != "" {
		AppConfig.Database.Name = name
	}
	if charset := firstEnvString("DATABASE_CHARSET", "DB_CHARSET"); charset != "" {
		AppConfig.Database.Charset = charset
	}
	if loc := firstEnvString("DATABASE_LOC", "DB_LOC"); loc != "" {
		AppConfig.Database.Loc = loc
	}
	if parseTime, ok := firstEnvBool("DATABASE_PARSE_TIME", "DB_PARSE_TIME"); ok {
		AppConfig.Database.ParseTime = parseTime
	}

	if s := viper.GetString("JWT_SECRET"); s != "" {
		AppConfig.JWT.Secret = s
	}
	if privateKey := viper.GetString("JWT_PRIVATE_KEY_PATH"); privateKey != "" {
		AppConfig.JWT.PrivateKeyPath = privateKey
	}
	if publicKey := viper.GetString("JWT_PUBLIC_KEY_PATH"); publicKey != "" {
		AppConfig.JWT.PublicKeyPath = publicKey
	}

	log.Println("✅ Configuration file loaded successfully!")
}

func firstEnvString(keys ...string) string {
	for _, key := range keys {
		if value := viper.GetString(key); value != "" {
			return value
		}
	}

	return ""
}

func firstEnvInt(keys ...string) int {
	for _, key := range keys {
		if value := viper.GetInt(key); value != 0 {
			return value
		}
	}

	return 0
}

func firstEnvBool(keys ...string) (bool, bool) {
	for _, key := range keys {
		if value := viper.GetString(key); value != "" {
			return viper.GetBool(key), true
		}
	}

	return false, false
}
