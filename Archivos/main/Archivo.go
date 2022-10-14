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

type Email struct {
	Directorio                string
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
}

func IngresarEmails() {
	usuarios, err := ioutil.ReadDir("../../enron_mail_20110402/maildir")
	if err != nil {
		log.Fatal(err)
	}
	for _, usuario := range usuarios {
		directorios, error := ioutil.ReadDir("../../enron_mail_20110402/maildir/" + usuario.Name())
		if error != nil {
			log.Fatal(error)
		}
		var emails []Email
		for _, directorio := range directorios {
			if directorio.IsDir() {
				emails = append(emails, ObtenerArchivos(("../../enron_mail_20110402/maildir/" + usuario.Name() + "/" + directorio.Name()))...)
			} else {
				emails = append(emails, Emails(("../../enron_mail_20110402/maildir/" + usuario.Name() + "/" + directorio.Name())))
			}
		}
		res := map[string]interface{}{"index": "emails", "records": emails}
		resJSON, err := json.Marshal(res)
		file, err := os.Create("ejemplo.json")
		if err != nil {
			panic(err)
		}
		_, err = file.WriteString(string(resJSON))
		if err != nil {
			panic(err)
		}
		file.Close()
		cmd := "curl http://localhost:4080/api/_bulkv2 -i -u admin:0208Mavl  --data-binary '@ejemplo.json'"
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
			fmt.Println(out)
		}

	}

	fmt.Println("termino")
}
func Emails(dir string) Email {

	var email Email

	archivoAbierto, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
	}
	email.Directorio = dir
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
		var prefijo string = strings.ToLower(nombresComoArreglo[0])
		if prefijo == "message-id" {
			email.Message_ID = nombresComoArreglo[1]
		} else if prefijo == "date" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 {
				email.Date = nombresComoArreglo[1]
				for j := 2; j < len(nombresComoArreglo); j++ {
					email.Date = email.Date + ":" + nombresComoArreglo[j]
				}
			}
		} else if prefijo == "from" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 {
				email.From = nombresComoArreglo[1]
			}
		} else if prefijo == "to" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 && email.X_FileName == "" {
				email.To = nombresComoArreglo[1]
				var val bool = true
				if email.X_FileName != "" {
					val = false
				}
				for val {
					if len(lines) > (i + 1) {
						comodin := strings.ToLower(strings.Split(lines[i+1], ":")[0])
						if comodin != "subject" {
							if comodin != "cc" {
								c := strings.Split(lines[(i+1)], ":")
								for j := 0; j < len(c); j++ {
									email.To = email.To + ":" + c[j]
								}
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
		} else if prefijo == "subject" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 && email.X_FileName == "" {
				email.Subject = nombresComoArreglo[1]
				for j := 2; j < len(nombresComoArreglo); j++ {
					email.Subject = email.Subject + ":" + nombresComoArreglo[j]
				}
				var val bool = true
				if email.X_FileName != "" {
					val = false
				}
				for val {
					if len(lines) > (i + 1) {
						comodin := strings.ToLower(strings.Split(lines[i+1], ":")[0])
						if comodin != "cc" {
							if comodin != "mime-version" {
								c := strings.Split(lines[(i+1)], ":")
								for j := 0; j < len(c); j++ {
									email.Subject = email.Subject + ":" + c[j]
								}
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
		} else if prefijo == "cc" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 && email.X_FileName == "" {
				email.Cc = nombresComoArreglo[1]
				var val bool = true
				if email.X_FileName != "" {
					val = false
				}
				for val {

					if len(lines) > (i + 1) {
						comodin := strings.ToLower(strings.Split(lines[i+1], ":")[0])
						if comodin != "mime-version" {
							if comodin != "content-type" {
								c := strings.Split(lines[(i+1)], ":")
								for j := 0; j < len(c); j++ {
									email.Cc = email.Cc + ":" + c[j]
								}
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
		} else if prefijo == "mime-version" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 {
				email.Mime_Version = nombresComoArreglo[1]
			}
		} else if prefijo == "content-type" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 {
				email.Content_Type = nombresComoArreglo[1]
			}
		} else if prefijo == "content-transfer-encoding" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 {
				email.Content_Transfer_Encoding = nombresComoArreglo[1]
			}
		} else if prefijo == "bcc" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 {
				email.Bcc = nombresComoArreglo[1]
				var val bool = true
				if email.X_FileName != "" {
					val = false
				}
				for val {
					if len(lines) > (i + 1) {
						comodin := strings.ToLower(strings.Split(lines[i+1], ":")[0])
						if comodin != "x-from" {
							if comodin != "x-bcc" {
								c := strings.Split(lines[(i+1)], ":")
								for j := 0; j < len(c); j++ {
									email.Bcc = email.Bcc + ":" + c[j]
								}
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
		} else if prefijo == "x-from" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 {
				email.X_From = nombresComoArreglo[1]
			}
		} else if prefijo == "x-to" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 && email.X_FileName == "" {
				email.X_To = nombresComoArreglo[1]
				var val bool = true
				if email.X_FileName != "" {
					val = false
				}
				for val {
					if len(lines) > (i + 1) {
						comodin := strings.ToLower(strings.Split(lines[i+1], ":")[0])
						if comodin != "x-cc" {
							if comodin != "x-bcc" {
								c := strings.Split(lines[(i+1)], ":")
								for j := 0; j < len(c); j++ {
									email.X_To = email.X_To + ":" + c[j]
								}
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
		} else if prefijo == "x-cc" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 && email.X_FileName == "" {
				email.X_cc = nombresComoArreglo[1]
				var val bool = true
				if email.X_FileName != "" {
					val = false
				}
				for val {
					if len(lines) > (i + 1) {
						comodin := strings.ToLower(strings.Split(lines[i+1], ":")[0])
						if comodin != "x-bcc" {
							if comodin != "x-folder" {
								c := strings.Split(lines[(i+1)], ":")
								for j := 0; j < len(c); j++ {
									email.X_cc = email.X_cc + ":" + c[j]
								}
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
		} else if prefijo == "x-bcc" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 {
				email.X_bcc = nombresComoArreglo[1]
				var val bool = true
				if email.X_FileName != "" {
					val = false
				}
				for val {
					if len(lines) > (i + 1) {
						comodin := strings.ToLower(strings.Split(lines[i+1], ":")[0])
						if comodin != "x-folder" {
							if comodin != "x-origin" {
								c := strings.Split(lines[(i+1)], ":")
								for j := 0; j < len(c); j++ {
									email.X_bcc = email.X_bcc + ":" + c[j]
								}
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
		} else if prefijo == "x-folder" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 {
				email.X_Folder = nombresComoArreglo[1]
			}
		} else if prefijo == "x-origin" && email.X_FileName == "" {
			if len(nombresComoArreglo) > 1 {
				email.X_Origin = nombresComoArreglo[1]
			}
		} else if prefijo == "x-filename" {
			if len(nombresComoArreglo) > 1 {
				email.X_FileName = nombresComoArreglo[1]
			}
		} else {
			for j := i; j < len(lines); j++ {
				email.Content = email.Content + lines[j] + `\n`
			}
			i = len(lines)
		}
		i++
	}
	return email
}
func ObtenerArchivos(dir string) []Email {

	var emails []Email
	archivos, error := ioutil.ReadDir(dir)
	for _, archivo := range archivos {
		if error != nil {
			log.Fatal(error)
		}
		if archivo.IsDir() {
			dir := dir + "/" + archivo.Name()
			emails = append(emails, ObtenerArchivos((dir))...)

		} else {
			emails = append(emails, Emails(dir+"/"+archivo.Name()))
		}
	}
	return emails
}
