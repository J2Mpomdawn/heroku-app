package ts

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"net/mail"
	"net/smtp"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"io/ioutil"
	"path/filepath"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/scorredoira/email"
)

type postime struct {
	h rune
	m bool
}

type post struct {
	//name...イベント名(開始日をもとに作成した番号を名前とする)
	Name string `gorm:"type:smallint unsigned"`
	//num...投稿時間
	Num int `gorm:"type:smallint unsigned"`
	//100~50000位
	One, Two, Three, Four, Five, Six string `gorm:"type:mediumint unsigned"`
}

type period struct {
	//番号
	Num string `gorm:"type:smallint unsigned PRIMARY KEY"`
	//name...イベント名
	Name string `gorm:"type:varchar(60)"`
	//period...開催期間
	Period string `gorm:"type:char(21)"`
	//times...データの数
	Times int `gorm:"type:smallint unsigned"`
	//done...整理状況
	//0:未使用, 1:記録済(不足), 2:記録済(充足)
	Done int `gorm:"type:tinyint default 0"`
}

var (
	bl                    bool
	m, fwt, twt           int
	poh                   rune
	tweet, eventname, pri string
	start, end            time.Time
	err                   error
	dates, twlog          []string = make([]string, 14), make([]string, 4)
	ru, r                 []rune
	krn, kk, kig          []int = make([]int, 6), make([]int, 6), make([]int, 6)
	poti                  postime
	pt                    post
)

//ツイート取得準備
func setconf() (api *anaconda.TwitterApi, v url.Values) {

	//apiの設定
	//本番か開発かで設定を変える
	/*if os.Getenv("PORT") == "" {
		if err = godotenv.Load("dev.env"); err != nil {
			fmt.Printf("--couldn't load env---\n%v\n", err)
		}
	}*/
	anaconda.SetConsumerKey(os.Getenv("ConsumerKey"))
	anaconda.SetConsumerSecret(os.Getenv("ConsumerSecret"))
	api = anaconda.NewTwitterApi(os.Getenv("AccessToken"), os.Getenv("AccessTokenSecret"))

	//とるのはテキストBOTさんの投稿、上から4個
	v = url.Values{}
	v.Set("screen_name", "imas_ml_td_t")
	v.Set("count", "4")
	return
}

//データベースの準備
func gormcore() *gorm.DB {

	//mysqlの設定
	//本番か開発かで設定を変える
	protocol := "tcp(" + os.Getenv("DB_HOSTNAME") + ":3306)"
	if os.Getenv("PORT") == "8080" {
		protocol = ""
	}
	db, err := gorm.Open("mysql", os.Getenv("DB_USERNAME")+
		":"+os.Getenv("DB_PASSWORD")+"@"+protocol+"/"+os.Getenv("DB_NAME"))
	if err != nil {
		panic(err.Error())
	}
	return db
}

//通知の処理
func tuti() {
label:
	for a, v := range ru {
		if a == len(ru)-14 {
			break
		}

		//"MILLION LIVE W"か"ミリコレ"ならm=8
		if reflect.DeepEqual(ru[a:a+4], []rune{12511, 12522, 12467, 12524}) || reflect.DeepEqual(ru[a:a+14], []rune{77, 73, 76, 76, 73, 79, 78, 32, 76, 73, 86, 69, 32, 87}) {
			m = 8
		}
		switch v {

		//"「"の場所を記録
		case 12300:
			poti.h = rune(a)

			//"「"と"」"の場所をもとにイベントの名前を取得
		case 12301:

			//既に記録されてたら以降の処理はスルー
			if string(ru[poti.h+1:a]) == eventname {
				continue

				//新イベならeventnameを上書き
			} else {
				eventname = string(ru[poti.h+1 : a])
			}

			//mをもとに処理を分岐
			if m == 8 {
				m = 9
			} else {
				m = 10
			}

			//":"なら色々やる
		case 58:

			//RTの場合は一回目の":"をスルー
			if poti.m {
				poti.m = false
				continue
			}

			//最終日を取得
			end, err = time.Parse("2006/1/2 15:04", string(ru[a+7:a+23]))

			if err != nil {
				fmt.Printf("--couldn't set LD---\n%v\n", err)
			}

			//"ミリコレ"か"WORKING"なら最終日まで待つ
			if m == 9 {
				ato := -time.Since(end.Add(-539 * time.Minute))
				fmt.Printf("イベント名「%s」\n", eventname)
				fmt.Printf("次のイベントまであと %v\n", ato)
				/*
					ほんまに実装するときはatoとfmt二つ消して
					time.sleep(ato)を
					time.Sleep(-time.Since(end.Add(-539 * time.Minute)))
					↑こうする

					なんかこのままでもいいような気がしてきた
				*/
				time.Sleep(ato)
				m = 11
				return

				//それ以外のイベントのとき
			} else if m == 10 {

				//開始日を取得
				start, err = time.Parse("2006/1/2 15", string(ru[a-13:a]))
				if err != nil {
					fmt.Printf("--couldn't set SD---\n%v\n", err)
				}

				//イベントの日付を記録
				dates = make([]string, 13)
				bb := 0
				for ; bb <= end.Day()-start.Day(); bb++ {
					day := start.AddDate(0, 0, bb).Day()
					day1 := day/10 + 48
					day2 := day%10 + 48
					if day < 10 {
						day1 = 48
						day2 = day + 48
					}
					dates[bb] = string([]byte{byte(day1), byte(day2)})
				}

				//fwtとtwtを記録
				fwt = start.Hour() << 1
				twt = ((bb-1)*24+end.Add(1*time.Minute).Hour())<<1 - fwt

				//2019/12/01を最初の日とする
				birth, err := time.Parse("2006/01/02", "2019/12/01")
				if err != nil {
					fmt.Printf("--couldn't set BD---\n%v\n", err)
				}

				//"イベント開始日は最初の日から何日目か"をprikeyに
				pri = strconv.FormatFloat((start.Add(-time.Duration(start.Hour())*time.Hour).Sub(birth) / 24).Hours(), 'f', 0, 64)

				//イベントの名前と期間とプリキーをセット
				pe := period{}
				pe.Num = pri
				pe.Name = eventname
				pe.Period = start.Format("2006/01/02") + "~" + end.Format("2006/01/02")
				pe.Times = twt

				//データベースに登録
				db := gormcore()
				//db.Table("list").CreateTable(&period{})
				db.Table("list").Save(&pe)
				db.Close()
			}

			break label
		}
	}
}

//ツイートの取得と加工
func gettweets(api *anaconda.TwitterApi, v url.Values) {

	//イベントの日付が記録されてなかったら
	if dates[0] == "" {

		//通知BOTさんの最新の投稿を取得
		v.Set("screen_name", "imas_ml_td_i")
		v.Set("count", "2")
		twii, err := api.GetUserTimeline(v)
		if err != nil {
			fmt.Printf("--couldn't get tweets---\n%v\n", err)
		}

		//改行コードを"\n"に統一
		//「折り返し」が入ってたら一個前の投稿を使う
		if strings.Contains(twii[0].FullText, "折り返し") {
			tweet = strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(twii[1].FullText)
		} else {
			tweet = strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(twii[0].FullText)
		}

		//rune配列に変換して通知の処理
		ru = []rune(tweet)
		poti.m = false
		tuti()

		//url.Valuesの設定を元に戻す
		v.Set("screen_name", "imas_ml_td_t")
		v.Set("count", "4")
	}

	//DB用意
	db := gormcore()
	//db.Table("datas").CreateTable(&post{})
	defer db.Close()

	//ツイートの取得
	tweets, err := api.GetUserTimeline(v)
	if err != nil {
		fmt.Printf("--couldn't get tweets---\n%v\n", err)
		return
	}

	//古い投稿から処理したいからカウントダウンでループ
	for a := 3; a >= 0; a-- {
		poti.m = false

		//取得したツイートがtwlogに記録されてるかどうか
		for b := 0; b < 4; b++ {
			if tweets[a].FullText == twlog[b] {
				poti.m = true
				break
			}
		}

		//ツイートをtwlogに記録
		twlog[3-a] = tweets[a].FullText

		//すでに記録されてたら以降の処理をスルー
		if poti.m {
			continue
		}

		//改行の統一、rune変換
		tweet = strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(tweets[a].FullText)
		ru = []rune(tweet)

		//RTかどうかで処理を分ける
		//RTじゃなかったら
		if tweets[a].RetweetedStatus == nil {
			m = 0
			kurai := 0
			for b, w := range ru {
				switch w {

				//"位"の数を数える
				case 20301:
					kurai++

				//":","(","\n"の場所を記録(5個ずつ)
				case 58:
					if m > 5 {
						continue
					}
					krn[m] = b
				case 40:
					if m > 5 {
						continue
					}
					kk[m] = b
				case 10:
					if m > 5 {
						continue
					}
					kig[m] = b
					m++

					//"#"の場所をもとに色々やる
				case 35:

					//何時のボーダーか取得
					poti.h = ru[b-6]*10 + ru[b-5] - 528

					//30分ならtrue
					if ru[b-3] == 48 {
						poti.m = false
					} else {
						poti.m = true
					}

					//poti.hを2倍にして、30分ならさらに1を足す
					poti.h <<= 1
					if poti.m {
						poti.h++
					}

					//最終結果の場合はmに1を足す(m=7)
					if ru[b-2] == 26524 {
						m++
					}

					//"位"の数だけループする
					for c := 0; c < kurai; c++ {

						//":"と"("もしくは"\n"からボーダーの数値を取得
						if m == 7 {
							r = []rune(string(ru[krn[c]+2 : kig[c]]))
						} else {
							r = []rune(string(ru[krn[c]+2 : kk[c]-1]))
						}

						//","を消す
						for l, x := range r {
							if x == 44 {
								r = append(r[:l], r[l+1:]...)
							}
						}

						//ptにイベント名と整形した数値をセット
						pt.Name = pri
						switch c {
						case 0:
							pt.One = string(r)
						case 1:
							pt.Two = string(r)
						case 2:
							pt.Three = string(r)
						case 3:
							pt.Four = string(r)
						case 4:
							pt.Five = string(r)
						case 5:
							pt.Six = string(r)
						}
						if err != nil {
							fmt.Printf("--couldn't convert AtoI---\n%v\n", err)
						}
					}

					//参加者不足カバー
					for c := 0; c < (6 - kurai); c++ {
						switch c {
						case 0:
							pt.Six = "0"
						case 1:
							pt.Five = "0"
						case 2:
							pt.Four = "0"
						case 3:
							pt.Three = "0"
						case 4:
							pt.Two = "0"
						case 5:
							pt.One = "0"
						}
					}

					//何日のボーダーか取得
					tweet = string(ru[b-9 : b-7])

					//最終日やったら
					if m == 7 {
						pt.Num = twt
						db.Table("datas").Save(&pt)

						//これ以降の処理はもういい
						return
					}

					//そうじゃなかったら適切な場所に記録
					for d, x := range dates {
						if tweet == x {
							pt.Num = 48*d + int(poti.h) - fwt
							db.Table("datas").Save(&pt)
						}
					}
				}
			}
			//RTやったら通知の処置
		} else {
			poti.m = true
			tuti()
		}
	}
	/*
		何か見たいものがあったらこの下に書く
	*/
}

//記録されてるイベントのリストを作成
func eventlist() (list string) {

	//記録されてるイベントの情報を取得
	pes := []period{}
	db := gormcore()
	db.Table("list").Select("num,name,period").Order("num desc").Where("num>2").Find(&pes)
	db.Close()

	//byte配列用意
	var listb []byte

	//listbにイベント情報を詰め込んでいく
	for _, v := range pes {

		//まずはnum
		b := []byte(v.Num)
		for _, w := range b {
			listb = append(listb, w)
		}

		//"//"を追加
		listb = append(listb, 47, 47)

		//次はname
		b = []byte(v.Name)
		for _, w := range b {
			listb = append(listb, w)
		}

		//"//"を追加
		listb = append(listb, 47, 47)

		//最後はperiod
		b = []byte(v.Period)
		for _, w := range b {
			listb = append(listb, w)
		}

		//",,"を追加
		listb = append(listb, 44, 44)
	}

	//完成したlistbを後ろ2文字を削って文字列にして返す
	return string(listb[:len(listb)-2])
}

//ローディングルーレット
func roulette() (src, alt string) {
	rand.Seed(time.Now().UnixNano())
	//０ ～ ("Intnの引数"-１) の整数
	n := rand.Intn(1)
	src = "/static/img/hinata" + strconv.Itoa(n) + ".gif"
	switch n {
	case 0:
		alt = "ローディング用のGIF。ラブリーフルーティアひなたがぴょんぴょん跳ねてる"
	case 1:
		//衣装の数だけ増やす
		alt = ""
	}
	return
}

//選択したイベントのグラフの素を作成
func graphinfo(ev, ra, he string) (sc int, ele string) {
	//ev...選択したイベント
	//ra...指定した順位
	//he...#graphの高さ

	//選択したイベントをコンマ区切りで取り出す
	se := strings.Split(ev, ",")

	//DB用意
	db := gormcore()
	defer db.Close()
	spt := []post{}
	ls := period{}

	//最終結果(達)のなかで一番でかいのを記録
	rm := float64(0)
	if len(se) != 0 {

		//比べる用
		comp := float64(0)

		for _, v := range se {

			//pt初期化
			pt = post{}

			//numの最大値を取得
			db.Table("datas").Select("max(num)as one").Where("name=?", v).Find(&pt)

			//300か348か396か
			db.Table("list").Select("times").Where("num=?", v).Find(&ls)
			name := "1"
			num := "348"
			switch ls.Times {
			case 300:
				name = "0"
				num = "300"
			case 396:
				name = "2"
				num = "396"
			}

			//そのnumの平均と比べて最終日予想
			db.Table("datas").Select("(select " + ra + " from datas where name=" + name + " and num=" + num + ")*((select " + ra + " from datas where name=" + v + " and num=" + pt.One + ")/(select " + ra + " from datas where name=" + name + " and num=" + pt.One + "))as two").Find(&pt)

			//開催期間が396より大きいとき
			if ls.Times > 396 {
				pt.Two = "6137039"
			}

			//float64に変換
			comp, err = strconv.ParseFloat(pt.Two, 64)
			if err != nil {
				fmt.Printf("--couldn't convert AtoF---\n%v\n", err)
			}

			//記録済みのrmより大きかったら上書き
			if comp > rm {
				rm = comp
			}
		}
	}

	//0のままやったら平均値を設定
	if rm == 0 {
		db.Table("datas").Select(ra + " as one").Where("name=0 and num=348").Find(&pt)
		rm, err = strconv.ParseFloat(pt.One, 64)
		if err != nil {
			fmt.Printf("--couldn't convert AtoF---\n%v\n", err)
		}
	}

	//y軸目盛用の数字
	/*
		目盛りは6個
		→sc*6
	*/
	sc = int(rm/300000+0.9) * 50000

	//#graphの高さをfloat64に変換
	hei, err := strconv.ParseFloat(he, 64)
	if err != nil {
		fmt.Printf("--couldn't convert AtoF---\n%v\n", err)
	}

	//割る
	rm = hei / rm

	//select句作成
	sel := ra + " as one,num"

	//被り確認用
	num := 0

	//イベント達の記録をまとめる箱
	boxs := make([][]float64, len(se))

	//グラフの文字列用バイト配列
	gra := make([]byte, 0, 32768)

	//イベント毎の処理
	for a, v := range se {

		//記録を取得
		db.Table("datas").Select(sel).Where("name=?", v).Order("num").Find(&spt)

		//num初期化
		num = -1

		//イベントの記録を入れる箱
		db.Table("list").Select("times").Where("num=?", v).Find(&ls)
		box := make([]float64, ls.Times+1)

		//「M0 0 」
		gra = append(gra, 77, 48, 32, 48, 32)

		//詰め込む
		for _, w := range spt {

			//指定したイベントの指定した順位のボーダー
			border, err := strconv.ParseFloat(w.One, 64)
			if err != nil {
				fmt.Printf("--couldn't convert AtoI---\n%v\n", err)
			}

			//被ってたらDBから削除
			if w.Num == num {
				db.Table("datas").Where("name=? and num=?", v, num).Limit(1).Delete(&w)
				continue
			}

			//記録
			box[w.Num] = border
			num = w.Num

			//graにデータを書き込んでいく
			//「L」
			gra = append(gra, 76)
			//"w.Num"
			gra = append(gra, []byte(strconv.FormatFloat(float64(w.Num)*2.51417108, 'f', 6, 64))...)
			//「 」
			gra = append(gra, 32)
			//"border"
			gra = append(gra, []byte(strconv.FormatFloat(border*rm, 'f', 6, 64))...)
			//「 」
			gra = append(gra, 32)
		}

		//boxを詰める
		boxs[a] = box

		//最後のときはコンマ追加しない
		if a == len(se)-1 {
			continue
		}

		//グラフのpathをコンマで区切る
		gra = append(gra, 44)
	}

	ele = string(gra)

	return
}

//xlsxに記録
func backup() (ok bool) {

	//db用意
	db := gormcore()
	//終わったら閉じる
	defer db.Close()

	//リストと記録
	lists := []period{}
	datas := []post{}

	//平均以外のリスト取得
	if err = db.Table("list").Select("num,times,done").Where("num>2 and done<>2").Find(&lists).Error; err != nil {
		fmt.Printf("--couldn't get lists---\n%v\n", err)
		return
	}

	//ファイル開く
	xf, err := excelize.OpenFile("datas.xlsx", excelize.Options{Password: os.Getenv("XlPassword")})
	if err != nil {
		fmt.Printf("--couldn't open the file---\n%v\n", err)
		return
	}

	//イベント毎に記録
	for _, list := range lists {

		//記録する列
		colname := "A"

		//未記録
		if list.Done == 0 {
			for a := 1; a < 7; a++ {
				nt := strconv.Itoa(list.Times*10 + a)
				if err = xf.InsertCol(nt, "A"); err != nil {
					fmt.Printf("--couldn't insert the col---\n%v\n", err)
					return
				}
				if err = xf.SetCellStr(nt, "A1", list.Num); err != nil {
					fmt.Printf("--couldn't set the str---\n%v\n", err)
					return
				}
			}
			db.Table("list").Where("num=?", list.Num).Update("done", 1)

			//不揃い記録済み
		} else if list.Done == 1 {

			//どこの列か探す
			sc, err := xf.SearchSheet(strconv.Itoa(list.Times*10+1), list.Num)
			if err != nil {
				fmt.Printf("--couldn't search the val---\n%v\n", err)
				return
			}
			colname, _, err = excelize.SplitCellName(sc[0])
			if err != nil {
				fmt.Printf("--couldn't divide the name---\n%v\n", err)
				return
			}
		} else {
			continue
		}

		//記録を取得して適当なマスに記録
		/*2秒待ってみる*/
		time.Sleep(2 * time.Second)
		db.Raw("select*from(select*from datas where name=? order by num)as A group by num", list.Num).Scan(&datas)

		for _, data := range datas {
			for a := 1; a < 7; a++ {
				nt := strconv.Itoa(list.Times*10 + a)
				nvs := ""
				switch a {
				case 1:
					nvs = data.One
				case 2:
					nvs = data.Two
				case 3:
					nvs = data.Three
				case 4:
					nvs = data.Four
				case 5:
					nvs = data.Five
				case 6:
					nvs = data.Six
				}
				nvi, err := strconv.Atoi(nvs)
				if err != nil {
					fmt.Printf("--couldn't convert AtoI---\n%v\n", err)
					return
				}
				if err = xf.SetCellInt(nt, colname+strconv.Itoa(data.Num+1), nvi); err != nil {
					fmt.Printf("--couldn't set the int---\n%v\n", err)
					return
				}
			}
		}

		//全部揃ったらもう使わん
		if len(datas) == list.Times {
			db.Table("list").Where("num=?", list.Num).Update("done", 2)
		}
	}

	//平均を計算して更新
	for a := 0; a < 3; a++ {
		e := 0
		switch a {
		case 0:
			pt.Name = "0"
			e = 302
		case 1:
			pt.Name = "1"
			e = 350
		case 2:
			pt.Name = "2"
			e = 398
		}
		for b := 2; b < e; b++ {
			pt.Num = b - 1
			for c := 1; c < 7; c++ {
				colname, err := excelize.ColumnNumberToName(6*a + c)
				if err != nil {
					fmt.Printf("--couldn't convert Num2Nam---\n%v\n", err)
					return
				}
				res, err := xf.CalcCellValue("a", colname+strconv.Itoa(b))
				if err != nil {
					if err.Error() == "#DIV/0!" {
						res = "0"
					} else {
						fmt.Printf("--couldn't calc formula---\n%v\n", err)
						return
					}
				}
				switch c {
				case 1:
					pt.One = res
				case 2:
					pt.Two = res
				case 3:
					pt.Three = res
				case 4:
					pt.Four = res
				case 5:
					pt.Five = res
				case 6:
					pt.Six = res
				}
			}
			db.Table("datas").Where("name=? and num=?", pt.Name, pt.Num).Update(&pt)
		}
	}

	//セーブ
	if err = xf.SaveAs("./tmp/datas.xlsx"); err != nil {
		fmt.Printf("--couldn't save the file---\n%v\n", err)
		return
	}
	return true
}

//メール送信
func send(tion bool) (ok bool) {

	//sub,body,from,to,server等設定
	subject := time.Now().Format("2006/1/2")
	body := func(opera bool) (style string) {
		if opera {
			style = "manual"
		} else {
			style = "auto"
		}
		return
	}(tion)
	from := "imymemine"
	fromAdd := os.Getenv("M_from")
	pass := os.Getenv("M_frompass")
	to0 := os.Getenv("M_to")
	servername := os.Getenv("M_servername")
	host := strings.Split(servername, ":")[0]

	//再構築失敗時
	if bl {
		subject = "えまーじぇんしー"
		body = "再構築になんか問題あり"
		from = "uuruurs"
	}

	//メッセージ作成
	ml := email.NewMessage(subject, body)
	ml.From = mail.Address{Name: from, Address: fromAdd}
	ml.To = []string{to0}
	if !bl {
		if err = ml.Attach("./tmp/datas.xlsx"); err != nil {
			fmt.Printf("--couldn't attach the file---\n%v\n", err)
			return
		}
	}

	//通信準備
	auth := smtp.PlainAuth("", fromAdd, pass, host)
	tlsconfig := &tls.Config{InsecureSkipVerify: true, ServerName: host}
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		fmt.Printf("--couldn't tls connection---\n%v\n", err)
		return
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Printf("--couldn't create a new client---\n%v\n", err)
		return
	}
	if err = c.Auth(auth); err != nil {
		fmt.Printf("--couldn't authenticate the client---\n%v\n", err)
		return
	}
	//通信開始
	if err = c.Mail(fromAdd); err != nil {
		fmt.Printf("--couldn't start send the mail---\n%v\n", err)
		return
	}
	if err = c.Rcpt(to0); err != nil {
		fmt.Printf("--couldn't specify the recipient---\n%v\n", err)
		return
	}
	wd, err := c.Data()
	if err != nil {
		fmt.Printf("--couldn't start send the message---\n%v\n", err)
		return
	}
	if _, err = wd.Write(ml.Bytes()); err != nil {
		fmt.Printf("--couldn't send the message---\n%v\n", err)
		return
	}
	//通信終了
	if err = wd.Close(); err != nil {
		fmt.Printf("--couldn't close the connection---\n%v\n", err)
		return
	}
	c.Quit()

	fmt.Println("ok")
	return true
}

////////////////////////////////////////
func send1(tion bool) (ok bool) {
	subject := time.Now().Format("2006/1/2")
	body := func(opera bool) (style string) {
		if opera {
			style = "manual"
		} else {
			style = "auto"
		}
		return
	}(tion)
	from := "imymemine"
	fromAdd := os.Getenv("M_from")
	pass := os.Getenv("M_frompass")
	to0 := os.Getenv("M_to")
	servername := os.Getenv("M_servername")
	host := strings.Split(servername, ":")[0]
	if bl {
		subject = "えまーじぇんしー"
		body = "再構築になんか問題あり"
		from = "uuruurs"
	}
	ml := email.NewMessage(subject, body)
	ml.From = mail.Address{Name: from, Address: fromAdd}
	ml.To = []string{to0}
	if !bl {
		if err = ml.Attach("./tmp/test.txt"); err != nil {
			fmt.Printf("--couldn't attach the file---\n%v\n", err)
			return
		}
	}
	auth := smtp.PlainAuth("", fromAdd, pass, host)
	tlsconfig := &tls.Config{InsecureSkipVerify: true, ServerName: host}
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		fmt.Printf("--couldn't tls connection---\n%v\n", err)
		return
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Printf("--couldn't create a new client---\n%v\n", err)
		return
	}
	if err = c.Auth(auth); err != nil {
		fmt.Printf("--couldn't authenticate the client---\n%v\n", err)
		return
	}
	if err = c.Mail(fromAdd); err != nil {
		fmt.Printf("--couldn't start send the mail---\n%v\n", err)
		return
	}
	if err = c.Rcpt(to0); err != nil {
		fmt.Printf("--couldn't specify the recipient---\n%v\n", err)
		return
	}
	wd, err := c.Data()
	if err != nil {
		fmt.Printf("--couldn't start send the message---\n%v\n", err)
		return
	}
	if _, err = wd.Write(ml.Bytes()); err != nil {
		fmt.Printf("--couldn't send the message---\n%v\n", err)
		return
	}
	if err = wd.Close(); err != nil {
		fmt.Printf("--couldn't close the connection---\n%v\n", err)
		return
	}
	c.Quit()
	fmt.Println("ok")
	return true
}
func dirwalk(dir string) (paths []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, "\n")
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, "\n", filepath.Join(dir, file.Name()))
	}
	return
}

///////////////////////////////////////////

//rebuild
func remake() (ok bool) {

	//db用意
	db := gormcore()

	//ファイル開く
	xf, err := excelize.OpenFile("datas.xlsx", excelize.Options{Password: os.Getenv("XlPassword")})
	if err != nil {
		fmt.Printf("--couldn't open the file---\n%v\n", err)
		return
	}

	//トランザクション開始
	tx := db.Begin()
	err = func(dbt *gorm.DB) error {

		//datas空
		if err = dbt.Exec("truncate table datas").Error; err != nil {
			fmt.Printf("--couldn't truncate the table---\n%v\n", err)
			return err
		}

		//pt初期化
		pt = post{}

		//300,348,396の３回
		ott := "300"
		ottn := 302
		for a := 0; a < 3; a++ {
			switch a {
			case 1:
				ott = "348"
				ottn = 350
			case 2:
				ott = "396"
				ottn = 398
			}

			//ott+"1"のA1から順にforで見て空にあたったら終わり(?)
			for b := 1; ; b++ {

				//列名
				colname, err := excelize.ColumnNumberToName(b)
				if err != nil {
					fmt.Printf("--couldn't convert Num2Nam---\n%v\n", err)
					return err
				}

				//イベント番号
				name, err := xf.GetCellValue(ott+"1", colname+"1")
				if err != nil {
					fmt.Printf("--couldn't get the val---\n%v\n", err)
					return err
				}

				//終わり
				if name == "" {
					break
				}

				pt.Name = name

				//continue用
				tobe := false

				for c := 2; c < ottn; c++ {
					for d := 1; d < 7; d++ {
						num, err := xf.GetCellValue(ott+strconv.Itoa(d), colname+strconv.Itoa(c))
						if err != nil {
							fmt.Printf("--couldn't get the val---\n%v\n", err)
							return err
						}

						//100位が空やったらスルー
						if num == "" && d == 1 {
							tobe = true
							continue
						}

						//ptに色々セット
						pt.Num = c - 1
						switch d {
						case 1:
							pt.One = num
						case 2:
							pt.Two = num
						case 3:
							pt.Three = num
						case 4:
							pt.Four = num
						case 5:
							pt.Five = num
						case 6:
							pt.Six = num
						}
					}
					if tobe {
						tobe = false
						continue
					}

					// insert
					if err = dbt.Table("datas").Save(&pt).Error; err != nil {
						fmt.Printf("--couldn't save the date---\n%v\n", err)
						return err
					}
				}
			}
		}
		return nil
	}(tx)

	//虎終了和閉扉
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
			ok = false
		}
		db.Close()
	}()
	return true
}

//管理する感じのやつ
func control(c *gin.Context) (ajax interface{}) {
	db := gormcore()
	defer db.Close()
	switch c.PostForm("f") {

	//tables取得
	case "0":
		tables, err := db.Raw("show tables").Rows()
		if err != nil {
			fmt.Println(err)
		}
		defer tables.Close()
		tbr := make([]rune, 0, 5)
		table := ""
		for tables.Next() {
			tables.Scan(&table)
			tbr = append(tbr, []rune(table)...)
			tbr = append(tbr, 44, 32)
		}
		return string(tbr)

		//records取得
	case "1":
		rcsb := make([]byte, 0, 32768)
		rcsb = append(rcsb, 227, 131, 172, 227, 130, 179, 227, 131, 188, 227, 131, 137, 230, 149, 176, 58, 32)
		s := c.PostForm("s")
		w := c.PostForm("w")
		if s == "" {
			s = "*"
		}
		if w != "" {
			w = " where " + w
		}
		n := c.PostForm("n")
		switch n {

		//from datas
		case "datas":
			rcs := []post{}
			db.Raw("select " + s + " from datas" + w).Scan(&rcs)
			rcsb = append(rcsb, []byte(strconv.Itoa(len(rcs)))...)
			rcsb = append(rcsb, 10, 110, 97, 109, 101, 44, 32, 110, 117, 109, 44, 32, 111, 110, 101, 44, 32, 116, 119, 111, 44, 32, 116, 104, 114, 101, 101, 44, 32, 102, 111, 117, 114, 44, 32, 102, 105, 118, 101, 44, 32, 115, 105, 120, 10)
			for _, v := range rcs {
				rcsb = append(rcsb, 10)
				rcsb = append(rcsb, []byte(v.Name)...)
				rcsb = append(rcsb, 44)
				rcsb = append(rcsb, []byte(strconv.Itoa(v.Num))...)
				rcsb = append(rcsb, 44)
				rcsb = append(rcsb, []byte(v.One)...)
				rcsb = append(rcsb, 44)
				rcsb = append(rcsb, []byte(v.Two)...)
				rcsb = append(rcsb, 44)
				rcsb = append(rcsb, []byte(v.Three)...)
				rcsb = append(rcsb, 44)
				rcsb = append(rcsb, []byte(v.Four)...)
				rcsb = append(rcsb, 44)
				rcsb = append(rcsb, []byte(v.Five)...)
				rcsb = append(rcsb, 44)
				rcsb = append(rcsb, []byte(v.Six)...)
			}
			return string(rcsb)

		//from list
		case "list":
			rcs := []period{}
			db.Raw("select " + s + " from list" + w).Scan(&rcs)
			rcsb = append(rcsb, []byte(strconv.Itoa(len(rcs)))...)
			rcsb = append(rcsb, 10, 110, 117, 109, 44, 32, 110, 97, 109, 101, 44, 32, 112, 101, 114, 105, 111, 100, 44, 32, 116, 105, 109, 101, 115, 44, 32, 100, 111, 110, 101, 10)
			for _, v := range rcs {
				rcsb = append(rcsb, 10)
				rcsb = append(rcsb, []byte(v.Num)...)
				rcsb = append(rcsb, 44)
				rcsb = append(rcsb, []byte(v.Name)...)
				rcsb = append(rcsb, 44)
				rcsb = append(rcsb, []byte(v.Period)...)
				rcsb = append(rcsb, 44)
				rcsb = append(rcsb, []byte(strconv.Itoa(v.Times))...)
				rcsb = append(rcsb, 44)
				rcsb = append(rcsb, []byte(strconv.Itoa(v.Done))...)
			}
			return string(rcsb)
		}

		//記録記録
	case "2":
		if backup() {
			if !send(true) {
				fmt.Println("send failed")
			}
		} else {
			fmt.Println("backup failed")
		}

		//手動取得
	case "3":
		api, v := setconf()
		gettweets(api, v)

		//リボーン
	case "4":
		bl = true
		fmt.Println("rebuilding...")
		if !remake() {
			if !send(true) {
				fmt.Println("send failed")
			}
			break
		}
		bl = false
		fmt.Println("success!")

		//一時保存テスト
	case "5":
		/*
			f, err := os.Create("./tmp/test.txt")
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()
			f.WriteString("test")
			if !send1(true) {
				fmt.Println("failed")
			}*/
		if err = os.Mkdir("tmp", 0777); err != nil {
			fmt.Println(err)
		}
		f, err := os.Create("./tmp/test.txt")
		if err != nil {
			fmt.Println(err)
		}
		f.WriteString("123\n456")
		if !send1(true) {
			fmt.Println("(*>△<)")
		}
		return dirwalk("../../")

		//caseを追加するときに分かりやすいように置いとく
	default:
		/*なし*/
	}
	return
}

func Run() {
	//現在時刻
	now := time.Now().Add(9 * time.Hour)

	if os.Getenv("PORT") == "" {
		if err = godotenv.Load("dev.env"); err != nil {
			fmt.Printf("--couldn't load env---\n%v\n", err)
		}
	}

	//ツイート取得の準備
	api, v := setconf()

	//master
	mas := os.Getenv("Master")

	//サーバーの準備
	r := gin.Default()
	r.LoadHTMLGlob("view/*.html")
	r.Static("/static", "./view/static")

	//トップ
	r.GET("/", func(c *gin.Context) {
		src, alt := roulette()
		c.HTML(http.StatusOK, "top.html", gin.H{"src": src, "alt": alt})
	})

	//HEADリクエスト
	r.HEAD("/")

	//いらん
	r.GET("room/:name", func(c *gin.Context) {
		c.HTML(http.StatusOK, "room.html", gin.H{"Name": c.Param("name")})
	})

	//イベントページ
	r.GET("events", func(c *gin.Context) {
		src, alt := roulette()

		if bl {
			//DB再構築中
			c.HTML(http.StatusOK, "room.html", gin.H{"Name": "メンテ中"})
		}

		c.HTML(http.StatusOK, "events.html", gin.H{"src": src, "alt": alt, "list": eventlist()})
	})

	//選択したイベントの情報
	r.POST("event", func(c *gin.Context) {
		sc, ele := graphinfo(c.PostForm("events"), c.PostForm("rank"), c.PostForm("height"))
		c.HTML(http.StatusOK, "event.html", gin.H{"selected": ele, "scale": sc})
	})

	//イベントポイント計算
	r.GET("calc1", func(c *gin.Context) {
		c.HTML(http.StatusOK, "calc1.html", gin.H{})
	})

	//仮実装
	r.GET("calc1/kari", func(c *gin.Context) {
		c.HTML(http.StatusOK, "calc0.9.html", gin.H{})
	})

	//○○計算
	r.GET("calc2", func(c *gin.Context) {
		c.HTML(http.StatusOK, "calc2.html", gin.H{})
	})

	//いらん
	r.GET("tweet", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ajax.html", gin.H{"tweet": 0, "coord": "M0 0 L75 150 L150 0"})
	})

	//管理
	r.GET(mas, func(c *gin.Context) {
		c.HTML(http.StatusOK, "care.html", gin.H{"time": now.Format("2006/01/02 15:04:05")})
	})

	//管理ajax
	r.POST("careajax", func(c *gin.Context) {
		c.HTML(http.StatusOK, "careajax.html", gin.H{"ajax": control(c)})
	})

	/*
		//download
		//今は使わんけどそのうち使うときのために残しとく
		r.GET("download/datas", func(c *gin.Context) {
			c.Writer.Header().Add("Content-Disposition", "attachment; filename=datas.xlsx")
			c.Writer.Header().Add("Content-Type", "application/octet-stream")
			c.File("datas.xlsx")
		})
	*/

	//30分毎にツイートの取得
	go func() {
		var (
			//1800000...30min
			//10000...10s
			wait, gap int = 1800000, 0
			wt        time.Duration
			start     time.Time
		)

		//開始時に一回やっとく
		gettweets(api, v)

		//一日やったらバックアップ作成
		if now.Day() == 1 {
			if backup() {
				if !send(false) {
					fmt.Println("send failed")
				}
			} else {
				fmt.Println("backup failed")
			}
		}

		g := []byte(time.Now().Format("05.0")[3:4])[0]
		if g < 53 {
			g += 5
		} else {
			g -= 5
		}
		for {
			start = time.Now()
			wt = time.Duration(wait-gap) * 1000000
			time.Sleep(wt)
			gettweets(api, v)
			//gapが大きいとsleepの時間はマイナス(0秒)になるけど、gapの値が正常になるまで0秒待つのが続くから、gapの値を正常にする処理が必要
			//とりあえずこんな感じにしとく
			if m == 11 {
				continue
			}
			if []byte(start.Format("05.0")[3:4])[0] == g {
				gap += 500
			}
			gap = int(time.Since(start))/1000000 - wait + gap
		}
	}()

	/*エラーが出たらサーバーが止まるからしっかり対策する*/
	/*
		投稿に25,000とか50,000位の記録がない場合がある
		「":"と"("もしくは"\n"からボーダーの数値を取得」のところでエラー
		ミスじゃなくて、単純に50000位まで人がおらんのが原因
	*/

	r.Run(":" + os.Getenv("PORT"))
	/*
		cmdで「set PORT=○○○○」を実行後に「http://localhost:○○○○/」にアクセスする
	*/
}
