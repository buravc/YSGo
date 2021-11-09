package godisApi

import (
	godis "YSGo/godis/server"
	godisApi "YSGo/godisApi/models"
	"YSGo/repository"
	"encoding/json"
	"net/http"
)

type GodisApiServer struct {
	godisRepository repository.Repository
	HttpServer      *http.Server
}

func NewApiServer() *GodisApiServer {
	godis := godis.Singleton()
	defaultRepository := &repository.DefaultRepository{GodisServer: godis}
	godisRepository := repository.Repository(defaultRepository)
	handler := http.NewServeMux()
	httpServer := &http.Server{Addr: "0.0.0.0:8090", Handler: handler}
	apiServer := &GodisApiServer{godisRepository: godisRepository, HttpServer: httpServer}

	handler.HandleFunc("/get", apiServer.get)
	handler.HandleFunc("/set", apiServer.set)
	handler.HandleFunc("/flush", apiServer.flush)

	httpServer.ListenAndServe()

	return apiServer
}

func (apiServer *GodisApiServer) get(w http.ResponseWriter, req *http.Request) {

	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := req.URL.Query()
	key := query.Get("key")

	if key == "" {
		response(w, godisApi.GetResult{Success: false, Value: ""})
		return
	}

	val, ok := apiServer.godisRepository.Get(key)

	response(w, godisApi.GetResult{Success: ok, Value: val})
}

func (apiServer *GodisApiServer) set(w http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(req.Body)

	var setReq godisApi.SetRequest
	err := decoder.Decode(&setReq)

	if err != nil {
		response(w, godisApi.SetResult{Success: false})
		return
	}

	apiServer.godisRepository.Set(setReq.Key, setReq.Value)

	response(w, godisApi.SetResult{Success: true})
}

func (apiServer *GodisApiServer) flush(w http.ResponseWriter, req *http.Request) {

	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	apiServer.godisRepository.Flush()

	response(w, godisApi.FlushResult{Success: true})
}

func response(w http.ResponseWriter, v interface{}) {

	jsn, err := json.Marshal(v)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		jsn, err = json.Marshal(godisApi.InternalError{Error: "An internal error has occurred!"})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsn)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsn)
}
