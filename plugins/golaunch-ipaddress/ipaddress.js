var electron = require('electron');
var app = electron.app;
var events = require('events');
var clipboard = require('clipboard');
var locip = require('ip');
var http = require('http');

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

var Plugin = function() {};

Plugin.prototype.__proto__ = events.EventEmitter.prototype;

Plugin.prototype.init = function init(metadata) {
  this.metadata = metadata;
};

Plugin.prototype.query = function query(query) {
  var self = this;

  if (query.startsWith("ipaddress")) {
    extip(function(err, eip) {
      var results = [];

      if (!err) {
        results.push({
          icon: self.metadata._icon,
          title: "External IP: " + eip,
          subtitle: "Copy to clipboard",
          score: -1,
          query: query,
          id: self.metadata.id,
          data: eip
        });
      }

      var localip = locip.address();
      results.push({
        icon: self.metadata._icon,
        title: "Local IP: " + localip,
        subtitle: "Copy to clipboard",
        score: -1,
        query: query,
        id: self.metadata.id,
        data: localip
      });

      self.emit('response', {
        'result': results
      });
    });
  }
};

Plugin.prototype.action = function action(action) {
  clipboard.writeText(action.queryResult.data);
};

module.exports = Plugin;
