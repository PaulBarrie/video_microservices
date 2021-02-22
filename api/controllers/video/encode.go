package video

import (
	"config"
	"controllers/user"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"models"
	"net/http"
	"net/smtp"
	"net/url"
	"os"

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
	vals := EmailValues{Username: usr.Username, Videoname: videoName}
	content, err := getHTMLContent("/go/src/api/controllers/video/emails_template/encode_beg.tmpl", vals)
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
	vals := EmailValues{Username: usr.Username, Videoname: videoName}
	content, err := getHTMLContent("/go/src/api/controllers/video/emails_template/encode_end.tmpl", vals)
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

func getHTMLContent(templateFile string, fields EmailValues) (string, error) {
	t, err := template.New("T").ParseFiles(templateFile)
	if err != nil {
		log.Println(err)
		return "", err
	}
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	if err := t.Execute(os.Stdout, fields); err != nil {
		log.Println("Error in exec template")
		log.Println(err)
	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = old
	return string(out[:]), nil

}

func sendEmail(content string, header map[string]string) error {
	// fromHeader := mail.Address{(*config.API.Smtp).Name, (*config.API.Smtp).Email}
	header["Content-Type"] = `text/html; charset="UTF-8"`
	// log.Println("CONTENT\n")
	// log.Println(content)
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject:" + header["Subject"] + "\n"
	message := subject + mime + "\n\n" + content
	// auth := smtp.PlainAuth("", (*config.API.Smtp).Username, (*config.API.Smtp).Password, (*config.API.Smtp).Addr)
	if err := smtp.SendMail((*config.API.Smtp).Addr, nil, (*config.API.Smtp).Email, []string{header["To"]}, []byte(message)); err != nil {
		log.Println("Error SendMail: ", err)
		return err
	}
	log.Println("Email Sent!")
	return nil
}

// func sendEmail(bMsg []byte, header map[string]string) error {
// 	header["Content-Type"] = `text/html; charset="UTF-8"`

// 	auth := smtp.PlainAuth("", (*config.API.Smtp).Username, (*config.API.Smtp).Password, (*config.API.Smtp).Host)
// 	tlsconfig := &tls.Config{
// 		InsecureSkipVerify: true,
// 		ServerName:         (*config.API.Smtp).Addr,
// 	}

// 	tlsDial, err := tls.Dial("tcp", (*config.API.Smtp).Addr, tlsconfig)
// 	if err != nil {
// 		log.Println("Error in Dial")
// 		log.Println(err)
// 		return err
// 	}

// 	c, err := smtp.NewClient(tlsDial, (*config.API.Smtp).Host)
// 	if err != nil {
// 		log.Println("Error in creating client")
// 		log.Println(err)
// 		return err
// 	}
// 	if err = c.Auth(auth); err != nil {
// 		log.Println("Error in auth")
// 		log.Println(err)
// 		return err
// 	}
// 	if err = c.Mail("contact@myyt.com"); err != nil {
// 		log.Println("Error in c.mail")
// 		log.Println(err)
// 		return err
// 	}
// 	if err = c.Rcpt(header["To"]); err != nil {
// 		log.Println("Error in addressing email")
// 		log.Println(err)
// 		return err
// 	}
// 	w, err := c.Data()
// 	if err != nil {
// 		log.Println("Error in creating data")
// 		log.Println(err)
// 		return err
// 	}
// 	_, err = w.Write(bMsg)
// 	if err != nil {
// 		log.Println("Error in write")
// 		log.Println(err)
// 		return err
// 	}
// 	c.Quit()
// 	return nil
// }

//EmailValues defines the custom fiels that will be sent in the emails
type EmailValues struct {
	Username  string
	Videoname string
}
