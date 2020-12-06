/*******************
  関数とか変数たち
*******************/

//省略
const dg = function(a) {
    return document.getElementById(a);
}//経験値,元気のラベルを押したときにキャレットを末尾に設定
  , forcus = function() {
    const l = dg(this.s)
      , ll = l.value.length;
    l.setSelectionRange(ll, ll);
}//元気最大値があがるレベルリスト
  , level = []
  , maxfuns = function() {
    const a = (b,c,d)=>{
        while (b < c) {
            b += d;
            level.push(b);
        }
    }
    ;
    a(0, 58, 2);
    a(57, 147, 3);
    a(146, 422, 4);
    a(421, 581, 5);
    a(580, 700, 6);
}//レベルに対する元気の最大値
  , maxfun = function(a) {
    return (a < 700) ? level.map(b=>a - b).filter(c=>c >= 0).length + 60 : 240;
}//次のレベルまでの経験値
  , nextlevel = function(a, b) {
    return (a - 1) * 100 + 50 - b;
}//元気ゲージ
  , fnbar = function() {
    const mf = maxfun(dg("lv").value)
      , par = 100 * dg("fn").value / mf;
    dg("mfn").textContent = mf;
    if (par <= 100) {
        dg("bar2").value = par;
        dg("bar3").value = 0;
        dg("fn").style.color = "#4d5759e3";
        dg("fn").style.textShadow = "0 0 1.5px #25452d";
    } else if (par <= 200) {
        dg("bar2").value = 100;
        dg("bar3").value = par - 100;
        dg("fn").style.color = "#d55962b8";
        dg("fn").style.textShadow = "0 0 1.5px #c34c56";
    } else {
        alert("out");
    }
}//経験値ゲージ
  , exbar = function() {
    const nx = nextlevel(dg("lv").value, 0)
      , par = 100 * dg("exp").value / nx;
    dg("mexp").textContent = nx;
    if (par < 100) {
        dg("bar1").value = par;
    } else {
        alert("out");
    }
}//半角
  , half = function(a) {
    return a.replace(/[！-～]/g, b=>String.fromCharCode(b.charCodeAt(0) - 0xFEE0));
}//3桁区切り
  , separate = function(a) {
    if (a === "") {
        alert("空っぽ");
        return "";
    }
    a = half(a).replace(/,/g, "").trim();
    if (!/^[+|-]?(\d*)(\.\d+)?$/.test(a)) {
        alert("数字を入力するべき");
        return a
    }
    return new Intl.NumberFormat().format(Math.round(a));
}//時計
  , clock = function() {
    setInterval(()=>{
        const date = new Date();
        dg("now").textContent = date.getFullYear() % 100 + "/" + (date.getMonth() + 1) + "/" + date.getDate() + " " + date.getHours() + ":" + date.getMinutes() + ":" + date.getSeconds();
    }
    , 1000);
};
//10,20,30の箱
let i1, i2, i3;
//元気の横のボタンを押したら使用アイテム欄を表示
const open = function() {
    i1 = +document.querySelector("#ener>:nth-child(2)").value;
    i2 = +document.querySelector("#ener>:nth-child(3)").value;
    i3 = +document.querySelector("#ener>:nth-child(4)").value;
    document.getElementById("ener").style.display = "block";
}//キャンセル
  , cancel = function() {
    document.querySelector("#ener>:nth-child(2)").value = i1;
    document.querySelector("#ener>:nth-child(3)").value = i2;
    document.querySelector("#ener>:nth-child(4)").value = i3;
    setTimeout(()=>{
        dg("ener").style.display = "none"
    }
    , 200);
}//使用アイテムリセット
  , clear = function() {
    document.querySelector("#ener>:nth-child(2)").value = 0;
    document.querySelector("#ener>:nth-child(3)").value = 0;
    document.querySelector("#ener>:nth-child(4)").value = 0;
}//使用アイテム保存
  , use = function() {
    i1 = +document.querySelector("#ener>:nth-child(2)").value;
    i2 = +document.querySelector("#ener>:nth-child(3)").value;
    i3 = +document.querySelector("#ener>:nth-child(4)").value;
    setTimeout(()=>{
        dg("ener").style.display = "none"
    }
    , 200);
}//ボタンの文字達
  , bvs = [{
    s: "le",
    b: false,
    t: "ライブ",
    f: "お仕事"
}, {
    s: "bc",
    b: false,
    t: "１倍",
    f: "２倍"
}, {
    s: "st",
    b: false,
    t: "シアター",
    f: "ツアー"
}, {
    s: "iN",
    b: false,
    t: "自然回復分を含める",
    f: "自然回復分を含めない"
}]
  , switchButton = function() {
    const bv = bvs[this.n];
    if (bv.b) {
        dg(bv.s).value = bv.t;
    } else {
        dg(bv.s).value = bv.f;
    }
    if (this.n === 2) {
        dg("bi").options[2] = null;
        const td = document.createElement("option");
        if (bv.b) {
            td.value = 3;
            td.appendChild(document.createTextNode("３倍"));
        } else {
            td.value = 4;
            td.appendChild(document.createTextNode("４倍"));
        }
        dg("bi").appendChild(td);
    }
    bv.b = !bv.b;
};

/*************
  関数の実行
*************/

//読み込み時
window.addEventListener("load", function() {
    clock();
    maxfuns();
    fnbar();
    exbar();
    dg("dia").value = separate(dg("dia").value);
}, false);

//forcus実行>経験値
dg("mexp").addEventListener("click", {
    handleEvent: forcus,
    s: "exp"
}, false);

//forcus実行>元気
dg("mfn").addEventListener("click", {
    handleEvent: forcus,
    s: "fn"
}, false);

//fn,ex bar実行>レベル
dg("lv").addEventListener("change", function() {
    fnbar();
    exbar();
});

//fnbar実行>元気
dg("fn").addEventListener("change", fnbar, false);

//exbar実行>経験値
dg("exp").addEventListener("change", exbar, false);

//separate実行
dg("dia").addEventListener("blur", function() {
    this.value = separate(this.value);
}, false);

//フォーカスしたらコンマ消す
dg("dia").addEventListener("focus", function() {
    this.value = this.value.replace(/,/g, "");
}, false);

//open実行
dg("fns").addEventListener("click", open, false);

//cancel実行
dg("ca").addEventListener("click", cancel, false);

//clear実行
dg("cl").addEventListener("click", clear, false);

//use実行
dg("us").addEventListener("click", use, false);

//switchButton>ライブ
dg("le").addEventListener("click", {
    handleEvent: switchButton,
    n: 0
}, false);

//switchButton>１倍
dg("bc").addEventListener("click", {
    handleEvent: switchButton,
    n: 1
}, false);

//switchButton>シアター
dg("st").addEventListener("click", {
    handleEvent: switchButton,
    n: 2
}, false);

//switchButton>含める
dg("iN").addEventListener("click", {
    handleEvent: switchButton,
    n: 3
}, false);

/********************/

/********************/
var lv, exp, fn, dia, tix, wh, wm, of, cosm, pl, st, mfn, mi, kai, ka, kaif = (a)=>{
    kai = ((b3 % 2 === 0) ? 30 : 60) * a;
    ka = kai / a;
}
, gtps, iv, cm, ivsm, mi, gte, ec, ei, ps, Lv, Exp, Fn, Mi, Mfn, Pl, Gte, Ps, Cosm, Ivsm, Cufn, fnm, cufn, par, pas = (a)=>parseInt(dg(a).value), dsa = (a,b)=>dg(a).setAttribute("value", b), dgt = (a,b)=>dg(a).textContent = b, b1 = 0, b2 = 0, b3 = 0, b4 = 0
, yoso = (a,b,c)=>{
    while (a < b) {
        a += c;
        level.push(a);
    }
}
;
var fil = (a)=>{
    hiki = level.map(function(b) {
        return a - b;
    });
    sa = hiki.filter(function(c) {
        return c >= 0;
    });
    return Math.min.apply({}, sa);
}
  , sfnm = (a)=>(a < 700) ? (level.indexOf(a - fil(a)) + 61) : 240
  , fng = ()=>{
    dg("mfn").textContent = sfnm(pas("lv"));
    par = 100 * pas("fn") / sfnm(pas("lv"));
    if (par <= 100) {
        dg("bar2").value = par;
        dg("bar3").value = 0;
        dg("fn").style.color = "#4d5759e3";
        dg("fn").style.textShadow = "0 0 1.5px #25452d";
    } else if (par <= 200) {
        dg("bar2").value = 100;
        dg("bar3").value = par - 100;
        dg("fn").style.color = "#d55962b8";
        dg("fn").style.textShadow = "0 0 1.5px #c34c56";
    } else {
        alert("out");
    }
}
  , nx = (a,b)=>(a - 1) * 100 + 50 - b

  , exg = ()=>{
    dg("mexp").textContent = nx(pas("lv"), 0);
    par = 100 * pas("exp") / nx(pas("lv"), 0);
    if (par < 100) {
        dg("bar1").value = par;
    } else {
        alert("out");
    }
}
  , kanm = (a)=>{
    if (a === '') {
        return '';
    }
    a = han(a).replace(/,/g, "").trim();
    if (!/^[+|-]?(\d*)(\.\d+)?$/.test(a)) {
        return a;
    }
    var b = Math.round(a);
    return new Intl.NumberFormat().format(b);
}
  , kesu = (c)=>{
    return c.replace(/,/g, "");
}
  , han = (c)=>{
    var d = c.replace(/[！-～]/g, function(e) {
        return String.fromCharCode(e.charCodeAt(0) - 0xFEE0);
    });
    return d;
}
  , rgtp = (a)=>{
    for (i = 0; i < a; i++) {
        let ra = Math.random();
        let rp = (ra > 0.24889867841) ? 48 : (ra > 0.09471365638) ? 84 : 65;
        rgtps += rp;
        //〇□△/回数、〇□70%/回数
    }
    return rgtps;
}
  , sinzo = (a,b)=>{
    let coltm = ((a + Mfn) < ka) ? kai / ka : Math.floor((a + Mfn) / ka);
    Cosm += coltm;
    Pl += ((a + Mfn) < ka) ? kai : 0;
    Mfn = (a + Mfn) % ka;
    let gti = Math.floor(st[0] * coltm) + Mi;
    let gtp = (cm === "l") ? st[1] * coltm : rgtp(coltm);
    rgtps = 0;
    let ivtm = Math.floor(gti / st[2]);
    Ivsm += ivtm;
    Mi = gti % st[2];
    let gtpp = st[3] * ivtm;
    Gte += (306 * (coltm + ivtm) + b);
    Ps += (gtp + gtpp);
    return [Pl, Mfn, Mi, Gte, Ps, Cosm, Ivsm];
}
  , cure = ()=>{
    for (; ; ) {
        if (nx(Lv, Exp) > Gte)
            break;
        Gte -= nx(Lv, Exp);
        Exp = 0;
        cufn += sfnm(Lv + 1);
        Lv++;
    }
    ;return cufn;
}
  , keisan = ()=>{
    //Lv,Exp,Fn,Mi,Mfn,Pl,Gte,Ps,Cosm,Ivsm,Cufn
    lv = pas("lv");
    Lv = lv;
    exp = pas("exp");
    Exp = exp;
    fn = pas("fn");
    Fn = fn,
    tix = pas("tix");
    let hm = dg("tm").value.split(":");
    wh = parseInt(hm[0], 10);
    wm = parseInt(hm[1], 10);
    of = 60 * ((pas("kika") - 1) * 24 - 15 + wh);
    iv = (dg("st").value === "シアター") ? "s" : "t";
    cm = (dg("le").value === "ライブ") ? "l" : "e";
    st = (iv === "s") ? [(ip = (cm === "l") ? 85 * ((b3 % 2 === 0) ? 1 : 2) : 59.5 * ((b3 % 2 === 0) ? 1 : 2)), ip, 180 * pas("bi"), 537 * pas("bi")] : [6 * ((b3 % 2 === 0) ? 1 : 2), 140 * ((b3 % 2 === 0) ? 1 : 2), 20 * pas("bi"), 720 * pas("bi")];
    kaif(10);
    mi = (iv === "s") ? (360 * pas("kika")) : (40 * pas("kika"));
    Mi = mi;
    mfn = ((wh + wm / 60 < 21) ? 12 * ((pas("kika") - 1) * 24 - 15) + (12 * wh + Math.floor(wm / 5)) : 2087) + ((cm === "l") ? 0 : tix);
    Mfn = mfn;
    //ec = 306/((b3%2===0)?1:2);
    //ei = 306/pas("bi");
    kabe = pl = gte = ps = cosm = ivsm = cufn = Pl = Gte = Ps = Cosm = Ivsm = 0;
    let cfs = []
      , so = []
      , soo = []
      , slc = (a,b)=>so.slice(a - 1, a)[0][b];

    //自然回復分計算---
    sinzo(Fn, of);
    for (i = 0; i < 4; i++) {
        cure();
        if (cufn + Mfn < ((b3 % 2 === 0) ? 30 : 60))
            break;
        sinzo(cufn, 0);
        cfs.push(cufn);
        cufn = 0;
    }
    dgt("ncl", Lv);
    dgt("ncp", Ps);
    //---ここまで
    //リセット---
    let pS = Ps, pL, lV;
    Pl = pl;
    Gte = gte;
    Ps = ps;
    Cosm = cosm;
    Ivsm = ivsm;
    Lv = lv;
    Mfn = mfn;
    Mi = mi;
    Exp = exp;
    //---ここまで

    if (pS === pas("gl")) {
        dgt("rsf", fn + mfn + cfs.reduce((a,b)=>a += b) + "(=自然回復分とレベルアップの分)")
    } else if (pS < pas("gl")) {
        //自然回復分より大きい---
        sinzo(Fn, of);
        while (Ps < pas("gl")) {
            cure();
            let iti = sinzo(cufn, 0);
            cufn = 0;
            soo.push(Pl);
            so.push(Lv, iti);
        }
        let io = soo.indexOf(Pl) * 2;
        pS = Ps;
        pL = Pl;
        lV = Lv;
        Lv = so.slice(io - 2, io - 1)[0];
        Pl = slc(io, 0);
        Mfn = slc(io, 1);
        Mi = slc(io, 2);
        Gte = slc(io, 3);
        Ps = slc(io, 4);
        Cosm = slc(io, 5);
        Ivsm = slc(io, 6);
        kaif(1);
        while (Ps <= pas("gl")) {
            cure();
            ni = sinzo(cufn, 0);
            cufn = 0;
        }
        if (Pl === pL) {
            Pl = pL;
            Ps = pS;
            Lv = lV;
        }
        //---ここまで
    }
    dgt("rsf", Pl);
    dgt("rsp", Ps);
    dgt("rsl", Lv);
    dgt("cos", Cosm);
    dgt("ivs", Ivsm);
    dgt("rmf", Mfn);
    dgt("rmi", Mi);

}

dg("ksan").onclick = keisan;