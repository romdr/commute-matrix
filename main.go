package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"googlemaps.github.io/maps"
	"gopkg.in/yaml.v2"
)

type Config struct {
	OriginAddresses      []string `yaml:"origins"`
	DestinationAddresses []string `yaml:"destinations"`
	Routes               struct {
		To   string `yaml:"to"`
		From string `yaml:"from"`
	}
	Cron []struct {
		DestinationAddress string   `yaml:"destination"`
		ToTimes            []string `yaml:"to_times"`
		FromTimes          []string `yaml:"from_times"`
	}
	ApiKey string `yaml:"apikey"`
}

func (config *Config) load() {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("ERROR: Reading config file: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Fatalf("ERROR: Parsing configuration: %s", err)
	}
}

func getDistanceMatrix(client *maps.Client, origins []string, destinations []string) *maps.DistanceMatrixResponse {
	request := &maps.DistanceMatrixRequest{
		Origins:       origins,
		Destinations:  destinations,
		DepartureTime: "now",
		Mode:          maps.TravelModeDriving,
		TrafficModel:  maps.TrafficModelBestGuess,
	}

	resp, err := client.DistanceMatrix(context.Background(), request)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	return resp
}

func printTable(distanceMatrixResponse *maps.DistanceMatrixResponse, origins, destinations []string, reverse bool) {
	title := "TO"
	if reverse {
		title = "FROM"
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetAutoWrapText(false)
	table.SetHeader(append([]string{fmt.Sprintf("%s %s", title, time.Now().Format("01/02/2006 15:04"))}, destinations...))

	for i, origin := range origins {
		tableRow := []string{origin}

		for j, _ := range destinations {

			from := i
			to := j
			if reverse {
				from = j
				to = i
			}

			tableRow = append(tableRow, fmt.Sprintf("%.1f", distanceMatrixResponse.Rows[from].Elements[to].DurationInTraffic.Minutes()))
		}

		table.Append(tableRow)
	}

	table.Render()
	fmt.Println("")
}

func main() {
	var config Config
	config.load()

	client, err := maps.NewClient(maps.WithAPIKey(config.ApiKey))
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	printTable(getDistanceMatrix(client, config.OriginAddresses, config.DestinationAddresses), config.OriginAddresses, config.DestinationAddresses, false)
	printTable(getDistanceMatrix(client, config.DestinationAddresses, config.OriginAddresses), config.OriginAddresses, config.DestinationAddresses, true)
}
