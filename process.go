package main

import(
    "net/http"
    "github.com/PuerkitoBio/goquery"
    "strings"
    "strconv"
    "log"
)

type File struct {
    head *goquery.Selection
    body *goquery.Selection
}

type Profile struct {
    Number int `json:"number"`
    Sex  string `json:"sex"`
    Nationality string `json:"nationality"`
    Age int `json:"age"`
    Path string `json: "path"`
    Date string `json: "date`
    Hospital string `json: "hospital"`
    Contacts int `json: "contacts"`
    Quarantine int `json: "quarantine"`
    Route []string `json:"route"`
}

func SelectorToData(f *File) Profile {
    profile := Profile{}
    f.head.Each(func(k int, s *goquery.Selection) {
        tag := s.Find("span").First().Text()
        context := s.Find("span").Last().Text()
        if strings.Contains(tag, "환자번호") {
            profile.Number, _ = strconv.Atoi(context)
        } else if strings.Contains(tag, "인적사항") {
            t1 := strings.Replace(context, "(", "", 1)
            t2 := strings.Replace(t1, ")", "", 1)
            t3 := strings.Replace(t2, ",", "", 1)
            t4 := strings.Replace(t3, "'", "", 1)
            text := strings.Split(t4, " ")
            profile.Sex = text[0]
            profile.Nationality = text[1]
            profile.Age, _ = strconv.Atoi(text[2])
        } else if strings.Contains(tag, "감염경로") {
            profile.Path = context
        } else if strings.Contains(tag, "확진일자") {
            t1 := strings.Join(strings.Split(context, "\n"), "")
            t2 := strings.Join(strings.Split(t1, "\t"), "")
            profile.Date = t2
        } else if strings.Contains(tag, "입원기관") {
            profile.Hospital = context
        } else if strings.Contains(tag, "접촉자수") {
            t1 := strings.Replace(context, "(", "", 1)
            t2 := strings.Replace(t1, ")", "", 1)
            t3 := strings.Split(t2, " ")
            profile.Contacts, _ = strconv.Atoi(t3[0])
            profile.Quarantine, _ = strconv.Atoi(t3[1])
        }
    })
    routes := make([]string, 0)
    f.body.Each(func(k int, s *goquery.Selection) {
        routes = append(routes, s.Text())
    })
    profile.Route = routes

    return profile
}

func GetPage(url string) *goquery.Document {
    res, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()
    if res.StatusCode != 200 {
        log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
    }
    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        log.Fatal(err)
    }
    return doc
}

func MakeFileArray(document *goquery.Document) []File {
    files := make([]File, 0)
    context := document.Find(".in_list")
    context.Each(func(i int, s *goquery.Selection) {
        s.Find(".onelist").Each(func(j int, s *goquery.Selection) {
            head := s.Find(".info_s li")
            body := s.Find(".info_mtxt .s_listin_dot li")
            file := File{head, body}
            files = append(files, file)
        })
    })
    return files
}

func GetData() []Profile {
    doc := GetPage("http://ncov.mohw.go.kr/bdBoardList.do?brdId=1&brdGubun=12")
    files := MakeFileArray(doc)
    output := make([]Profile, 0)
    for _, v := range files {
        o := SelectorToData(&v)
        output = append(output, o)
    }
    return output
}

func Filter(data []Profile, index string) []Profile {
    output := make([]Profile, 0)
    for _, v := range data {
        for _, r := range v.Route {
            if strings.Contains(r, index) {
                output = append(output, v)
                break;
            }
        }
    }

    return output
}
