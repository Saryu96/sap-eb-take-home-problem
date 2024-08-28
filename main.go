package main

import (
	"encoding/csv"
	"html/template"
	"net/http"
	"os"
)

// Trail represents a single trail with various attributes
type Trail struct {
	Name       string
	AccessType string
	Restrooms  bool
	Picnic     bool
	Fishing    bool
	Address    string
	Fee        bool
	BikeRack   bool
	Difficulty string
}

// mapDifficulty converts the trail class code to a human-readable difficulty level
func mapDifficulty(class string) string {
	switch class {
	case "T1":
		return "Easy"
	case "T2":
		return "Moderate"
	case "T3":
		return "Difficult"
	case "T4":
		return "Most Difficult"
	default:
		return "No"
	}
}

// loadTrailsData reads trail data from a CSV file and maps it to a slice of Trail structs
func loadTrailsData(filename string) ([]Trail, error) {
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var trails []Trail
	// Iterate over the records, starting from the second row
	for _, record := range records[1:] {
		// Map each record to a Trail struct
		trails = append(trails, Trail{
			Name:       record[29],
			AccessType: record[6],
			Restrooms:  record[1] == "Yes",
			Picnic:     record[2] == "Yes",
			Fishing:    record[3] == "Yes",
			Address:    record[8],
			Fee:        record[9] == "Yes",
			BikeRack:   record[10] == "Yes",
			Difficulty: mapDifficulty(record[7]),
		})
	}

	return trails, nil
}

// filterTrails filters the trail data based on user input criteria
func filterTrails(trails []Trail, address, difficulty string, restrooms, picnic, fishing, fee, bikeRack bool) []Trail {
	var filtered []Trail
	for _, trail := range trails {
		// Apply the filters based on user input
		if (address == "" || trail.Address == address) &&
			(difficulty == "" || trail.Difficulty == difficulty) &&
			(!restrooms || trail.Restrooms) &&
			(!picnic || trail.Picnic) &&
			(!fishing || trail.Fishing) &&
			(!fee || trail.Fee) &&
			(!bikeRack || trail.BikeRack) {
			filtered = append(filtered, trail)
		}
	}
	return filtered
}

// handleTrails processes the HTTP request
func handleTrails(w http.ResponseWriter, r *http.Request) {
	// Parse user input from the request
	r.ParseForm()
	address := r.FormValue("address")
	difficulty := r.FormValue("difficulty")
	restrooms := r.FormValue("restrooms") == "on"
	picnic := r.FormValue("picnic") == "on"
	fishing := r.FormValue("fishing") == "on"
	fee := r.FormValue("fee") == "on"
	bikeRack := r.FormValue("bikerack") == "on"

	// Load trail data from the CSV file
	trails, err := loadTrailsData("BoulderTrailHeads.csv")
	if err != nil {
		http.Error(w, "Error loading trails data", http.StatusInternalServerError)
		return
	}

	// Filter the trails based on user input
	filteredTrails := filterTrails(trails, address, difficulty, restrooms, picnic, fishing, fee, bikeRack)

	// Render the filtered trails using the HTML template
	tmpl := template.Must(template.ParseFiles("trails.html"))
	tmpl.Execute(w, filteredTrails)
}

// main initializes the HTTP server and routes
func main() {
	// Register the /trails route
	http.HandleFunc("/trails", handleTrails)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the HTTP server on port 8080
	http.ListenAndServe(":8080", nil)
}
