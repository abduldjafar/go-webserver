package services

import (
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type Services interface {
	PostTokenToKafkaCLient(token string, topic string, url string)
	PostPathToKafkaClient(name string, fullpathname string, idxGroup string, topic string, idxTotal string, idxNumber string, file multipart.File, path string, url string, HostUrl string, idx_method string)
}

type idxservices struct {
}

func (*idxservices) PostTokenToKafkaCLient(token string, topic string, url string) {
	payload := strings.NewReader("{\n\t\"token\":\"" + token + "\",\n\t\"topic\":\"" + topic + "\"\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	log.Println(string(body))
}
func (*idxservices) PostPathToKafkaClient(filename string, fullpathname string, idxGroup string, topic string, idxTotal string, idxNumber string, file multipart.File, path string, url string, HostUrl string, idx_method string) {

	filename = strings.ToLower(filename)
	fullpathname = strings.Replace(fullpathname, "\\u00", "=", 0)

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

			// payload := strings.NewReader("{\n\t\"name\":\"" + name + "\",\n\t\"topic\":\"" + topic + "\"\n}")
			payload := strings.NewReader("{\n\t\"name\":\"" + fullpathname + "\",\n\t\"topic\":\"" + topic + "\",\n\t\"idxGroup\":\"" + idxGroup + "\",\n\t\"idxTotal\":" + idxTotal + ",\n\t\"idx_method\":\"" + idx_method + "\",\n\t\"path\":\"" + HostUrl + "\",\n\t\"filename\":\"" + filename + "\",\n\t\"idxNumber\":" + idxNumber + "\n\t\n}")

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

func ImplementServices() Services {
	return &idxservices{}
}
