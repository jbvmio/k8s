package k8s

import (
	"github.com/tidwall/gjson"
)

// findSelfLinks parses a json response for the specified named resource and returns any selfLinks if found
func findSelfLinks(raw []byte, name string) []string {
	search := string(`items.#[metadata.name%"*` + name + `*"]#.metadata.selfLink`)
	results := gjson.GetManyBytes(raw, search)
	var selfLinks []string
	for _, url := range results {
		url.ForEach(func(k, v gjson.Result) bool {
			selfLinks = append(selfLinks, v.Raw)
			return true
		})
	}
	return selfLinks
}

// returnSelfLinks parses a json response and returns any selfLinks
func returnSelfLinks(raw []byte) []string {
	search := string(`items.#.metadata.selfLink`)
	results := gjson.GetManyBytes(raw, search)
	var selfLinks []string
	for _, url := range results {
		url.ForEach(func(k, v gjson.Result) bool {
			selfLinks = append(selfLinks, v.Raw)
			return true
		})
	}
	return selfLinks
}

// parseFor parses a json response for the specified named resource and returns the corresponding json for the resource
func parseFor(raw []byte, name string) []string {
	search := string(`items.#[metadata.name%"*` + name + `*"]#`)
	results := gjson.GetManyBytes(raw, search)
	var json []string
	for _, url := range results {
		url.ForEach(func(k, v gjson.Result) bool {
			json = append(json, v.Raw)
			return true
		})
	}
	return json
}

func keepAsBytes(json []byte, result gjson.Result) []byte {
	var raw []byte
	if result.Index > 0 {
		raw = json[result.Index : result.Index+len(result.Raw)]
	} else {
		raw = []byte(result.Raw)
	}
	return raw
}