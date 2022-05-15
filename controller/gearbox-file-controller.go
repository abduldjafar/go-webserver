package controller

import (
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
	"github.com/gogearbox/gearbox"
)

type gearboxFile struct {
	config *config.Configuration
}

func (g *gearboxFile) SetupConfig(config *config.Configuration) {
	g.config = config
	log.Println(g.config)
}

func (g *gearboxFile) GenerateFileToken() interface{} {
	return func(ctx gearbox.Context) {
		token := g.CreateToken("idx", 7)
		topic := ctx.Param("topic")

		g.sendTokenTokafkaCLient(token, topic)

		ctx.SendJSON(map[string]interface{}{
			"token": token,
		})

	}
}
func (g *gearboxFile) Get(download_path string) interface{} {
	return func(ctx gearbox.Context) {
		fileName := ctx.Param("filename")

		token := ctx.Param("token")

		_, err := g.Validate(token, "idx")

		if err != nil {
			ctx.SendJSON(map[string]interface{}{
				"error": err.Error(),
				"code":  500,
			})
		}

		targetPath := filepath.Join(download_path, fileName)

		if _, err := os.Stat(targetPath); err != nil {
			ctx.SendJSON(map[string]interface{}{
				"error": err,
				"code":  500,
			})
		}

		log.Println(targetPath)
		//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
		ctx.Context().Response.Header.Set("Content-Description", "File Transfer")
		ctx.Context().Response.Header.Set("Content-Transfer-Encoding", "binary")
		ctx.Context().Response.Header.Set("Content-Disposition", "attachment; filename="+fileName)
		ctx.Context().Response.Header.Set("Content-Type", "application/octet-stream")
		ctx.Context().SendFile(targetPath)

	}
}
func (g *gearboxFile) Create(path string) interface{} {
	return func(ctx gearbox.Context) {
		initConfig := g.config

		form, err := ctx.Context().Request.MultipartForm()
		if err != nil {
			ctx.SendJSON(map[string]interface{}{
				"error": err.Error(),
				"code":  500,
			})
		}

		fileHeader := form.File["file"][0]
		idxgroup := form.Value["idxgroup"][0]
		idxtotal := form.Value["idxtotal"][0]
		idxnumber := form.Value["idxnumber"][0]

		file, err := fileHeader.Open()
		if err != nil {
			ctx.SendJSON(map[string]interface{}{
				"error": err.Error(),
				"code":  500,
			})
		}

		go g.sendTokafkaCLient(fileHeader.Filename, idxgroup, initConfig.Kafka.Topic, idxtotal, idxnumber, file, path)

	}
}

func (g *gearboxFile) sendTokafkaCLient(name string, idxGroup string, topic string, idxTotal string, idxNumber string, file multipart.File, path string) {
	initConfig := g.config
	name = strings.ToLower(name)
	filename := name
	log.Println("processing file " + filename)

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
			payload := strings.NewReader("{\n\t\"name\":\"" + name + "\",\n\t\"topic\":\"" + topic + "\",\n\t\"idxGroup\":\"" + idxGroup + "\",\n\t\"idxTotal\":" + idxTotal + ",\n\t\"idxNumber\":" + idxNumber + "\n\t\n}")

			req, _ := http.NewRequest("POST", url, payload)

			req.Header.Add("Content-Type", "application/json")

			res, _ := http.DefaultClient.Do(req)

			if res != nil {

				defer res.Body.Close()
			}
			log.Println("success save file " + filename)

		}

	}

}

func (g *gearboxFile) sendTokenTokafkaCLient(token string, topic string) {
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

func (g *gearboxFile) CreateToken(secret string, day time.Duration) string {

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

func (g *gearboxFile) Validate(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	return token, err
}

func GearboxImplFileController() FileController {
	return &gearboxFile{}
}
