$.base64.utf8encode = true;

function getCookie(name) {
  var arr = document.cookie.match(new RegExp("(^| )" + name + "=([^;]*)(;|$)"));
  if (arr !== null) return unescape(arr[2]);
  return null;
}
loc = window.location;
var im_ws = new IM.NewIm('ws://im.dev.jxzy.com/im', true);
im_ws.tdata = "";
im_ws.on("error", function(ev) {
  console.log("onerror");
});

im_ws.on("connect", function(ev) {
  im_ws.emit("li", {
    "token": "b66dc44cd9882859d84670604ae276e6-ec6e5337-8852-41ee-9953-135158a7a1ac",
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
  // im_ws.sms(["U-202996"], 0, "message->这是中文");
}