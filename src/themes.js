var events = require('events');
var path = require('path');
var fs = require('fs');
var domain = require('domain');
var concat = require('concat-stream');

function ThemeManager() {
  var self = this;
  var model = {
    themes: {}
  };

  var loadFromPath = function(dirPath) {
    var d = domain.create();
    d.on('error', function() {
      console.log(arguments);
    });
    d.run(function() {
      fs.createReadStream(path.join(dirPath, 'style.css'), 'utf8').pipe(concat(function(data) {
        model.themes[path.basename(dirPath)] = data;
      }));
    });
  };

  this.init = function() {
    // webContents.on('dom-ready', function() {
    //   webContents.send("theme-change", "HEYYY");
    // });

    var themesDir = path.join(__dirname, '../', 'themes');
    fs.readdir(themesDir, function(err, files) {
      for (var x = 0; x < files.length; x++) {
        var dirPath = path.join(themesDir, files[x]);
        var stats = fs.statSync(dirPath);
        if (!stats.isDirectory()) {
          continue;
        }

        loadFromPath(dirPath);
      }
    });
  };

  this.get = function(name) {
    if (!name) {
      name = "dracula";
    }
    return model.themes[name.toLowerCase()];
  };
}

ThemeManager.prototype.__proto__ = events.EventEmitter.prototype;

module.exports = {
  ThemeManager: ThemeManager
};
