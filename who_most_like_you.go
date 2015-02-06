package who_most_like_you

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func init() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/locki", Locki)
	http.HandleFunc("/", root)
}

type FQLPost struct {
	Name string `json:"name"`
	Set  []struct {
		PostID string `json:"post_id"`

		UserID string `json:"user_id"`

		Name string `json:"name"`
		Sex  string `json:"sex"`
		UID  string `json:"uid"`
	} `json:"fql_result_set"`
}

func Locki(w http.ResponseWriter, r *http.Request) {

	query := `{"post_ids":"SELECT post_id FROM stream WHERE source_id=me() AND likes.count>0 LIMIT 5000",` +
		`"uids":"SELECT user_id FROM like WHERE post_id IN (SELECT post_id FROM #post_ids)",` +
		`"like_ids":"SELECT name,sex,uid FROM user WHERE uid IN (SELECT user_id FROM #uids)"}`
	accessToken := "CAACEdEose0cBAM1n9ZA41phlVgj0YcIQKv2nBfC0VHd4IUcthZBJc6ZCSlUjmWNY288FC7lPC0n05ZBZC1ImIbg5J5xQyFlLat2He6Lu4I5xwnrMRKAoZBNm8kUTLaa8mPTQgQAvAZBPBhh9CZBdQrUtTTd8iT0NGWOfZAU0ziNK5np2AkQqUtRI08BDF3pihCjwWD28LQFZB0KMBnZA2aasawd"
	v := url.Values{
		"access_token": []string{accessToken},
		"q":            []string{query},
	}
	req, err := http.NewRequest("GET", "https://graph.facebook.com/fql?"+v.Encode(), nil)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	if resp.StatusCode != 200 {
		w.Write(body)
		return
	}
	jresp := struct {
		Data []FQLPost `json:"data"`
	}{}
	if err := json.Unmarshal(body, &jresp); err != nil {
		log.Printf("body: %s, err: %v", body, err)
		return
	}
	log.Printf("%+v", jresp)

	rr := make(map[string]int)
	for _, p := range jresp.Data[1].Set {
		rr[p.UserID] = rr[p.UserID] + 1
	}

}

var rootTmpl = template.Must(template.ParseFiles("templates/index.html"))

func root(w http.ResponseWriter, r *http.Request) {
	page := struct {
		IsLogin bool
	}{
		IsLogin: true,
	}
	rootTmpl.Execute(w, page)
}
