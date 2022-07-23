package cron

import (
	idxauth "go-webserver/auth"
	"go-webserver/config"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-co-op/gocron"
)

type CronController interface {
	SetupConfig(config *config.Configuration)
	Scheduler()
}

type cronStruct struct {
	config *config.Configuration
}

type Token struct {
	*jwt.StandardClaims
}

var (
	auths idxauth.Auth = idxauth.ImplAuthService()
)

func (c *cronStruct) SetupConfig(config *config.Configuration) {
	c.config = config
}
func (c *cronStruct) Scheduler() {
	log.Println("cron will start every day at " + c.config.Kafka.CronTime)
	s := gocron.NewScheduler(time.Local)
	_, err := s.Every(1).Day().At(c.config.Kafka.CronTime).Do(func() {
		c.sendToken()
	})

	if err != nil {
		log.Println("error creating scheduler for token")
	}

	s.StartAsync()

	log.Println("scheduler for creating token started")
}

func (c *cronStruct) sendToken() {
	initConfig := c.config

	url := initConfig.Kafka.UrlProducer + "/cron"
	token := auths.CreateToken("idx", 7)
	topic := initConfig.Kafka.TokenTopic

	payload := strings.NewReader("{\n\t\"token\":\"" + token + "\",\n\t\"topic\":\"" + topic + "\"\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	body, _ := ioutil.ReadAll(res.Body)

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println("error creating token :" + string(body))
	} else {
		log.Println("succes create token")
	}

}

func CronScheduler() CronController {
	return &cronStruct{}
}
