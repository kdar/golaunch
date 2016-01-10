var events = require('events');
var fuzzy = require('../../sdk/js/fuzzy');
var sdk = require('../../sdk/js/sdk');
var path = require('path');
var plugin = require('../../sdk/js/plugin');

var Plugin = function() {
  var p = this;

  p.client = new plugin.Client();

  p.commands = {
    "restart": {
      title: "Restart GoLaunch",
      subtitle: "Restarts the application immediately."
    },
    "quit": {
      title: "Quit GoLaunch",
      subtitle: "Quits the program for good!"
    },
    "exit": {
      title: "Exit GoLaunch",
      subtitle: "Exits the program for good!"
    }
  };

  sdk.imageFileToEmbed(path.join(__dirname, "images", "exit.png")).then(function(img) {
    p.commands["quit"].icon = img;
    p.commands["exit"].icon = img;
  }, function(err) {
    console.error(err);
  });

  sdk.imageFileToEmbed(path.join(__dirname, "images", "restart.png")).then(function(img) {
    p.commands["restart"].icon = img;
  }, function(err) {
    console.error(err);
  });
};

Plugin.prototype.__proto__ = events.EventEmitter.prototype;

Plugin.prototype.init = function init(metadata) {
  this.metadata = metadata;
};

Plugin.prototype.query = function query(query) {
  var p = this;

  var results = [];
  for (var key in p.commands) {
    var match = fuzzy.match(query, key);
    if (match.success) {
      results.push({
        icon: p.commands[key].icon,
        title: p.commands[key].title,
        subtitle: p.commands[key].subtitle,
        score: match.score,
        query: query,
        id: p.metadata.id,
        data: key
      });
    }
  }

  if (results.length > 0) {
    p.client.call("queryresults", results);
    return;
  }

  this.client.call("noqueryresults", null);
};

Plugin.prototype.action = function action(action) {
  if (action.queryResult.data == "restart") {
    function restart() {
      const args = process.argv.slice(1);
      //const out = fs.openSync('./out.log', 'a');
      //const err = fs.openSync('./out.log', 'a');
      var child = spawn(process.argv[0], args, {
        detached: true,
        //stdio: [ 'ignore', out, err ]
      });
      child.unref();
      app.quit();
    }

    this.client.call("eval", "("+restart.toString()+")()");
  } else if (action.queryResult.data == "quit" || action.queryResult.data == "exit") {
    this.client.call("eval", "app.quit();");
  }
};

var server = new plugin.Server();
server.register(new Plugin());
server.serve();
