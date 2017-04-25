package main

import (
	"bytes"
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"errors"
	"io/ioutil"
)

type stacksmithClient struct {
	apiKey		string
	stacksUrl	string
}

type Stack struct {
	Name 		string 		`json:"name"`
	Components 	[]Component 	`json:"components"`
	Flavour 	string		`json:"flavour"`
}

type Component struct {
	Id 	string	`json:"id"`
	Version string	`json:"version"`
}

type Output struct {
	Dockerfile	string	`json:"dockerfile"`
}

type Body struct {
	Id 		string	`json:"id"`
	Stack_url  	string	`json:"stack_url"`
}

func newStacksmithClient(apiKey string, stacksUrl string) (c *stacksmithClient) {
	c = &stacksmithClient {
		apiKey: 	apiKey,
		stacksUrl: 	stacksUrl,
	}
	return c
}

func (c *stacksmithClient) createStack(s Stack) (body *Body, err error) {
	stackJsonBuffer := new(bytes.Buffer)
	json.NewEncoder(stackJsonBuffer).Encode(s)

	res, err := http.Post(c.stacksUrl, "application/json; charset=utf-8", stackJsonBuffer)
	if err != nil {
		log.Println("Unable to create a new http request", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		errorBody, _ := ioutil.ReadAll(res.Body)
		log.Println("Wrong status code")
		return nil, errors.New(string(errorBody))
	}

	json.NewDecoder(res.Body).Decode(&body)
	return body, nil
}

func (c *stacksmithClient) getStack(id string) (output *Output, err error) {

	stackUrl := c.stacksUrl + "/" + id
	res, err := c.httpRequest(stackUrl)
	if err != nil {
		log.Println("Something went wrong trying to do the http request", err)
		return nil, err
	}
	var body struct {
		// httpbin.org sends back key/value pairs, no map[string][]string
		Id 		string	`json:"id"`
		Output  	Output	`json:"output"`
	}

	json.NewDecoder(res.Body).Decode(&body)
	return &body.Output, nil
}

func (c * stacksmithClient) httpRequest(url string) (res *http.Response, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Add("api_key", c.apiKey)
	req.URL.RawQuery = q.Encode()
	res, err = client.Do(req)

	log.Println("Url requested", req.URL.String())
	if err != nil {
		log.Println("Unable to create a new http request", err)
		return nil, err
	}

	if res.StatusCode != 200 {
		defer res.Body.Close()
		errorBody, _ := ioutil.ReadAll(res.Body)
		log.Println("Wrong status code")
		return nil, errors.New(string(errorBody))
	}

	return res, nil
}

func (c *stacksmithClient) getDockerfile(id string) (body string, err error) {

	stackOutput, err := c.getStack(id);
	if err != nil {
		log.Println("Something went wrong calling getDockerfile", err)
		return
	}

	dockerfileUrl := stackOutput.Dockerfile
	res, err := c.httpRequest(dockerfileUrl)
	if err != nil {
		log.Println("Something went wrong calling httpRequest", err)
		return "", err
	}
	defer res.Body.Close()
	dockerfileContent, _ := ioutil.ReadAll(res.Body)
	fmt.Printf(string(body))

	return string(dockerfileContent), nil
}