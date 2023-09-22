package twilio

import (
	"goVoice/internal/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go/client"
)

func GetValidator (cfg *config.Config) client.RequestValidator {
	return client.NewRequestValidator(cfg.TwilioAuthToken)
}

func RequireValidTwilioSignature(validator *client.RequestValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := "https://some--digits.ngrok.io" + c.Request.URL.Path // FIXME temp path
		signatureHeader := c.Request.Header.Get("X-Twilio-Signature")
		params := make(map[string]string)
		c.Request.ParseForm()
		for k, v := range c.Request.PostForm {
			params[k] = v[0]
		}

		if !validator.Validate(url, params, signatureHeader) {
			log.Printf("Incomming request not originating from Twilio. URL: %s, Params: %v, Signature: %s", url, params, signatureHeader)
			c.Copy().AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}