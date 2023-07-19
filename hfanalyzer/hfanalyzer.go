package hfanalyzer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	AgreementModelURL      = "https://huggingface.co/%s"
	AgreementDatasetURL    = "https://huggingface.co/datasets/%s"
	RawModelFileURL        = "https://huggingface.co/%s/raw/%s/%s"
	RawDatasetFileURL      = "https://huggingface.co/datasets/%s/raw/%s/%s"
	LfsModelResolverURL    = "https://huggingface.co/%s/resolve/%s/%s"
	LfsDatasetResolverURL  = "https://huggingface.co/datasets/%s/resolve/%s/%s"
	JsonModelsFileTreeURL  = "https://huggingface.co/api/models/%s/tree/%s/%s"
	JsonDatasetFileTreeURL = "https://huggingface.co/api/datasets/%s/tree/%s/%s"
)

var (
	RequiresAuth = false
	AuthToken    = ""
)

type hfmodel struct {
	Type        string `json:"type"`
	Oid         string `json:"oid"`
	Size        int    `json:"size"`
	Path        string `json:"path"`
	IsDirectory bool
	IsLFS       bool

	AppendedPath    string
	SkipDownloading bool
	FilterSkip      bool
	DownloadLink    string
	Lfs             *hflfs `json:"lfs,omitempty"`
}
type hflfs struct {
	Oid_SHA265  string `json:"oid"` // in lfs, oid is sha256 of the file
	Size        int64  `json:"size"`
	PointerSize int    `json:"pointerSize"`
}

func Analyze(ModelDatasetName string, Storage string, ModelBranch string, token string) error {
	getFileList(ModelDatasetName, Storage, ModelBranch, token)
	return nil
}

func getFileList(ModelDatasetName string, Storage string, ModelBranch string, token string) error {
	//we have to check first if model dataset name exists, then try join Storage/modelname, then try online
	foundLocal := false
	localPath := ""
	paths := []string{ModelDatasetName, path.Join(Storage, ModelDatasetName)} // Replace with the actual paths you want to check

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			foundLocal = true
			localPath = path
			break
		} else if os.IsNotExist(err) {
			continue
		} else { // another weird error
			return err
		}
	}

	if foundLocal { //pull the files locally
		fmt.Printf("\nFound Model Locally at %s", localPath)

		files, err := ioutil.ReadDir(localPath)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			fmt.Println(file.Name())
		}
	} else {
		fmt.Printf("\nModel Not found locally, try hugging face online")
		// initial parameters
		branch := ModelBranch
		JsonTreeVaraible := JsonModelsFileTreeURL
		if token != "" {
			RequiresAuth = true
			AuthToken = token
		}
		AgreementURL := fmt.Sprintf(AgreementModelURL, ModelDatasetName)
		JsonFileListURL := fmt.Sprintf(JsonTreeVaraible, ModelDatasetName, branch, "")
		//end of initial parameters

		client := &http.Client{}
		req, err := http.NewRequest("GET", JsonFileListURL, nil)
		if err != nil {
			return err
		}
		if RequiresAuth {
			// Set the authorization header with the Bearer token
			bearerToken := AuthToken
			req.Header.Add("Authorization", "Bearer "+bearerToken)
		}
		resp, err := client.Do(req)
		if err != nil {
			// fmt.Println("Error:", err)
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 401 && RequiresAuth == false {
			return fmt.Errorf("\nThis Repo requires access token, generate an access token form huggingface, and pass it using flag: -t TOKEN")
		}
		if resp.StatusCode == 403 {
			return fmt.Errorf("\nYou need to manually Accept the agreement for this model/dataset: %s on HuggingFace site, No bypass will be implemeted", AgreementURL)
		}
		// Read the response body into a byte slice
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// fmt.Println("Error:", err)
			return err

		}
		jsonFilesList := []hfmodel{}
		err = json.Unmarshal(content, &jsonFilesList)
		if err != nil {
			return err
		}
		fmt.Println(jsonFilesList)

	}
	return nil
}
