var IM = (function() {
    function IM(url, recon) {
        this.url = url;
        this.EV = {};
        this.recon = recon;
        this.times = 0;
        this.tdata = "";
        this.closed = true;
        this.showlog = true;
        this.mrs = {};
        this.conn();
    }
    IM.WIN_SEQ = "^-^";
    IM.prototype.log = function() {
        if (!this.showlog) {
            return;
        }
        switch (arguments.length) {
            case 1:
                console.log(arguments[0]);
                break;
            case 2:
                console.log(arguments[0], arguments[1]);
                break;
            case 3:
                console.log(arguments[0], arguments[1], arguments[2]);
                break;
            case 4:
                console.log(arguments[0], arguments[1], arguments[2], arguments[3]);
                break;
            default:
                console.log(arguments[0]);
                break;
        }
    };
    IM.prototype.conn = function() {
        var tar = this;
        this.ws = new WebSocket(this.url);
        this.ws.onerror = function(ev) {
            tar.onerror(ev);
        };
        this.ws.onopen = function(ev) {
            tar.onopen(ev);
        };
        this.ws.onmessage = function(ev) {
            tar.onmessage(ev);
        };
        this.ws.onclose = function(ev) {
            tar.onclose(ev);
        };
        this.log("connecting->" + this.url);
    };
    IM.prototype.onerror = function(ev) {
        this.log("onerror->", ev);
        if (this.EV.error) {
            this.EV.error(ev);
        }
        if (this.EV.err) {
            this.EV.err(ev);
        }
    };
    IM.prototype.onopen = function(ev) {
        this.log("onopen->", this.url);
        if (this.EV.connect) {
            this.EV.connect(ev);
        }
        if (this.EV.open) {
            this.EV.open(ev);
        }
        this.times = 0;
        this.closed = false;
    };
    IM.prototype.onmessage = function(ev) {
        this.tdata += ev.data;
        if (this.tdata.substr(this.tdata.length - 1) != "\n") {
            return;
        }
        var tdata = this.tdata;
        var cmds = tdata.split("^-^");
        this.tdata = "";
        if (cmds.length < 2) {
            this.log("receive invalid data: " + tdata);
            return;
        }
        var args = JSON.parse(cmds[1]);
        if (cmds[0] == "m") {
            var tim = this;
            setTimeout(function() {
                tim.emit("mr", {
                    i: args.i,
                    a: args.a,
                });
            }, 1000);
            if (this[args.i]) {
                return;
            }
            args.c = $.base64.atob(args.c, true);
            this[args.i] = true;
        }
        this.on_(cmds[0], args);
    };
    IM.prototype.onclose = function(ev) {
        this.log("onclose->", ev);
        this.closed = true;
        if (this.EV.close) {
            this.EV.close(ev);
        }
        this.log("ws is closed..");
        if (this.recon) {
            this.log("ws will reconnect after " + (this.times * 100) + " ms");
            var tim = this;
            setTimeout(function() {
                tim.conn();
            }, this.times * 300);
            this.times++;
        }
    };
    IM.prototype.on = function(name, func) {
        this.EV[name] = func;
    };
    IM.prototype.on_ = function(name, args) {
        if (this.EV[name]) {
            this.EV[name](args);
        }
    };
    IM.prototype.emit = function(name, args) {
        this.ws.send(name + IM.WIN_SEQ + JSON.stringify(args) + "\n");
    };
    //send text message.
    IM.prototype.sms = function(r, t, c) {
        if (!r || r.length < 1 || t === undefined || !c) {
            this.log("sms args error", r, t, c);
            return;
        }
        this.emit("m", {
            r: r,
            t: t,
            c: $.base64.btoa(c),
        });
    };
    IM.prototype.sms2 = function(m) {
        if (!m.r || m.r.length < 1 || m.t === undefined || !m.c) {
            this.log("sms args error", m);
            return;
        }
        this.emit("m", {
            r: m.r,
            t: m.t,
            c: $.base64.btoa(m.c),
        });
    };
    IM.NewIm = function(url, recon) {
        return new IM(url, recon);
    };
    window.IM = IM;
    return IM;
})();