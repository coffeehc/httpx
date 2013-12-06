package web

import "net/http"

func fileHandler(w http.ResponseWriter, req *http.Request) {
	entity := GetCacheEntity(req.RequestURI)
	if entity != nil {
		//....
	} else {
	}
}
