package middleware

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rayark/zin"
)

// HMACSHA1Signer returns a middleware wrapper to add hmac signing string in
// response header
func HMACSHA1Signer(hmacHeaderKey, nounceHeaderKey string, secret []byte) zin.Middleware {
	return func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			deferWriter := NewDeferWriter(w)
			defer deferWriter.WriteAll()

			key := secret
			if nounceHeaderKey != "" {
				nounceInHex := r.Header.Get(nounceHeaderKey)
				nounce, err := hex.DecodeString(nounceInHex)
				if err == nil {
					key = append(key, nounce...)
				}
			}

			h(deferWriter, r, p)
			hmacSignature := generateSignature(deferWriter.Bytes(), key)
			deferWriter.Header().Set(hmacHeaderKey, hmacSignature)
		}
	}
}

func generateSignature(msg, key []byte) string {
	h := hmac.New(sha1.New, key)
	h.Write(msg)
	return "sha1=" + hex.EncodeToString(h.Sum(nil))
}
