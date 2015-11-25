var electron = require('electron');
var app = electron.app;
var path = require('path');
var spawn = require('child_process').spawn;
var fs = require('fs');
var toml = require('toml');
var concat = require('concat-stream');
var events = require('events');
var domain = require('domain');
var sdk = require('../sdk/js/sdk');

var PluginManager = function() {
  var self = this;
  var model = {
    plugins: {}
  };

  // for queries
  var lastQuery = null;
  var queryResults = [];

  function response(data) {
    // process.nextTick(function () {
    //   self.emit('plugin_response', data);
    // });

    // if (!data.result || data.result.length == 0) {
    //   return;
    // }

    queryResults.push.apply(queryResults, data.result);

    queryResults.sort(function(a, b) {
      if (a.score == -1) {
        return -1;
      } else if (b.score == -1) {
        return 1;
      }
      return b.score - a.score;
    });

    process.nextTick(function () {
      self.emit('query-results', queryResults);
    });
  };

  var pluginRequest = function(data) {
    for (var key in model.plugins) {
      directPluginRequest(data, model.plugins[key]);
    }
  };

  var directPluginRequest = function(data, plugin) {
    switch (plugin.type) {
    case 'stdio':
      var _process = plugin._process;
      process.nextTick(function () {
        _process.stdin.write(JSON.stringify(data));
      });
      break;
    case 'js':
      var _object = plugin._object;
      process.nextTick(function () {
        _object.emit('request', data);
      });
      break;
    }
  };

  var loadFromPath = function(dirPath) {
    var d = domain.create();
    d.on('error', function(err) {
      console.error(err);
    });
    d.run(function() {
      fs.createReadStream(path.join(dirPath, 'plugin.toml'), 'utf8').pipe(concat(function(data) {
        var parsed = toml.parse(data);

        if (parsed.icon) {
          parsed._icon = sdk.imageFileToEmbedSync(path.join(dirPath, parsed.icon));
        }

        parsed._appdata = path.join(app.getPath("userData"), parsed.name);

        switch (parsed.type) {
        case 'stdio':
          var plugin = spawn(path.join(dirPath, parsed.main), parsed.arguments || [], {
            cwd: dirPath
          });
          plugin.stdout.setEncoding('utf8');
          plugin.stdin.setEncoding('utf8');
          plugin.stderr.setEncoding('utf8');

          // plugin.on('error', function() {
          //   console.log("ui:", arguments);
          // });
          //
          plugin.stdout.on('data', function(data) {
            response(JSON.parse(data));
          });

          plugin.stderr.on('data', function(data) {
            process.stdout.write("backend: " + data);
          });

          plugin.stdin.write(JSON.stringify({
            method: "init",
            params: parsed
          }));

          parsed._process = plugin;

          model.plugins[parsed.id] = parsed;
          break;
        case 'js':
          var cls = require(path.join(dirPath, parsed.main));
          parsed._object = new cls();

          parsed._object.on('response', function(data) {
            response(data);
          });

          parsed._object.emit('request', {
            "method": "init",
            "params": parsed
          });

          model.plugins[parsed.id] = parsed;
          break;
        }
      }));
    });
  };

  this.init = function() {
    var pluginsDir = path.join(__dirname, '../', 'plugins');
    fs.readdir(pluginsDir, function(err, files) {
      for (var x = 0; x < files.length; x++) {
        var dirPath = path.join(pluginsDir, files[x]);
        var stats = fs.statSync(dirPath);
        if (!stats.isDirectory()) {
          continue;
        }

        loadFromPath(dirPath);
      }
    });
  };

  this.pluginQuery = function(query) {
    if (lastQuery != query) {
      lastQuery = query;
      queryResults = [];
      process.nextTick(function () {
        pluginRequest({
          "method": "query",
          "params": query
        });
      });
    }
  };

  this.pluginAction = function(data) {
    process.nextTick(function () {
      directPluginRequest({
        "method": "action",
        "params": data
      }, model.plugins[data.queryResult.id]);
    });
  };

  this.shutdown = function() {
    for (var x = 0; x < model.plugins.length; x++) {
      switch (model.plugins[x].type) {
      case 'stdio':
        model.plugins[x]._process.kill('SIGINT');
        break;
      }
    }
  };
};

PluginManager.prototype.__proto__ = events.EventEmitter.prototype;

module.exports = {
  PluginManager: PluginManager
};
