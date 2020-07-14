package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
)

// ESLogger function for ElastciSearch log handler
func ESLogger(c echo.Context) error {
	body := echo.Map{}
	if err := c.Bind(&body); err != nil {
		return err
	}

	indexPattern := os.Getenv("INDEX_PATTERN")

	if indexPattern == "" {
		indexPattern = "kong-2006-01-02"
	}

	currentTime := time.Now().Format(indexPattern)

	esHost := os.Getenv("ES_HOST")
	if esHost == "" {
		esHost = "127.0.0.1"
	}

	esPort := os.Getenv("ES_PORT")
	if esPort == "" {
		esPort = "9200"
	}

	tr := &http.Transport{}
	client := &http.Client{Transport: tr}

	reqBody, err := json.Marshal(body)

	// Call the api
	resp, err := client.Post(
		fmt.Sprintf("http://%s:%s/%s/_doc", esHost, esPort, currentTime),
		"application/json",
		bytes.NewBuffer(reqBody),
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	// res := map[string]string{"currentTime": currentTime, "bar": "ok", "es_host": esHost, "es_port": esPort}
	return c.JSON(http.StatusOK, result)
}
