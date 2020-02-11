package main

import(
    "fmt"
    "github.com/gorilla/mux"
    "encoding/json"
    "net/http"
    "os"
)

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/v1/search", callback)
    router.HandleFunc("/v1/search/{city}", callback)
    router.HandleFunc("/v1/search/{city}/{district}", callback)
    http.Handle("/", headerMiddleWare(router))
    err := http.ListenAndServe(":8001", nil)
    if err != nil {
        fmt.Fprintf(os.Stderr, "http error", err)
    }
}

func callback(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    var output string
    switch len(vars) {
    case 0:
        data := GetData()
        output = MakeJSONString(data)
    case 1:
        data := Filter(GetData(), vars["city"])
        output = MakeJSONString(data)
    case 2:
        data := Filter(Filter(GetData(), vars["city"]), vars["district"])
        output = MakeJSONString(data)
    }

    w.Write([]byte(output))
}

func MakeJSONString(data []Profile) string {
    jsonByte, err := json.Marshal(data)
    if err != nil {
        fmt.Println("MakeJSONString error: ", err)
    }
    jsonString := string(jsonByte)
    return jsonString
}

func headerMiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w,r)
	})
}


