package redirect

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/DevJHansen/tinyurl.com.na/pkg/config"
	"github.com/go-resty/resty/v2"
	"github.com/ua-parser/uap-go/uaparser"
)

type IPInfoResponse struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}

func ProcessUrl(u string) (string, error) {
	var link = u
	validString := regexp.MustCompile(`^[^\s]+\.[^\s]+$`)

	if !validString.MatchString(link) {
		return "", errors.New("invalid link provided")
	}

	if strings.HasPrefix(link, "www.") {
		link = strings.Replace(link, "www.", "https://", 1)
	}

	if strings.HasPrefix(link, "http://") {
		link = strings.Replace(link, "http://", "https://", 1)
	}

	if !strings.HasPrefix(link, "https://") {
		link = "https://" + link
	}

	parsedURL, err := url.Parse(link)
	if err != nil || parsedURL.Host == "" {
		return "", err
	}

	validHostname := regexp.MustCompile(`^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !validHostname.MatchString(parsedURL.Host) {
		return "", errors.New("invalid hostname")
	}

	return link, nil
}

func GenerateUID() (string, error) {
	bytes := make([]byte, 6)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	uid := base64.URLEncoding.EncodeToString(bytes)

	uid = strings.TrimRight(uid, "=")
	return uid, nil
}

func IsBot(r *http.Request) bool {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		// If User-Agent is missing, it's suspicious
		return true
	}

	userAgent = strings.ToLower(userAgent)

	botSignatures := []string{
		"bot",
		"crawl",
		"slurp",
		"spider",
		"mediapartners",
		"baiduspider",
		"bingbot",
		"duckduckbot",
		"googlebot",
		"yandexbot",
		"facebot",
		"ia_archiver",
		"python-urllib",
		"headless",
		"wget",
		"curl",
		"java",
		"phantomjs",
		"selenium",
		"headlesschrome",
		"puppeteer",
		"jsdom",
		"nodejs",
	}

	for _, signature := range botSignatures {
		if strings.Contains(userAgent, signature) {
			return true
		}
	}

	return false
}

func LogRedirectAnalytics(r *http.Request, key string, app *config.Config) {
	parser := uaparser.NewFromSaved()
	client := parser.Parse(r.UserAgent())
	flagged := IsBot(r)
	ip := getIPFromRequest(r)
	country := getCountryFromIP(ip, (*app.Envs)["IP_INFO_API_KEY"])

	analyticsData := map[string]interface{}{
		"target_key":  key,
		"user_agent":  r.UserAgent(),
		"referrer":    r.Referer(),
		"ip_address":  ip,
		"device_type": client.Device.Family,
		"os":          client.Os.Family,
		"flagged":     flagged,
		"country":     country,
	}

	_, _, err := app.Supabase.From("redirect_analytics").Insert(analyticsData, false, "", "", "").Execute()

	if err != nil {
		log.Printf("Error logging analytics for key %s: %v", key, err)
	}
}

func getIPFromRequest(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	return ip
}

func getCountryFromIP(ip string, apiKey string) string {
	client := resty.New()

	var result IPInfoResponse
	resp, err := client.R().
		SetQueryParam("ip", ip).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey)).
		SetResult(&result).
		Get("https://ipinfo.io")

	if err != nil || resp.IsError() {
		log.Println("Failed to get country from IP:", err)
		return ""
	}

	return result.Country
}
