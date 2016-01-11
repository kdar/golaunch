var events = require('events');
var ncp = require("copy-paste");
var locip = require('ip');
var http = require('http');
var plugin = require('../../sdk/js/plugin');

var extip = function(cb) {
  var services = [
    'http://curlmyip.com/',
    'http://icanhazip.com/',
    'http://myexternalip.com/raw'
  ];

  var done = false;
  var errors = [];
  var requests;

  var abort = function(requests) {
    process.nextTick(function() {
      requests.forEach(function(request) {
        request.abort();
      });
    });
  };

  var onResponse = function(err, ip) {
    if (done) {
      return;
    }
    if (err) {
      errors.push(err);
    }
    if (ip) {
      done = true;
      abort(requests); //async
      return cb(null, ip);
    }
    if (errors.length === services.length) {
      done = true;
      abort(requests); //async
      return cb(errors, null);
    }
  };

  requests = services.map(function(service) {
    var req = http.request(service, function(res) {
      res.on('data', function(data) {
        onResponse(null, data);
      });
    });
    req.on('error', function(e) {
      onResponse(e.message, null);
    });
    req.end();
    return req;
  });
};

var Plugin = function() {
  this.client = new plugin.Client();
};

Plugin.prototype.__proto__ = events.EventEmitter.prototype;

Plugin.prototype.init = function init(metadata) {
  this.metadata = metadata;
};

Plugin.prototype.query = function query(query) {
  var p = this;

  if (query.startsWith("ipaddress")) {
    extip(function(err, eip) {
      var results = [];

      if (!err) {
        results.push({
          icon: p.metadata._icon,
          title: "External IP: " + eip,
          subtitle: "Copy to clipboard",
          score: -1,
          query: query,
          id: p.metadata.id,
          data: String(eip).trim()
        });
      }

      var localip = locip.address();
      results.push({
        icon: p.metadata._icon,
        title: "Local IP: " + localip,
        subtitle: "Copy to clipboard",
        score: -1,
        query: query,
        id: p.metadata.id,
        data: localip
      });

      p.client.call("queryresults", results);
    });
    return;
  }

  this.client.call("noqueryresults", null);
};

Plugin.prototype.action = function action(action) {
  ncp.copy(action.queryResult.data);
};

var server = new plugin.Server();
server.register(new Plugin());
server.serve();
