package organs

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

type postDatas struct {
	Num   int `json:"num"`
	One   int `json:"one"`
	Two   int `json:"two"`
	Three int `json:"three"`
	Four  int `json:"four"`
	Five  int `json:"five"`
	Six   int `json:"six"`
}

type jpost struct {
	Name  int         `json:"name"`
	Datas []postDatas `json:"datas"`
}

type period struct {
	//番号
	Num string `gorm:"type:smallint unsigned PRIMARY KEY" json:"num"`
	//name...イベント名
	Name string `gorm:"type:varchar(60)" json:"name"`
	//period...開催期間
	Period string `gorm:"type:char(21)" json:"period"`
	//times...データの数
	Times int `gorm:"type:smallint unsigned" json:"times"`
	//done...整理状況
	//0:未使用, 1:記録済(不足), 2:記録済(充足)
	Done int `gorm:"type:tinyint default 0" json:"done"`
}

var (
	m, fwt, twt    int
	bl             bool
	err            error
	poti           postime
	db             *gorm.DB
	twiapi         *anaconda.TwitterApi
	twiv_t, twiv_i = url.Values{"screen_name": {"imas_ml_td_t"}, "count": {"4"}}, url.Values{"screen_name": {"imas_ml_td_i"}, "count": {"1"}}
	ename, pri     string
	dates, twlog   = make([]string, 14), make([]string, 4)
	krn, kk, kig   = make([]int, 6), make([]int, 6), make([]int, 6)
	pt             post
)

//データベースの準備
func gormcore() {

	//mysqlの設定
	//本番か開発かで設定を変える
	protocol := "tcp(" + os.Getenv("DB_HOSTNAME") + ":3306)"
	if os.Getenv("PORT") == "8080" {
		protocol = ""
	}
	db, err = gorm.Open("mysql", os.Getenv("DB_USERNAME")+
		":"+os.Getenv("DB_PASSWORD")+"@"+protocol+"/"+os.Getenv("DB_NAME"))
	if err != nil {
		fmt.Printf("--couldn't connect the DB---\n%v\n", err)
	}
}

//ツイート取得準備
func setconf() {

	//api
	anaconda.SetConsumerKey(os.Getenv("ConsumerKey"))
	anaconda.SetConsumerSecret(os.Getenv("ConsumerSecret"))
	twiapi = anaconda.NewTwitterApi(os.Getenv("AccessToken"), os.Getenv("AccessTokenSecret"))

	//テキストBOTさんの投稿、上から4個
	/*Twiv_t = url.Values{"screen_name": {"imas_ml_td_t"}, "count": {"4"}}

	//通知BOTさんの投稿用
	Twiv_i = url.Values{"screen_name": {"imas_ml_td_i"}, "count": {"2"}}*/
}
