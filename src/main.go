package main
import (
	"fmt"
	//import the Paho Go MQTT library
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"./example"
	"os"
	"path/filepath"
	"log"
	"strings"
	"net/http"
	"io/ioutil"
	"sync"
	"io"
	"os/exec"
	"net/url"
	"encoding/json"
	"reflect"
)

type receiveMessage struct {
	topic string
	payload []byte
}

type httpresult struct{
	Artist string `json:"Artist"`
	Song string `json:"Song"`
	Coverurl string `json:"Cover"`
	Songurl string `json:"Url"`
	NoticeOrNot bool `json:"NoticeOrNot"`
}

//type mqttResponse struct{
//	Result httpresult  `json:"Result"`
//	Response string `json:"Response"`
//}

type cmdMsg struct{
	Cmd string `json:"cmd"`
	Progress int `json:"progress,omitempty"`
	Volume int `json:"volume,omitempty"`
}

type Audiostate struct{
	Playstate string
	Playprogress int

	PlayParams httpresult

	Playvolume int

	VoiceboxConnect bool
	OrderText string
}

var choke = make(chan receiveMessage)
//var curPath string = getCurrentDirectory()
var queryurl = "http://10.134.142.140:8080/Query"+"?"+"query"+"="
var wait = &sync.WaitGroup{}
var broker = "tcp://127.0.0.1:61613"
var clientId = "test"
var mytopic = "yipai_test/topic"
var songtopic = mytopic+"/song"
var cmdtopic = mytopic+"/cmd"
var MAudiostate Audiostate

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	choke <- receiveMessage{msg.Topic(),msg.Payload()}

}

func checkFileIsExist(filename string) (bool) {
	var exist = true;
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false;
	}
	return exist;
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
func httpGet(queryurl string,text string) string {

	params := url.Values{"query":{text}}

	client := &http.Client{}
	resp, err := client.PostForm(queryurl,params)
	check(err)

	defer resp.Body.Close()
	body, err3 := ioutil.ReadAll(resp.Body)
	check(err3)
	//
	fmt.Println(string(body))
	return string(body)
}



func main() {


	defer func(){ // 必须要先声明defer，否则不能捕获到panic异常
		if err:=recover();err!=nil{
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
	}()

	//opts := MQTT.NewClientOptions().AddBroker("test.amber-link.com:1883")
	opts := MQTT.NewClientOptions().AddBroker(broker)
	opts.SetClientID(clientId)
	opts.SetDefaultPublishHandler(f)
	opts.SetPassword("password")
	opts.SetUsername("admin")



	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	go receiveMQTT(&c)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}else{
		fmt.Println("SUCCESSFUL CONNECT")
		fmt.Println(token)
	}

	if token := c.Subscribe(mytopic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println("FAILED TO SUBSCRIBE")
		fmt.Println(token.Error())
	}else{
		fmt.Println("SUCCESSFUL SUBSCRIBE")
		fmt.Println(token)
	}

	http.Handle("/", http.FileServer(http.Dir("/tmp/static/")))
	//http.HandleFunc("/", responseVoice)
	// 设置监听的端口
	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	wait.Wait()
}

func receiveMQTT(client *MQTT.Client){
	wait.Add(1)
	defer wait.Done()
	for true {
		msg := <-choke

		var f *os.File
		var text string
		var songresult httpresult
		if string(msg.payload) == "end" {
			if (checkFileIsExist("test.raw")) {

				cmd := exec.Command("/usr/local/ffmpeg/bin/ffmpeg", "-f", "s16le", "-ar", "16000", "-ac", "1", "-acodec", "pcm_s16le", "-i", "test.raw", "test.wav")
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					fmt.Printf("Run returns: %s\n", err)
				}
				text = example.WavToText("test.wav")
				fmt.Print(text)
				os.Remove("test.raw")
				os.Remove("test.wav")

				songresult = explaintext(text,client)
				songresult.NoticeOrNot = true
				jsonstring, _ := json.Marshal(songresult)
				(*client).Publish(songtopic, 0, false, jsonstring)
			}
		} else if string(msg.payload) == "cancel" {
			if (checkFileIsExist("test.raw")) {
				os.Remove("test.raw")
			}
			if (checkFileIsExist("test.wav")) {
				os.Remove("test.wav")
			}
		} else if (len(msg.payload)==1280){
			f, _ = os.OpenFile("test.raw", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			n, _ :=io.WriteString(f, string(msg.payload))
			f.Close()
			println(n)
		} else {
			(*client).Publish(cmdtopic,0,false, msg.payload)
			println(string(msg.payload))
			var cmdresult cmdMsg
			json.Unmarshal(msg.payload,&cmdresult)
			example.PrintStruct(reflect.TypeOf(cmdresult),reflect.ValueOf(cmdresult),1)
			switch cmdresult.Cmd {
			case "PROGRESSCHANGE":
				MAudiostate.Playprogress = cmdresult.Progress
			case "VOLUMECHANGE":
				MAudiostate.Playvolume = cmdresult.Volume
			case "PLAY":
				MAudiostate.Playstate = "PLAYING"
			case "PAUSE":
				MAudiostate.Playstate = "PAUSED"
			case "NEXT":
				songresult = explaintext(MAudiostate.PlayParams.Song,client)
				songresult.NoticeOrNot = false
				jsonstring, _ := json.Marshal(songresult)
				(*client).Publish(songtopic, 0, false, jsonstring)
			default:
				println("default")
			}
			example.PrintStruct(reflect.TypeOf(MAudiostate),reflect.ValueOf(MAudiostate),7)
		}
	}


}


func explaintext(text string,client *MQTT.Client)(httpresult){
	result := httpGet(queryurl,text)
	var songresult httpresult
	json.Unmarshal([]byte(result),&songresult)
	if(songresult.Song!=""&&songresult.Artist!="") {
		MAudiostate.PlayParams = songresult
		MAudiostate.Playstate = "PLAYING"
		MAudiostate.Playprogress = 0
		MAudiostate.OrderText = text
		example.TextToSpeech("正在为您播放 "+songresult.Artist+" "+songresult.Song, "/tmp/static/result.wav")
		if(checkFileIsExist("/tmp/static/result.mp3")) {
			os.Remove("/tmp/static/result.mp3")
		}
		cmd2 := exec.Command("/usr/local/ffmpeg/bin/ffmpeg",  "-i", "/tmp/static/result.wav", "/tmp/static/result.mp3")
		cmd2.Stderr = os.Stderr
		if err := cmd2.Run(); err != nil {
			fmt.Printf("Run returns: %s\n", err)
		}


	}else{
		example.TextToSpeech("抱歉没有找到这首歌","/tmp/static/result.wav")
		if(checkFileIsExist("/tmp/static/result.mp3")) {
			os.Remove("/tmp/static/result.mp3")
		}
		cmd2 := exec.Command("/usr/local/ffmpeg/bin/ffmpeg",  "-i", "/tmp/static/result.wav", "/tmp/static/result.mp3")
		cmd2.Stderr = os.Stderr
		if err := cmd2.Run(); err != nil {
			fmt.Printf("Run returns: %s\n", err)
		}
	}
	return songresult
}


func responseVoice(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	// 这些信息是输出到服务器端的打印信息
	//fmt.Println("request map:", r.Form)
	//fmt.Println("path", r.URL.Path)
	//fmt.Println("scheme", r.URL.Scheme)
	//fmt.Println(r.Form["url_long"])
	//
	//for k, v := range r.Form {
	//	fmt.Println("key:", k)
	//	fmt.Println("val:", strings.Join(v, ";"))
	//}

	// 这个写入到w的信息是输出到客户端的
	if(r.URL.Path=="/") {

		w.Header().Set("Content-Disposition","attachment; filename=\"result.wav\"")
		w.Header().Set("content-type", "audio/wav")
		w.Header().Set("content-type", "audio/wav")
		fin, err := os.Open("/tmp/static/result.wav")
		defer fin.Close()
		check(err)
		fd, _ := ioutil.ReadAll(fin)
		w.Write(fd)
	}else {
		http.NotFound(w, r)
	}
}
