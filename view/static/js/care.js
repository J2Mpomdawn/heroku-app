/*
ローディング時の表示非表示
*/
window.addEventListener("load",function(){
  document.getElementById("loading").remove();
},false)

/*
ajaxの処理
*/
const aj = document.getElementById("records");
var tbs = function(response){
      document.getElementById("tables").innerHTML=response;
    },
    rcs = function(response){
      document.getElementById("records").innerHTML=response;
      aj.firstElementChild.style.whiteSpace="pre";
      if (aj.firstElementChild.scrollHeight>200) {
        aj.firstElementChild.style.height="200px";
        aj.firstElementChild.style.overflowY="scroll";
      }
    },
    drs = function(response){
      document.getElementById("directries").innerHTML=response;
    },
    ajf = function(){
      const ajax = new XMLHttpRequest;
      ajax.open("POST","careajax",true);
      ajax.setRequestHeader("content-type", "application/x-www-form-urlencoded;charset=UTF-8");
      let n = this.arg;
      ajax.addEventListener("load",function(){
        switch (n) {
          case 0:
            tbs(this.responseText);
            break;
          case 1:
            if (/^[0-9A-Za-z]+$/.test(document.getElementById("tablename").value)) {
              rcs(this.responseText);
            } else {
              alert("入力ミス");
              return
            }
            break;
          case 5:
          drs(thi.responseText);
            break;
          default:
            console.log(this);
        }
      },false);
      ajax.send("f="+n+"&n="+document.getElementById("tablename").value+"&s="+document.getElementById("select").value+"&w="+document.getElementById("where").value);
    };

/*
table取得
*/
document.getElementById("looktables").addEventListener("click",{handleEvent:ajf,arg:0},false);

/*
record取得
*/
document.getElementById("lookrecords").addEventListener("click",{handleEvent:ajf,arg:1},false);

/*
tableリセット
*/
document.getElementById("getdatas").addEventListener("click",{handleEvent:ajf,arg:2},false);

/*mysqlテスト*/
document.getElementById("dbreq").addEventListener("click",{handleEvent:ajf,arg:3},false);

//time test
document.getElementById("tes").addEventListener("click",function(){
  let date1=new Date(document.getElementById("h3").innerText);
  console.log(date1);
  let date2=new Date();
  console.log(date2);
  console.log(Math.round((date2-date1)/360000)/10,"時間");
},false);

/*reborn 仮*/
document.getElementById("reborn").addEventListener("click",{handleEvent:ajf,arg:4},false);

/*tmp test*/
document.getElementById("cretes").addEventListener("click",{handleEvent:ajf,arg:5},false);
