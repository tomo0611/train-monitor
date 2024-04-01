package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*
	1 特急
	2 快速急行
	3 急行
	4 準急
	5 区間準急
	6 普通
	7 区間快速
	8 区間急行
	9 区間快速急行
	10  一般貸切
	11  特急貸切
*/

type Dest struct {
	Text string `json:"text"`
	Code string `json:"code"`
	Line string `json:"line"`
}

type Train struct {
	No           string `json:"no"`
	Pos          string `json:"pos"`
	Direction    int    `json:"direction"`
	Nickname     string `json:"nickname"`
	Type         string `json:"type"`
	DisplayType  string `json:"displayType"`
	Dest         Dest   `json:"dest"`
	Via          string `json:"via"`
	DelayMinutes int    `json:"delayMinutes"`
	TypeChange   string `json:"typeChange"`
	NumberOfCars int    `json:"numberOfCars"`
}

func main() {

	//【A】近鉄奈良線 (大阪難波～近鉄奈良)
	resp, err := http.Get("https://tid.kintetsu.co.jp/LocationHtml/trainlocationinfo01.html?innerLink=true")
	if err != nil {
		fmt.Println("http get error")
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("html parse error")
		panic(err)
	}

	stations := make(map[string]string)
	doc.Find(".stationinfo").Each(func(_ int, s *goquery.Selection) {
		station_id := s.Find(".routetype").Text()
		station_name := strings.TrimSpace(s.Find(".station").Text())
		station_name = strings.Replace(station_name, "　", "", 1)
		stations[station_id] = station_name
	})

	doc.Find(".scrollArea-box").Each(func(_ int, s *goquery.Selection) {
		trainNum := s.Find(".image").Length()
		for i := 0; i < trainNum; i++ {
			// 列車種別
			trainType := s.Find(".traintype").Eq(i).Text()
			trainTypeNum := "0"
			sss, exists := s.Find(".traintype").Eq(i).Attr("class")
			if exists {
				trainTypeNum = strings.Replace(sss, "traintype traintype", "", 1)
			}
			fmt.Println("Type: " + trainType)
			fmt.Println("Type Num: " + trainTypeNum)
			// 行き先
			station := strings.TrimSpace(s.Find(".station").Eq(i).Text())
			station = strings.TrimSpace(strings.Split(station, "\n")[strings.Count(station, "\n")])
			// 例)特急ひのとりでは、行先欄に「特急ひのとり」が含まれるというとんでもない実装をしているのでそれの対策をしている
			fmt.Println("Bound for : " + station)
			// 列車番号
			train_no := s.Find(".trainno").Eq(i).Text()
			fmt.Println("Train No: " + train_no)
			// 編成数
			train_length := s.Find(".trainno").Eq(i).Next().Text()
			train_length = strings.Replace(train_length, "両編成", "", 1)
			train_length_int, err := strconv.Atoi(train_length)
			if err != nil {
				fmt.Println("Train Length Parse Error")
				panic(err)
			}
			fmt.Println("Train Length: " + train_length)
			// 遅れ状況
			delay_str := s.Find(".delay").Find(".rightbox").Eq(i).Text()
			delay := 0
			if delay_str == "遅れなし" {
				delay = 0
			} else {
				delay, err = strconv.Atoi(strings.Replace(strings.Replace(delay_str, "遅れ：", "", 1), "分", "", 1))
				if err != nil {
					fmt.Println("Delay Parse Error")
					panic(err)
				}
			}
			fmt.Println("Delay: " + strconv.Itoa(delay))
			// 現在位置
			current_location := s.Find(".border").Find(".rightbox").Eq(i).Text()
			fmt.Println("Current Location: " + current_location + "\n")

			train := Train{
				No:           train_no,
				Pos:          current_location,
				Direction:    0,
				Nickname:     "",
				Type:         trainTypeNum,
				DisplayType:  trainType,
				Dest:         Dest{Text: station, Code: "", Line: "kin_n003"},
				Via:          "",
				DelayMinutes: delay,
				TypeChange:   "",
				NumberOfCars: train_length_int,
			}
			b, err := json.Marshal(train)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(b))
		}
	})
}
