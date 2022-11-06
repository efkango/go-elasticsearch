package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	baseURL string
}

func NewClient(host string) *Client {
	return &Client{host}
}

func (c *Client) CheckHealth() error {
	response, err := http.Get(c.baseURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed check Elasticsearch health: %v", err)
	}

	log.Println("debug health check response: ", string(responseBody))

	return nil
}

func (c *Client) CreateIndex() error {
	body := `
	{
		"mappings": {
			"properties": {
				"id": {
					"type": "integer"
				},
				"name": {
					"type": "text"
				},
				"address": {
					"type": "text"
				},
				"salary": {
					"type": "float"
				}
			}
		}
	}
	`

	req, err := http.NewRequest("PUT", c.baseURL+"/employee", strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to make a create index request: %v", err)
	}

	httpClient := http.Client{}
	req.Header.Add("Content-type", "application/json")
	response, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make a http call to create an index: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed read create index response: %v", err)
	}

	log.Println("debug create index response: ", string(responseBody))

	return nil
}

func (c *Client) InsertData(e *Employee) error {
	body, _ := json.Marshal(e)

	id := strconv.Itoa(e.Id)
	req, err := http.NewRequest("PUT", c.baseURL+"/employee/_doc/"+id, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to make a insert data request: %v", err)
	}

	httpClient := http.Client{}
	req.Header.Add("Content-type", "application/json")
	response, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make a http call to insert data: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed read insert data response: %v", err)
	}

	log.Println("debug insert data response: ", string(responseBody))

	return nil
}

func (c *Client) SeedingData(idStart, n int) error {
	for i := idStart; i < n; i++ {
		if err := c.InsertData(&Employee{
			Id:      i,
			Name:    "person" + strconv.Itoa(i),
			Address: "address" + strconv.Itoa(i),
			Salary:  float64(i * 100),
		}); err != nil {
			return fmt.Errorf("failed seeding data with id %d: %v", i, err)
		}
	}

	return nil
}

func (c *Client) UpdateData(e *Employee) error {
	body, _ := json.Marshal(map[string]*Employee{
		"doc": e,
	})

	id := strconv.Itoa(e.Id)
	req, err := http.NewRequest("POST", c.baseURL+"/employee/_update/"+id, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to make a update data request: %v", err)
	}

	httpClient := http.Client{}
	req.Header.Add("Content-type", "application/json")
	response, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make a http call to update data: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed read update data response: %v", err)
	}

	log.Println("debug update data response: ", string(responseBody))

	return nil
}

func (c *Client) DeleteData(id int) error {

	req, err := http.NewRequest("DELETE", c.baseURL+"/employee/_doc/"+strconv.Itoa(id), nil)
	if err != nil {
		return fmt.Errorf("failed to make a delete data request: %v", err)
	}

	httpClient := http.Client{}
	req.Header.Add("Content-type", "application/json")
	response, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make a http call to delete data: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed read delete data response: %v", err)
	}

	log.Println("debug delete data response: ", string(responseBody))

	return nil
}

func (c *Client) SearchData(keyword string) ([]*Employee, error) {
	query := fmt.Sprintf(`
	{
		"query": {
			"match": {
				"name": "%s"
			}
		}
	}
	`, keyword)

	req, err := http.NewRequest("GET", c.baseURL+"/employee/_search", strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("failed to make a search data request: %v", err)
	}

	httpClient := http.Client{}
	req.Header.Add("Content-type", "application/json")
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make a http call to search data: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read insert data response: %v", err)
	}

	var searchHits SearchHits
	if err := json.Unmarshal(responseBody, &searchHits); err != nil {
		return nil, fmt.Errorf("failed read unmarshal data response: %v", err)
	}

	var employees []*Employee
	for _, hit := range searchHits.Hits.Hits {
		employees = append(employees, hit.Source)
	}

	return employees, nil
}
