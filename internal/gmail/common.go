package gmail

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/robsztal/FanMyMail/internal/pages"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

func (h *Handlers) HandleLogin(w http.ResponseWriter, r *http.Request, oauthConf *oauth2.Config, oauthStateString string) {
	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("Parsing url failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Ctx(r.Context()).Info().Str("url", URL.String()).Msg("Parsed url")
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	URL.RawQuery = parameters.Encode()
	parsedURL := URL.String()
	log.Ctx(r.Context()).Debug().Str("url", parsedURL).Msg("parsed url")
	http.Redirect(w, r, parsedURL, http.StatusTemporaryRedirect)
}

func (h *Handlers) HandleMain(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(pages.IndexPage))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
