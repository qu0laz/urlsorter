package urlsorter

import (
	"fmt"
	"net/url"
	"strings"
)

func SchemeHostSplit(x string) (string, string) {
	n, err := url.Parse(x)
	if err != nil {
		fmt.Println("error wirh scheme host parsing", err)
	}
	ns := n.Scheme
	nh := n.Host
	c := ns + "://" + nh
	return c, nh
}

func LinkCheck(links []string, scheme string, host string) ([]string, []string) {
	var URLS []string
	var urlsIQS []string
	for _, n := range links {
		if n != "" {
			y := NotAnAssetBool(n)
			if y == true {
				urlNoQ, urlQ := URLScrub(n, scheme, host)
				URLS = append(URLS, urlNoQ)
				urlsIQS = append(urlsIQS, urlQ)
			}
		}
	}
	uniqueURLs := UniqueSlice(URLS)
	uniqueURLQuerys := UniqueSlice(urlsIQS)
	return uniqueURLs, uniqueURLQuerys
}

func IntExtSort(seedURL string, links []string) ([]string, []string) {
	var intL []string
	var extL []string
	for _, x := range links {
		if strings.Contains(x, seedURL) {
			intL = append(intL, x)
			//fmt.Println("internal", x)
		} else {
			extL = append(extL, x)
			//fmt.Println("external", x)
		}
	}
	return intL, extL
}

func UniqueSlice(s []string) []string {
	Map := make(map[string]int)
	for _, val := range s {
		Map[val] = 1
		//	fmt.Println("FOUND", val)
	}
	Slice := make([]string, 0)
	for v := range Map {
		Slice = append(Slice, v)
	}
	return Slice
}

func NotAnAssetBool(url string) bool {
	c := []string{"/amp/", ";", "wp-json", "feed", "tel:", "mailto:", ":void"}
	s := []string{".oembed", ".css", ".js", ".jpeg", ".jpg", ".pdf", ".ico", ".svg", ".png", ".xml", ".woff2", ".ttf", ".otf", ".xlsx", ".csv", ".xls", ".zip", ".gif", ".psd", ".mp3", ".mp4", ".m4a", ".doc", ".docx", ".bak", ".gz", ".apk"}
	for _, i := range c {
		if strings.Contains(url, i) {
			return false
		}
	}
	for _, i := range s {
		if strings.Contains(url, i) {
			return false
		}
	}
	return true
}

func URLScrub(u string, scheme string, host string) (string, string) {
	var URL, urlIQS string
	//handles fringe cases
	u = RemoveSpaces(u)
	//e.g. bitly.com/about /n
	u = RemoveNewline(u)
	//e.g. .bitly.com/about
	if strings.HasPrefix(u, ".") {
		u = RemoveFirstChar(u)
	}

	//removes anything after an on page #id
	//e.g. bitly.com/about ->#contact
	u = RemoveHash(u)
	//Adds schem/host to home relative link
	//e.g. "/" ->https://bitly.com
	if u == "/" {
		c := scheme + "://" + host
		URL = c
		return URL, ""
	}

	//Catches All standard https://www config
	if strings.HasPrefix(u, "//") || strings.HasPrefix(u, "https://") || strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://www") || strings.HasPrefix(u, "http://www") {
		URL = RemoveParams(u)
		urlIQS = u
		//fmt.Println("<<<<FIRED-1",u)
		//catches /bitly.com/about
	} else if strings.HasPrefix(u, "/") {
		x := CheckURLAppendSchemeHostIfNeeded(scheme, host, u)
		URL = RemoveParams(x)
		urlIQS = x
		//fmt.Println("<<<<FIRED-2",u)
		// Catches links that have no slashes and other malformed hrefs
	} else {
		//fmt.Println("<<<<FIRED-3", u)
		c := fmt.Sprintf("%s://%s\n", scheme, host)
		if !strings.HasPrefix(u, "/") {
			c = c + "/" + u
		}
		URL = RemoveParams(c)
		urlIQS = c
	}
	//URL = ClipTrailingSlashes(URL)
	//urlIQS = ClipTrailingSlashes(urlIQS)

	return URL, urlIQS
}

func RemoveFirstChar(input string) string {
	if len(input) <= 1 {
		return ""
	}
	return input[1:]
}
func RemoveSpaces(i string) string {
	return strings.Replace(i, " ", "", -1)
}

func RemoveNewline(i string) string {
	return strings.Replace(i, "\n", "", -1)
}

func RemoveParams(u string) string {
	v, err := url.Parse(u)
	if err != nil {
		//fmt.Println("error on", u)
		return ""
	}
	v.RawQuery = ""
	return v.String()
}
func RemoveHash(u string) string {
	if strings.Contains(u, "#") {
		split := strings.Split(u, "#")
		s := split[0]
		return s
	} else {
		return u
	}
}
func CheckURLAppendSchemeHostIfNeeded(scheme string, host string, link string) string {
	//contains url
	v := new(url.URL)
	v.Scheme = scheme
	v.Host = host

	baseDomain := fmt.Sprintf("%s://%s\n", v.Scheme, v.Host)
	//If Link contains domain already OR contains :// andassumes there is already a domain somewhere in the string
	if strings.HasPrefix(link, baseDomain) || strings.Contains(link, "://") {
		//link is Clean
		return link
		//If its a dummy link and contains a new line char return nothing
	} else if strings.Contains(link, "\n") {
		return ""
		//If it's a normal relative link add the domain
	} else if strings.HasPrefix(link, "/") {
		return baseDomain + link
		//If its a link with no beginning "/" fringe case
	} else {
		return baseDomain + "/" + link
	}
}
