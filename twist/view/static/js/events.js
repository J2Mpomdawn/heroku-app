/*
画面サイズによって改行を追加する関数
背景gifサイズを変更する機能も追加
*/
function listrans(){
  let en = document.getElementsByClassName("ename"),
      pe = document.getElementsByClassName("period"),
     ena,
     per,
       b;
  const w = (window.innerWidth < 640) ? true : false,
       br = (en[0].innerHTML.indexOf("\n") === -1) ? true : false ;
  for (let i = 0; i < en.length; i++) {
    if (br) {
      if (w) {
        ena = en[i].innerHTML;
        b = ena.indexOf("～");
        en[i].innerHTML = ena.slice(0,b)+"\n"+ena.slice(b);
        per = pe[i].innerHTML;
        b = per.indexOf("~");
        pe[i].innerHTML = per.slice(0,b)+"\n"+per.slice(b);
      }
    } else {
      if (!w) {
        ena = en[i].innerHTML;
        b = ena.indexOf("～");
        en[i].innerHTML = ena.slice(0,b-1)+ena.slice(b);
        per = pe[i].innerHTML;
        b = per.indexOf("~");
        pe[i].innerHTML = per.slice(0,b-1)+per.slice(b);
      }
    }
  }
  const back = document.getElementById("back");
  if (back===null){
    return;
  }
  back.width=0;
  back.height=0;
  back.width=document.getElementsByClassName("table")[0].getBoundingClientRect().width;
  back.height=document.getElementsByClassName("table")[0].getBoundingClientRect().height;
}

/*
ローディング時に"table"の中身をイベント達のテーブルに書き換える
*/
window.addEventListener("load", function(){
  const idt = document.getElementById("table"),
       list = idt.innerHTML.split(",,"),
      table = document.createElement("table"),
        trh = document.createElement("tr");
  idt.innerHTML="";
  table.classList.add("table");
  for (let i = 0; i < 3; i++) {
    const th = document.createElement("th");
    th.scope = "col";
    switch (i) {
      case 0:
        th.innerHTML = "レ";
        break;
      case 1:
        th.innerHTML = "イベント名";
        break;
      default:
        th.innerHTML = "期間";
    }
    trh.appendChild(th);
  }
  const thd = document.createElement("thead"),
        tbd = document.createElement("tbody");
  thd.appendChild(trh);
  table.appendChild(thd);
  for (let i = 0; i < list.length; i++) {
    const trb = document.createElement("tr"),
         info = list[i].split("//"),
          td1 = document.createElement("td"),
           ch = document.createElement("input");
    ch.setAttribute("type","checkbox");
    ch.setAttribute("name","check");
    ch.setAttribute("value",info[0]);
    td1.appendChild(ch);
    trb.appendChild(td1);
    const td2 = document.createElement("td");
    td2.classList.add("ename");
    td2.innerHTML = info[1];
    trb.appendChild(td2);
    const td3 = document.createElement("td");
    td3.classList.add("period");
    td3.innerHTML = info[2];
    trb.appendChild(td3);
    tbd.appendChild(trb);
  }
  table.appendChild(tbd);
  idt.appendChild(table);
  listrans();
  const img = document.createElement("img");
  img.setAttribute("id","back");
  img.setAttribute("src","/static/img/1041uuu1.gif");
  img.setAttribute("width",document.getElementsByClassName("table")[0].getBoundingClientRect().width);
  img.setAttribute("height",document.getElementsByClassName("table")[0].getBoundingClientRect().height);
  idt.children[0].insertBefore(img,thd);
}, false);

/*
グラフ描画用の関数
*/
const gra = document.getElementById("graph"),
   figure = document.getElementById("figure");
function graph(){
  const datas=document.getElementById("datas").innerHTML.split(",");
  document.getElementById("datas").remove();
  for (let i=0;i<datas.length;i++){
    const pathn = document.createElementNS("http://www.w3.org/2000/svg","path");
    /*
    色をイベント毎に設定したい
    赤青黄色。。。みたいに適当に順番に設定してもいい
    この色を下に表示するイベント名にも使う
    */
    const evecolo = "#"+("000000"+(75915*Math.floor(222*(i/(datas.length+1)))-i*75915).toString(16)).slice(-6);
    pathn.setAttribute("stroke",evecolo);
    pathn.setAttribute("d",datas[i]);
    datas[i]="";
    gra.children[1].appendChild(pathn);
  }

  const sc=document.getElementById("scal").innerHTML;
  document.getElementById("scal").remove();
}

/*
グラフの最大サイズをもとにtransformする関数
*/
function retrans(){
  let w, h, wm=0, hm=0;
  const chc = gra.children[0].childElementCount;
  gra.children[0].removeAttribute("style");
  for (let i = 0; i < chc; i++) {
    w = gra.children[0].children[i].getBoundingClientRect().width;
    if (w > wm) {
      wm = w;
    }
    h = gra.children[0].children[i].getBoundingClientRect().height;
    if (h > hm) {
      hm = h;
    }
  }
  gra.children[0].style.transform="scaleX("+(gra.getBoundingClientRect().width/wm)+") scaleY("+(gra.getBoundingClientRect().height/hm)+")";
  gra.children[1].children[0].setAttribute("d","M0 "+(gra.clientHeight-parseFloat(window.getComputedStyle(figure).height))/2+" v"+window.getComputedStyle(figure).height.slice(0,-2));
  gra.children[1].children[1].setAttribute("d","M"+(parseFloat(window.getComputedStyle(figure).width)-figure.getBoundingClientRect().width+2)/2+" 0 h"+(figure.getBoundingClientRect().width-2));
}

/*
選択されたイベントのグラフを描写
*/
const che = document.getElementsByName("check");
function check(){
  const arr = [];
  for (let i = 0; i < che.length; i++) {
    if (che[i].checked) arr.push(che[i].value);
  }
  if (!arr.length) return;
  const ajax = new XMLHttpRequest;
  ajax.open("POST","/event",true);
  ajax.setRequestHeader("content-type", "application/x-www-form-urlencoded;charset=UTF-8");
  ajax.onload=function (){
    gra.innerHTML=ajax.responseText;
    graph();
    retrans();
    gra.nextElementSibling.id="grad";
  }
  ajax.send("events="+arr+"&rank="+document.getElementById("rank").value+"&height="+gra.getBoundingClientRect().height.toFixed(5));
}

/*
ローディング時のgifを消す
*/
window.addEventListener("load",function(){
  setTimeout(
    function(){
      document.getElementById("loading").remove();
    },"500");
},false);

document.getElementById("bt").onclick=check;

/*
色んなサイズを画面サイズをもとに調整
*/
(function () {
  let timer,
       renl = () => {listrans();retrans();};
  window.addEventListener("resize",function(){
    clearTimeout(timer);
    timer = setTimeout(renl,300);
  },false);
}());

/*
スクロールしたら背景gifをリスタート
*/
(function () {
  let timer,
        res = () => {
          const back = document.getElementById("back");
          back.id="";
          setTimeout(function(){
            back.id="back";
          },10);
        };
  window.addEventListener("scroll",function(){
    clearTimeout(timer);
    timer = setTimeout(res,500);
  },false);
  document.getElementById("table").addEventListener("scroll",function(){
    clearTimeout(timer);
    timer = setTimeout(res,500);
  },false);
}());


/*
色テスト
*/
function clr(){
  const ct = document.getElementById("ct"),
        ar = [];
  for (let i = 0; i < 222; i++){
    /*const cc = "#"+("000000"+(86037*i).toString(16)).slice(-6),*/
    const cc = 75915*i;
    ar.push(cc);
  }
  /*
  arに色のリストが格納されてる
  222個格納されてる
  */
  /*
  指定したイベントの数がn個やったら
    222* (i/n+1)
  で、iは1～nの整数でループさせる
  */
  let ind = 0;
  for (let n = 1; n < 10; n++){
    const dvv = document.createElement("div");
    dvv.style.marginBottom = "10px";
    for (let i = 1; i <= n; i++){
      const dv = document.createElement("div");
      dv.style.width = "300px";
      dv.style.height = "20px";
      dv.style.backgroundColor="#"+("000000"+(75915*Math.floor(222*(i/(n+1)))-n*75915).toString(16)).slice(-6);
      dvv.appendChild(dv)
    }
    ct.appendChild(dvv);
  }
}
document.getElementById("bbbb").onclick=clr;
