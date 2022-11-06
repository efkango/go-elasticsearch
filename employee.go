package main

type Employee struct {
	Id      int     `json:"id"`
	Name    string  `json:"name,omitempty"`
	Address string  `json:"address,omitempty"`
	Salary  float64 `json:"salary,omitempty"`
}

type SearchHits struct {
	Hits struct {
		Hits []*struct {
			Source *Employee `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
