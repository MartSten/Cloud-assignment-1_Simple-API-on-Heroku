package main

import (
	"fmt"
	"net/http"
	"strings"
	"encoding/json"
	"io/ioutil"
	"os"
)


/**
Function for handling of the url.
 */
func urlHandler(w http.ResponseWriter, r *http.Request){

	//framed code form: https://stackoverflow.com/questions/31622052/how-to-serve-up-a-json-response-using-go/31622112
	//-------------------------------------------------------------
	w.Header().Set("Content-Type", "application/json")	//sets content-type so that user clients will know to expect json
	//-------------------------------------------------------------

	//Function variables
	var apiContributor string	//The url for the contributors api
	var apiLanguage string	//The url for the languages api
	var url string = string(r.URL.Path)	//gets the url fom r (the http request) and casts it to string

	var projectData = strings.Split(url, "/")	//gets the name and owner of a github project
	if len(projectData) != 6 {	//if the incorrect url is requested: send error.
		http.Error(w, "incorrect url. needs 5 parameters", 400)
		return
	}

	apiContributor= "https://api." + projectData[3] + "/repos/" + projectData[4] + "/" + projectData[5] + "/contributors"

	apiLanguage= "https://api." + projectData[3] + "/repos/" + projectData[4] + "/" + projectData[5] + "/languages"

	pOwner := projectData[4]	//project owner
	theProject := projectData[5]	//project name
	var topCom string	//top contributor
	var amountOfComs int	//top contributor's amount of commits

	theContributors := getContributor(apiContributor)
	if len(theContributors) < 1 {	//if there were errors getting the contributor data
		fmt.Println("Could not find the correct api (contributor) data. Is the url correct?")
		topCom = "No committers found"
		amountOfComs = 0
	} else {	//If there were no errors getting the contributor data
		topCom = theContributors[0].Committer
		amountOfComs = theContributors[0].Commits
	}

	theLanguages := getLanguage(apiLanguage)
	if len(theLanguages.Language) < 1{
		fmt.Println("Could not find the correct api (language) data. Is the url correct?")
		theLanguages.Language = append(theLanguages.Language, "No languages fund")
	}

	//Struct for sending the final json payload
	theResponse := finalStruct{
		Project: theProject,
		Owner: pOwner,
		Committer: topCom,
		Commits: amountOfComs,
		Language: theLanguages.Language,
	}

	json.NewEncoder(w).Encode(theResponse)

}


/**
GET /repos/:owner/:repo/languages
A struct that holds the repo's language data
 */

type languages struct {
	Language []string
}

/**
function getLanguage()
The function gets a github repo's languages based on an url to the github api
PARAM: The function takes an url to a github repo's language api
RETURN: On success the function returns an Language slice containing languages. On failure the function returns an empty Language{}
 */
func getLanguage(urlToAPI string) languages {

	resp, err := http.Get(urlToAPI)	//Gets the api's content based on url
	if err != nil{
		//http.Error(writer, http.StatusText(400), 400)
	}

	body, err := ioutil.ReadAll(resp.Body)

	//Done with request. Closing it.
	defer resp.Body.Close()

	var a interface{}	//create interface a

	lang := languages{}

	//Process languages by unmarshaling them, and puts them inn to a
	err =json.Unmarshal(body, &a)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not find the correct api (languages) data. Is the url correct?")
		return languages{}
	}
	//fmt.Println(language)

	m := a.(map[string]interface{})		//map a
	//for each a - put it inn to languages through lang.Language
	for k := range m {
		lang.Language = append(lang.Language, k)
	}
	//fmt.Println(lang)
	if lang.Language[0] == "message" {
		fmt.Println("Could not find the correct api (languages) data. Is the github-url correct?")
		return languages{}
	}

	return lang
}

/**
GET /repos/:owner/:repo/contributors
A struct that holds the repo's contributor data
based on the github api
 */
type contributors []struct {
	Committer string `json:"login"`
	Commits int `json:"contributions"`
}

/**
function getContributor()
the function gets a github repo's contributors from the github api.
it then processes them and returns them as an array
PARAM: the function takes the url for the github api as it's parameter
RETURN: on success it returns an array of contributors. On failure it returns an empty contributor
 */
func getContributor(urlToAPI string) contributors {

	resp, err := http.Get(urlToAPI)
	if err != nil{
		//http.Error(writer, http.StatusText(400), 400)
	}

	body, err := ioutil.ReadAll(resp.Body)

	//Done with request. Closing it.
	defer resp.Body.Close()

	//Process json fom the api
	topCon := contributors{}
	err =json.Unmarshal(body, &topCon)
	if err != nil {		//if the json can't be placed in contributors, see if it can be placed in apiError
		fmt.Println(err)
		errorAPI := apiError{}
		err2 := json.Unmarshal(body, &errorAPI)
		if err2 != nil {	//if the json can't be placed in apiError
				fmt.Println(err2)
			return contributors{}
			}
		fmt.Println("Could not find the correct api data. Is the url correct?")
		return contributors{}
	}
	//fmt.Println(topCon)
	return topCon
}

/**
apiError struct
Holds data fom the github api's error message
Used when the script fails to get data fom a given api url
 */
type apiError struct{
	Message string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

/**
The struct finalStruct holds the data which is to be presented as json to the user
 */
type finalStruct struct{
	Project string
	Owner string
	Committer string
	Commits int
	Language []string
}


/**
Function main
Does main things
 */
func main() {
	port := os.Getenv("port")
	http.HandleFunc("/projectinfo/v1/", urlHandler)
	//getContributor("https://api.github.com/repos/apache/camel/contributors")
	//getLanguage("https://api.github.com/repos/apache/camel/languages")
	fmt.Println(port)
	//fmt.Println(os.Environ())
	http.ListenAndServe(port, nil)

}
