package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"log"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/api",
		Index,
	},
	Route{
		"CreateApp",
		"POST",
		"/api/apps",
		CreateApp,
	},
	Route{
		"AppIndex",
		"GET",
		"/api/apps",
		AppIndex,
	},
	Route{
		"GetApp",
		"GET",
		"/api/apps/{instanceId}",
		GetApp,
	},
	Route{
		"DeleteApp",
		"DELETE",
		"/api/apps/{instanceId}",
		DeleteApp,
	},
}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc

		router.
		Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	log.Println(w, "Welcome!\n")
}

func AppIndex(w http.ResponseWriter, r *http.Request) {
	res, err := getInstances("")
	if err != nil {
		log.Println("Unable to get instances", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	type awsInstance struct {
		Id 	 *string	`json:"id"`
		PublicIp *string	`json:"publicIp"`
	}
	var instances []awsInstance

	for _, res := range res.Reservations {
		for _, i := range res.Instances {
			awsInstance := awsInstance {
				PublicIp: i.PublicIpAddress,
				Id: i.InstanceId,
			}
			instances = append(instances, awsInstance)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(instances); err != nil {
		panic(err)
	}
}

func GetApp(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	res, err := getInstances(params["instanceId"])
	if err != nil {
		log.Println("Unable to get instances", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	type awsInstance struct {
		Id 	 *string	`json:"id"`
		PublicIp *string	`json:"publicIp"`
	}
	var instances []awsInstance

	for _, res := range res.Reservations {
		for _, i := range res.Instances {
			awsInstance := awsInstance {
				PublicIp: i.PublicIpAddress,
				Id: i.InstanceId,
			}
			instances = append(instances, awsInstance)
		}
	}

	if len(instances) < 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(instances[0]); err != nil {
		panic(err)
	}
}

func CreateApp(w http.ResponseWriter, r *http.Request) {

	createStackBody, err := ssClient.createStack(nodeStack)
	if err != nil {
		log.Println("Unable to createStack", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	log.Println("Successfully create stack")

	var userdata string
	if dockerfile, err := ssClient.getDockerfile(createStackBody.Id); err != nil {
		log.Println("Unable to get Dockerfile", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	} else {
		log.Println("Successfully retrieved Dockerfile", dockerfile)
		userdata, err = GenerateCloudConfig(dockerfile)
		if err != nil {
			log.Println("Unable to generate CloudConfig", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}


	if res, err := createInstance(userdata); err != nil {
		log.Println("Unable to createStack", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	} else {
		log.Println("Successfully created instance")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(res.Instances[0]); err != nil {
			panic(err)
		}
	}
}

func DeleteApp(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, err := deleteInstance(params["instanceId"])
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		log.Println("Unable to delete Instance", err)
	}
	log.Println("Successfully deleted instance")
	w.WriteHeader(http.StatusNoContent)
}