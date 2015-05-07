$.base64.utf8encode = true;

function getCookie(name) {
  var arr = document.cookie.match(new RegExp("(^| )" + name + "=([^;]*)(;|$)"));
  if (arr !== null) return unescape(arr[2]);
  return null;
}
loc = window.location;
var im_ws = new IM.NewIm('ws://' + loc.host + "/ws", true);
im_ws.tdata = "";
im_ws.on("error", function(ev) {
  console.log("onerror");
});

im_ws.on("connect", function(ev) {
  im_ws.emit("li", {
    "token": "abc",
  });
});
im_ws.on("m", function(m) {
  console.log(m);
});
im_ws.on("close", function() {
  clearInterval(im_ws.timer);
});
im_ws.on("li", function(arg) {
  if (arg.code !== 0) {
    console.error("login error->", arg);
  } else {
    console.log("login success->", arg);
    im_ws.timer = setInterval(msg, 3000);
  }
});

function msg() {
  if (im_ws.closed) {
    return;
  }
  im_ws.sms(["S-Robot-0"], 0, "message->这是中文");
}