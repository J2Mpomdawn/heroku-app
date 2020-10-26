package ts

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
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
	//done...整理確認
	//0...未整理, 1...整理済, 2...使用済
	Done int `gorm:"type:tinyint default 0"`
}

var (
	n, m, fwt, twt        int
	poh                   rune
	tweet, eventname, pri string
	start, end            time.Time
	err                   error
	dates, twlog          []string = make([]string, 14), make([]string, 4)
	ru, r                 []rune
	krn, kk, kig          []int = make([]int, 6), make([]int, 6), make([]int, 6)
	poti                  postime
	pt                    post
	pts                   [][]post = make([][]post, 14)
)

//ツイート取得準備
func setconf() (api *anaconda.TwitterApi, v url.Values) {

	//apiの設定
	anaconda.SetConsumerKey(os.Getenv("ConsumerKey"))
	anaconda.SetConsumerSecret(os.Getenv("ConsumerSecret"))
	api = anaconda.NewTwitterApi("AccessToken", "AccessTokenSecret")

	//とるのはテキストBOTさんの投稿、上から4個
	v = url.Values{}
	v.Set("screen_name", "imas_ml_td_t")
	v.Set("count", "4")

	//ボーダーを記録する箱を用意
	for a := 0; a < 14; a++ {
		p := make([]post, 48)
		pts[a] = p
	}
	return
}

//データベースの準備
func gormcore() *gorm.DB {

	//mysqlの設定
	//本番か開発かで設定を変える
	protocol := "tcp(" + os.Getenv("DB_HOSTNAME") + ":3306)"
	if os.Getenv("PORT") == "" {
		if err = godotenv.Load("dev.env"); err != nil {
			fmt.Printf("--couldn't load env---\n%v\n", err)
		}
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
				db.Table("list").CreateTable(&period{})
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
			fmt.Println(err)
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
	db.Table("datas").CreateTable(&post{})
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

						//記録整理
						//未整理のイベントを列挙
						ds := []period{}
						db.Table("list").Select("num").Where("done=?", 0).Find(&ds)

						spt := []post{}
						//被りチェック用
						n := 0

						//イベント毎に被りを調べてあったら消していく
						for _, x := range ds {
							db.Table("datas").Select("num").Where("name=?", x.Num).Find(&spt)

							//n初期化
							n = -1

							for _, y := range spt {

								//かぶってたら消す
								if y.Num == n {
									db.Table("datas").Where("name=? and num=?", x.Num, n).Limit(1).Delete(&y)
									continue
								}
								n = y.Num
							}

							//済
							db.Table("list").Where("num=?", x.Num).Update("done", 1)
						}

						//平均に使うものを列挙
						//開催期間が８日
						//整理済み未使用(１)

						db.Table("list").Select("num").Where("times=348 or done=1").Find(&ds)

						//使えるのがあったら
						if l := len(ds); l != 0 {

							//列挙したものをもとにwhere句を作る
							wherb := make([]byte, 0, 128)

							//「num=? and(」を追加
							wherb = append(wherb, 110, 117, 109, 61, 63, 32, 97, 110, 100, 40)

							for _, x := range ds {

								//「name=」を追加
								wherb = append(wherb, 110, 97, 109, 101, 61)

								//"x.Num"を追加
								wherb = append(wherb, []byte(x.Num)...)

								//「 or 」を追加
								wherb = append(wherb, 32, 111, 114, 32)
							}

							//後ろ４文字を消して「)」を追加
							wherb = append(wherb[:len(wherb)-4], 41)

							//文字列に変換
							wher := string(wherb)

							//個数を文字列に変換
							ls := strconv.Itoa(l)

							//平均を計算して記録
							pt.Name = "0"
							for c := 1; c < 349; c++ {

								//cを文字列に変換
								cs := strconv.Itoa(c)

								//平均取得
								db.Table("datas").Select("(sum(one)+(select one from datas where name=0 and num="+cs+")*(select times from list where num=0))/((select times from list where num=0)+"+ls+") as one,(sum(two)+(select two from datas where name=0 and num="+cs+")*(select times from list where num=0))/((select times from list where num=0)+"+ls+") as two,(sum(three)+(select three from datas where name=0 and num="+cs+")*(select times from list where num=0))/((select times from list where num=0)+"+ls+") as three,(sum(four)+(select four from datas where name=0 and num="+cs+")*(select times from list where num=0))/((select times from list where num=0)+"+ls+") as four,(sum(five)+(select five from datas where name=0 and num="+cs+")*(select times from list where num=0))/((select times from list where num=0)+"+ls+") as five,(sum(six)+(select six from datas where name=0 and num="+cs+")*(select times from list where num=0))/((select times from list where num=0)+"+ls+") as six").Where(wher, c).Find(&pt)

								//平均記録
								//整数型やから四捨五入される
								pt.Num = c
								db.Table("datas").Where("name=0 and num=?", c).Update(&pt)
							}

							//使用済み＆times上書き
							db.Exec("update list set times=(select times from(select times from list where num=0)as t)+? where num=0", l)
							for d, x := range wherb {
								if x == 40 {
									db.Table("list").Where(strings.NewReplacer("name", "num").Replace(string(wherb[d+1:len(wherb)-1]))).Update("done", 2)
									break
								}
							}
						}

						//これ以降の処理はもういい
						return
					}

					//そうじゃなかったら適切な場所に記録
					for d, x := range dates {
						if tweet == x {
							pt.Num = 48*d + int(poti.h) - fwt
							/*db.Table("datas").Save(&pt)*/

							/*
								DB操作のテスト用に使う
								そのときは上の１行をコメントアウトする
								(別にせんでもいいならそのまま)
							*/
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
	db.Table("list").Select("num,name,period").Order("num desc").Where("num!=0").Find(&pes)
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

			//そのnumの平均と比べて最終日予想
			db.Table("datas").Select("(select " + ra + " from datas where name=0 and num=348)*((select " + ra + " from datas where name=" + v + " and num=" + pt.One + ")/(select " + ra + " from datas where name=0 and num=" + pt.One + "))as two").Find(&pt)

			//開催期間が348じゃないとき
			if db.Table("list").Select("times").Where("num=?", v).Find(&ls); ls.Times != 348 {

				//開催期間が300のとき
				if ls.Times == 300 {
					if pt.One == "300" {
						db.Table("datas").Select(ra+" as two").Where("name=? and num=?", v, pt.One).Find(&pt)
					} else {
						db.Table("datas").Select(ra + " as two").Where("name=0 and num=400").Find(&pt)
					}

					//開催期間が348より大きいとき
				} else {
					if pt.One == "396" {
						db.Table("datas").Select(ra+" as two").Where("name=? and num=?", v, pt.One).Find(&pt)
					} else {
						db.Table("datas").Select(ra + " as two").Where("name=0 and num=496").Find(&pt)
					}
				}
			}

			//float64に変換
			comp, err = strconv.ParseFloat(pt.Two, 64)
			if err != nil {
				fmt.Printf("\n--couldn't convert AtoF---\n%v\n", err)
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
			fmt.Printf("\n--couldn't convert AtoF---\n%v\n", err)
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
		fmt.Printf("\n--couldn't convert AtoF---\n%v\n", err)
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

//管理する感じのやつ
func control(f, n, s, w string) (ajax interface{}) {
	db := gormcore()
	defer db.Close()

	switch f {

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
		if s == "" {
			s = "*"
		}
		if w != "" {
			w = " where " + w
		}
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
	}
	return
}

func Run() {

	//ツイート取得の準備
	api, v := setconf()

	//サーバーの準備
	r := gin.Default()
	r.LoadHTMLGlob("view/*.html")
	r.Static("/static", "./view/static")

	//ローディングGIF用の文字列
	var src, alt string

	//トップ
	r.GET("/", func(c *gin.Context) {
		src, alt = roulette()
		c.HTML(http.StatusOK, "top.html", gin.H{"src": src, "alt": alt})
	})

	//いらん
	r.GET("room/:name", func(c *gin.Context) {
		c.HTML(http.StatusOK, "room.html", gin.H{"Name": c.Param("name")})
	})

	//イベントページ
	r.GET("events", func(c *gin.Context) {
		src, alt = roulette()
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

	//○○計算
	r.GET("calc2", func(c *gin.Context) {
		c.HTML(http.StatusOK, "calc2.html", gin.H{})
	})

	//いらん
	r.GET("tweet", func(c *gin.Context) {
		n++
		c.HTML(http.StatusOK, "ajax.html", gin.H{"tweet": n, "coord": "M0 0 L75 150 L150 0"})
	})

	//DB管理
	r.GET("control", func(c *gin.Context) {
		c.HTML(http.StatusOK, "control.html", gin.H{})
	})

	//tables
	r.POST("controlajax", func(c *gin.Context) {
		c.HTML(http.StatusOK, "controlajax.html", gin.H{"ajax": control(c.PostForm("f"), c.PostForm("n"), c.PostForm("s"), c.PostForm("w"))})
	})

	//30分毎にツイートの取得
	//サーバー建てるのと並行してやってもらう
	//ちゃんと30分毎になるように処理時間を測って微調整もする
	go func() {
		var (
			//1800000...30min
			//10000...10s
			wait, gap int = 1800000, 0
			wt        time.Duration
			start     time.Time
		)
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

	if port := os.Getenv("PORT"); port != "" {
		r.Run(":" + os.Getenv("PORT"))
	} else {
		r.Run(":8080")
	}
	/*
		cmdで「set PORT=○○○○」を実行後に「http://localhost:○○○○/」にアクセスする
	*/
}
