package main

import (
	"log"
	"time"

	"github.com/andreandradecosta/rpimonitor/db"
	"github.com/andreandradecosta/rpiserver/server"
	"github.com/namsral/flag"
)

var (
	buildInfo string
)

func main() {
	log.Printf("Build info: %s", buildInfo)

	config := flag.String("config", "", "Config file path")
	location := flag.String("LOCATION", "America/Sao_Paulo", "Host location")
	host := flag.String("HOST", "", "Domain")
	httpPort := flag.String("PORT", "8080", "HTTP port")
	httpsPort := flag.String("HTTPS_PORT", "", "HTTPS port")
	isDev := flag.Bool("IS_DEVELOPMENT", false, "Is Dev Env.")
	cert := flag.String("CERT", "", "Certification path")
	key := flag.String("KEY", "", "Private Key path")
	redisHost := flag.String("REDIS_HOST", "localhost:6379", "Redis host:port")
	redisPasswd := flag.String("REDIS_PASSWD", "", "Redis password")
	mongoURL := flag.String("MONGO_URL", "localhost", "mongodb://user:pass@host:port/database")

	flag.Parse()

	log.Println("Starting server...")
	if *config != "" {
		log.Println("Using ", *config)
	}
	hostLocation, err := time.LoadLocation(*location)
	if err != nil {
		panic(err)
	}

	db := db.NewDB(*mongoURL, *redisHost, *redisPasswd)

	s := &server.HTTPServer{
		Host:         *host,
		HTTPPort:     *httpPort,
		HTTPSPort:    *httpsPort,
		IsDev:        *isDev,
		Cert:         *cert,
		Key:          *key,
		RedisPool:    db.RedisPool,
		MongoSession: db.MongoSession,
		HostLocation: hostLocation,
	}
	s.Start()
}
