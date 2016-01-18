const electron = require('electron');
const ipcRenderer = electron.ipcRenderer;

function Client() {
  this.windowID = electron.remote.getCurrentWindow().id;
};

Client.prototype.call = function call(method, params) {
  // process.send({
  //   'method': method,
  //   'params': params
  // });
  ipcRenderer.send("plugin-" + this.windowID, {
    'method': method,
    'params': params
  });
};

function Server() {
  this.windowID = electron.remote.getCurrentWindow().id;
};

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
  ipcRenderer.on("plugin-" + this.windowID, function(event, m) {
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
