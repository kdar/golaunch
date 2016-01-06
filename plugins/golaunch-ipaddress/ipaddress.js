var electron = require('electron');
var app = electron.app;
var events = require('events');
var clipboard = require('clipboard');
var locip = require('ip');
var extip = require('external-ip')();

var Plugin = function() {};

Plugin.prototype.__proto__ = events.EventEmitter.prototype;

Plugin.prototype.init = function init(metadata) {
  this.metadata = metadata;
};

Plugin.prototype.query = function query(query) {
  var self = this;

  if (query.startsWith("ipaddress")) {
    extip(function (err, eip) {
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
