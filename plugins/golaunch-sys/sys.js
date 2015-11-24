var electron = require('electron');
var app = electron.app;
var events = require('events');
var fuzzy = require('../../sdk/js/fuzzy');
var sdk = require('../../sdk/js/sdk');
var spawn = require('child_process').spawn;
var path = require('path');

var Plugin = function() {
  var self = this;

  var commands = {
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
    commands["quit"].image = img;
    commands["exit"].image = img;
  }, function(err) {
    console.error(err);
  });

  sdk.imageFileToEmbed(path.join(__dirname, "images", "restart.png")).then(function(img) {
    commands["restart"].image = img;
  }, function(err) {
    console.error(err);
  });

  this.on('request', function(data) {
    if (data.method == "init") {
      self.metadata = data.params;
    } else if (data.method == "query") {
      var results = [];
      for (var key in commands) {
        var match = fuzzy.match(data.params, key);
        if (match.success) {
          results.push({
            image: commands[key].image,
            title: commands[key].title,
            subtitle: commands[key].subtitle,
            score: match.score,
            query: data.params,
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
    } else if (data.method == "action") {
      if (data.params.queryResult.data == "restart") {
        const args = process.argv.slice(1);
        //const out = fs.openSync('./out.log', 'a');
        //const err = fs.openSync('./out.log', 'a');
        var child = spawn(process.argv[0], args, {
          detached: true,
          //stdio: [ 'ignore', out, err ]
        });
        child.unref();
        app.quit();
      } else if (data.params.queryResult.data == "quit" || data.params.queryResult.data == "exit") {
        app.quit();
      }
    }
  });
};

Plugin.prototype.__proto__ = events.EventEmitter.prototype;

module.exports = Plugin;
