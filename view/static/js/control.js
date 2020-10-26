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
    }
    ajf = function(){
      const ajax = new XMLHttpRequest;
      ajax.open("POST","controlajax",true);
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
document.getElementById("resettables").addEventListener("click",{handleEvent:ajf,arg:2},false)
