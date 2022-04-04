package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello World!")
}

func createReqBody(fileBytes []byte) (string, io.Reader, error) {

	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf) // body writer

	fw1, _ := bw.CreateFormFile("file", "file")
	fw1.Write(fileBytes)

	bw.Close() //write the tail boundry
	return bw.FormDataContentType(), buf, nil
}

func ScanFile(w http.ResponseWriter, r *http.Request) {
	//fileType := r.PostFormValue("type")
	file, _, err := r.FormFile("file")
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}
	fmt.Println(string(fileBytes))

	client := &http.Client{}

	contType, reader, err := createReqBody(fileBytes)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", scanFileAPIURL, reader)
	req.Header.Add("Content-Type", contType)
	req.Header.Set("x-apikey", "11f5291a3ed4fd31bbe5cb0521e87272bf61a1054b056bf819ebfb0fc36931ee")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(400)
		msg := "Something wrong with response of scan api"
		w.Write([]byte(msg))
		fmt.Println(msg)
	}
	fmt.Println(string(body))

	scanURLAPIResultResultParam := ScanURLAPIResultResultParam{}
	if err := json.Unmarshal(body, &scanURLAPIResultResultParam); err != nil {
		w.WriteHeader(400)
		msg := "Something wrong with json"
		w.Write([]byte(msg))
		fmt.Println(msg)
	}

	result, err := getReport(scanURLAPIResultResultParam.Data.ID)
	if err != nil {
		w.WriteHeader(400)
		msg := "Something wrong with get report api " + err.Error()
		w.Write([]byte(msg))
		fmt.Println(msg)
	}
	fmt.Println(result)
	jsonResult, _ := json.Marshal(result)
	w.Write(jsonResult)
	return
}

func ScanURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In FUNC ScanURL")

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(400)
		msg := "Something wrong with read request"
		w.Write([]byte(msg))
		fmt.Println(msg)
		return
	}

	scanParam := ScanURLParam{}
	if err := json.Unmarshal(body, &scanParam); err != nil {
		w.WriteHeader(400)
		msg := "Something wrong with json"
		w.Write([]byte(msg))
		fmt.Println(msg)
		return
	}

	result, err := getResult(scanParam.Url)
	if err == nil {
		jsonResult, _ := json.Marshal(result)
		w.Write(jsonResult)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", scanURLAPIURL, strings.NewReader(fmt.Sprintf("url=%s", scanParam.Url)))
	if err != nil {
		w.WriteHeader(400)
		msg := "Something wrong with new request to scan api"
		w.Write([]byte(msg))
		fmt.Println(msg)
		return
	}
	req.PostForm = url.Values{}

	req.Header.Set("x-apikey", "11f5291a3ed4fd31bbe5cb0521e87272bf61a1054b056bf819ebfb0fc36931ee")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(400)
		msg := "Something wrong with response of scan api"
		w.Write([]byte(msg))
		fmt.Println(msg)
		return
	}
	fmt.Println(string(body))
	scanURLAPIResultResultParam := ScanURLAPIResultResultParam{}
	if err := json.Unmarshal(body, &scanURLAPIResultResultParam); err != nil {
		w.WriteHeader(400)
		msg := "Something wrong with json"
		w.Write([]byte(msg))
		fmt.Println(msg)
		return
	}

	reportResult, err := getReport(scanURLAPIResultResultParam.Data.ID)
	if err != nil {
		w.WriteHeader(400)
		msg := "Something wrong with get report api " + err.Error()
		w.Write([]byte(msg))
		fmt.Println(msg)
		return
	}
	reportResult.Data.Attributes.Results.URL = scanParam.Url
	fmt.Println(result)

	err = saveResult(reportResult)
	if err != nil {
		w.WriteHeader(400)
		msg := "Something wrong save result"
		w.Write([]byte(msg))
		fmt.Println(msg)
		return
	}

	jsonResult, _ := json.Marshal(reportResult)
	w.Write(jsonResult)
	return
	//w.Write([]byte(reportResult))

}

func getReport(reportID string) (*GetReportAPIResultResultParam, error) {
	client := &http.Client{}
	getReportAPIURLWithID := getReportAPIURL + reportID
	req, err := http.NewRequest("GET", getReportAPIURLWithID,
		strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	for {
		req.Header.Set("x-apikey", "11f5291a3ed4fd31bbe5cb0521e87272bf61a1054b056bf819ebfb0fc36931ee")
		resp, err := client.Do(req)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		param := GetReportAPIResultResultParam{}
		if err := json.Unmarshal(body, &param); err != nil {
			return nil, err
		}
		if param.Data.Attributes.Status == "completed" {
			//result, err := json.Marshal(param.Data.Attributes.Results)
			//if err != nil {
			//	return nil, err
			//}
			return &param, nil
		}
		time.Sleep(2 * time.Second)
	}
	return nil, nil
}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	data := checkHistoryCache()
	if data != nil {
		jsonData, _ := json.Marshal(data)
		w.Write(jsonData)
		return
	}
	data, _ = getAllData()
	cacheData(data)

	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
	return
}
