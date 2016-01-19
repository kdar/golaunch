const electron = require('electron');
const ipcRenderer = electron.ipcRenderer;

var windowID = electron.remote.getCurrentWindow().id;

process.on('uncaughtException', function uncaughtException(err) {
  ipcRenderer.send("plugin-error-" + windowID, err.stack);
});

function Client() {};

Client.prototype.call = function call(method, params) {
  // process.send({
  //   'method': method,
  //   'params': params
  // });
  ipcRenderer.send("plugin-" + windowID, {
    'method': method,
    'params': params
  });
};

function Server() {};

Server.prototype.register = function register(plugin) {
  this.p = plugin;
};

Server.prototype.serve = function serve() {
  var s = this;
  // process.on('message', function(m) {
  //   switch (m.method) {
	// 	case "init":
	// 		s.p.init(m.params)
  //     break;
	// 	case "query":
	// 		s.p.query(m.params)
  //     break;
	// 	case "action":
	// 		s.p.action(m.params)
  //     break;
	// 	}
  // });
  ipcRenderer.on("plugin-" + windowID, function(event, m) {
    switch (m.method) {
		case "init":
			s.p.init(m.params)
      break;
		case "query":
			s.p.query(m.params)
      break;
		case "action":
			s.p.action(m.params)
      break;
		}
  });
};

module.exports = {
  Client: Client,
  Server: Server
};
