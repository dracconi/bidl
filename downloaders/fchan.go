package downloaders

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type fchanThread struct {
	Posts []struct {
		Id  int    `json:"no"`
		Tim int    `json:"tim",omiempty`
		Ext string `json:"ext",omiempty`
	} `json:"posts"`
}

func fchanThreadApi(m map[string]string) ([]RemoteImage, error) {
	board := m["board"]
	threadId := m["no"]
	resp, err := http.Get("https://a.4cdn.org/" + board + "/thread/" + threadId + ".json")
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)
	var thread fchanThread
	err = dec.Decode(&thread)
	if err != nil {
		return nil, err
	}

	var urls []RemoteImage

	for _, v := range thread.Posts {
		if v.Tim != 0 { // omiempty should leave 0 when there's no such field
			urls = append(urls, RemoteImage{
				Remote: "https://i.4cdn.org/" + board + "/" + strconv.Itoa(v.Tim) + v.Ext,
				Local:  board + "/" + threadId + "/" + strconv.Itoa(v.Tim) + v.Ext})
		}
	}
	return urls, nil
}
