package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Re_Email struct {
	Header  string
	From    string
	Sent    string
	To      string
	Date    string
	CC      string
	Subject string
	Content string

	re any
}
type Email struct {
	//User                      string
	//Directorio                string
	Message_ID                string
	Date                      string
	From                      string
	To                        string
	Subject                   string
	Cc                        string
	Mime_Version              string
	Content_Type              string
	Content_Transfer_Encoding string
	Bcc                       string
	X_From                    string
	X_To                      string
	X_cc                      string
	X_bcc                     string
	X_Folder                  string
	X_Origin                  string
	X_FileName                string
	Content                   string
	re                        Re_Email
}

func IngresarEmails() {

	usuarios, err := ioutil.ReadDir("../../enron_mail_20110402/maildir")
	if err != nil {
		log.Fatal(err)
	}

	for _, usuario := range usuarios {
		directorios, error := ioutil.ReadDir("../../enron_mail_20110402/maildir/" + usuario.Name())
		if error != nil {
			//log.Fatal(error)
		}
		var emails []Email
		for _, directorio := range directorios {

			archivos, error := ioutil.ReadDir("../../enron_mail_20110402/maildir/" + usuario.Name() + "/" + directorio.Name())
			if error != nil {
				//log.Fatal(error)
			}

			for _, archivo := range archivos {

				if error != nil {
					log.Fatal(archivo)
				}
				//emails.PushBack(Emails(usuario.Name(), directorio.Name(), archivo.Name()))
				emails = append(emails, Emails(usuario.Name(), directorio.Name(), archivo.Name()))

			}

		}
		res := map[string]interface{}{"index": "email", "records": emails}
		resJSON, err := json.Marshal(res)
		file, err := os.Create("ejemplo.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = file.WriteString(string(resJSON))
		if err != nil {
			panic(err)
		}
		fmt.Println("cambio de directorio")
		cmd := "curl http://localhost:4080/api/_bulkv2 -i -u admin:0208Mavl  --data-binary '@ejemplo.json'"
		out := string(Cmd(cmd, true))
		//out := string(Cmd(cmd,false))
		fmt.Println(out)
		//

	}
	fmt.Println("termino")
}
func Emails(usuario string, directorio string, archivos string) Email {

	var email Email
	//email.User = usuario
	//email.Directorio = directorio + "/" + archivos

	archivoAbierto, err := os.Open("../../enron_mail_20110402/maildir/" + usuario + "/" + directorio + "/" + archivos)
	if err != nil {
		//log.Fatal(err)
	}
	//fmt.Println("../enron_mail_20110402/maildir/" + usuario.Name() + "/" + directorio.Name() + "/" + archivo.Name())
	//fmt.Println("archivoAbierto")
	fileScanner := bufio.NewScanner(archivoAbierto)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	archivoAbierto.Close()
	var i int = 0
	for i < len(lines) {

		nombresComoArreglo := strings.Split(lines[i], ":")
		switch nombresComoArreglo[0] {
		case "Message-ID":
			email.Message_ID = nombresComoArreglo[1]
		case "Date":
			if len(nombresComoArreglo) > 1 {
				email.Date = nombresComoArreglo[1]
			}
		case "From":
			if len(nombresComoArreglo) > 1 {
				email.From = nombresComoArreglo[1]
			}

		case "To":
			if len(nombresComoArreglo) > 1 && email.X_FileName == "" {
				email.To = nombresComoArreglo[1]
				if strings.Split(lines[i+1], ":")[0] != "Subject" {
					if strings.ToLower(strings.Split(lines[i+1], ":")[0]) != "cc" {

						var val bool = true
						for val {
							if email.X_FileName != "" {
								val = false
							}
							if len(lines) > (i + 1) {
								if strings.Split(lines[i+1], ":")[0] != "Subject" {
									if strings.ToLower(strings.Split(lines[i+1], ":")[0]) != "cc" {
										//fmt.Println(strings.Split(lines[(i+1)], ":")[0])
										email.To = email.To + strings.Split(lines[(i+1)], ":")[0]
										i++

									} else {
										val = false
									}
								} else {
									val = false
								}
							}
						}
					}

				}
			}

		case "Subject":
			if len(nombresComoArreglo) > 1 {
				if nombresComoArreglo[1] == "Re" || nombresComoArreglo[1] == "RE" {
					if len(nombresComoArreglo) > 2 {
						email.Subject = "RE: " + nombresComoArreglo[2]
					}
				} else if nombresComoArreglo[1] == "FW" || nombresComoArreglo[1] == "fw" {
					if len(nombresComoArreglo) > 2 {
						email.Subject = "FW: " + nombresComoArreglo[2]
					}
				} else {
					email.Subject = nombresComoArreglo[1]
				}
			}
		case "Cc":
			if len(nombresComoArreglo) > 1 {
				email.Cc = nombresComoArreglo[1]
			}
		case "Mime-Version":
			if len(nombresComoArreglo) > 1 {
				email.Mime_Version = nombresComoArreglo[1]
			}
		case "Content-Type":
			if len(nombresComoArreglo) > 1 {
				email.Content_Type = nombresComoArreglo[1]
			}
		case "Content-Transfer-Encoding":
			if len(nombresComoArreglo) > 1 {
				email.Content_Transfer_Encoding = nombresComoArreglo[1]
			}
		case "Bcc":
			if len(nombresComoArreglo) > 1 {
				email.Bcc = nombresComoArreglo[1]
			}
		case "X-From":
			if len(nombresComoArreglo) > 1 {
				email.X_From = nombresComoArreglo[1]
			}
		case "X-To":
			if len(nombresComoArreglo) > 1 && email.X_FileName == "" {
				email.X_To = nombresComoArreglo[1]
				if strings.ToLower(strings.Split(lines[i+1], ":")[0]) != "x-cc" {
					if strings.ToLower(strings.Split(lines[i+1], ":")[0]) != "x-bcc" {

						var val bool = true
						for val {
							if email.X_FileName != "" {
								val = false
							}
							if len(lines) > (i + 1) {
								if strings.ToLower(strings.Split(lines[i+1], ":")[0]) != "x-cc" {
									if strings.ToLower(strings.Split(lines[i+1], ":")[0]) != "x-bcc" {
										//fmt.Println(strings.Split(lines[(i+1)], ":")[0])
										email.X_To = email.To + strings.Split(lines[(i+1)], ":")[0]
										i++

									} else {
										val = false
									}
								} else {
									val = false
								}
							}
						}
					}

				}
			}

		case "X-cc":
			if len(nombresComoArreglo) > 1 {
				email.X_cc = nombresComoArreglo[1]
			}
		case "X-bcc":
			if len(nombresComoArreglo) > 1 {
				email.X_bcc = nombresComoArreglo[1]
			}
		case "X-Folder":
			if len(nombresComoArreglo) > 1 {
				email.X_Folder = nombresComoArreglo[1]
			}
		case "X-Origin":
			if len(nombresComoArreglo) > 1 {
				email.X_Origin = nombresComoArreglo[1]
			}
		case "X-FileName":
			if len(nombresComoArreglo) > 1 {
				email.X_FileName = nombresComoArreglo[1]
			}
		case " -----Original Message-----":
			email.re = addEmail(i, "../../enron_mail_20110402/maildir/"+usuario+"/"+directorio+"/"+archivos)

			i = len(lines)
		case "-----Original Message-----":
			email.re = addEmail(i, "../../enron_mail_20110402/maildir/"+usuario+"/"+directorio+"/"+archivos)

			i = len(lines)
		case "----- Original Message-----":
			email.re = addEmail(i, "../../enron_mail_20110402/maildir/"+usuario+"/"+directorio+"/"+archivos)

			i = len(lines)
		case "--------- Inline attachment follows ---------":
			email.re = addEmail(i, "../../enron_mail_20110402/maildir/"+usuario+"/"+directorio+"/"+archivos)

			i = len(lines)
		default:
			if len(strings.Split(nombresComoArreglo[0], "---------------------- Forwarded by")) == 1 {
				email.Content = email.Content + `\n` + nombresComoArreglo[0]
			} else {
				email.re = addEmail(i, "../../enron_mail_20110402/maildir/"+usuario+"/"+directorio+"/"+archivos)

				i = len(lines)
			}

		}
		i++

	}

	return email

}

func addEmail(j int, dir string) Re_Email {

	archivoAbierto, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(archivoAbierto)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	archivoAbierto.Close()
	var email Re_Email
	if len(strings.Split(lines[j], "---------------------- Forwarded by")) == 1 {
		email.Header = lines[j]
	} else {
		email.Header = lines[j] + lines[j+1]
		j++
	}
	j++
	for i := j; i < len(lines); i++ {
		nombresComoArreglo := strings.Split(lines[i], ":")
		switch strings.ToLower(nombresComoArreglo[0]) {
		case "from":
			if len(nombresComoArreglo) > 1 {
				email.From = nombresComoArreglo[1]
			}
		case "to":
			if len(nombresComoArreglo) > 1 {
				email.To = nombresComoArreglo[1]
			}
		case "sent":
			if len(nombresComoArreglo) > 1 {
				email.Sent = nombresComoArreglo[1]
			}
		case "cc":
			if len(nombresComoArreglo) > 1 {
				email.CC = nombresComoArreglo[1]
			}
		case "date":
			if len(nombresComoArreglo) > 1 {
				email.Date = nombresComoArreglo[1]
			}
		case "subject":
			if len(nombresComoArreglo) > 1 {
				if nombresComoArreglo[1] == "Re" || nombresComoArreglo[1] == "RE" {
					if len(nombresComoArreglo) > 2 {
						email.Subject = "RE: " + nombresComoArreglo[2]
					}
				} else if nombresComoArreglo[1] == "FW" || nombresComoArreglo[1] == "fw" {
					if len(nombresComoArreglo) > 2 {
						email.Subject = "FW: " + nombresComoArreglo[2]
					}
				} else {
					email.Subject = nombresComoArreglo[1]
				}
			}
		case " -----original message-----":
			email.re = addEmail(i, dir)

			i = len(lines)
		case "-----original message-----":
			email.re = addEmail(i, dir)

			i = len(lines)
		case "----- original message-----":
			email.re = addEmail(i, dir)

			i = len(lines)
		case "--------- inline attachment follows ---------":
			email.re = addEmail(i, dir)

			i = len(lines)
		default:

			if len(strings.Split(nombresComoArreglo[0], "---------------------- Forwarded by")) == 1 {
				email.Content = email.Content + `\n` + nombresComoArreglo[0]
			} else {
				i = len(lines)
			}

		}
	}

	return email
}

func Cmd(cmd string, shell bool) []byte {
	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			panic("some error found")
		}
		return out
	} else {
		out, err := exec.Command(cmd).Output()
		if err != nil {
			panic("some error found")
		}
		return out
	}
}
