var electron = require('electron');
var app = electron.app;
var events = require('events');
var fuzzy = require('../../sdk/js/fuzzy');
var sdk = require('../../sdk/js/sdk');
var spawn = require('child_process').spawn;
var path = require('path');

var Plugin = function() {
  var self = this;

  self.commands = {
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
    self.commands["quit"].icon = img;
    self.commands["exit"].icon = img;
  }, function(err) {
    console.error(err);
  });

  sdk.imageFileToEmbed(path.join(__dirname, "images", "restart.png")).then(function(img) {
    self.commands["restart"].icon = img;
  }, function(err) {
    console.error(err);
  });
};

Plugin.prototype.__proto__ = events.EventEmitter.prototype;

Plugin.prototype.init = function init(metadata) {
  this.metadata = metadata;
};

Plugin.prototype.query = function query(query) {
  var self = this;

  var results = [];
  for (var key in self.commands) {
    var match = fuzzy.match(query, key);
    if (match.success) {
      results.push({
        icon: self.commands[key].icon,
        title: self.commands[key].title,
        subtitle: self.commands[key].subtitle,
        score: match.score,
        query: query,
        id: self.metadata.id,
        data: key
      });
    }
  }

  if (results.length > 0) {
    self.emit('response', {
      'result': results
    });
  }
};

Plugin.prototype.action = function action(action) {
  if (action.queryResult.data == "restart") {
    const args = process.argv.slice(1);
    //const out = fs.openSync('./out.log', 'a');
    //const err = fs.openSync('./out.log', 'a');
    var child = spawn(process.argv[0], args, {
      detached: true,
      //stdio: [ 'ignore', out, err ]
    });
    child.unref();
    app.quit();
  } else if (action.queryResult.data == "quit" || action.queryResult.data == "exit") {
    app.quit();
  }
};

module.exports = Plugin;
