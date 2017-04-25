package main

import (
	"text/template"
	"encoding/base64"
	"log"
	"bytes"
	"github.com/Masterminds/sprig"
)

type cloudConfigFile struct {
	Dockerfile string
	Domain 	string
}

func base64UserData(userDataFileContent string) (userdata string, err error) {
	if userDataFileContent != "" {
		userdata = base64.StdEncoding.EncodeToString([]byte(userDataFileContent))
	}
	return
}

func populateTemplates(ccf cloudConfigFile) (populatedFile string, err error) {
	t := template.Must(
		template.New("cloud-config.yml.tmpl").Funcs(sprig.TxtFuncMap()).ParseGlob(("cloud-config.yml.tmpl")))

	var buffer bytes.Buffer
	err = t.Execute(&buffer, ccf)
	if err != nil {
		log.Fatalf("template execution: %s", err)
		return
	}

	populatedFile = buffer.String()
	return
}

func GenerateCloudConfig(dockerfileContent string) (userData string, err error){
	ccf := cloudConfigFile{
		Dockerfile: dockerfileContent,
		Domain: "ghost.listo.com",
	}
	populatedFile, err := populateTemplates(ccf)
	if err != nil {
		log.Fatalf("Something went wrong calling populateTemplates", err)
		return
	}

	userData, err = base64UserData(populatedFile)
	if err != nil {
		log.Fatalf("Something went wrong calling base64UserData", err)
		return
	}

	return
}