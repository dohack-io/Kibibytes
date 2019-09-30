package googleMapsUrlGenerator

//https://developers.google.com/maps/documentation/urls/guide#directions-action
import (
	"Kibibytes/utils"
	"fmt"
	"net/url"
)

func ToHome(id int64) string {
	var query string
	user := utils.GetUser(id)

	if user.Location != "" {
		destination := url.QueryEscape(user.Location)
		if user.Travelmode == "" {
			query = fmt.Sprintf("%s&destination=%s", "https://www.google.com/maps/dir/?api=1", destination)
		} else {
			query = fmt.Sprintf("%s&destination=%s&travelmode=%s", "https://www.google.com/maps/dir/?api=1", destination, user.Travelmode)
		}
		return query
	} else {
		return "Sorry we can't find your home"
	}

}

func FromTo(origin string, destination string, id int64) string {
	var query string

	user := utils.GetUser(id)
	origin = url.QueryEscape(origin)
	destination = url.QueryEscape(destination)

	if user.Travelmode == "" {
		query = fmt.Sprintf("%s&origin=%s&destination=%s", "https://www.google.com/maps/dir/?api=1", origin, destination)
	} else {
		query = fmt.Sprintf("%s&origin=%s&destination=%s&travelmode=%s", "https://www.google.com/maps/dir/?api=1", origin, destination, user.Travelmode)
	}
	return query
}

func Find(query string) string {
	var link string

	query = url.QueryEscape(query)
	link = fmt.Sprintf("%s&query=%s", "https://www.google.com/maps/search/?api=1", query)

	return link
}
