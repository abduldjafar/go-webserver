package cron

import (
	"fmt"
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

func (c *cronStruct) SetupConfig(config *config.Configuration) {
	c.config = config
}
func (c *cronStruct) Scheduler() {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(7).Week().Do(func() {
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
	token := CreateToken("idx", 7)
	topic := initConfig.Kafka.TokenTopic

	payload := strings.NewReader("{\n\t\"token\":\"" + token + "\",\n\t\"topic\":\"" + topic + "\"\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

func CronScheduler() CronController {
	return &cronStruct{}
}

func CreateToken(secret string, day time.Duration) string {

	expiresAt := time.Now().Add(time.Hour * 24 * day).Unix()
	tk := &Token{

		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println(err)
	}

	return tokenString

}
