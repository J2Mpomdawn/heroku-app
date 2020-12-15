package organs

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//管理する感じのやつ
func control(c *gin.Context) (ajax interface{}) {
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
		tweektweets()

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

		//けすのめんｄ
	case "5":
		//

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
		if err := godotenv.Load("dev.env"); err != nil {
			fmt.Printf("--couldn't load env---\n%v\n", err)
		}
	}

	//mysqlの準備
	gormcore()
	defer db.Close()

	//ツイート取得の準備
	setconf()

	//master
	mas := os.Getenv("Master")

	//サーバーの準備
	r := gin.Default()
	r.LoadHTMLGlob("view/html/*.html")
	r.Static("view", "./view")

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
		tweektweets()

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
			tweektweets()
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
