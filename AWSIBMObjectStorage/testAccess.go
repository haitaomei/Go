package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/haitaomei/GoUtil/AWSIBMObjectStorage/awsauth"
)

var AccessKeyID = "a05d7fd864fe44d1a643c82e113f491d"
var SecretAccessKey = "d118d6b590a4ac15f0e83c545d307a48c952b904615f5476"

func main() {
	url := "https://s3.eu-geo.objectstorage.softlayer.net/"
	bucketName := "haitao-oddk"
	res, _ := listBucketContents(url, bucketName)
	objs := getObjectsWithinDir("tm15ccac9a6-306c-4e7c-bfd7-6a6c69e68556", res)

	for _, obj := range *objs {
		suc, err := deleteObject(url + bucketName + "/" + encode(obj))
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Delete", obj, "success=", suc)
	}
}

func encode(uri string) string {
	return strings.Replace(url.PathEscape(uri), "%2F", "/", -1)
}

func listBucketContents(url string, bucketName string) (string, error) {
	url += bucketName + "/"
	client := new(http.Client)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	awsauth.Sign(req, awsauth.Credentials{
		AccessKeyID:     AccessKeyID,
		SecretAccessKey: SecretAccessKey,
	})

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		bodyString := string(bodyBytes)
		return bodyString, nil
	}
	return "", nil
}

func getObjectsWithinDir(dir string, res string) *[]string {
	i := 0
	j := 0
	var objs []string
	for i < len(res)-5 {
		if res[i:i+5] == "<Key>" {
			for j = i + 5; j < len(res)-6; j++ {
				if res[j:j+6] == "</Key>" {
					break
				}
			}
			obj := res[i+5 : j]
			if obj[0:len(dir)] == dir {
				objs = append(objs, obj)
			}
		}
		i++
	}
	return &objs
}

func deleteObject(uri string) (bool, error) {
	client := new(http.Client)
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return false, err
	}
	awsauth.Sign(req, awsauth.Credentials{
		AccessKeyID:     AccessKeyID,
		SecretAccessKey: SecretAccessKey,
	})
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		// fmt.Println("Delete", uri, "sucessed")
		return true, nil
	} else {
		// bodyBytes, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// bodyString := string(bodyBytes)
		// fmt.Println(bodyString)
		return false, err
	}
	return true, nil
}
