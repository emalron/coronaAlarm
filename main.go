package main

import(
    "fmt"
    "log"
    "github.com/PuerkitoBio/goquery"
    "os"
    "net/http"
    "strings"
)

type File struct {
    head *goquery.Selection
    body *goquery.Selection
}

func PrintCon(f *File) {
    fmt.Println(f.head.Text())
    f.body.Each(func(i int, s *goquery.Selection) {
        fmt.Println(s.Text())
    })
}

type Profile struct {
    title string
    route []string
}

func SelectorToData(f *File) *Profile {
    profile := Profile{}
    profile.title = f.head.Text()
    profile.route = make([]string, 0)
    f.body.Each(func(i int, s *goquery.Selection) {
        s.Children().Each(func(n int, so *goquery.Selection) {
            profile.route = append(profile.route, so.Text())
        })
    })
    return &profile
}

func main() {
    args := os.Args[1:]

    // f, err := os.Open("a.html")
    res, err := http.Get("http://ncov.mohw.go.kr/bdBoardList.do?brdId=1&brdGubun=12")
    if err != nil {
        log.Fatal(err)
    }
    // defer f.Close()
    defer res.Body.Close()
    if res.StatusCode != 200 {
        log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
    }
    // doc, err := goquery.NewDocumentFromReader(f)
    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        log.Fatal(err)
    }

    indexes := make([]int, 0)
    body := doc.Find("h5").Siblings()
    body.Each(func(i int, s *goquery.Selection) {
        if s.HasClass("s_title_in3") {
            indexes = append(indexes, i)
        }
    })

    files := make([]File, 0)
    for i:=0; i<len(indexes)-1; i++ {
        file := File{body.Eq(indexes[i]), body.Slice(indexes[i]+1, indexes[i+1])}
        files = append(files, file)
    }
    fi := File{body.Eq(indexes[len(indexes)-1]), body.Slice(indexes[len(indexes)-1]+1, body.Size())}
    files = append(files, fi)

    output := make([]*Profile, 0)
    for _, v := range files {
        o := SelectorToData(&v)
        output = append(output, o)
    }

    for _, v := range output {
        for _, r := range v.route {
            if strings.Contains(r, args[0]) {
                fmt.Println(r)
            }
        }
    }
}
