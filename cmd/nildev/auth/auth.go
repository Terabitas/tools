package auth

import (
	"net/http"
	"sync"

	"io"

	"github.com/nildev/lib/log"
	"github.com/skratchdot/open-golang/open"
)

var wg sync.WaitGroup
var accessToken AccessToken

var ret = `
<html>
	<head>
		<script type="text/javascript">
		//<![CDATA[
			function parse_url(url) {
				var match = url.match(/^(http|https|ftp)?(?:[\:\/]*)([a-z0-9\.-]*)(?:\:([0-9]+))?(\/[^?#]*)?(?:\?([^#]*))?(?:#(.*))?$/i);
				var ret   = new Object();

				ret['protocol'] = '';
				ret['host']     = match[2];
				ret['port']     = '';
				ret['path']     = '';
				ret['query']    = '';
				ret['fragment'] = '';

				if(match[1]){
					ret['protocol'] = match[1];
				}

				if(match[3]){
					ret['port']     = match[3];
				}

				if(match[4]){
					ret['path']     = match[4];
				}

				if(match[5]){
					ret['query']    = match[5];
				}

				if(match[6]){
					ret['fragment'] = match[6];
				}

				return ret;
			}

			var url_parts = parse_url(location.href);

			var protocol  = url_parts['protocol'];
			var host      = url_parts['host'];
			var port      = url_parts['port'];
			var path      = url_parts['path'];
			var query     = url_parts['query'];
			var fragment  = url_parts['fragment'];

			var xmlhttp=new XMLHttpRequest();
			xmlhttp.open("GET", 'http://127.0.0.1:7777/store?'+fragment);
			xmlhttp.onreadystatechange = function() {}
			xmlhttp.send({});

			window.close()

		//]]>
		</script>
	<head>
	<body>
	</body>
</html>
`

func processOAuth(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, ret)
}

func persistToken(w http.ResponseWriter, r *http.Request) {
	//	fmt.Printf("Got token [%s] \n", r.URL.Query().Get("access_token"))
	//	fmt.Printf("Got scopes [%s] \n", r.URL.Query().Get("scopes"))
	//	fmt.Printf("Got token_type [%s] \n", r.URL.Query().Get("token_type"))
	accessToken = AccessToken(r.URL.Query().Get("access_token"))
	wg.Done()
}

// Auth user against selected provider and return AccessToken
func Auth(provider string) AccessToken {
	port := "7777"
	go func() {
		log.Infof("Starting HTTP server on 0.0.0.0:%s...", port)
		http.HandleFunc("/oauth2", processOAuth)
		http.HandleFunc("/store", persistToken)

		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.WithField("error", err).Fatal("Unable to create listener")
		}
	}()

	if provider == "bitbucket.org" {
		log.Infof("Please authenticate yourself within [%s] provider in opened browser window ...", provider)
		wg.Add(1)
		open.Start("https://bitbucket.org/site/oauth2/authorize?client_id=FUGTa554rDNASAzdBg&response_type=token")
	}
	wg.Wait()

	return accessToken
}
