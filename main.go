package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jasonlvhit/gocron"
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

func printToDestinations(client *maps.Client, config *Config) {
	printTable(getDistanceMatrix(client, config.OriginAddresses, config.DestinationAddresses), config.OriginAddresses, config.DestinationAddresses, false)
}

func printFromDestinations(client *maps.Client, config *Config) {
	printTable(getDistanceMatrix(client, config.DestinationAddresses, config.OriginAddresses), config.OriginAddresses, config.DestinationAddresses, true)
}

func toMilitaryTime(s string) string {
	t, err := time.Parse("03:04pm", s)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	return t.Format("15:04")
}

func main() {
	// Load the configuration
	var config Config
	config.load()

	// Instantiate the Maps API client
	client, err := maps.NewClient(maps.WithAPIKey(config.ApiKey))
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	// Parse arguments
	nowPtr := flag.Bool("now", false, "Print the commute matrix now (instead of scheduling it)")
	flag.Parse()

	if *nowPtr {
		printToDestinations(client, &config)
		printFromDestinations(client, &config)
	} else {
		fmt.Println("Commute matrix scheduled...")
		for _, schedule := range config.Cron {
			for _, toTime := range schedule.ToTimes {
				gocron.Every(1).Day().At(toMilitaryTime(toTime)).Do(printToDestinations, client, &config)
			}

			for _, fromTime := range schedule.FromTimes {
				gocron.Every(1).Day().At(toMilitaryTime(fromTime)).Do(printFromDestinations, client, &config)
			}
		}

		<-gocron.Start()
	}
}
