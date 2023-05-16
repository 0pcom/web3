package main

import (
	"encoding/base64"
	"fmt"
	cc "github.com/ivanpirog/coloredcobra"
	"strings"
	"github.com/bitfield/script"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"github.com/skycoin/skywire-utilities/pkg/cipher"
)

const name = "web3"

const shcmd = `/usr/bin/bash -c`

var (
	// vars that contain generated html pages
	mainhtml          *string
	// html snippets
	htmlstart   string = "<!doctype html><html lang=en><head></head><body style='background-color:black;color:white;'>\n<style type='text/css'>\npre {\n  font-family:Courier New;\n  font-size:10pt;\n}\n.af_line {\n  color: gray;\n  text-decoration: none;\n}\n.column {\n  float: left;\n  width: 30%;\n  padding: 10px;\n}\n.row:after {\n  content: '';\n  display: table;\n  clear: both;\n}\n</style>\n<pre>"
	n0          string = "<a id='top' class='anchor' aria-hidden='true' href='#top'></a>"
	n1          string = "  <a href='/'>home</a>"
	n13         string = "\n<br>\n"
	navlinks    string = n1 + n13
	htmltoplink string = "<a href='#top'>top of page</a>\n"
	htmlend     string = "</pre></body></html>"
	htmlstyle   string = "<style>\npre {\n  font-family:Courier New;\n  font-size:10pt;\n}\n.af_line {\n  color: gray;\n  text-decoration: none;\n}\n.column {\n  float: left;\n  width: 30%;\n  padding: 10px;\n}\n.row:after {\n  content: '';\n  display: table;\n  clear: both;\n}\n</style>\n"
	bodystyle   string = "<body style='background-color:black;color:white;'>\n"
)


func main() {
	Execute()
}

func init() {
	//	rootCmd.AddCommand(
		//		myCmd,
		//	)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	defaultport, err := strconv.Atoi(os.Getenv("WEBPORT"))
	if err != nil {
		defaultport = 8080
	}
	rootCmd.Flags().IntVarP(&webPort, "port", "p", defaultport, "port to serve on - env WEBPORT="+os.Getenv("WEBPORT"))
}

var rootCmd = &cobra.Command{
	Use:   name,
	Short: "web application template",
	Long: `
	┬ ┬┌─┐┌┐   ┌┬┐┬ ┬┬─┐┌─┐┌─┐
	│││├┤ ├┴┐   │ ├─┤├┬┘├┤ ├┤
	└┴┘└─┘└─┘   ┴ ┴ ┴┴└─└─┘└─┘`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf(scriptExec(fmt.Sprintf(`%s "lsof -ti tcp:%s | xargs kill -9"`, shcmd, strconv.Itoa(webPort))))
		Server()
	},
}

func Execute() {
	cc.Init(&cc.Config{
		RootCmd:       rootCmd,
		Headings:      cc.HiBlue + cc.Bold, //+ cc.Underline,
		Commands:      cc.HiBlue + cc.Bold,
		CmdShortDescr: cc.HiBlue,
		Example:       cc.HiBlue + cc.Italic,
		ExecName:      cc.HiBlue + cc.Bold,
		Flags:         cc.HiBlue + cc.Bold,
		//FlagsDataType: cc.HiBlue,
		FlagsDescr:      cc.HiBlue,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Failed to execute command: ", err)
	}
}

var (
	webPort int
)

func mainhtmlfunc() {
	l := "<!doctype html><html lang=en><head><title>Web3</title></head><body style='background-color:black;color:white;'>\n<style type='text/css'>\npre {\n  font-family:Courier New;\n  font-size:10pt;\n}\n.af_line {\n  color: gray;\n  text-decoration: none;\n}\n.column {\n  float: left;\n  width: 30%;\n  padding: 10px;\n}\n.row:after {\n  content: '';\n  display: table;\n  clear: both;\n}\n</style>\n<pre>"
	l += navlinks
	l += fmt.Sprintf("\n%s\n",
		sex("source "+name+".sh ;  _dayscalc "))
	l += fmt.Sprintf("%s",
		sex("source "+name+".sh ;  _rainbowcal "))
	l += "</body></html>"
	mainhtml = &l
}

func sex(cmd string) string {
	return scriptExecHTML(fmt.Sprintf(`%s "%s"`, shcmd, cmd))
}

func scriptExecHTML(cmd string) string {
	res, err := script.Exec(cmd).String()
	if err != nil {
		res += fmt.Sprintf("<br><p style='color:red'>error during script.Exec:\n<br> %v\n<br></p>command:\n<br>\n%s\n<br>\n%s", err, cmd, res)
	}
	return res
}
func scriptExec(cmd string) string {
	fmt.Printf("executing command: \n %s", cmd)
	res, err := script.Exec(cmd).String()
	if err != nil {
		res = fmt.Sprintf("error during script.Exec:\n %v\nCommand:\n%s\nResult:\n%s\n", err, cmd, res)
	}
	return res
}

func Server() {
//	fmt.Println("generating in-memory html")
//	go mainhtmlfunc()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	r1 := gin.Default()
	r1.GET("/", func(c *gin.Context) {
  c.Writer.Header().Set("Server", "")
	c.Writer.WriteHeader(http.StatusOK)
	mainhtmlfunc()
	c.Writer.Write([]byte(*mainhtml))
		return
	})

	r1.GET("/:pkport/:file", func(c *gin.Context) {
		var publicKey string
		var dmsgPort string
		file := c.Param("file")		// Validate Skywire public key
		result := strings.SplitN(c.Param("pkport"), ":", 2)
		if len(result) == 2 {
			publicKey = result[0]
			dmsgPort = result[1]
			} else {
				errorMessage := fmt.Sprintf("Invalid Skywire <public-key>:<port>")
				c.Writer.Header().Set("Server", "")
				c.Writer.WriteHeader(http.StatusBadRequest)
				c.Writer.Write([]byte(errorMessage))
				return
		}
		pk := cipher.PubKey{}
		err := pk.Set(publicKey)
		if err != nil {
			errorMessage := fmt.Sprintf("Invalid Skywire public key: %s", err.Error())
			c.Writer.Header().Set("Server", "")
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.Writer.Write([]byte(errorMessage))
			return
		}
		// Handle valid public key
		cmd := fmt.Sprintf(`bash -c 'dmsgget -n dmsg://%s:%s/%s'`, pk, dmsgPort, file)
		res, err := script.Exec(cmd).String()
		if err != nil {
			errorMessage := fmt.Sprintf("error during script.Exec:\n%v\ncommand:\n%s\nresult:\n%s", err, cmd, res)
			c.Writer.Header().Set("Server", "")
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.Writer.Write([]byte(errorMessage))
			return
		}
		c.Writer.Header().Set("Server", "")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte(res))
	})



	faviconBase64 := `AAABAAEAICAAAAEAIACoEAAAFgAAACgAAAAgAAAAQAAAAAEAIAAAAAAAABAAACIuAAAiLgAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAC4BIwAAAAoBnVRgGb5mdEjQcH9d2HeCXdl5gV7Mcnhoy3J3a9B2emjc
gnxg34d7XsJ1bUGiZGAq0YVyXOGQdV7ikXVd5JJ3XtSKcVy6fWtKmGdmNIBYXhP///8AAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACZVWEAUDg1Ardhc1bRboLA13KF69x1iPXeeIfv3XmF
6d99hN/gf4Pb4oGD2+WIgd3jiH/p1oN2yK1tZFbIgGyE24h1yeCPdODVim6ty4RtfcOAbmK1dmlq
xoNpcN6VbXP0pnZy5ptyaceFYh7//+AAAAAAAAAAAAAAAAAAAAAAAJVQXgBRMDMDtl9zXMRne6K8
ZHWKtmJwgLZjcHS2YHBuumVycMVvdmzJdHRqx3dxZrtvam63cWtoomtkHo5aSQOnYWAcsG1gIalx
YkrQhnGc65h52PqigOz8pID2/Kh/+v6tfvr8r3n676Zyvd+YbS//140AwoVkAFtKNQBJQysCXU40
A3c7TwWtV29YxmV+otVribfcb43J4HON2eN2juXedon54HqI/eF+hf3fgYH21358ybBraXKaY1wT
mFpbAKBiXhPKgG9bwHdtcsl/bWnYi3CA5ZV0rfOfe+L7p374/qt+//+xe//+sXr/8qhzvsuNXyTk
n2sA/8PzAK9ZdFC7XHx4pVNsOc5kh9LibJT/6HCY/+tzl//td5X/73qV//F+lP/xgJT/5n6K9NJ3
fK2+bW9qzXl0kcp9cH6ybWUyrm1jS9uJdLrzlIT07pOA3d+Ld7TYiW+ByoJoaM+IZ4Hqm3S39aZ4
7P6ve//9sHr+6qNuk6pyWgqMQF4ex1mIwNRfkP61VXh7y2GGx+Zqmf/pbJv/6m+Z/+t0l//teJX/
63qR/tp1hdm4Zm+AwWp1f+J+h9DqhIn+z3l3ndN8eaHbhHq0yX1tcfOUhPP9mYn//puI//qcg/7x
mXzj55V3scyEa3zPiWlo55pwm/eodeLrom+wtn9fHLJSeknTWJPp312c/8NZhLa1V3eJ3WeU/Ohq
nP/obJr/6nKX/+R0kf3LbH65rFtrdcxsf7jofI/18YSR/+qDivi1aGmM4IN/w+2LhfnLeXJ33Yd6
uvmXh//9m4j//56H//+ghv/+oYT/+aF/9+ybd9PXjW6JxoRqU7l8YSi+jmECtEx9bdZXl/3hWqD/
01uS7aNMbXzNYYjY5WqZ/+Zsmf/ebZD7vmJ3sK1bbHvXb4fU63uS//B+lP/zgZX/5H6I6bNlaXXp
hIfh946N/9qAfMjFdW568pKE9PyXif/+mon//52H//6ghP/6oIH68Z1649qPcKq9fmZAq3FdCrF5
XwCvR3x91lOZ/+BWof/bWpr/ulN/rLFUdZrbZZL83GaS+bhaeaWmU2yS1myJ5uZzk//kdY//43iM
/+Z7jP/UdoDByHN1cumDiPv0io7/64iG87pta33hh33I+ZOK//mVif/xln/t7JV7utSHcIHKgmll
2YxvgeeWdbvfkXGoxoNiRqxIeG3QUZX831Oh/99WoP/OWI7mpU1uXK9QdYK2WHd4jUpeQKhYbaG+
X3mdxWR9hMlofnfDaXd4w2x1erdpbUm3bWtByXN2nNZ6fLbhgoHPyXV2g816dEXdhH2b24V6msN3
a3PEeGt95Y18qvKZfd76oID+/6SE//ykgP/ckXGkp01yQctQkuLbUZ//1lOZ/cBShcWcSGcuAAAA
AXY4TRG3W3dfvlt+lcZfga3IY4HE0GmFyNZuhsvNbX7ItmZuUMVueIjNb32pyG55j8J1cGCiXmAe
/6GWAIxeTAi3dWVd3IV81PSUhvj8mof//5+F//+hhP//o4X/+aR7/dySbHeDPFkVu0yGqb5Lic2t
R3uKnEZraZNLXxwAAAAAl0lkRNFgjPDhZJn/5Wia/+hrmv/pbpn/7HOZ/+Bxj+e2YXBx3naI6Ox+
kP/gfIfovXBuhJ9iWhx0Rz4Cr29hF7NvY1vNfm+H4453zPWYgvf9m4b//6GD//+ig//xnHnbzodo
OP+CuACGP1slnkRwcrlJhsi8ToXbm0dpY6tOdH6dSWlqxlqGxeBhmv/mZpz/52ic/+lsnP/pb5n/
2myMwbVfcYLkdY/23HWGz7lmbnPIb3eR1Hh9xrtrb1nEdG+L4oWA6Nd+erDEeGtwyXxuf+aPfMn3
moL9+ZyB/N+QdIN+UUQDjENfAFk4NgWqSXeGx0+P+7BHfK6rR3ig01SX/cdUis+jS2x1y12I0+Jl
mv/lZpv/52qa/+Vsl//KZoGnwWZ4m9Btgc6yYGxqzG5+uOh9jfzwg5D/0nV+sMdydZLxiov/9Y6L
/+yLhe/bhHmswnZrataGcpLZiHSnv3dnHtKDcQBZITcArk6AAIo7YRalSXVnl0JpeMZRjOfcU6D/
2lWc/8FUhsebR2iDy1yJ0+Jkmv/mZp3/4WmW+LJccYutXW2XrFtsc85sf8bqeZH+8H2V//GBk//Z
eYHPynN2ceuEiv33i4//946N//SQiP/ih3/bt3BqR5hfUwySX04BmWFTAAAAAAAAAAAAOE8UAKBE
cACZRGsvtkqBscpOkufSUpfy0laU/7tSgc6eS2h6yFyHxd5imP/YZpDroFJpUpFPX0TIZn+85XSR
/+13l//vepb/8H2U/996hua1ZWt35oGI7u+Hi/bpiITn3oR8utB8dme0a2kUx3dyAAAAAAAAAAAA
AAAAAAAAAAAAAAAAJgAgAA4AEAGFOVsNnD1wKaFDcVC9S4fJ2FSa/8NTh9iiSW2Bull6mrtdenyN
SVwNwWF+feBrk/vrcZn/7HWY/+13lv/vepb/5nqN6rdlbXPUd33Vw210Z8RxcSiybGUQb0ZAA3ZK
RAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACNOWIAAAAAArtJh4DZU53/3Fad
/81WkN+fS2pTAAAAAQAAAAHKZISe5GqY/+htmf/qcZf/7HSX/+54l//oeZDtuGdveMNvc7/Fb3Qa
yHF2AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFsbQgDi
YqIArUZ8VNBRlfXeVaD/2leb/7lTf5d6PFIHFxkDA7hedWPUZI3f42uW/+lvmf/rcpj/7XWY/+Z2
kfOyYm6BrGVoeqtkZxCtZmkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAJZDZgCRQmIdwk2LwtpPn//cU5//yFOMyJ9PZ0rDW4KavFl8lKVSa3e6XHig
0GaH4eFskvrocZb/23GK8aZeaGaJVlgdhFFUAohUVwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAEHAP+B4wCqSXZWzE6T7tdTmv+7T4KutU58odpa
mP/cXpn8zl2L5bZVebGiTmmJsFhxkcxpga+2YHGKg0dRFJxTYAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAo0hvAIA/Uwa1TH5+
xVGL6KVIcYzGU4rJ3lme/+Jdn//jYJ3/4GKZ/9Zij+/OY4e8sVpyi4tNWCWNTloAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAizxeAIo7XQyaSWpRoUhwWctTkODcVp3/3lid/+Fenf/iYJz/4GSZ/9djkPC8XHuckVFc
HJVSXwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKRKcgCjSnEVrkp5c71OhLvCUojcyFWL3cZahsTB
XX+QuVl5S49GXA6rVG8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAG4zSAAvFh0Dj0Jg
D5NAZBuYQWccnlBmEYBMSwaSUVsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA////
//////////////////AAAH/gAAAH4AAAA4AAQAGAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAQA
AgAAAIAAAACAAAABwAAAAfAAAAfwAAAP/gAAf/8AAH//AAB//4AB//+AA///wAP///AH///4H///
//////////////8=`
	faviconBuffer, _ := base64.StdEncoding.DecodeString(faviconBase64)

	r1.GET("/favicon.ico", func(c *gin.Context) {
		_, _ = c.Writer.WriteString(string(faviconBuffer))
	})

	go func() {
		fmt.Printf("listening on http://127.0.0.1:%d using gin router\n", webPort)
		r1.Run(fmt.Sprintf(":%d", webPort))
		wg.Done()
	}()

	wg.Wait()
}
