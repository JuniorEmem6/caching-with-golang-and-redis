package main

import (
	"caching/main/data"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	data.Connect()
	data.RedisClient()
	http.HandleFunc("/", product)

	fmt.Println("Server started at port 3000")
	http.ListenAndServe(":3000", nil)
}

func product(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		var product data.Product

		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&product); err != nil && err != io.EOF {
			panic(err)
		}

		err := data.InsertProduct(&product)
		if err != nil {
			data, _ := json.Marshal(`{"status": "false", "message": "Erroring creating product"}`)
			res.Write(data)
		} else {
			res.Write([]byte(`{"status": "true", "message": "Product inserted successfully"}`))
		}

	} else if req.Method == "GET" {
		id := strings.TrimPrefix(req.URL.Path, "/")
		newId, _ := strconv.Atoi(id)
		result := data.GetCache(newId)
		fmt.Println(result)
		if result == "Cache miss" {
			name, price, description, err := data.SelectProduct(newId)

			if err != nil {
				res.Write([]byte(`{"status": "false", "message": "Product not found"}`))

			} else {
				// fmt.Println(name, description, price)
				data.SetCache(newId, name, price, description)
				res.Write([]byte(`{"status": "true", "message": "Product retrieved successfully from database"}`))
			}
		} else {
			res.Write([]byte(`{"status": "true", "message": "Product retrieved successfully from cache"}`))

		}

	}
}
