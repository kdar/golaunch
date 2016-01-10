var ncp = require("copy-paste");
var math = require('mathjs');
var plugin = require('../../sdk/js/plugin');

var Plugin = function() {
  this.client = new plugin.Client();
};

Plugin.prototype.init = function init(metadata) {
  this.metadata = metadata;
};

Plugin.prototype.query = function query(query) {
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
    this.client.queryResults([{
      icon: this.metadata._icon,
      title: "" + answer,
      subtitle: "Copy this answer to clipboard",
      score: -1,
      query: query,
      id: this.metadata.id
    }]);
  } catch (e) {}
};

Plugin.prototype.action = function action(action) {
  ncp.copy(action.queryResult.title);
};

var server = new plugin.Server();
server.register(new Plugin());
server.serve();
