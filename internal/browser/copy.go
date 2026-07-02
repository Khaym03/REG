package browser

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func CopyCookies(src *rod.Browser, dst *rod.Browser) error {
	currentCookies, err := src.GetCookies()
	if err != nil {
		return err
	}

	var params []*proto.NetworkCookieParam
	for _, cookie := range currentCookies {
		sourcePortCopy := cookie.SourcePort

		param := &proto.NetworkCookieParam{
			Name:         cookie.Name,
			Value:        cookie.Value,
			Domain:       cookie.Domain,
			Path:         cookie.Path,
			Secure:       cookie.Secure,
			HTTPOnly:     cookie.HTTPOnly,
			SameSite:     cookie.SameSite,
			Expires:      cookie.Expires,
			Priority:     cookie.Priority,
			SameParty:    cookie.SameParty,
			SourceScheme: cookie.SourceScheme,
			SourcePort:   &sourcePortCopy,
			PartitionKey: cookie.PartitionKey,
		}

		params = append(params, param)
	}

	return dst.SetCookies(params)
}
