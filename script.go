//Golang script made by GabrielCtz

package main

import (
    "fmt"
    "encoding/json"
    "net/http"
	"os"
)

type MojangData struct {
    UUID string `json:"uuid"`
}

type LabyNetData struct {
    User struct {
        UUID     string `json:"uuid"`
        Username string `json:"username"`
    } `json:"user"`
    NameHistory []struct {
        Username   string `json:"username"`
        ChangedAt  string `json:"changed_at"`
    } `json:"name_history"`
    Badges []struct {
        Name       string `json:"name"`
        ReceivedAt string `json:"received_at"`
    } `json:"badges"`
}

func UUID(pseudo string) (string, error) {
    url := "https://api.ashcon.app/mojang/v2/user/" + pseudo
    req, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer req.Body.Close()

    var data MojangData
    err = json.NewDecoder(req.Body).Decode(&data)
    if err != nil {
        return "", err
    }

    return data.UUID, nil
}

func Labynet(pseudo string) {
    UUID, err := UUID(pseudo)

    if err != nil {
        fmt.Println("[!] This user doesn't exists...")
        return
    }

    headers := map[string]string{
        "User-Agent": "Mozilla/5.0 (compatible; my-laby-discord-bot/1.0; +Niklas#4822)",
    }
    req, err := getWithHeaders("https://laby.net/api/user/"+UUID+"/get-snippet", headers)
    if err != nil {
        fmt.Println(err)
        return
    }

    var data LabyNetData
    err = json.NewDecoder(req.Body).Decode(&data)
    if err != nil {
        fmt.Println(err)
        return
    }

    user := data.User.UUID
    username := data.User.Username
	fmt.Println(" [\n | Results for", username)
    fmt.Printf(" | UUID : %s\n | Pseudo : %s\n", user, username)
	fmt.Println(" |")

    for i := 0; i < 25; i++ {
        headers := map[string]string{
            "User-Agent": "Mozilla/5.0 (compatible; my-laby-discord-bot/1.0; +Niklas#4822)",
        }
        res, err := getWithHeaders("https://laby.net/api/user/"+UUID+"/get-snippet", headers)
        if err != nil {
            fmt.Println(err)
            continue
        }

        var data LabyNetData
        err = json.NewDecoder(res.Body).Decode(&data)
        if err != nil {
            fmt.Println(err)
            continue
        }

        if i < len(data.NameHistory) {
            prevname := data.NameHistory[i].Username
            pdate := data.NameHistory[i].ChangedAt
            if pdate != "" {
				fmt.Printf(" | [%d] History Pseudo : %s ~ %s\n", i, prevname, pdate)            }
        }

        if i < len(data.Badges) {
            badges := data.Badges[i].Name
            btime := data.Badges[i].ReceivedAt
            if btime != "" {
                fmt.Printf(" | %s ~ %s\n", badges, btime)
            }
        }
    }
	fmt.Println(" [")
}

func getWithHeaders(url string, headers map[string]string) (*http.Response, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    for key, value := range headers {
        req.Header.Set(key, value)
    }

    return client.Do(req)
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("[-] go run script.go <pseudo>")
        return
    }

    pseudo := os.Args[1]
    Labynet(pseudo)
}
