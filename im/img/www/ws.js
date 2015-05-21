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
    "token": "69657ec87da52096e4f8e354ec940404-b0fe5f7c-dff6-4e53-8a9b-0b78083fa799",
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
  im_ws.emit("ur", {});
});

function msg() {
  if (im_ws.closed) {
    return;
  }
  // im_ws.sms(["U-202996"], 0, "message->这是中文");
}