package main

import (
	"os"
	"os/exec"
	"encoding/base64"
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "gcloud tag plugin"
	app.Usage = "gcloud tag plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "source_tag",
			Usage:  "source tag",
			EnvVar: "PLUGIN_SOURCE_TAG",
		},				
		cli.StringFlag{
			Name:   "dest_tag",
			Usage:  "destination tag",
			EnvVar: "PLUGIN_DEST_TAG",
		},		
		cli.StringSliceFlag{
			Name:     "repositories",
			Usage:    "repositories to tag",
			EnvVar:   "PLUGIN_REPOSITORIES",
		},
		cli.StringFlag{
			Name:   "service_key",
			Usage:  "GCP service key",
			EnvVar: "PLUGIN_SERVICE_KEY",
		},	
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func decodeServiceKey(encodedKey string) []byte {
	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)

	if err != nil {
		logrus.Error(err)
	}			

	logrus.Info("service key decoded")	

    return decodedKey
}

func writeServiceKey(decodedKey []byte) {
    f, err := os.Create("service_key.json")

    if err != nil {
        logrus.Fatal(err)
    }

    defer f.Close()

    _, err2 := f.Write(decodedKey)

    if err2 != nil {
        logrus.Fatal(err2)
    }

	logrus.Info("service key written")	
}

func authenticate(decodedKey []byte) {
	var serviceKey map[string]interface{}

	json.Unmarshal([]byte(decodedKey), &serviceKey)

	var clientEmail string = serviceKey["client_email"].(string)

	var projectId string = serviceKey["project_id"].(string)

	_, err := exec.Command("gcloud", "auth", "activate-service-account", clientEmail, "--key-file=service_key.json", "--project=" + projectId).Output()	

	if err != nil {
		logrus.Error(err)
	}		
	
	logrus.Info("authenticated " + clientEmail)	
}

func run(c *cli.Context) {

	decodedKey := decodeServiceKey(c.String("service_key"))

	writeServiceKey(decodedKey)

	authenticate(decodedKey)

	for _, repo := range c.StringSlice("repositories") {
		var src string = repo + ":" + c.String("source_tag")

		var dest string = repo + ":" + c.String("dest_tag")

		_, err := exec.Command("gcloud", "container", "images", "add-tag", src, dest).Output()
		
		logrus.Info(c)

		if err != nil {
			logrus.Error(err)
		}		

		logrus.Info(src + " >> " + dest)	
	}	
}