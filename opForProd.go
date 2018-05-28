package main
import (
"os"
"os/exec"
"fmt"
"io/ioutil"
"bufio"
"strings"
"net/http"
"strconv"
"time"
"runtime"
"encoding/json"
)

var (
_NEEDLES []string
_PREFIX []string
_SUFIX []string
_HTML string
_PORT = "555"
_PAGE = "IsItOpForProd"
_PAGEDELIVERED = false
)


type Line struct{
	Content string
	Location string
	Index int} 

func init() {
	_PREFIX = []string{
		"", "'", "\"", "//", "/*", "#", 
	}
	_SUFIX = []string{
		"", "'", "\"",
	}
	b, err := ioutil.ReadFile("./list.json") // just pass the file name
    if err != nil {
    	fmt.Println("list.json not found")
        return
    }
    json.Unmarshal(b, &_NEEDLES)
}

func main(){

	arg := "";
	var files []string;
	var lines []Line;


	if len(os.Args) < 2  {
		fmt.Println( "NO ARGS" )
		return 
	}else{
		arg = os.Args[1]
	}


	if b, r := isDirectory(arg); r != nil{//FAIL
		fmt.Println(arg , "NOT A FILE NOR A FODLER")
		return
	}else if !b{//FILE
		files = []string{ arg }
	}else if b{//DIRECTORY
		arg = directorySeparator(arg)
		files = explore(arg)
	}

	for _, f := range files{
		lines=append(lines,readFile(f)...)
	}

	displayResult(lines)
}

func isDirectory(path string) (bool, error) {
    if  fileInfo, err := os.Stat(path); err == nil{
    	return fileInfo.IsDir(),nil
    }else{
    	 return false, err
    }
 }
   
func directorySeparator(path string) (newpath string) {
	newpath =  path
	if   path[len(path)-1:] != string (os.PathSeparator){
		newpath += string (os.PathSeparator) 
	}
	return 
}

func explore(loc string)(contents []string){
   files, _ := ioutil.ReadDir(loc)
   for _, f := range files {
   		str := loc + f.Name()
         if boolean, _ := isDirectory( str ); boolean{
         	str = directorySeparator(str)
         		contents = append(contents,explore(str)... )
         	}else{
         		contents = append(contents , str )
         	}
    }
   return}

func readFile(loc string)(contents []Line){
 file, _ := os.Open(loc)
  defer file.Close()
  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)

  i:=1
  for scanner.Scan() {
  	s_bfr := scanner.Text();
  	if detectNeedle(s_bfr){
  		contents = append(contents, Line{s_bfr,loc,i})
  	}
  	i++
  }
  return}

func detectNeedle(line string) bool{
	for _, needle := range _NEEDLES{
		for _, pre := range _PREFIX{
			for _, su := range _SUFIX{
			  	if strings.Contains(line, pre+needle+su){//EqualFold
			  		  return true
			  	}
			} 
		} 
	}
  	return false}

func displayBrowser(){
	switch runtime.GOOS {
	case "linux":
	    exec.Command("xdg-open", "http://localhost:"+_PORT+"/"+_PAGE+"/").Start()
	case "darwin":
	    exec.Command("open", "http://localhost:"+_PORT+"/"+_PAGE+"/").Start()
	 case "windows":
	 	exec.Command(`C:\Windows\System32\rundll32.exe`, "url.dll,FileProtocolHandler", "http://localhost:"+_PORT+"/"+_PAGE+"/").Start()
	}}

func handleWebServer(){
	http.HandleFunc("/"+_PAGE+"/", servResult)
	http.ListenAndServe(":"+_PORT, nil)
}

func formateHTMl(lines []Line){
	_HTML = "<div style='text-align:center'>"
	if len(lines)==0{
		_HTML += "<p>IT SEEMS OP FOR PROD</p>"
		return
	}else{
		_HTML +="<p>"+ strconv.FormatInt(int64(len(lines)), 10)   +" verification required for prod</p><br><br>"
	}

	for _, f := range lines{
		_HTML+= "<p>"+f.Location +":"+ strconv.FormatInt(int64(f.Index), 10)+ "</p>"
		_HTML+= "<p style='color:lightcoral'>"+f.Content + "</p><br>"
	}
	_HTML+= "</div>"
}

func servResult(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, _HTML)
    _PAGEDELIVERED=true
}


func displayResult(lines []Line){
	formateHTMl(lines)

	go handleWebServer()
	displayBrowser()
	for !_PAGEDELIVERED{ time.Sleep(5 * time.Second)}
} 
