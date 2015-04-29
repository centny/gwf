$.base64.utf8encode = true;

function getCookie(name) {
  var arr = document.cookie.match(new RegExp("(^| )" + name + "=([^;]*)(;|$)"));
  if (arr != null) return unescape(arr[2]);
  return null;
}
loc = window.location;
var im_ws = new WebSocket('ws://' + loc.host + "/ws");
im_ws.tdata = "";
im_ws.onerror = function(event) {
  console.log("onerror");
};

im_ws.onopen = function(event) {
  console.log("onopen");
  send("li^-^" + JSON.stringify({
    "token": "abc",
  }));
};

im_ws.onmessage = function(ev) {
  im_ws.tdata += ev.data;
  if (im_ws.tdata.substr(im_ws.tdata.length - 1) != "\n") {
    return;
  }
  var tdata = im_ws.tdata;
  var cmds = tdata.split("^-^");
  im_ws.tdata = "";
  if (cmds.length < 2) {
    console.log("receive invalid data:" + tdata);
    return;
  }
  var args = JSON.parse(cmds[1]);
  if (cmds[0] == "m") {
    args.c = $.base64.atob(args.c, true)
  }
  console.log(cmds[0], args);
};

function send(data) {
  console.log("sending->" + data);
  im_ws.send(data + "\n");
};

function msg() {
  send("m^-^" + JSON.stringify({
    r: ["S-Robot-0"],
    c: $.base64.btoa("message->这是中文"),
  }));
}
setInterval(msg, 3000);