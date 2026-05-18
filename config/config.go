package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server        ServerConfig
	DB            DBConfig
	JWT           JWTConfig
	Redis         RedisConfig
	Elasticsearch ESConfig
	AWS           AWSConfig
	Email         EmailConfig
	Ollama        OllamaConfig
	Firebase      FirebaseConfig
}

type ServerConfig struct {
	Port        string
	ContextPath string
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}

type JWTConfig struct {
	SecretKey string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type ESConfig struct {
	Host     string
	Username string
	Password string
}

type AWSConfig struct {
	AccessKey       string
	SecretKey       string
	Region          string
	S3Bucket        string
	CloudfrontDomain string
}

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type OllamaConfig struct {
	Host string
}

type FirebaseConfig struct {
	CredentialsPath string
}

var AppConfig Config

func Load() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			log.Fatalf("설정 파일 읽기 실패: %v", err)
		}
		log.Println(".env 파일 없음, 환경변수 사용")
	}

	AppConfig = Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "8080"),
			ContextPath: getEnv("CONTEXT_PATH", "/puppynote"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			Name:     getEnv("DB_NAME", "puppynote"),
			Username: getEnv("DB_USERNAME", ""),
			Password: getEnv("DB_PASSWORD", ""),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", ""),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Elasticsearch: ESConfig{
			Host:     getEnv("ES_HOST", ""),
			Username: getEnv("ES_USERNAME", ""),
			Password: getEnv("ES_PASSWORD", ""),
		},
		AWS: AWSConfig{
			AccessKey:        getEnv("AWS_ACCESS_KEY", ""),
			SecretKey:        getEnv("AWS_SECRET_KEY", ""),
			Region:           getEnv("AWS_REGION", "ap-northeast-2"),
			S3Bucket:         getEnv("S3_BUCKET", "puppynote-dev"),
			CloudfrontDomain: getEnv("CLOUDFRONT_DOMAIN", ""),
		},
		Email: EmailConfig{
			Host:     getEnv("EMAIL_HOST", "smtp.gmail.com"),
			Port:     587,
			Username: getEnv("EMAIL_USERNAME", ""),
			Password: getEnv("EMAIL_PASSWORD", ""),
		},
		Ollama: OllamaConfig{
			Host: getEnv("OLLAMA_HOST", "http://localhost:11434"),
		},
		Firebase: FirebaseConfig{
			CredentialsPath: getEnv("FIREBASE_CREDENTIALS_PATH", "firebase-credentials.json"),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val := viper.GetString(key); val != "" {
		return val
	}
	return defaultVal
}
