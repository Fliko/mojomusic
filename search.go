package main

import (
	"context"
	"fmt"
	"os/exec"

	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

type ytResult struct {
	id        string
	title     string
	channel   string
	thumbnail []byte
}

func handleError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %x", msg, err))
	}
}

// print results  of  search list
func printSearchListResults(response *youtube.SearchListResponse) ytResult {
	result := ytResult{}
	for _, item := range response.Items {
		fmt.Println(item)
		result = ytResult{
			id:      item.Id.VideoId,
			title:   item.Snippet.Title,
			channel: item.Snippet.ChannelTitle,
		}
	}
	url := fmt.Sprintf("https://youtu.be/%s", result.id)
	fmt.Println(url)
	out, err := exec.Command("youtube-dl", "-x", "--audio-quality", "0", "--exec", "./dca.sh", url).Output()
	fmt.Println(string(out))
	handleError(err, "youtube-dl")
	//thumbnail, err := exec.Command("youtube-dl", "--get-thumbnail", url).Output()
	//handleError(err, "get thumbnail")
	//result.thumbnail = thumbnail

	return result
}

func searchListByKeyword(service *youtube.Service, part string, maxResults int64, q string, typeArgument string) ytResult {
	call := service.Search.List(part)
	if maxResults != 0 {
		call = call.MaxResults(maxResults)
	}
	if q != "" {
		call = call.Q(q)
	}
	if typeArgument != "" {
		call = call.Type(typeArgument)
	}
	response, err := call.Do()
	handleError(err, "")
	if len(response.Items) == 0 {
		fmt.Println("Nothing Found")
		return ytResult{}
	}
	return printSearchListResults(response)
}

func ytSearch(search string) ytResult {
	// AIzaSyAdznmKD2m9a0VGXei2nRQO2nTA6ZhB8sY
	service, err := youtube.NewService(context.Background(), option.WithAPIKey("AIzaSyAdznmKD2m9a0VGXei2nRQO2nTA6ZhB8sY"))
	handleError(err, "new service")
	search = search + " hd hq"
	res := searchListByKeyword(service, "snippet", 1, search, "video")
	return res
}
