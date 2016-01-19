var plugin = require('../../sdk/js/plugin');
var events = require('events');
var fuzzy = require('../../sdk/js/fuzzy');
var sdk = require('../../sdk/js/sdk');
var path = require('path');

var Plugin = function() {
  var p = this;

  p.client = new plugin.Client();

  p.commands = {
    "restart": {
      title: "Restart GoLaunch",
      subtitle: "Restarts the application immediately.",
      icon: "restart.png"
    },
    "settings": {
      title: "GoLaunch Settings",
      subtitle: "Configure GoLaunch and plugins.",
      icon: "settings.png"
    },
    "quit": {
      title: "Quit GoLaunch",
      subtitle: "Quits the program for good!",
      icon: "exit.png"
    },
    "exit": {
      title: "Exit GoLaunch",
      subtitle: "Exits the program for good!",
      icon: "exit.png"
    }
  };
};

Plugin.prototype.__proto__ = events.EventEmitter.prototype;

Plugin.prototype.init = function init(metadata) {
  var p = this;
  this.metadata = metadata;

  for (var key in p.commands) {
    if (p.commands[key].icon) {
      !function outer(key) {
        sdk.imageFileToEmbed(path.join(__dirname, "images", p.commands[key].icon)).then(function(img) {
          p.commands[key].icon = img;
        }, function(err) {
          p.client.call("log", err);
          console.error(err);
        });
      }(key);
    }
  }
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
  } else if (action.queryResult.data == "settings") {
    this.client.call("opensettings", null);
  }
};

var server = new plugin.Server();
server.register(new Plugin());
server.serve();
