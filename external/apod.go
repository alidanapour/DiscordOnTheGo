// This package is for RESTful API calls to NASA's Astronomy Photo of the Day
// api.

package external

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	term "example.com/terminal"
	"github.com/bwmarrin/discordgo"
)

type Apod struct {
	Date        string `json:"date"`
	Explanation string `json:"explanation"`
	Title       string `json:"title"`
	ImgURL      string `json:"url"`
	Copyright   string `json:"copyright"`
}

const dateLayout = "2006-01-02"

var APOD Apod
var currDate time.Time

func init() {
	currDate = time.Now()
	APOD.Date = "2000-01-01"
	//getApod() //TODO: uncomment when needed for actual run
}

// This function does the API call and updates the Astronomy Photo of the Day data
// into APOD. Credit: https://www.soberkoder.com/consume-rest-api-go/
func getApod() error {

	if APOD.Date != currDate.Format(dateLayout) {
		term.Print(term.HTTPGET, "https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY")
		resp, err := http.Get("https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY")
		if err != nil {
			term.Print(term.ERROR, "APOD Error: GET failed")
			return err
		}

		// Convert response into []byte slice and store into APOD struct
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &APOD)
	}
	return nil
}

// Main function. This returns a discordgo MessageEmbed filled with APOD values, error otherwise
func ApodRequest() (*discordgo.MessageEmbed, error) {
	if getApod() == nil {
		embed := &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       0x003366,
			Title:       APOD.Title,
			Description: APOD.Explanation,
			Image: &discordgo.MessageEmbedImage{
				URL: APOD.ImgURL,
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Copyright: " + APOD.Copyright,
			},
		}
		return embed, nil
	}
	return nil, errors.New("Unable to get APOD, most likely reached api call limit.")
}
