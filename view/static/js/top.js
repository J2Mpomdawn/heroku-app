/*
文字を入れてsendボタンを押したら/room/○○にとぶ
最終的には消す
*/
function move() {
  //入力されたものを取得
  const name = document.getElementById("user");
  //空じゃなかったら
  if (name.value != "") {
    window.location = "/room/" + name.value;
  }
}

document.getElementById("nb").onclick = move;

/*
ローディング時の表示非表示
*/
window.addEventListener("load",function(){
  document.getElementById("loading").remove();
},false)

/*
リンクにカーソルを合わせたら説明文の色を変える
*/
const links = document.getElementById("links");
function changeColor(e){
  const target = e.target;
  if (target.tagName == "A"){
    let expl = document.getElementById(target.href.substr(22)+"Explanation");
    if (e.type == "mouseover"){
      /*奈緒カラー*/
      expl.style.color = "#788bc5";
      expl.style.transition = "0.5s";
    }
    if (e.type == "mouseout"){
      expl.style.color = "";
      expl.style.transition = "0.2s";
    }
  }
}

links.onmouseover = links.onmouseout = changeColor;
