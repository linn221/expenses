package routing

import (
	"fmt"
	"net/http"
)

type myMux struct {
	*http.ServeMux
}

func (mux *myMux) HandleFunc(pattern string, f func(http.ResponseWriter, *http.Request)) {
	fmt.Println(pattern)
	mux.ServeMux.HandleFunc(pattern, f)
}
