package godisApi

import (
	godis "YSGo/godis/server"
	godisApi "YSGo/godisApi/models"
	"YSGo/repository"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type GodisApiServer struct {
	godisRepository repository.Repository
	HttpServer      *http.Server
}

func NewApiServer() *GodisApiServer {
	log.Println("Starting GodisApiServer...")

	godis := godis.Singleton()
	defaultRepository := &repository.DefaultRepository{GodisServer: godis}
	godisRepository := repository.Repository(defaultRepository)

	port := os.Getenv("GodisApiPort")

	if port == "" {
		port = "8090"
		log.Printf("'GodisApiPort' env variable was not set. Falling back to default: %s\n", port)
	}

	num, err := strconv.ParseInt(port, 0, 16)

	if (num < 1024 || num > 99999) || err != nil {
		log.Fatal("Failed to parse 'GodisApiPort' env variable! Port number must be between 1024 and 99999.")
	}

	endpoint := fmt.Sprintf("0.0.0.0:%s", port)

	handler := http.NewServeMux()
	httpServer := &http.Server{Addr: endpoint, Handler: handler}
	apiServer := &GodisApiServer{godisRepository: godisRepository, HttpServer: httpServer}

	handler.HandleFunc("/get", apiServer.get)
	handler.HandleFunc("/set", apiServer.set)
	handler.HandleFunc("/flush", apiServer.flush)

	httpServer.ListenAndServe()

	log.Println("GodisApiServer started successfully!")

	return apiServer
}

func (apiServer *GodisApiServer) get(w http.ResponseWriter, req *http.Request) {

	reqUUID := pseudo_uuid(req.RemoteAddr)

	logRequestReceived(reqUUID, "/get", req)

	if req.Method != "GET" {
		log.Printf("Logged by %s : Http method is not GET. Current request method: %s\n", reqUUID, req.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := req.URL.Query()
	key := query.Get("key")

	if key == "" {
		log.Printf("Logged by %s : Key parameter is empty.\n", reqUUID)
		response(w, godisApi.GetResult{Success: false, Value: ""}, reqUUID)
		return
	}

	val, ok := apiServer.godisRepository.Get(key)

	log.Printf("Logged by %s : Request successful. {Key: \"%s\", Value: \"%s\", Ok: %t}", reqUUID, key, val, ok)
	response(w, godisApi.GetResult{Success: ok, Value: val}, reqUUID)
}

func (apiServer *GodisApiServer) set(w http.ResponseWriter, req *http.Request) {

	reqUUID := pseudo_uuid(req.RemoteAddr)

	logRequestReceived(reqUUID, "/set", req)

	if req.Method != "POST" {
		log.Printf("Logged by %s : Http method is not POST. Current request method: %s\n", reqUUID, req.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		response(w, godisApi.SetResult{Success: false}, reqUUID)
		return
	}

	log.Printf("Logged by %s : Request Body: %s\n", reqUUID, string(body))

	var setReq godisApi.SetRequest
	err = json.Unmarshal(body, &setReq)

	if err != nil {
		log.Printf("Logged by %s : Could not convert body to json. Body: %s\n", reqUUID, string(body))
		response(w, godisApi.SetResult{Success: false}, reqUUID)
		return
	}

	apiServer.godisRepository.Set(setReq.Key, setReq.Value)

	log.Printf("Logged by %s : Request successful.", reqUUID)

	response(w, godisApi.SetResult{Success: true}, reqUUID)
}

func (apiServer *GodisApiServer) flush(w http.ResponseWriter, req *http.Request) {

	reqUUID := pseudo_uuid(req.RemoteAddr)

	logRequestReceived(reqUUID, "/flush", req)

	if req.Method != "GET" {
		log.Printf("Logged by %s : Http method is not GET. Current request method: %s\n", reqUUID, req.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	apiServer.godisRepository.Flush()

	log.Printf("%s : Request successful.", reqUUID)
	response(w, godisApi.FlushResult{Success: true}, reqUUID)
}

func response(w http.ResponseWriter, v interface{}, reqUUID string) {

	w.Header().Set("Content-Type", "application/json")

	jsn, err := json.Marshal(v)

	if err != nil {
		log.Printf("Logged by %s : An error occurred while converting result to json. Responding with internal error!\n", reqUUID)
		jsn, err = json.Marshal(godisApi.InternalError{Error: "An internal error has occurred!"})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsn)
		return
	}

	log.Printf("Logged by %s : Responding with: %s\n", reqUUID, jsn)
	w.WriteHeader(http.StatusOK)
	w.Write(jsn)
}

func pseudo_uuid(ipAddr string) (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X--%s", b[0:4], b[4:6], b[6:8], b[8:10], b[10:], ipAddr)

	return
}

func logRequestReceived(reqUUID, endpoint string, req *http.Request) {
	reqRemoteAddr := req.RemoteAddr
	reqUrl := req.URL
	reqMethod := req.Method
	reqAccept := req.Header.Get("accept")
	reqUserAgent := req.UserAgent()

	log.Printf("Request received. Endpoint: %s\nRequest ID: %s\nRemote Address: %s\nUrl: %s\nMethod: %s\nAccept: %s\nUser Agent: %s\n", endpoint, reqUUID, reqRemoteAddr, reqUrl, reqMethod, reqAccept, reqUserAgent)
}
