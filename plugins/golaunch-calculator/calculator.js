var clipboard = require('clipboard');
var events = require('events');
var math = require('mathjs');

var Plugin = function() {
  var self = this;

  this.on('request', function(data) {
    if (data.method == "init") {
      self.metadata = data.params;
    } else if (data.method == "query") {
      try {
        var answer = math.eval(data.params);
        if (typeof(answer) == "function") {
          return;
        }

        // don't care if the answer is only words
        if (/^[A-Za-z]+$/.test(answer)) {
          return;
        }

        //console.dir(answer);
        self.emit('response', {
          'result': [{
            image: self.metadata._icon,
            title: "" + answer,
            subtitle: "Copy this answer to clipboard",
            score: -1,
            query: data.params,
            id: self.metadata.id
          }]
        });
      } catch (e) {}
    } else if (data.method == "action") {
      clipboard.writeText(data.params.queryResult.title);
    }
  });
};

Plugin.prototype.__proto__ = events.EventEmitter.prototype;

module.exports = Plugin;
