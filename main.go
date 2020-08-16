package QW50b25pLUN6YXBsaWNraQ_

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	http.ListenAndServe(":8080", r)
}
