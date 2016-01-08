var clipboard = require('clipboard');
var events = require('events');
var math = null;

var Plugin = function() {};

Plugin.prototype.__proto__ = events.EventEmitter.prototype;

Plugin.prototype.init = function init(metadata) {
  this.metadata = metadata;
};

Plugin.prototype.query = function query(query) {
  // delay load math module
  if (math == null) {
    math = require('mathjs');
  }

  try {
    var answer = math.eval(query);
    if (typeof(answer) == "function") {
      return;
    }

    // don't care if the answer is only words
    if (/^[A-Za-z]+$/.test(answer)) {
      return;
    }

    //console.dir(answer);
    this.emit('response', {
      'result': [{
        icon: this.metadata._icon,
        title: "" + answer,
        subtitle: "Copy this answer to clipboard",
        score: -1,
        query: query,
        id: this.metadata.id
      }]
    });
  } catch (e) {}
};

Plugin.prototype.action = function action(action) {
  clipboard.writeText(action.queryResult.title);
};

module.exports = Plugin;
