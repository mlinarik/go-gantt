package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
)

// Chart represents a Gantt chart
type Chart struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	StartYear  int        `json:"startYear"`
	StartQ     int        `json:"startQuarter"`
	EndYear    int        `json:"endYear"`
	EndQ       int        `json:"endQuarter"`
	Categories []Category `json:"categories"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// Category represents a grouping of tasks
type Category struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Tasks []Task `json:"tasks"`
}

// Task represents a single item in the Gantt chart
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartYear   int    `json:"startYear"`
	StartQ      int    `json:"startQuarter"`
	EndYear     int    `json:"endYear"`
	EndQ        int    `json:"endQuarter"`
	Color       string `json:"color,omitempty"`
}

// ChartStore manages the collection of charts
type ChartStore struct {
	charts map[string]*Chart
}

// NewChartStore creates a new chart store
func NewChartStore() *ChartStore {
	return &ChartStore{
		charts: make(map[string]*Chart),
	}
}

// Add adds a new chart to the store
func (s *ChartStore) Add(chart *Chart) {
	if chart.ID == "" {
		chart.ID = uuid.New().String()
	}
	now := time.Now()
	chart.CreatedAt = now
	chart.UpdatedAt = now

	// Ensure all categories and tasks have IDs
	for i := range chart.Categories {
		if chart.Categories[i].ID == "" {
			chart.Categories[i].ID = uuid.New().String()
		}
		for j := range chart.Categories[i].Tasks {
			if chart.Categories[i].Tasks[j].ID == "" {
				chart.Categories[i].Tasks[j].ID = uuid.New().String()
			}
		}
	}

	s.charts[chart.ID] = chart
}

// Get retrieves a chart by ID
func (s *ChartStore) Get(id string) *Chart {
	return s.charts[id]
}

// GetAll returns all charts
func (s *ChartStore) GetAll() []*Chart {
	charts := make([]*Chart, 0, len(s.charts))
	for _, chart := range s.charts {
		charts = append(charts, chart)
	}
	return charts
}

// Update updates an existing chart
func (s *ChartStore) Update(chart *Chart) {
	if existing := s.charts[chart.ID]; existing != nil {
		chart.CreatedAt = existing.CreatedAt
	}
	chart.UpdatedAt = time.Now()

	// Ensure all categories and tasks have IDs
	for i := range chart.Categories {
		if chart.Categories[i].ID == "" {
			chart.Categories[i].ID = uuid.New().String()
		}
		for j := range chart.Categories[i].Tasks {
			if chart.Categories[i].Tasks[j].ID == "" {
				chart.Categories[i].Tasks[j].ID = uuid.New().String()
			}
		}
	}

	s.charts[chart.ID] = chart
}

// Delete removes a chart
func (s *ChartStore) Delete(id string) {
	delete(s.charts, id)
}

// Save persists the charts to a file
func (s *ChartStore) Save(filename string) error {
	data, err := json.MarshalIndent(s.charts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// Load reads charts from a file
func (s *ChartStore) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, not an error
		}
		return err
	}

	return json.Unmarshal(data, &s.charts)
}
