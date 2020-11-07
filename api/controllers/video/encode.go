package video

import (
	"bytes"
	"config"
	"controllers/user"
	"fmt"
	"io/ioutil"
	"log"
	"models"
	"net/http"
	"net/mail"
	"net/smtp"
	"net/url"
	"text/template"

	"github.com/gorilla/mux"
)

func sendEncodeRequest(bucketName string, format int64, filename string, r *http.Request) {
	log.Println("Sending encode request")
	//Prepare values for sending emails
	uid := mux.Vars(r)["id"]
	user := user.ReqUserByID(uid)
	videoName := r.FormValue("name")

	formData := url.Values{
		"bucket_name": {bucketName},
		"format":      {fmt.Sprintf("%d", format)},
		"filename":    {filename},
	}
	if err := sendEmailEncoding(user, videoName); err != nil {
		log.Println(err)
		return
	}
	resp, err := http.PostForm("http://video_encoder:3001/encode", formData)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if err := sendEmailEncoded(user, videoName); err != nil {
		log.Println(err)
		return
	}
	log.Printf("Request result: %s \n", string(body))
}

func sendEmailEncoding(usr models.User, videoName string) error {
	header := make(map[string]string)
	header["To"] = usr.Email
	header["Subject"] = "[Do not reply]: your video is uploading"
	content, err := getHTMLContent("/go/src/api/controllers/video/emails_template/encoding_email.tmpl", emailValues{usr.Username, videoName})
	if err != nil {
		log.Println(err)
		return nil
	}
	if err = sendEmail(content, header); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func sendEmailEncoded(usr models.User, videoName string) error {
	header := make(map[string]string)
	header["To"] = usr.Email
	header["Subject"] = "[Do not reply]: your video is uploaded"
	content, err := getHTMLContent("/go/src/api/controllers/video/emails_template/encoded_email.tmpl", emailValues{usr.Username, videoName})
	if err != nil {
		log.Println(err)
		return nil
	}
	if err = sendEmail(content, header); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func getHTMLContent(templateFile string, fields emailValues) ([]byte, error) {
	t := template.New("action")
	t, err := t.ParseFiles(templateFile)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	var content bytes.Buffer
	if err = t.Execute(&content, fields); err != nil {
		log.Println(err)
		return []byte{}, err
	}
	return content.Bytes(), nil

}

func sendEmail(content []byte, header map[string]string) error {

	header["Content-Type"] = `text/html; charset="UTF-8"`
	fromHeader := mail.Address{(*config.Api.Postfix).Name, (*config.Api.Postfix).Email}
	from := fromHeader.String()
	header["From"] = from

	bMsg := content
	c, err := smtp.Dial((*config.Api.Postfix).Addr)
	if err != nil {
		log.Println(err)
		return err
	}
	defer c.Close()
	if err = c.Mail(from); err != nil {
		log.Println(err)
		return err
	}
	w, err := c.Data()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = w.Write(bMsg)
	if err != nil {
		log.Println(err)
		return err
	}
	err = c.Quit()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type emailValues struct {
	username  string
	videoname string
}
