package subsonic

import (
	"encoding/xml"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mdlayher/wavepipe/api"
	"github.com/mdlayher/wavepipe/data"

	"github.com/gorilla/context"
	"github.com/mdlayher/goset"
	"github.com/unrolled/render"
)

// IndexesContainer represents a Subsonic indexes container
type IndexesContainer struct {
	XMLName xml.Name `xml:"indexes,omitempty"`

	LastModified int64   `xml:"lastModified,attr"`
	Indexes      []Index `xml:"index"`
}

// Index represents an alphabetical Subsonic index
type Index struct {
	XMLName xml.Name `xml:"index"`

	Name string `xml:"name,attr"`

	// This is not actually an artist, but that's the way Subsonic treats it
	Artists []Artist `xml:"artist"`
}

// GetIndexes is used in Subsonic to return an alphabetical index of folders and IDs
func GetIndexes(res http.ResponseWriter, req *http.Request) {
	// Retrieve render
	r := context.Get(req, api.CtxRender).(*render.Render)

	// Create a new response container, build indexes container
	c := newContainer()
	c.Indexes = &IndexesContainer{
		// TODO: replace with actual last scan time
		LastModified: time.Now().Unix(),
	}

	// Fetch list of all folders, ordered alphabetically
	folders, err := data.DB.AllFoldersByTitle()
	if err != nil {
		log.Println(err)
		r.XML(res, 200, ErrGeneric)
		return
	}

	// Use a set to track indexes which already exist
	indexSet := set.New()

	// Iterate all folders and begin building indexes
	indexes := make([]Index, 0)
	for i, f := range folders {
		// Skip folder with ID 1, the root media folder
		if f.ID == 1 {
			continue
		}

		// Get the initial character of the folder title
		char := string(f.Title[0])

		// Create the index if it doesn't already exist
		if indexSet.Add(char) {
			indexes = append(indexes, Index{Name: char})
		}

		// Add this folder to the index at the current position
		indexes[i].Artists = append(indexes[i].Artists, Artist{
			Name: f.Title,
			// Since Subsonic and wavepipe have different data models, we get around
			// the ID restriction by adding a prefix describing what this actually is
			ID: "folder_" + strconv.Itoa(f.ID),
		})
	}

	// Add indexes, write response
	c.Indexes.Indexes = indexes
	r.XML(res, 200, c)
}
