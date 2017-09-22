package main

import "encoding/json"
import "fmt"
import "io/ioutil"
import "net/http"
import "os"
import "time"

type Member struct {
	MembershipId   string
	DisplayName    string
	MembershipType int
}

type MemberResponse struct {
	Response []Member
}

type ProfileResponse struct {
	Response CharacterInfo
}

type CharacterInfo struct {
	CharacterProgressions struct {
		Data map[string]map[string]map[string]int
	}
	Characters struct {
		Data map[string]map[string]interface{}
	}
}

type BungieClient struct {
	Client *http.Client
	APIKey string
}

func NewBungieClient(apiKey string) *BungieClient {
	return &BungieClient{
		Client: &http.Client{Timeout: 60 * time.Second},
		APIKey: apiKey,
	}
}

func (b *BungieClient) NewRequest(method string, url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("%v", err)
	}
	req.Header.Set("X-API-Key", b.APIKey)
	return req
}

func (b *BungieClient) Do(req *http.Request) (*http.Response, error) {
	return b.Client.Do(req)
}

func (b *BungieClient) Get(url string) (*http.Response, error) {
	req := b.NewRequest("GET", url)
	return b.Client.Do(req)
}

const root = "https://bungie.net/Platform/Destiny2/"

func GetMembershipId(client *BungieClient, membershipType int, username string) string {
	membership_path := "SearchDestinyPlayer/%d/%s/"
	url := fmt.Sprintf(root+membership_path, membershipType, username)

	resp, err := client.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("%v", err)
	}

	var memRes MemberResponse
	err = json.NewDecoder(resp.Body).Decode(&memRes)
	if err != nil {
		fmt.Printf("%v", err)
	}
	member := memRes.Response[0]
	return member.MembershipId
}

func GetProfile(client *BungieClient, membershipType int, memberId string) {
	// We need URL parameters this time. components='200,202'
	profile_path := "%d/Profile/%s"
	url := fmt.Sprintf(root+profile_path, membershipType, memberId)

	// Because we have to add query information, we can't use the BungieClient.Get method.
	req := client.NewRequest("GET", url)
	q := req.URL.Query()
	q.Add("components", "200,202")
	req.URL.RawQuery = q.Encode()

	fmt.Printf("Query string: %s", req.URL.String())

	fmt.Printf("About to request\n")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Printf("About to decode\n")
	var profResponse ProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&profResponse)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("Raw body: %s\nEnd of Raw\n", bodyString)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%+v\nEnd of Parsed\n", profResponse)
}

func main() {
	api_key := os.Getenv("BUNGIE_API_KEY")
	client := NewBungieClient(api_key)
	mem_id := GetMembershipId(client, 2, "guubu")
	fmt.Printf("%s\n", mem_id)
	GetProfile(client, 2, mem_id)
}
