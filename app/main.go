package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"io/ioutil"
	"io"
	"log"
	"path/filepath"
	"encoding/json"
	"strings"
	"github.com/zalando/go-keyring"
	"github.com/Jeffail/gabs"
	"bytes"
)

var token_variable = ""
var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
//var logfile
var keychain_app_service = "github-forkrefresh1"
var keychain_username = "dmore1"
var app_log_file = "app.log"

//log to std out
func setup_logging(){

	//logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	var logfile, err  = os.Create(app_log_file)

    if err != nil {
     log.Fatal(err)
    }

    defer logfile.Close()

    log.SetOutput(logfile)
    //aiming to also log to console
    log.SetOutput(os.Stdout)
    log.SetOutput(os.Stderr)

}

//this method stores your secret on the OS keychain. 
func store_secret_on_keychain(token string ){

	service := keychain_app_service
    user := keychain_username
    password := token
    //if you want to inject it onto or from an env var...
    //password := token
    //os.Setenv("GITHUB_TOKEN",password)
    //password = os.Getenv("GITHUB_TOKEN")
  
    // set password
    err := keyring.Set(service, user, password)
    if err != nil {
        log.Fatal(err)
    }
}

//this method retrieves secret from the keychain. 
func retrieve_secret_from_keychain() (string){

	service := keychain_app_service
    user := keychain_username

	// get password
    secret, err := keyring.Get(service, user)
    if err != nil {
        log.Fatal(err)
    }

    return secret
}


func main() {

	//setup log to file
	setup_logging()

	//keychain app service from env var.
	keychain_app_service = os.Getenv("KEYCHAIN_APP_SERVICE")
	if (len(keychain_app_service) == 0 ) {
        //log.Fatal(err)
        os.Exit(1)
    }
	log.Println("keychain_app_service:" + keychain_app_service)

	//keychain username from env var
	keychain_username = os.Getenv("KEYCHAIN_USERNAME")
	if (len(keychain_username) == 0){
		//log.Fatal(err)
		os.Exit(1)
	}
	log.Println("keychain_username:" + keychain_username)

	//uncomment to store your secret o keychain
	//store_secret_on_keychain("GITHUB_TOKEN")

	var github_token_env_var = os.Getenv("GITHUB_TOKEN")
	if (len(github_token_env_var) == 0){
		//pull it from the keychain instead.
		token_variable = retrieve_secret_from_keychain()
	}else{
		//use passed in token if non-empty string
		token_variable = github_token_env_var
	}

	var github_gist = os.Getenv("REPOS_GIST")

	var arr []string
	var err error
	//var getErr = nil 
	//var readErr = nil 
	
	if (len(github_gist) > 0){

		log.Println("downloading repos list from gist")
		log.Println("" + github_gist)
		var url = github_gist
		
		var gistClient = http.Client{
			Timeout: time.Second * 2, // Timeout after 2 seconds
		}

		var req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Println("unable to pull down yer gist")
			log.Fatal(err)
			os.Exit(1)
		}

		req.Header.Set("User-Agent", "test")

		var res, getErr = gistClient.Do(req)
		if getErr != nil {
			log.Println("error http get when pulling down gist")
			log.Fatal(getErr)
			os.Exit(1)
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		var body, readErr = ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Println("error reading the gist")
			log.Fatal(readErr)
			os.Exit(1)
		}
		
		var err3 = json.Unmarshal([]byte(body), &arr)
		if err3 != nil {
	      fmt.Println("error3:", err3)
	      os.Exit(1)
	    }
	}else{
		/** using repos_repo.json file **/
		//refactor to pull dynamically.
		//file must be json array not json
		absPath, _ := filepath.Abs("repos_repo.json")
		f, err := os.Open(absPath)
		if err != nil {
		    log.Fatal(err)
		}
		log.Println("Successfully Opened repos_repo.json")
		defer f.Close()

		//unmarshall
		byteValue, _ := ioutil.ReadAll(f)   
		err3 := json.Unmarshal([]byte(byteValue), &arr)
		if err3 != nil {
	      fmt.Println("error3:", err3)
	      os.Exit(1)
	    }
	}

	
	//var result map[string]interface{} not working. 
	//in memory test works fine also
	/**
	var dataJson = `[
    	"dmore/aws-vault-local-os-keychain-mfa",
    	"dmore/aws-multi-region-cicd-with-terraform"
	]`
	err2 := json.Unmarshal([]byte(dataJson), &arr)
	if err2 != nil {
      fmt.Println("error2:", err2)
      os.Exit(1)
    }
    **/
	
	
    log.Printf("Unmarshaled: %s", arr)
    //loop through
    for i := 0; i < len(arr); i++ {
		
		var reponame = string(arr[i])
		reponame = strings.TrimSuffix(reponame, "/")
		reponame = strings.TrimPrefix(reponame, "/")
		fmt.Println(reponame)
		
		//LOOP HERE EACH ELEMENT

	    var ret = ""
	    var branch = ""
	    //grab branches first so we know what branch name we need in advance...
		branch, err = fork_get_query_branch(reponame)
		if err != nil {
			log.Fatalln(err)
			//ret2,err2 := call("", "main", "", "POST")
			//if err2 != nil {
			//	log.Fatalln(err)
			//	//continue
			//}
			//404
			log.Println("notok1")
			continue
		}
		//has the branch been found though 404?
		if strings.Contains(string(ret), "Not Found") {
			log.Println("branch not found...")
			log.Println(string(ret))

			log.Println("Branch Not Found so we will skip this one ")
			log.Println("branch not there....repo does not exist...")
			continue
		}else{
			log.Println("branch found ok")
		}
	    //relevant branch call
		ret,err = fork_refresh_call(branch, reponame, "POST")
		if err != nil {
			log.Fatalln(err)
			//ret2,err2 := call("", "main", "", "POST")
			//if err2 != nil {
			//	log.Fatalln(err)
			//	//continue
			//}
			//404
			log.Println("notok2")
			continue
		}
		//don't print it too much content
		log.Println("- - - - - - -- - - - - -- - -- - -")
		log.Println(string(ret))
		log.Println("- - - - - - -- - - - - -- - -- - -")
		
		if strings.Contains(string(ret), "Not Found") {
			log.Println("checking...")
			log.Println(string(ret))

			log.Println("Not Found found on stringifed response => [main]")
			//main call
			ret,err = fork_refresh_call("main", reponame, "POST")
			if err != nil {
				log.Fatalln(err)
				//ret2,err2 := call("", "main", "", "POST")
				//if err2 != nil {
				//	log.Fatalln(err)
				//	//continue
				//}
			}
			log.Println("notok3")
			continue
		}else{
			log.Println("ok")
			continue
		}
		
	}
	//exit after looping repo names
	log.Println("Done")

	//run as a cron
	os.Exit(0)

}


func fork_get_query_branch(reponame string) (string, error) {

	reponame = strings.TrimSuffix(reponame, "/")
	reponame = strings.TrimPrefix(reponame, "/")

	httpposturl := "https://api.github.com/repos/" + reponame + "/branches"
	fmt.Println("url: " + string(httpposturl))
	request, err := http.NewRequest("GET", httpposturl, nil)
	if err != nil {
	    log.Fatal(err)
	}
	//request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Accept", "Accept: application/vnd.github+json")
	request.Header.Set("Authorization", "token " + token_variable)
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
	    log.Fatal(err)
	}
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		log.Fatalln(err)
		return "nil", err
	}

	//fmt.Println("response :", response.Errorf)
	log.Println("response Status:", response.Status)
	log.Println("response Body:", string(b))

	if strings.Contains(string(b), "Not Found") {
		return "nil", err
	}

	var objmap []map[string]interface{}
	if err := json.Unmarshal(b, &objmap); err != nil {
    	log.Fatal(err)
    	//branch is not there, so blows up
    	return "unable_to_unmarshall_branch_response", err
	}
	if (len(objmap) == 0){
		return "failed to log", err
	}
	log.Println("Unmarshalled")
	fmt.Println(objmap[0]["name"])

	var return_branch = ""
	for k, v := range objmap[0] {
	    switch c := v.(type) {
	    case string:
	    	if k == "name" {
	    		return_branch = string(c)
	    		fmt.Printf("Item %q is a string, containing %q\n", k, c)
	    	}
	        
	    case float64:
	        //fmt.Printf("Looks like item %q is a number, specifically %f\n", k, c)
	        continue
	    default:
	        //fmt.Printf("Not sure what type item %q is, but I think it might be %T\n", k, c)
	        continue
	    }
	}
	fmt.Println("return_branch is " + return_branch)
	return string(return_branch), nil
}

func fork_refresh_call(branch string, reponame string, method string) (string, error) {
	//now that we know the branch name in advance we can use that instead of this.
	
	jsonObj := gabs.New()
	// or gabs.Wrap(jsonObject) to work on an existing map[string]interface{}

	jsonObj.Set("" + branch, "branch")

	jsonOutput := jsonObj.String()

	fmt.Println(jsonObj.String())
	fmt.Println(jsonObj.StringIndent("", "  "))

	var jsonStr = []byte(jsonOutput)
   
	reponame = strings.TrimSuffix(reponame, "/")
	reponame = strings.TrimPrefix(reponame, "/")

	httpposturl := "https://api.github.com/repos/" + reponame + "/merge-upstream"
	//fmt.Println("url: %s", (string)httpposturl)
	request, err := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonStr))
	if err != nil {
	    log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Set("Authorization", "token " + token_variable)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
	    log.Fatal(err)
	}
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
		return "nil", err
	}

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Body:", string(b))
	return string(b), nil
	//return fmt.Println(string(b))
}
