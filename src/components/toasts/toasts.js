var electron = require('electron');
var BrowserWindow = electron.BrowserWindow;
var ipc = electron.ipcMain;

var serialize = function(obj, prefix) {
  var str = [];
  for(var p in obj) {
    if (obj.hasOwnProperty(p)) {
      var k = prefix ? prefix + "[" + p + "]" : p, v = obj[p];
      str.push(typeof v == "object" ?
        serialize(v, k) :
        encodeURIComponent(k) + "=" + encodeURIComponent(v));
    }
  }
  return str.join("&");
};

var show = function(currentWindow, opts) {
	opts.timeout = opts.timeout || 5000;

	var self = this;
	this.window = new BrowserWindow({
		width: opts.width,
		title: opts.title,
		skipTaskbar: true,
		alwaysOnTop: true,
		frame: false,
		show: false
	});

	var pos = currentWindow.getPosition();
	var display = electron.screen.getDisplayNearestPoint({x:pos[0], y:pos[1]});

	var htmlFile = opts.htmlFile || 'file://' + __dirname + '/default.html?' + serialize(opts);
	this.window.loadURL(htmlFile);

	this.window.webContents.on('did-finish-load', function() {
		if (self.window) {
			var width = self.window.getSize()[0];
			var height = self.window.getSize()[1];
			self.window.setPosition(display.workAreaSize.width - width - 4, display.workAreaSize.height - height);
			self.window.show();
		}
	});
};

var Toasts = function(){
	return this;
};

Toasts.prototype.init = function(currentWindow) {
	var self = this;
	ipc.on('golaunch-toast', function(event, opts) {
	  show.call(self, currentWindow, opts);
	});
};

module.exports = Toasts;
