package discogs

import (
	"fmt"
	"net/http"

	"github.com/irlndts/go-apirequest"
)

const (
	discogsAPI = "https://api.discogs.com/"
)

type Options struct {
	Currency  string
	UserAgent string
	Token     string
}

// Client is a Discogs client for making Discogs API requests.
type Client struct {
	Release *ReleaseService
	Master  *MasterService
	Artist  *ArtistService
	Label   *LabelService
	Search  *SearchService
}

// NewClient returns a new Client.
func NewClient(o *Options) (*Client, error) {
	base := apirequest.New().Client(&http.Client{}).Base(discogsAPI)
	if o.UserAgent != "" {
		base.Set("User-Agent", o.UserAgent)
	}

	cur, err := currency(o.Currency)
	if err != nil {
		return nil, err
	}

	// set token, it's required for some queries like search
	if o.Token != "" {
		base.Set("Authorization", "Discogs token="+o.Token)
	}

	return &Client{
		Release: newReleaseService(base.New(), cur),
		Artist:  newArtistService(base.New()),
		Label:   newLabelService(base.New()),
		Master:  newMasterService(base.New()),
		Search:  newSearchService(base.New()),
	}, nil
}

// currency validates currency for marketplace data.
// Defaults to the authenticated users currency. Must be one of the following:
// USD GBP EUR CAD AUD JPY CHF MXN BRL NZD SEK ZAR
func currency(c string) (string, error) {
	switch c {
	case "USD", "GBP", "EUR", "CAD", "AUD", "JPY", "CHF", "MXN", "BRL", "NZD", "SEK", "ZAR":
		return c, nil
	default:
		return "", fmt.Errorf("%v\n", "Invalid currency abbreviation.")
	}
	return "USD", nil
}
