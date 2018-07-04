package main

import (
	"net/http"
	"fmt"
	"net/url"
	"strings"
	"io/ioutil"
	"encoding/json"
)

func main() {

	http.HandleFunc("/openlink", JDRender)

	http.ListenAndServe(":8989", nil);
}

func JDRender(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	jdUrl := r.Form.Get("url")
	_, err := url.ParseRequestURI(jdUrl)
	if err != nil || !strings.Contains(jdUrl, "jd.com") {
		errInfo := `{"errcode":-1,"message":"不支持的资源域名 Only => *.jd.com"}`
		fmt.Fprintf(w, errInfo)
		return
	} else {
		u := strings.Split(jdUrl, ":")
		jdRurl := url.QueryEscape(u[1])
		fmt.Print()
		reqUrl := "https://api.jd.com/routerjson?access_token=d485f611-7d1f-49d5-97a9-d5ac5450c659&app_key=8E76F960FE21B40A83B167F22223759C&method=jingdong.wxsq.mjgj.link.GetOpenLink&v=2.0&sign=35D5AA4F312FD08D8070B951105A8D9C&timestamp=2018-07-03%2014%3A54%3A45"
		data := url.Values{"jump": {"0"}, "rurl": {url.QueryEscape("http://dc2.jd.com/auto.php?service=transfer&type=pms&to=" + jdRurl + "&openlink=1")}}
		body := strings.NewReader(data.Encode())
		clt := http.Client{}
		resp, err := clt.Post(reqUrl, "application/x-www-form-urlencoded", body)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		content, err := ioutil.ReadAll(resp.Body)
		respBody := string(content)
		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(respBody), &dat); err == nil {
			r, _ := json.Marshal(dat["jingdong_wxsq_mjgj_link_GetOpenLink_responce"])
			if err := json.Unmarshal(r, &dat); err == nil {
				r, _ := json.Marshal(dat["open_link_result"])
				fmt.Fprintf(w, string(r))
			}
		}
	}
	fmt.Print("\r\n访问者ip：" + r.Header.Get("X-Real-Ip") + strings.Split(r.RemoteAddr, ":")[0])

}
