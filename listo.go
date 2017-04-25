package main
import (
	"net/http"
	"log"

)

var (
	address = ":8080"
	stacksmithApiKey = ""
	stacksmithStacksUrl = "https://stacksmith.bitnami.com/api/v1/stacks"

	node = Component {
		Id: "node",
		Version: "6.9.5",
	}
	nodeStack = Stack {
		Name:	"Node stack for listo",
		Components: []Component{node},
		Flavour: "node-base",
	}

	ssClient = newStacksmithClient(stacksmithApiKey, stacksmithStacksUrl)

)

func main() {
	router := NewRouter()
	log.Println("Listening on:", address)
	log.Fatal(http.ListenAndServe(address, router))
}

//REQUEST URL
//curl -X POST --header "Content-Type: application/json" --header "Accept: application/json" -d "{
//\"name\": \"Test Node Stack\",
//\"components\": [{
//\"id\": \"node\",
//\"version\": \"7.9.0\"
//}],
//\"flavor\": \"node-base\"
//}" "https://stacksmith.bitnami.com/api/v1/stacks?api_key=ff42ef80e10a3514f203d9bd3362b33ba45e9e92e63d3392067dbe51c69de5fd"

//RESPONSE
//{
//"id": "0yw0r0j",
//"stack_url": "https://stacksmith.bitnami.com/api/v1/stacks/u6befo4"
//}

//REQUEST stack_url
//RESPONSE
//get field dockerfile
//{
//"output": {"dockerfile": https://stacksmith.bitnami.com/api/v1/stacks/u6befo4.dockerfile}
//}

// REQUEST dockerfile