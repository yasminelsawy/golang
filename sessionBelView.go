package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
	//"github.com/gorilla/mux"

	cors "github.com/heppu/simple-cors"
)

var (
	WelcomeMessage = "Welcome, Do u want a used car or a new car ?"

	// sessions = {
	//   "uuid1" = Session{...},
	//   ...
	// }
	sessions = map[string]Session{}

	processor = sampleProcessor
	//CarTypeMap     = make(map[string]string)
	//ChatStage      = 0
	//ChatStageUsed  = 0
	//ChatStageInt   = 0
	//WelcomeMessage = "Welcome, Do u want a used car or a new car ?"
	//User           Person
	//url            = "https://egypt.yallamotor.com/new-cars"
	//urlUsed        = "https://eg.hatla2ee.com/en/car"
	//CarTypeUsedMap = make(map[string]string)
	//CarMakeUsedMap = make(map[string]string)
	//CarVersionMap  = make(map[string]string)
	//contentY       = []string{}
	//contentK       = []string{}
	//errorNoCar     = ""
)

type (
	Car struct {
		Model   string `json:"Model,omitempty"`
		Type    string `json:"Type,omitempty"`
		YearMax int64  `json:"YearMax,omitempty"`
		YearMin int64  `json:"YearMin,omitempty"`
		//from the user input(USED BAS)
		PriceMax float64 `json:"PriceMax,omitempty"`
		//from the user input(USED BAS)
		PriceMin float64 `json:"PriceMin,omitempty"`
		Used     bool    `json:"Used,omitempty"`
		Overview string  `json:"Overview,omitempty"`
		//the final price of the version
		Price   float64 `json:"PriceMax,omitempty"` //the car price will me only 1 value not range   remove 2 price values and let it be 1 value price
		Version string  `json:"Version,omitempty"`
	}

	Person struct {
		UUID      string `json:"UUID,omitempty"`
		CarPerson Car    `json:"Cars,omitempty"`
	}
)

type (
	// Session Holds info about a session
	Session map[string]interface{}

	// JSON Holds a JSON object
	JSON map[string]interface{}

	// Processor Alias for Process func
	Processor func(session Session, message string) (string, error)
)

func sampleProcessor(session Session, message string) (string, error) {
	// Make sure a history key is defined in the session which points to a slice of strings

//var message string
	// Process the received message
	//fmt.Printf("%#v\n", session)
	ChatStageInt, ifAssignedInt := session["ChatStageInt"].(int)
	fmt.Println(ChatStageInt)
	fmt.Println(ifAssignedInt)
	if !ifAssignedInt {
		ChatStageInt = 0
	}

	switch ChatStageInt {
	case 0:
		message = checkOnWelcomeMessage(message,session )
		fmt.Println("CASE 0", message)
	}

	ChatStage, ifAssigned := session["ChatStage"].(int)

	if !ifAssigned {
		ChatStage = 0
	}

	fmt.Println(ChatStage)
	switch ChatStage {
	case 1:
		message = ModelCar(session/*uuid*/)
	//	fmt.Println("ChatStage --> CASE 1", message)
	case 2:
		message = checkType(message,session/*uuid*/)
		// fmt.Println("ChatStage --> CASE 2", message)
	case 3:
		message = TypeCar(session/*uuid*/)
		//fmt.Println("ChatStage --> CASE 3", message)
	case 4:
		message = checkOnCarType(message,session/*uuid*/)
		//fmt.Println("ChatStage --> CASE 4", message)
	case 5:
		message = VersionCar(session)
		//fmt.Println("ChatStage --> CASE 5", message)
	case 6:
		message = checkVersion(message,session)
		//fmt.Println("ChatStage --> CASE 6", message)
	case 7:
		message = finalresult(session)
		//fmt.Println("ChatStage --> CASE 7", message)

	}
	ChatStageUsed, ifAssigned := session["ChatStageUsed"].(int)

	if !ifAssigned {
		ChatStageUsed = 0
	}

	session["ChatStageUsed"] = ChatStageUsed

	switch ChatStageUsed {

	case 1:
		message = ModelCarUsed(session)
		//fmt.Println("ChatStageUsed --> CASE 1", message)
	case 2:
		message = checkMakeused(message,session)
		//fmt.Println("ChatStageUsed --> CASE 2", message)
	case 3:
		message = TypeCarUsed(session)
		//fmt.Println("ChatStageUsed --> CASE 3", message)
	case 4:
		message = checkOnCarTypeUsed(message,session)
		//fmt.Println("ChatStageUsed --> CASE 4", message)
	case 5:
		message = checkOnMinPriceUsed(message,session)
		//fmt.Println("ChatStageUsed --> CASE 5", message)
	case 6:
		message = checkOnMaxPriceUsed(message,session)
		//fmt.Println("ChatStageUsed --> CASE 6", message)
	case 7:
		message = checkOnMinYearUsed(message,session)
		//fmt.Println("ChatStageUsed --> CASE 7", message)
	case 8:
		message = checkOnMaxYearUsed(message,session)
		//fmt.Println("ChatStageUsed --> CASE 8", message)

	}

	_, historyFound := session["history"]
	if !historyFound {
		session["history"] = []string{}
	}

	// Fetch the history from session and cast it to an array of strings
	history, _ := session["history"].([]string)

	// Make sure the message is unique in history
	for _, m := range history {
		if strings.EqualFold(m, message) {
			return "", fmt.Errorf("You've already ordered %s before!", message)
		}
	}

	// Add the message in the parsed body to the messages in the session
	history = append(history, message)

	// Form a sentence out of the history in the form Message 1, Message 2, and Message 3
	l := len(history)
	wordsForSentence := make([]string, l)
	copy(wordsForSentence, history)
	if l > 1 {
		wordsForSentence[l-1] = "and " + wordsForSentence[l-1]
	}
	sentence := strings.Join(wordsForSentence, ", ")

	// Save the updated history to the session
	session["history"] = history

	return message, nil
	return fmt.Sprintf("So, you want %s! What else?", strings.ToLower(sentence)), nil
}

// withLog Wraps HandlerFuncs to log requests to Stdout
func withLog(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := httptest.NewRecorder()
		fn(c, r)
		log.Printf("[%d] %-4s %s\n", c.Code, r.Method, r.URL.Path)

		for k, v := range c.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(c.Code)
		c.Body.WriteTo(w)
	}
}

// writeJSON Writes the JSON equivilant for data into ResponseWriter w
func writeJSON(w http.ResponseWriter, data JSON) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// ProcessFunc Sets the processor of the chatbot
func ProcessFunc(p Processor) {
	processor = p
}

// handleWelcome Handles /welcome and responds with a welcome message and a generated UUID
func handleWelcome(w http.ResponseWriter, r *http.Request) {
	// Generate a UUID.
	hasher := md5.New()
	hasher.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	uuid := hex.EncodeToString(hasher.Sum(nil))


	session := Session{}

	session["CarTypeMap"] = make(map[string]string)
	session["ChatStage"] = 0
	session["ChatStageUsed"] = 0
	session["ChatStageInt"] = 0
	var User Person
	session["User"] = User
	session["CarTypeUsedMap"] = make(map[string]string)
	session["CarMakeUsedMap"] = make(map[string]string)
	session["CarVersionMap"] = make(map[string]string)
	session["contentY"] = []string{}
	session["contentK"] = []string{}
	session["errorNoCar"] = ""
	session["url"] = "https://egypt.yallamotor.com/new-cars"
	session["urlUsed"] = "https://eg.hatla2ee.com/en/car"

	// Create a session for this UUID
	sessions[uuid] = session
	//User := session["User"].(Person)

	//User.UUID = uuid
	//session["User"]=User
	//fmt.Println(User)
	w.Header().Add("Authorization", uuid)
	fmt.Println(w.Header())
	// Write a JSON containg the welcome message and the generated UUID
	writeJSON(w, JSON{
		"uuid":    uuid,
		"message": WelcomeMessage,
	})
}

func checkOnWelcomeMessage(data string,session Session/*uuid string*/) string {
	//var session Session = sessions[uuid] 

	fmt.Println(data)
	ChatStage, ifAssigned := session["ChatStage"].(int)

	if !ifAssigned {
		ChatStage = 0
	}
	ChatStageUsed, ifAssignedUsed := session["ChatStageUsed"].(int)

	if !ifAssignedUsed {
		ChatStageUsed = 0
	}

	ChatStageInt, ifAssignedInt := session["ChatStageInt"].(int)

	if !ifAssignedInt {
		ChatStageInt = 0
	}

	if strings.Contains(strings.ToLower(data), "new") {
		ChatStageInt=3
		ChatStage++
		session["ChatStage"] = ChatStage
		session["ChatStageInt"] = ChatStageInt
		//fmt.Printf("%#v\n", session)

		//sessions[uuid]=session
		return "Nfsak feh car type eh?"
		//return ModelCar()

	} else if strings.Contains(strings.ToLower(data), "used") {
		User := session["User"].(Person)
		User.CarPerson.Used = true
		ChatStageUsed++
		ChatStage = -1
		ChatStageInt++
		session["ChatStageUsed"] = ChatStageUsed
		session["ChatStage"] = ChatStage
		session["ChatStageInt"] = ChatStageInt
		session["User"] = User

		return "Nfsak feh car type eh?"
		//later add if condition law both are working
	} else {
		return "mosh fahem 3ayz new waala used??"
	}

}

//yasmin search for the types in a certain model(given by gina) NEW
func searchForTypeFillMap(session Session /*uuid string*/) {
	//	url:="https://egypt.yallamotor.com/new-cars/kia"
	//isModel := false
	//var session Session = sessions[uuid] 

	errorNoCar := session["errorNoCar"].(string)
	url, ok := session["url"].(string)
	if !ok {
		session["url"] = "https://egypt.yallamotor.com/new-cars"
	}

	resp, err := http.Get(url)
	ChatStage, ifAssigned := session["ChatStage"].(int)
	fmt.Println(ChatStage)
	fmt.Println(ifAssigned)
	if !ifAssigned {
		ChatStage = 0
	}

	if err != nil {
		fmt.Println("error failed \"" + url)
		return
	}
	if resp.StatusCode == 404 {
		errorNoCar = "No Cars Found"
		session["errorNoCar"] = errorNoCar
		ChatStage--
		session["ChatStage"] = ChatStage
		url = "https://egypt.yallamotor.com/new-cars"
		session["url"] = url
		//sessions[uuid]=session
		return
	}
	ChatStage++
	session["ChatStage"] = ChatStage
	doc := resp.Body
	defer doc.Close()
	//toknize my doc
	tokenizedDoc := html.NewTokenizer(doc)
	User := session["User"].(Person)
	/*given the tokenizerDoc the html is tokenized by repeatedly calling .Next()
	which parses the nest token and return its 1) TYPE 2)ERROR */
	for {
		typeToken := tokenizedDoc.Next()

		switch typeToken {
		case html.ErrorToken:
			return

		case html.StartTagToken:
			{
				tag := tokenizedDoc.Token()
				tobeCompared := "/new-cars/" + User.CarPerson.Type + "/"
				isA := tag.Data == "a"
				if isA {
					CarTypeMap := session["CarTypeMap"].(map[string]string)

					for _, a := range tag.Attr {

						if a.Key == "href" && strings.Contains(a.Val, tobeCompared) {
							tokenizedDoc.Next()
							fmt.Println("ana gow3a")

							mapKey := (tokenizedDoc.Token()).Data
							fmt.Println(mapKey)
							if strings.ToLower(mapKey) != "view detail" && !strings.Contains(strings.ToLower(mapKey), "price in egypt") {

								//CarTypeMap[mapKey] = a.Val
								CarTypeMap[mapKey] = a.Val
							}

							break
						}
					}
					session["CarTypeMap"] = CarTypeMap
					//sessions[uuid]=session

					//isModel=false

				}
			}
		}
	}

}

//yasmin NEW
func TypeCar(session Session/*uuid string*/) string {
	//var session Session = sessions[uuid] 

	errorNoCar := session["errorNoCar"].(string)
	searchForTypeFillMap(session)
	if errorNoCar == "" {
		message := "Aayez anhy model?\n"
		CarTypeMap := session["CarTypeMap"].(map[string]string)
		for key, _ := range CarTypeMap {
			if key != ""{
			message += "" + key + "\n"
				
			}
		}
		return message

	}
	return errorNoCar

}

//yasmin NEW
func checkOnCarType(message string,session Session/*uuid string*/) string {
	//var session Session = sessions[uuid] 

	fmt.Println(message)
	CarTypeMap := session["CarTypeMap"].(map[string]string)

	url, ok := session["url"].(string)
	if !ok {
		session["url"] = "https://egypt.yallamotor.com/new-cars"
	}

	ChatStage, ifAssigned := session["ChatStage"].(int)
	fmt.Println(ChatStage)
	fmt.Println(ifAssigned)
	if !ifAssigned {
		ChatStage = 0
	}

	for key, Val := range CarTypeMap {
		if strings.Contains(strings.ToLower(message), strings.ToLower(key)) {
			ChatStage++
			session["ChatStage"] = ChatStage
			url = "https://egypt.yallamotor.com" + Val
			session["url"] = url
			fmt.Println(url)
			User := session["User"].(Person)
			User.CarPerson.Model = strings.ToLower(key)
			session["User"] = User
			fmt.Println("User.CarPerson.Type")
			fmt.Println(User.CarPerson.Type)
			fmt.Println("User.CarPerson.Model")
			fmt.Println(User.CarPerson.Model)
			fmt.Println(User)
			//retturn 3la to karims method
			//sessions[uuid]=session
			return VersionCar(session)
		}
	}
	session["ChatStage"] = ChatStage

	return "Laa sorry mosh 3andi !"

}

//gina NEW stage 1
func websrapLABEL(session Session/*uuid string*/) string {
	//var session Session = sessions[uuid] 

	url, ok := session["url"].(string)
	if !ok {
		session["url"] = "https://egypt.yallamotor.com/new-cars"
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error failed \"" + url)
		return "error"
	}
	ChatStage, ifAssigned := session["ChatStage"].(int)
	fmt.Println(ChatStage)
	fmt.Println(ifAssigned)
	if !ifAssigned {
		ChatStage = 0
	}

	ChatStage++
	session["ChatStage"] = ChatStage

	contentY := session["contentY"].([]string)

	stingToreturn := ""
	isAnchor := false
	//bytes, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("HTML:\n\n", string(bytes))
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return "error"
		case tt == html.StartTagToken:
			tag := z.Token()

			if tag.Data == "label" {
				isAnchor = true
			}
			//count := 5

			if isAnchor {

				if tt = z.Next(); tt == html.TextToken {

					label := z.Token().Data

					//label = label[0:5]
					if label != "Select Make" && label != "Select Model" && label != "Select Year" && label != "From (EGP)" && label != "To (EGP)" {

						contentY = append(contentY, label)
						stingToreturn += "   " + label+"\n"
						//fmt.Println(content)
						//fmt.Println(len(label))
						session["contentY"] = contentY
						//sessions[uuid]=session
						if len(contentY) == 49 {
							return stingToreturn
						}

					}

				}

			}

		}
	}
}

//gina New
func ModelCar(session Session/*uuid string*/) string {
//var session Session = sessions[uuid] 

	message2 := websrapLABEL(session/*uuid*/ )
	message := "Aayez anhy Type?\n"

	return message + "\n" + message2

}
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

//gina new to go to stage 2
func checkType(message string,session Session/*uuid string*/) string {
	//var session Session = sessions[uuid] 

	errorNoCar := session["errorNoCar"].(string)
	contentY := session["contentY"].([]string)
	url:=session["url"].(string)
	fmt.Println(contentY)
	errorNoCar = ""
	session["errorNoCar"] = errorNoCar

	ChatStage, ifAssigned := session["ChatStage"].(int)
	User := session["User"].(Person)
	fmt.Println(ChatStage)
	fmt.Println(ifAssigned)
	if !ifAssigned {
		ChatStage = 0
	}

	for _, a := range contentY {
		if strings.Contains(strings.ToLower(message), strings.ToLower(a)) {

			url = url + "" + "/" + "" + strings.ToLower(a)
			session["url"] = url
			fmt.Println(url)
			ChatStage++
			session["ChatStage"] = ChatStage
			User.CarPerson.Type = strings.ToLower(a)
			session["User"] = User
			//return on yasmin method
		//	sessions[uuid]=session
			return TypeCar(session/*uuid*/)
		}
	}
	return "laa sorry mosh 3andi"
}

//Kareem NEW VERsions to show to the user should be called by yasmins function

func versions(session Session) {

	CarVersionMap, ifAssigned := session["CarVersionMap"].(map[string]string)
	if !ifAssigned {
		session["CarVersionMap"] = make(map[string]string)
	}

	url, ok := session["url"].(string)
	if !ok {
		session["url"] = "https://egypt.yallamotor.com/new-cars"
	}

	resp, err := http.Get(url)

	ChatStage, ifAssigned := session["ChatStage"].(int)
	User := session["User"].(Person)

	fmt.Println(ChatStage)
	fmt.Println(ifAssigned)
	if !ifAssigned {
		ChatStage = 0
	}

	if err != nil {
		return
	}
	ChatStage++
	session["ChatStage"] = ChatStage

	//stingToreturn := ""
	doc := resp.Body

	defer resp.Body.Close()

	tokenizedDoc := html.NewTokenizer(doc)
	urlToCompare := "/new-cars/" + User.CarPerson.Type + "/"
	fmt.Println("urlToCompare")
	fmt.Println(urlToCompare)

	/*given the tokenizerDoc the html is tokenized by repeatedly calling .Next()
	which parses the nest token and return its 1) TYPE 2)ERROR */
	for {
		typeToken := tokenizedDoc.Next()

		switch typeToken {
		case html.ErrorToken:
			return

		case html.StartTagToken:
			{
				tag := tokenizedDoc.Token()

				isA := tag.Data == "a"
				if isA {
					fmt.Println(tag.Attr)

					for _, a := range tag.Attr {

						if a.Key == "href" && strings.Contains(a.Val, urlToCompare) {
							fmt.Println(tokenizedDoc.Next())
							mapKey := (tokenizedDoc.Token()).Data
							if strings.Contains(strings.ToLower(mapKey), strings.ToLower(User.CarPerson.Model)) && !strings.Contains(strings.ToLower(mapKey), "price in egypt") {
								CarVersionMap[mapKey] = a.Val

							}
							break
						}
					}

				}
			}
		}
	}
	session["CarVersionMap"] = CarVersionMap

}

//kareem function to fill the map
func VersionCar(session Session) string {
	versions(session)
	message := "Aayez anhy Version ya?\n"

	CarVersionMap, ifAssigned := session["CarVersionMap"].(map[string]string)
	if !ifAssigned {
		session["CarVersionMap"] = make(map[string]string)
	}

	for key, _ := range CarVersionMap {
		//fmt.Println(key)
		message += "" + key + "\n"
	}

	return message

}

//kareem NEW check on the version tht the user entered
func checkVersion(message string,session Session) string {

	CarVersionMap, ifAssigned := session["CarVersionMap"].(map[string]string)
	if !ifAssigned {
		session["CarVersionMap"] = make(map[string]string)
	}
	//fmt.Println(CarVersionMap)

	url, ok := session["url"].(string)
	if !ok {
		session["url"] = "https://egypt.yallamotor.com/new-cars"
	}

	ChatStage, ifAssigned := session["ChatStage"].(int)
	User := session["User"].(Person)

	//fmt.Println(ChatStage)
	//fmt.Println(ifAssigned)
	if !ifAssigned {
		ChatStage = 0
	}

	for key, Val := range CarVersionMap {
		fmt.Println(strings.ToLower(key))
		fmt.Println(strings.ToLower(message))
		if strings.Contains(strings.ToLower(message), strings.ToLower(key)) {
			ChatStage++
			session["ChatStage"] = ChatStage
			url = "https://egypt.yallamotor.com" + Val
			
			session["url"] = url
			User.CarPerson.Version = strings.ToLower(key)
			//fmt.Println(User.CarPerson.Version)
			//fmt.Println(User)
			session["User"] = User
			//retturn 3la to karims method
			return finalresult(session)
		}

	}
	return "mesh fahem bardo 3ayez anhy version meen dool"

}

//kareem final result which will return the url and the price
func finalresult(session Session) string {
	url, ok := session["url"].(string)
	if !ok {
		session["url"] = "https://egypt.yallamotor.com/new-cars"
	}

	price := ""
	resp, err := http.Get(url)

	ChatStage, ifAssigned := session["ChatStage"].(int)
	User := session["User"].(Person)

	fmt.Println(ChatStage)
	fmt.Println(ifAssigned)
	if !ifAssigned {
		ChatStage = 0
	}

	if err != nil {
		return "error" + url
	}
	defer resp.Body.Close()
	ChatStage++
	session["ChatStage"] = ChatStage

	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return "error"
		case tt == html.StartTagToken:
			tag := z.Token()

			isAnchor := tag.Data == "div"

			if isAnchor {

				tt = z.Next()
				if tt == html.TextToken {

					div := z.Token().Data

					ver := strings.Contains(div, "EGP")

					if ver {
						re := regexp.MustCompile("[0-9]+")
						x := re.FindAllString(div, -1) //list of the price

						for i := range x {
							price += x[i]
						}

						newFloat, _ := strconv.ParseFloat(price, 64)
						User.CarPerson.Price = newFloat //to set the price
						message := "price of the car is " + price + "\n" + "in order to view the car details check out this link  " + url
						session["User"] = User
						return message

					}

				}

			}

		}
	}
}

//yasmin USED
func TypeCarUsedFillMap(session Session) {

	isModel := false
	urlUsed := session["urlUsed"].(string)
	CarTypeUsedMap := session["CarTypeUsedMap"].(map[string]string)

	resp, err := http.Get(urlUsed)
	if err != nil {
		fmt.Println("error failed \"" + urlUsed)
		return
	}
	ChatStageUsed, ifAssigned := session["ChatStageUsed"].(int)
	fmt.Println(ChatStageUsed)
	fmt.Println(ifAssigned)
	if !ifAssigned {
		ChatStageUsed = 0
	}

	ChatStageUsed++
	session["ChatStageUsed"] = ChatStageUsed
	doc := resp.Body

	defer doc.Close()
	//toknize my doc
	tokenizedDoc := html.NewTokenizer(doc)

	/*given the tokenizerDoc the html is tokenized by repeatedly calling .Next()
	which parses the nest token and return its 1) TYPE 2)ERROR */
	for {
		typeToken := tokenizedDoc.Next()

		switch typeToken {
		case html.ErrorToken:
			return

		case html.StartTagToken:
			{
				tag := tokenizedDoc.Token()
				if tag.Data == "select" {
					for _, a := range tag.Attr {
						if a.Key == "id" && a.Val == "model" {
							isModel = true
							//fmt.Println(isModel)
							break
						} else {
							isModel = false

						}
					}
				}

				if tag.Data == "option" && isModel == true {
					for _, a := range tag.Attr {
						if a.Key == "value" && a.Val != "" {
						//	fmt.Println(a)
							tokenizedDoc.Next()
							mapKey := (tokenizedDoc.Token()).Data
							CarTypeUsedMap[mapKey] = a.Val
							break
						}

					}

				}

			}
		}
	}

	session["CarTypeUsedMap"] = CarTypeUsedMap
}

//Yasmin USED
func TypeCarUsed(session Session) string {
	TypeCarUsedFillMap(session)
	CarTypeUsedMap := session["CarTypeUsedMap"].(map[string]string)
	message := "Aayez anhy model?\n"

	for key, _ := range CarTypeUsedMap {
		message += "  " + key + " \n"

	}
	return message

}

//yasmin
func checkOnCarTypeUsed(message string ,session Session) string {
	CarTypeUsedMap := session["CarTypeUsedMap"].(map[string]string)

	for key, Val := range CarTypeUsedMap {
		//fmt.Println(strings.ToLower(key))
		//fmt.Println(strings.ToLower(message))
		if strings.Contains(strings.ToLower(message), strings.ToLower(key)) {

			ChatStageUsed, ifAssigned := session["ChatStageUsed"].(int)
			User := session["User"].(Person)

			fmt.Println(ChatStageUsed)
			fmt.Println(ifAssigned)
			if !ifAssigned {
				ChatStageUsed = 0
			}

			ChatStageUsed++
			session["ChatStageUsed"] = ChatStageUsed

			urlUsed := session["urlUsed"].(string)

			urlUsed += "" + "&model=" + "" + Val
			session["urlUsed"] = urlUsed

			fmt.Println(urlUsed)
			User.CarPerson.Model = key
			session["User"] = User
			fmt.Println(User)
			//retturn 3la to karims method
			return "oly min. range bat3k .... "
		}
	}

	return "Laa sorry mosh 3andi !"

}

//gina USED
func websrapLABELused(session Session) {
	// /search?make=22
	isMake := false
	urlUsed := session["urlUsed"].(string)
	CarMakeUsedMap := session["CarMakeUsedMap"].(map[string]string)

	resp, err := http.Get(urlUsed)
	if err != nil {
		fmt.Println("error failed \"" + urlUsed)
		return
	}
	ChatStageUsed, ifAssigned := session["ChatStageUsed"].(int)
	fmt.Println(ChatStageUsed)
	fmt.Println(ifAssigned)
	if !ifAssigned {
		ChatStageUsed = 0
	}

	ChatStageUsed++
	session["ChatStageUsed"] = ChatStageUsed
	doc := resp.Body

	defer doc.Close()
	//toknize my doc
	tokenizedDoc := html.NewTokenizer(doc)

	/*given the tokenizerDoc the html is tokenized by repeatedly calling .Next()
	which parses the nest token and return its 1) TYPE 2)ERROR */
	for {
		typeToken := tokenizedDoc.Next()

		switch typeToken {
		case html.ErrorToken:
			return

		case html.StartTagToken:
			{
				tag := tokenizedDoc.Token()
				if tag.Data == "select" {
					for _, a := range tag.Attr {
						if a.Key == "id" && a.Val == "make" {
							isMake = true
							fmt.Println(isMake)
							break
						} else {
							isMake = false

						}
					}
				}

				if tag.Data == "option" && isMake == true {
					for _, a := range tag.Attr {
						if a.Key == "value" && a.Val != "" {
							fmt.Println(a)
							tokenizedDoc.Next()
							mapKeyMake := (tokenizedDoc.Token()).Data
							CarMakeUsedMap[mapKeyMake] = a.Val
							break
						}

					}

				}

			}
		}
	}
	session["CarMakeUsedMap"] = CarMakeUsedMap
}

//gina USEd
func ModelCarUsed(session Session) string {
	websrapLABELused(session)
	message := "Aayez anhy Type?\n"
	CarMakeUsedMap := session["CarMakeUsedMap"].(map[string]string)

	for key, _ := range CarMakeUsedMap {
		message += "" + key + "\n"
	}
	return message

}

//gina Used
func checkMakeused(message string,session Session) string {
	CarMakeUsedMap := session["CarMakeUsedMap"].(map[string]string)

	for Key, Val := range CarMakeUsedMap {
		if strings.Contains(strings.ToLower(message), strings.ToLower(Key)) {
			ChatStageUsed, ifAssigned := session["ChatStageUsed"].(int)
			User := session["User"].(Person)

			fmt.Println(ChatStageUsed)
			fmt.Println(ifAssigned)
			if !ifAssigned {
				ChatStageUsed = 0
			}

			ChatStageUsed++
			session["ChatStageUsed"] = ChatStageUsed
			urlUsed := session["urlUsed"].(string)

			urlUsed += "" + "/search?make=" + "" + Val
			session["urlUsed"] = urlUsed
			fmt.Println(urlUsed)

			User.CarPerson.Type = Key
			session["User"] = User
			fmt.Println(User)
			return TypeCarUsed(session )

		}
	}
	return "Laa sorry mosh 3andi !"

}

//Kareem price Min USED
func checkOnMinPriceUsed(message string,session Session) string {
	User := session["User"].(Person)

	price := ""
	re := regexp.MustCompile("[0-9]+")
	x := re.FindAllString(message, -1) //list of the price

	for i := range x {
		price += x[i]
	}

	newInt, _ := strconv.ParseFloat(price, 64)
	User.CarPerson.PriceMin = newInt //to set the price min
	session["User"] = User
	ChatStageUsed, ifAssigned := session["ChatStageUsed"].(int)
	fmt.Println(ChatStageUsed)
	fmt.Println(ifAssigned)
	if !ifAssigned {
		ChatStageUsed = 0
	}

	ChatStageUsed++
	session["ChatStageUsed"] = ChatStageUsed
	urlUsed := session["urlUsed"].(string)

	urlUsed += "" + "&priceMin=" + "" + price
	session["urlUsed"] = urlUsed
	return "oli range max beta3k Kam??"

}

//Kareem Price Max Used
func checkOnMaxPriceUsed(message string,session Session) string {
	price := ""
	User := session["User"].(Person)

	re := regexp.MustCompile("[0-9]+")
	x := re.FindAllString(message, -1) //list of the price

	for i := range x {
		price += x[i]
	}

	newInt, _ := strconv.ParseFloat(price, 64)
	User.CarPerson.PriceMax = newInt //to set the price min
	session["User"] = User
	ChatStageUsed, ifAssigned := session["ChatStageUsed"].(int)
	fmt.Println(ChatStageUsed)
	fmt.Println(ifAssigned)
	if !ifAssigned {
		ChatStageUsed = 0
	}

	ChatStageUsed++
	session["ChatStageUsed"] = ChatStageUsed
	urlUsed := session["urlUsed"].(string)

	urlUsed += "" + "&priceMax=" + "" + price
	session["urlUsed"] = urlUsed

	return "oli men sant kan by  min??"
	//return "Inorder to view El cars check out this link "+""+urlUsed

}
func checkOnMaxYearUsed(message string,session Session) string {
	year := ""
	re := regexp.MustCompile("[0-9]+")
	User := session["User"].(Person)

	x := re.FindAllString(message, -1) //list of the price

	for i := range x {
		year += x[i]
	}

	newInt, _ := strconv.ParseInt(year, 0, 64)
	User.CarPerson.YearMax = newInt //to set the price min
	session["User"] = User
	ChatStageUsed, ifAssigned := session["ChatStageUsed"].(int)
	fmt.Println(ChatStageUsed)
	fmt.Println(ifAssigned)
	if !ifAssigned {
		ChatStageUsed = 0
	}

	ChatStageUsed++
	session["ChatStageUsed"] = ChatStageUsed
	urlUsed := session["urlUsed"].(string)

	urlUsed += "" + "&dateMax=" + "" + year
	session["urlUsed"] = urlUsed

	return "Inorder to view El cars check out this link " + "" + urlUsed

}

func checkOnMinYearUsed(message string,session Session) string {
	year := ""
	re := regexp.MustCompile("[0-9]+")
	User := session["User"].(Person)

	x := re.FindAllString(message, -1) //list of the price

	for i := range x {
		year += x[i]
	}

	newInt, _ := strconv.ParseInt(year, 0, 64)
	User.CarPerson.YearMin = newInt //to set the price min
	session["User"] = User
	ChatStageUsed, ifAssigned := session["ChatStageUsed"].(int)
	fmt.Println(ChatStageUsed)
	fmt.Println(ifAssigned)
	if !ifAssigned {
		ChatStageUsed = 0
	}

	ChatStageUsed++
	session["ChatStageUsed"] = ChatStageUsed
	urlUsed := session["urlUsed"].(string)

	urlUsed += "" + "&dateMin=" + "" + year
	session["urlUsed"] = urlUsed

	return "oli men sant kan by max??"

}

func handleChat(w http.ResponseWriter, r *http.Request) {
	// Make sure only POST requests are handled
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Make sure a UUID exists in the Authorization header
	uuid := r.Header.Get("Authorization")
	if uuid == "" {
		http.Error(w, "Missing or empty Authorization header.", http.StatusUnauthorized)
		return
	}

	// Make sure a session exists for the extracted UUID
	session, sessionFound := sessions[uuid]
	if !sessionFound {
		http.Error(w, fmt.Sprintf("No session found for: %v.", uuid), http.StatusUnauthorized)
		return
	}

	// Parse the JSON string in the body of the request
	data := JSON{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("Couldn't decode JSON: %v.", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Make sure a message key is defined in the body of the request
	_, messageFound := data["message"]
	if !messageFound {
		http.Error(w, "Missing message key in body.", http.StatusBadRequest)
		return
	}

	// Process the received message
	message, err := processor(session, data["message"].(string))
	if err != nil {
		http.Error(w, err.Error(), 422 /* http.StatusUnprocessableEntity */)
		return
	}

	writeJSON(w, JSON{
		"message": message,
	})

}

//yrab
// handle Handles /
func handle(w http.ResponseWriter, r *http.Request) {
	body :=
		"<!DOCTYPE html><html><head><title>Chatbot</title></head><body><pre style=\"font-family: monospace;\">\n" +
			"Available Routes:\n\n" +
			"  GET  /welcome -> handleWelcome\n" +
			"  POST /chat    -> handleChat\n" +
			"  GET  /        -> handle        (current)\n" +
			"</pre></body></html>"
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintln(w, body)
}

// Engage Gives control to the chatbot
func Engage(addr string) error {
	// HandleFuncs
	mux := http.NewServeMux()
	mux.HandleFunc("/welcome", withLog(handleWelcome))
	mux.HandleFunc("/chat", withLog(handleChat))
	mux.HandleFunc("/", withLog(handle))

	// Start the server
	return http.ListenAndServe(addr, cors.CORS(mux))
}

func main() {
	/*router := mux.NewRouter()

	router.HandleFunc("/welcome", WelcomeEndPoint).Methods("GET")

	router.HandleFunc("/chat", handleChat).Methods("POST")
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8099",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("SERVER ON RUNNING")
	log.Fatal(srv.ListenAndServe())*/

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	err := Engage(":" + port)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("SERVER ON RUNNING")
	}

}
