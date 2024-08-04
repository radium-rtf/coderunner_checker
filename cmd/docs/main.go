package main

import (
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	"net/http"

	"github.com/radium-rtf/coderunner_checker/internal/config"
)

func main() {
	cfg, err := config.LoadDocs()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(cfg)

	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), swaggerHandler); err != nil {
		log.Fatalln(err)
	}
}
