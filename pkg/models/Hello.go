package models

type Hello struct {
	Name      string     `json:"name"`
	Port      string     `json:"port"`
	IP        string     `json:"ip"`
	Label     string     `json:"label"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	URL       string   `json:"url"`
	Protected bool     `json:"protected"`
	Methods   []string `json:"methods"`
}
