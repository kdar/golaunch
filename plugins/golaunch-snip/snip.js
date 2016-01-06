// Minicap -captureregselect -clipimage -exit

var electron = require('electron');
var app = electron.app;
var events = require('events');
var spawn = require('child_process').spawn;
var os = require('os');
var path = require('path');

var Plugin = function() {};

Plugin.prototype.__proto__ = events.EventEmitter.prototype;

Plugin.prototype.init = function init(metadata) {
  this.metadata = metadata;
};

Plugin.prototype.query = function query(query) {
  var self = this;
  
  if (query.startsWith("snip")) {
    self.emit('response', {
      'result': [{
        icon: self.metadata._icon,
        title: "Screen Snipping",
        subtitle: "Snip screen to clipboard",
        score: -1,
        query: query,
        id: self.metadata.id
      }]
    });
  }
};

Plugin.prototype.action = function action(action) {
  var platform = os.platform();
  var arch = os.arch();
  switch (platform + arch) {
    case "win32ia32":
    case "win32x64":
      var child = spawn(path.join(__dirname, "platform", "win32-ia32", "MiniCap.exe"), ['-captureregselect', '-clipimage', '-exit'], {
        detached: true
      });
      child.unref();
    break;
  }
};

module.exports = Plugin;
