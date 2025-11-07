package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

var (
	store    *ChartStore
	storeMux sync.RWMutex
	dataFile = "charts.json"
)

func main() {
	// Initialize the store
	store = NewChartStore()
	if err := store.Load(dataFile); err != nil {
		log.Printf("Could not load data file: %v", err)
	}

	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/charts", getChartsHandler).Methods("GET")
	api.HandleFunc("/charts", createChartHandler).Methods("POST")
	api.HandleFunc("/charts/{id}", getChartHandler).Methods("GET")
	api.HandleFunc("/charts/{id}", updateChartHandler).Methods("PUT")
	api.HandleFunc("/charts/{id}", deleteChartHandler).Methods("DELETE")
	api.HandleFunc("/charts/{id}/export/svg", exportSVGHandler).Methods("GET")
	api.HandleFunc("/charts/{id}/export/png", exportPNGHandler).Methods("GET")
	api.HandleFunc("/charts/{id}/export/pdf", exportPDFHandler).Methods("GET")

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getChartsHandler(w http.ResponseWriter, r *http.Request) {
	storeMux.RLock()
	charts := store.GetAll()
	storeMux.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(charts)
}

func createChartHandler(w http.ResponseWriter, r *http.Request) {
	var chart Chart
	if err := json.NewDecoder(r.Body).Decode(&chart); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storeMux.Lock()
	store.Add(&chart)
	store.Save(dataFile)
	storeMux.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chart)
}

func getChartHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	storeMux.RLock()
	chart := store.Get(id)
	storeMux.RUnlock()

	if chart == nil {
		http.Error(w, "Chart not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chart)
}

func updateChartHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var chart Chart
	if err := json.NewDecoder(r.Body).Decode(&chart); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chart.ID = id
	storeMux.Lock()
	store.Update(&chart)
	store.Save(dataFile)
	storeMux.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chart)
}

func deleteChartHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	storeMux.Lock()
	store.Delete(id)
	store.Save(dataFile)
	storeMux.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func exportSVGHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	storeMux.RLock()
	chart := store.Get(id)
	storeMux.RUnlock()

	if chart == nil {
		http.Error(w, "Chart not found", http.StatusNotFound)
		return
	}

	svg, err := GenerateSVG(chart)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating SVG: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"chart-%s.svg\"", id))
	w.Write([]byte(svg))
}

func exportPNGHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	storeMux.RLock()
	chart := store.Get(id)
	storeMux.RUnlock()

	if chart == nil {
		http.Error(w, "Chart not found", http.StatusNotFound)
		return
	}

	pngData, err := GeneratePNG(chart)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating PNG: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"chart-%s.png\"", id))
	w.Write(pngData)
}

func exportPDFHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	storeMux.RLock()
	chart := store.Get(id)
	storeMux.RUnlock()

	if chart == nil {
		http.Error(w, "Chart not found", http.StatusNotFound)
		return
	}

	pdfData, err := GeneratePDF(chart)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating PDF: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"chart-%s.pdf\"", id))
	w.Write(pdfData)
}
