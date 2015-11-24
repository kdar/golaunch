var electron = require('electron');
var app = electron.app;
var events = require('events');
var clipboard = require('clipboard');
var locip = require('ip');
var extip = require('external-ip')();

var Plugin = function() {
  var self = this;

  this.on('request', function(data) {
    if (data.method == "init") {
      self.metadata = data.params;
    } else if (data.method == "query") {
      if (data.params.startsWith("ipaddress")) {
        extip(function (err, eip) {
          var results = [];

          if (!err) {
            results.push({
              image: self.metadata._icon,
              title: "External IP: " + eip,
              subtitle: "Copy to clipboard",
              score: -1,
              query: data.params,
              id: self.metadata.id,
              data: eip
            });
          }

          var localip = locip.address();
          results.push({
            image: self.metadata._icon,
            title: "Local IP: " + localip,
            subtitle: "Copy to clipboard",
            score: -1,
            query: data.params,
            id: self.metadata.id,
            data: localip
          });

          self.emit('response', {
            'result': results
          });
        });
      }
    } else if (data.method == "action") {
      clipboard.writeText(data.params.queryResult.data);
    }
  });
};

Plugin.prototype.__proto__ = events.EventEmitter.prototype;

module.exports = Plugin;
