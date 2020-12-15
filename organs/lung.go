package organs

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/mail"
	"net/smtp"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/scorredoira/email"
)

//通知の処理
func tuti(tweet []rune) {
label:
	for a, z := range tweet {
		if a == len(tweet)-14 {
			break
		}

		//"MILLION LIVE W"か"ミリコレ"ならm=8
		if reflect.DeepEqual(tweet[a:a+4], []rune{12511, 12522, 12467, 12524}) || reflect.DeepEqual(tweet[a:a+14], []rune{77, 73, 76, 76, 73, 79, 78, 32, 76, 73, 86, 69, 32, 87}) {
			m = 8
		}
		switch z {

		//"「"の場所を記録
		case 12300:
			poti.h = rune(a)

			//"「"と"」"の場所をもとにイベントの名前を取得
		case 12301:

			//既に記録されてたら以降の処理はスルー
			if string(tweet[poti.h+1:a]) == ename {
				continue

				//新イベならeventnameを上書き
			} else {
				ename = string(tweet[poti.h+1 : a])
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
			end, err := time.Parse("2006/1/2 15:04", string(tweet[a+7:a+23]))

			if err != nil {
				fmt.Printf("--couldn't set LD---\n%v\n", err)
			}

			//"ミリコレ"か"WORKING"なら最終日まで待つ
			if m == 9 {
				ato := -time.Since(end.Add(-539 * time.Minute))
				fmt.Printf("イベント名「%s」\n", ename)
				fmt.Printf("次のイベントまであと %v\n", ato)
				time.Sleep(ato)
				m = 11
				return

				//それ以外のイベントのとき
			} else {

				//開始日を取得
				start, err := time.Parse("2006/1/2 15", string(tweet[a-13:a]))
				if err != nil {
					fmt.Printf("--couldn't set SD---\n%v\n", err)
				}

				//イベントの日付を記録
				dates = make([]string, 14)
				b := 0
				for ; b <= end.Day()-start.Day(); b++ {
					day := start.AddDate(0, 0, b).Day()
					day1 := day/10 + 48
					day2 := day%10 + 48
					if day < 10 {
						day1 = 48
						day2 = day + 48
					}
					dates[b] = string([]byte{byte(day1), byte(day2)})
				}

				//fwtとtwtを記録
				fwt = start.Hour() << 1
				twt = ((b-1)*24+end.Add(1*time.Minute).Hour())<<1 - fwt

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
				pe.Name = ename
				pe.Period = start.Format("2006/01/02") + "~" + end.Format("2006/01/02")
				pe.Times = twt

				//既にあるかチェック
				che := struct{ Name string }{}
				db.Table("list").Select("name").Where("name=?", ename).Order("num desc").First(&che)

				//なかったら登録
				if che.Name == "" {
					db.Table("list").Save(&pe)
				}
			}
			break label
		}
	}
}

//ツイートの取得と加工
func tweektweets() {

	//イベントの日付が記録されてなかったら
	if dates[0] == "" {

		//通知BOTさんの最新の投稿を取得
		tweet, err := twiapi.GetUserTimeline(twiv_i)
		if err != nil {
			fmt.Printf("--couldn't get tweets---\n%v\n", err)
		}

		//改行コードを"\n"に統一
		twee := strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(tweet[0].FullText)

		//rune配列に変換して通知の処理
		tweer := []rune(twee)
		poti.m = false
		tuti(tweer)
	}

	//ツイートの取得
	tweet, err := twiapi.GetUserTimeline(twiv_t)
	if err != nil {
		fmt.Printf("--couldn't get tweets---\n%v\n", err)
		return
	}

	//古い投稿から処理したいからカウントダウンでループ
	for a := 3; a >= 0; a-- {
		poti.m = false

		//取得したツイートがtwlogに記録されてるかどうか
		for b := 0; b < 4; b++ {
			if tweet[a].FullText == twlog[b] {
				poti.m = true
				break
			}
		}

		//ツイートをtwlogに記録
		twlog[3-a] = tweet[a].FullText

		//すでに記録されてたら以降の処理をスルー
		if poti.m {
			continue
		}

		//改行の統一、rune変換
		tweer := []rune(strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(tweet[a].FullText))

		//RTじゃなかったら
		if tweet[a].RetweetedStatus == nil {
			m = 0
			kurai := 0
			for b, z := range tweer {
				switch z {

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
					poti.h = tweer[b-6]*10 + tweer[b-5] - 528

					//30分ならtrue
					if tweer[b-3] == 48 {
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
					if tweer[b-2] == 26524 {
						m++
					}

					//"位"の数だけループする
					for c := 0; c < kurai; c++ {
						var twer []rune

						//":"と"("もしくは"\n"からボーダーの数値を取得
						if m == 7 {
							twer = []rune(string(tweer[krn[c]+2 : kig[c]]))
						} else {
							twer = []rune(string(tweer[krn[c]+2 : kk[c]-1]))
						}

						//","を消す
						for d, y := range twer {
							if y == 44 {
								twer = append(twer[:d], twer[d+1:]...)
							}
						}

						//ptにイベント名と整形した数値をセット
						pt.Name = pri
						switch c {
						case 0:
							pt.One = string(twer)
						case 1:
							pt.Two = string(twer)
						case 2:
							pt.Three = string(twer)
						case 3:
							pt.Four = string(twer)
						case 4:
							pt.Five = string(twer)
						case 5:
							pt.Six = string(twer)
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

					//最終日やったら
					if m == 7 {
						pt.Num = twt
						db.Table("datas").Save(&pt)

						//これ以降の処理はもういい
						return
					}

					//何日のボーダーか取得
					wd := string(tweer[b-9 : b-7])

					//そうじゃなかったら適切な場所に記録
					for c, y := range dates {
						if wd == y {
							pt.Num = 48*c + int(poti.h) - fwt
							db.Table("datas").Save(&pt)
						}
					}
				}
			}

			//RTやったら通知の処置
		} else {
			poti.m = true
			tuti(tweer)
		}
	}
}

//記録されてるイベントのリストを作成
func eventlist() string {

	//記録されてるイベントの情報を取得
	pes := []period{}
	db.Table("list").Select("num,name,period").Order("num desc").Where("num>2").Find(&pes)

	//byte配列用意
	var listb []byte

	//listbにイベント情報を詰め込んでいく
	for _, z := range pes {

		//まずはnum
		nnpb := []byte(z.Num)
		for _, y := range nnpb {
			listb = append(listb, y)
		}

		//"//"を追加
		listb = append(listb, 47, 47)

		//次はname
		nnpb = []byte(z.Name)
		for _, y := range nnpb {
			listb = append(listb, y)
		}

		//"//"を追加
		listb = append(listb, 47, 47)

		//最後はperiod
		nnpb = []byte(z.Period)
		for _, y := range nnpb {
			listb = append(listb, y)
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
	src = "/view/img/hinata" + strconv.Itoa(n) + ".gif"
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

	pts := []post{}
	ls := period{}

	//最終結果(達)のなかで一番でかいのを記録
	rm := float64(0)
	if len(se) != 0 {

		//比べる用
		comp := float64(0)

		for _, z := range se {

			//pt初期化
			pt = post{}

			//numの最大値を取得
			db.Table("datas").Select("max(num)as one").Where("name=?", z).Find(&pt)

			//300か348か396か
			db.Table("list").Select("times").Where("num=?", z).Find(&ls)
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
			db.Table("datas").Select("(select " + ra + " from datas where name=" + name + " and num=" + num + ")*((select " + ra + " from datas where name=" + z + " and num=" + pt.One + ")/(select " + ra + " from datas where name=" + name + " and num=" + pt.One + "))as two").Find(&pt)

			//開催期間が396より大きいとき
			if ls.Times > 396 {
				/*仮*/
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
	for a, z := range se {

		//記録を取得
		db.Table("datas").Select(sel).Where("name=?", z).Order("num").Find(&pts)

		//num初期化
		num = -1

		//イベントの記録を入れる箱
		db.Table("list").Select("times").Where("num=?", z).Find(&ls)
		box := make([]float64, ls.Times+1)

		//「M0 0 」
		gra = append(gra, 77, 48, 32, 48, 32)

		//詰め込む
		for _, y := range pts {

			//指定したイベントの指定した順位のボーダー
			border, err := strconv.ParseFloat(y.One, 64)
			if err != nil {
				fmt.Printf("--couldn't convert AtoI---\n%v\n", err)
			}

			//被ってたらDBから削除
			if y.Num == num {
				db.Table("datas").Where("name=? and num=?", z, num).Limit(1).Delete(&y)
				continue
			}

			//記録
			box[y.Num] = border
			num = y.Num

			//graにデータを書き込んでいく
			//「L」
			gra = append(gra, 76)
			//"w.Num"
			gra = append(gra, []byte(strconv.FormatFloat(float64(y.Num)*2.51417108, 'f', 6, 64))...)
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
	if err = os.Mkdir("tmp", 0777); err != nil {
		fmt.Printf("--couldn't make the dir---\n%v\n", err)
		return
	}
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
		if err := ml.Attach("./tmp/datas.xlsx"); err != nil {
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
	if err := c.Auth(auth); err != nil {
		fmt.Printf("--couldn't authenticate the client---\n%v\n", err)
		return
	}
	//通信開始
	if err := c.Mail(fromAdd); err != nil {
		fmt.Printf("--couldn't start send the mail---\n%v\n", err)
		return
	}
	if err := c.Rcpt(to0); err != nil {
		fmt.Printf("--couldn't specify the recipient---\n%v\n", err)
		return
	}
	wd, err := c.Data()
	if err != nil {
		fmt.Printf("--couldn't start send the message---\n%v\n", err)
		return
	}
	if _, err := wd.Write(ml.Bytes()); err != nil {
		fmt.Printf("--couldn't send the message---\n%v\n", err)
		return
	}
	//通信終了
	if err := wd.Close(); err != nil {
		fmt.Printf("--couldn't close the connection---\n%v\n", err)
		return
	}
	c.Quit()

	fmt.Println("ok")
	return true
}

//rebuild
func remake() (ok bool) {

	//ファイル開く
	xf, err := excelize.OpenFile("datas.xlsx", excelize.Options{Password: os.Getenv("XlPassword")})
	if err != nil {
		fmt.Printf("--couldn't open the file---\n%v\n", err)
		return
	}

	//トランザクション開始
	tx := db.Begin()
	err = func() error {

		//datas空
		if err := tx.Exec("truncate table datas").Error; err != nil {
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
					if err = tx.Table("datas").Save(&pt).Error; err != nil {
						fmt.Printf("--couldn't save the date---\n%v\n", err)
						return err
					}
				}
			}
		}
		return nil
	}()

	//虎終了和閉扉
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
			ok = false
		}
	}()
	return true
}
