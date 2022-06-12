package controller

import (
	"errors"
	"fmt"
	"go-webserver/config"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Token struct {
	*jwt.StandardClaims
}

type ginFile struct {
	config *config.Configuration
}

func (g *ginFile) SetupConfig(config *config.Configuration) {
	g.config = config
}

func (g *ginFile) GenerateFileToken() interface{} {
	return func(ctx *gin.Context) {
		token := g.CreateToken("idx", 7)
		topic := ctx.Param("topic")

		g.sendTokenTokafkaCLient(token, topic)

		ctx.JSON(http.StatusOK, gin.H{"token": token})
	}
}
func (g *ginFile) Get(download_path string) interface{} {
	return func(ctx *gin.Context) {
		fileName := ctx.Param("filename")
		token := ctx.Param("token")

		_, err := g.Validate(token, "idx")
		if err != nil {
			ctx.String(403, "Not Authorized")
			return
		}

		targetPath := filepath.Join(download_path, fileName)
		//This ckeck is for example, I not sure is it can prevent all possible filename attacks - will be much better if real filename will not come from user side. I not even tryed this code
		if !strings.HasPrefix(filepath.Clean(targetPath), download_path) {
			ctx.String(403, "Look like you attacking me")
			return
		}

		if _, err := os.Stat(targetPath); errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does not exist
			ctx.String(404, "file not found")
			return
		}

		//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
		ctx.Header("Content-Description", "File Transfer")
		ctx.Header("Content-Transfer-Encoding", "binary")
		ctx.Header("Content-Disposition", "attachment; filename="+fileName)
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.File(targetPath)

	}
}
func (g *ginFile) Create(path string) interface{} {
	return func(c *gin.Context) {
		initConfig := g.config

		idxgroup := c.Request.FormValue("idxgroup")
		idxtotal := c.Request.FormValue("idxtotal")
		idxnumber := c.Request.FormValue("idxnumber")

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
			return
		}

		filename := header.Filename

		go g.sendTokafkaCLient(filename, idxgroup, initConfig.Kafka.Topic, idxtotal, idxnumber, file, path)

		c.JSON(http.StatusOK, gin.H{"status": "ok"})

	}
}

func (g *ginFile) sendTokafkaCLient(name string, idxGroup string, topic string, idxTotal string, idxNumber string, file multipart.File, path string) {
	initConfig := g.config

	filename := name
	out, err := os.Create(path + "/" + filename)
	if err != nil {
		log.Println(err)
	} else {
		defer out.Close()
		_, err = io.Copy(out, file)

		if err != nil {
			log.Println(err)
		} else {
			url := initConfig.Kafka.UrlProducer + "/produces"
			name = initConfig.Kafka.HostUrl + name

			// payload := strings.NewReader("{\n\t\"name\":\"" + name + "\",\n\t\"topic\":\"" + topic + "\"\n}")
			payload := strings.NewReader("{\n\t\"name\":\"" + name + "\",\n\t\"topic\":\"" + topic + "\",\n\t\"idxGroup\":\"" + idxGroup + "\",\n\t\"idxTotal\":" + idxTotal + ",\n\t\"path\":" + initConfig.Kafka.HostUrl + "\",\n\t\"filename\":" + filename + ",\n\t\"idxNumber\":" + idxNumber + "\n\t\n}")

			req, _ := http.NewRequest("POST", url, payload)

			req.Header.Add("Content-Type", "application/json")

			res, _ := http.DefaultClient.Do(req)

			if res != nil {

				defer res.Body.Close()
			}

		}

	}

}

func (g *ginFile) sendTokenTokafkaCLient(token string, topic string) {
	initConfig := g.config

	url := initConfig.Kafka.UrlProducer + "/cron"

	payload := strings.NewReader("{\n\t\"token\":\"" + token + "\",\n\t\"topic\":\"" + topic + "\"\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}

func (g *ginFile) CreateToken(secret string, day time.Duration) string {

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

func (g *ginFile) Validate(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	return token, err
}

func (g *ginFile) GetWithQuery(download_path string) interface{} {
	return nil
}
func GinImplFileController() FileController {
	return &ginFile{}
}
