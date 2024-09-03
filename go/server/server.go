package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/filter", filter)
	router.GET("/status", status)

	router.Run()
}

type Entry struct {
	Name     string  `json:"name"`
	Language string  `json:"language"`
	Id       string  `json:"id"`
	Bio      string  `json:"bio"`
	Version  float64 `json:"version"`
}

func status(context *gin.Context) {
	response := []byte("Ok")
	context.Data(200, "plain/text", response)
}

func filter(context *gin.Context) {
	context.Writer.Header().Set("content-type", "application/json")

	if err := filterJsonArray(context.Request.Body, context.Writer, filterEntryByVersion); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func filterEntryByVersion(entry Entry) bool {
	return entry.Version > 5
}

type Filter[T any] func(item T) bool

func filterJsonArray[T any](reader io.Reader, writer io.Writer, filter Filter[T]) error {
	decoder := json.NewDecoder(reader)

	if _, err := writer.Write([]byte("[")); err != nil {
		return err
	}

	// Get rid of the first opening bracket: "["
	if _, err := decoder.Token(); err != nil {
		return err
	}

	firstItem := true

	for decoder.More() {
		var item T

		if err := decoder.Decode(&item); err != nil {
			return err
		}

		if filter(item) {
			itemJson, err := json.Marshal(item)
			if err != nil {
				return err
			}

			if !firstItem {
				if _, err := writer.Write([]byte(",")); err != nil {
					return err
				}
			} else {
				firstItem = false
			}

			if _, err := writer.Write(itemJson); err != nil {
				return err
			}
		}
	}

	if _, err := writer.Write([]byte("]")); err != nil {
		return err
	}

	return nil
}
