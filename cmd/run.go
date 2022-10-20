package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/showwin/speedtest-go/speedtest"
	"github.com/spf13/cobra"
)

var host, username, password, client string
var saver bool

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a speedtest",
	Run: func(cmd *cobra.Command, args []string) {
		if host == "" || username == "" || password == "" || client == "" {
			log.Fatal("results-host, username, password & client are required flags")
		}

		user, err := speedtest.FetchUserInfo()
		if err != nil {
			fmt.Println("Warning: Cannot fetch user information. http://www.speedtest.net/speedtest-config.php is temporarily unavailable.")
		}

		servers, err := speedtest.FetchServers(user)
		if err != nil {
			log.Fatal("Cannot fetch servers", err)
		}

		targets, err := servers.FindServer(nil)
		if err != nil {
			log.Fatal("Cannot find server", err)
		}

		if len(targets) < 1 {
			log.Fatal("No server found")
		}

		server := targets[0]

		err = server.PingTest()
		if err != nil {
			log.Fatal("Ping test failed", err)
		}

		err = server.DownloadTest(saver)
		if err != nil {
			log.Fatalf("Download test failed: %s", err)
		}
		err = server.UploadTest(saver)
		if err != nil {
			log.Fatalf("Upload test failed: %s", err)
		}

		payload := map[string]interface{}{
			"id":      server.ID,
			"name":    server.Name,
			"country": server.Country,
			"sponsor": server.Sponsor,
			"lat":     server.Lat,
			"lon":     server.Lon,

			"dl_speed": server.DLSpeed,
			"ul_speed": server.ULSpeed,
			"latency":  server.Latency,

			"client": client,
		}

		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			log.Fatal("Cannot marshal result data to JSON", err)
		}

		client := http.Client{}
		req, err := http.NewRequest("POST", host, bytes.NewBuffer(jsonBytes))
		if err != nil {
			log.Fatal("Cannot create request", err)
		}

		req.SetBasicAuth(username, password)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Cannot send request", err)
		}

		if resp.StatusCode != 200 {
			log.Fatal("Request failed", resp.StatusCode)
		}
	},
}

func init() {
	runCmd.Flags().StringVar(&host, "results-host", "", "where the results are posted")
	runCmd.Flags().StringVar(&username, "username", "", "username for results host")
	runCmd.Flags().StringVar(&password, "password", "", "password for results host")
	runCmd.Flags().StringVar(&client, "client", "", "name of the client")
	runCmd.Flags().BoolVar(&saver, "saver-mode", false, "run the low data test")

	rootCmd.AddCommand(runCmd)
}
