package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"

	"github.com/radium-rtf/coderunner_checker/internal/config"
	"github.com/swaggo/swag"
)

func init() {
	docsPATH := os.Getenv("DOCS_CONFIG_PATH")
	file, err := os.ReadFile(docsPATH)
	if err != nil {
		panic(err)
	}

	cfg, err := config.LoadDocs()
	if err != nil {
		panic(err)
	}

	template := string(file)

	m := make(map[string]any)
	err = json.NewDecoder(strings.NewReader(template)).Decode(&m)
	if err != nil {
		panic(err)
	}
	m["basePath"] = "/"
	m["host"] = cfg.HttpENDPOINT

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(m)
	if err != nil {
		panic(err)
	}

	swagger := &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  buf.String(),
	}

	swag.Register(swagger.InstanceName(), swagger)
}
