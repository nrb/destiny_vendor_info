package main

import "encoding/json"
import "fmt"
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

func (b *BungieClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("%v", err)
	}
	req.Header.Set("X-API-Key", b.APIKey)
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

	fmt.Printf("About to request\n")
	resp, err := client.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Printf("About to decode\n")
	var profResponse ProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&profResponse)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%+v", profResponse)
}

func main() {
	api_key := os.Getenv("BUNGIE_API_KEY")
	client := NewBungieClient(api_key)
	mem_id := GetMembershipId(client, 2, "guubu")
	fmt.Printf("%s\n", mem_id)
	GetProfile(client, 2, mem_id)
}
