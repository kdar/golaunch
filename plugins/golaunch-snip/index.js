// Minicap -captureregselect -clipimage -exit

var events = require('events');
var spawn = require('child_process').spawn;
var os = require('os');
var path = require('path');
var plugin = require('../../sdk/js/plugin');

var Plugin = function() {
  this.client = new plugin.Client();
};

Plugin.prototype.__proto__ = events.EventEmitter.prototype;

Plugin.prototype.init = function init(metadata) {
  this.metadata = metadata;
};

Plugin.prototype.query = function query(query) {
  var self = this;

  if (query.startsWith("snip")) {
    this.client.call("queryresults", [{
      icon: self.metadata._icon,
      title: "Screen Snipping",
      subtitle: "Snip screen to clipboard",
      score: -1,
      query: query,
      id: self.metadata.id
    }]);

    return;
  }

  this.client.call("noqueryresults", null);
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

var server = new plugin.Server();
server.register(new Plugin());
server.serve();
