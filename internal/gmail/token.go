package gmail

import (
	"net/http"

	"github.com/robsztal/FanMyMail/cmd/config"
	"github.com/rs/zerolog/log"
)

var (
	oauthStateStringGl = "state-token"
)

type Handlers struct {
	cfg config.Config
}

func NewHandlers(cfg config.Config) Handlers {
	return Handlers{cfg: cfg}
}

func (h *Handlers) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	h.HandleLogin(w, r, h.cfg.OAuthCfg, oauthStateStringGl)
}

func (h *Handlers) HandleCallBackFromGoogle(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msgf("callback")
	state := r.FormValue("state")
	if state != oauthStateStringGl {
		log.Debug().Msgf("invalid oauth state, expected " + oauthStateStringGl + ", got " + state + "\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")
	if code == "" {
		log.Debug().Msgf("code not found")
		if _, err := w.Write([]byte("Code Not Found to provide AccessToken..\n")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			if _, err := w.Write([]byte("User has denied Permission..")); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	} else {
		token, err := h.cfg.OAuthCfg.Exchange(r.Context(), code)
		if err != nil {
			log.Error().Msgf("oauthConfGl.Exchange() failed with " + err.Error() + "\n")
			return
		}
		tokenCh <- token
		log.Debug().Msgf("TOKEN>> AccessToken>> " + token.AccessToken)
		log.Debug().Msgf("TOKEN>> Expiration Time>> " + token.Expiry.String())
		log.Debug().Msgf("TOKEN>> RefreshToken>> " + token.RefreshToken)
		return
	}
}
