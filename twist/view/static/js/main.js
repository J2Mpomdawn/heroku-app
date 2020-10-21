/*
10秒後にリロードする
リロードしたらまた10秒後にやることになるから、
実質　10秒毎にリロードする

function rl(){
  window.location="/";
  //window.location.reload();
}

setTimeout("rl()",10000);
*/

/*
id=testのdivタグに"<p>test</p>"を追加する


var newtweet = document.getElementById("test");

function add(){
  newtweet.insertAdjacentHTML("afterbegin","<p>test</p>");
}
*/
/*
10秒毎に関数を実行する
コールバック関数やから非同期

function loop(max,i){
  if(i<max){
    add();
    setTimeout(function(){loop(max,++i)},10000);
  }
}
loop(10,0);
*/

//document.getElementById("bt").onclick=add;

function regularly(){
  let ajax = new XMLHttpRequest;
  ajax.open("GET","/tweet",true);
  ajax.onload=function (){
    document.getElementById("test").innerHTML=ajax.responseText;
  }
  ajax.send(null);
}

(function loop(){
  regularly();
  setTimeout(function(){loop()},10000);
}());
