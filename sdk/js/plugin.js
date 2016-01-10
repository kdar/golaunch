function Client() {};

// Client.prototype.queryResults = function queryResults(results) {
//   process.send({
//     'result': results
//   });
// };

Client.prototype.call = function call(method, params) {
  // var params = null;
  // if (arguments.length == 2) {
  //   params = arguments[1];
  // } else if (arguments.length > 2) {
  //   params = arguments.splice(1);
  // }

  process.send({
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
  process.on('message', function(m) {
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
