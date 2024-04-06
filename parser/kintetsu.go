package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tomo0611/train-monitor/models"
)

func ParseData(data []byte) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	if err != nil {
		fmt.Println("html parse error")
		panic(err)
	}

	stations := make(map[string]models.Station)
	doc.Find(".stationinfo").Each(func(_ int, s *goquery.Selection) {
		station_id := s.Find(".routetype").Text()
		station_name := strings.TrimSpace(s.Find(".station").Text())
		station_name = strings.Replace(station_name, "　", "", 1)
		stations[station_id] = models.Station{
			Code:   station_id,
			Name:   station_name,
			NameEn: station_name, // 英語名はないのでとりあえず日本語名を入れておく
		}
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
			//fmt.Println("Type: " + trainType)
			//fmt.Println("Type Num: " + trainTypeNum)
			// 行き先
			station := strings.TrimSpace(s.Find(".station").Eq(i).Text())
			station = strings.TrimSpace(strings.Split(station, "\n")[strings.Count(station, "\n")])
			// 例)特急ひのとりでは、行先欄に「特急ひのとり」が含まれるというとんでもない実装をしているのでそれの対策をしている
			//fmt.Println("Bound for : " + station)
			// 列車番号
			train_no := s.Find(".trainno").Eq(i).Text()
			//fmt.Println("Train No: " + train_no)
			// 編成数
			train_length := s.Find(".trainno").Eq(i).Next().Text()
			train_length = strings.Replace(train_length, "両編成", "", 1)
			train_length_int, err := strconv.Atoi(train_length)
			if err != nil {
				fmt.Println("Train Length Parse Error")
				panic(err)
			}
			//fmt.Println("Train Length: " + train_length)
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
			//fmt.Println("Delay: " + strconv.Itoa(delay))
			// 現在位置
			current_location := s.Find(".border").Find(".rightbox").Eq(i).Text()
			//fmt.Println("Current Location: " + current_location + "\n")

			train := models.Train{
				No:           train_no,
				Pos:          current_location,
				Direction:    0,
				Nickname:     "",
				Type:         trainTypeNum,
				DisplayType:  trainType,
				Dest:         models.Dest{Text: station, Code: "", Line: "kin_n003"},
				Via:          "",
				DelaySeconds: delay * 60,
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
