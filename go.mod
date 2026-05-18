module github.com/PuppyNote/puppynote-server-golang

go 1.22

require (
	firebase.google.com/go/v4 v4.14.1
	github.com/aws/aws-sdk-go-v2 v1.26.0
	github.com/aws/aws-sdk-go-v2/config v1.27.7
	github.com/aws/aws-sdk-go-v2/credentials v1.17.7
	github.com/aws/aws-sdk-go-v2/service/s3 v1.53.0
	github.com/elastic/go-elasticsearch/v8 v8.12.1
	github.com/gin-gonic/gin v1.9.1
	github.com/prometheus/client_golang v1.19.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/redis/go-redis/v9 v9.5.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/viper v1.18.2
	golang.org/x/crypto v0.22.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gorm.io/driver/mysql v1.5.6
	gorm.io/gorm v1.25.9
)
