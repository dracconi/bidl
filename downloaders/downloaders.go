package downloaders

import (
	"errors"
	"net/url"
	"regexp"
)

func zipMatches(headers []string, matched []string) map[string]string {
	zipped := make(map[string]string)
	for i, v := range headers {
		if v != "" {
			zipped[v] = matched[i]
		}
	}
	return zipped
}

type rule struct {
	host    string
	re_path *regexp.Regexp
	handler func(map[string]string) ([]RemoteImage, error)
}

var rules []rule

func getNamedReData(str string, re *regexp.Regexp) map[string]string {
	return zipMatches(re.SubexpNames(), re.FindStringSubmatch(str))
}

func addRule(host string, re_path string, handler func(map[string]string) ([]RemoteImage, error)) {
	rules = append(rules, rule{
		host:    host,
		re_path: regexp.MustCompile(re_path),
		handler: handler})
}

func InitRules() {
	addRule("boards.4channel.org",
		`^/(?P<board>\w+)/thread/(?P<no>\d+)`,
		fchanThreadApi)
}

type RemoteImage struct {
	Remote string
	Local  string
}

func GetImUrls(url *url.URL) ([]RemoteImage, error) {
	for _, v := range rules {
		if v.host == url.Host {
			return v.handler(getNamedReData(url.Path, v.re_path))
		}
	}
	return nil, errors.New("Host didn't match to any rules.")
}
