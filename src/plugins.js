const electron = require('electron');
const app = electron.app;
const ipcMain = electron.ipcMain;
const BrowserWindow = electron.BrowserWindow;
// const remote = electron.remote;
const path = require('path');
const spawn = require('child_process').spawn;
const fs = require('fs');
const toml = require('toml');
const concat = require('concat-stream');
const events = require('events');
const domain = require('domain');
const sdk = require('../sdk/js/sdk');
const child_process = require('child_process');

var elapsedTime = function(hrtime) {
  var precision = 3; // 3 decimal places
  var elapsed = process.hrtime(hrtime)[1] / 1000000; // divide by a million to get nano to milli
  return process.hrtime(hrtime)[0] + "s" + elapsed.toFixed(precision) + "ms";
}

var PluginManager = function() {
  var self = this;
  var model = {
    plugins: {}
  };

  // for queries
  var lastQuery = null;
  var queryResults = [];
  var queryProgress = {
    count: 0,
    current: 0
  }

  // a function that is called when we get some sort of plugin data
  function pluginData(data) {
    // if (!data.result || data.result.length == 0) {
    //   return;
    // }

    switch (data.method) {
    // This is here so plugins not using JS can run JS code.
    case 'eval':
      eval(data.params);
      break;

    case 'log':
      console.log(data.params);
      break;

    case 'noqueryresults':
      queryProgress.current += 1;
      self.emit('progress-update', queryProgress);
      break;

    case 'queryresults': // just query results
      queryProgress.current += 1;
      self.emit('progress-update', queryProgress);

      // console.log(data.params);
      queryResults.push.apply(queryResults, data.params);

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

      break;

    default:
      console.log("unrecognized plugin data:", data);
    }
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
      var _process = plugin._process;
      process.nextTick(function () {
        _process.send(data);
      });
      break;
    case 'electron':
      var _process = plugin._process;
      _process.webContents.send("plugin-" + _process.id, data);
    }
  };

  var loadFromPath = function(dirPath) {
    var d = domain.create();
    d.on('error', function(err) {
      console.error(err.stack);
    });
    d.run(function() {
      fs.createReadStream(path.join(dirPath, 'plugin.toml'), 'utf8').pipe(concat(function(data) {
        var parsed = toml.parse(data);

        var startTime = process.hrtime();
        console.log(parsed.name + ": loading...");

        if (parsed.enabled === false) {
          return;
        }

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
            pluginData(JSON.parse(data));
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
          queryProgress.count += 1;
          break;
        case 'js':
          var child = child_process.fork(path.join(dirPath, parsed.main));

          child.on('message', function(m) {
            pluginData(m);
          });

          child.send({
            method: "init",
            params: parsed
          });

          parsed._process = child;

          model.plugins[parsed.id] = parsed;
          queryProgress.count += 1;
          break;
        case 'electron':
          var win = new BrowserWindow({
        		show: false
        	});
        	win.loadURL(`file://` + path.join(dirPath, parsed.main));

          ipcMain.on("plugin-" + win.id, function(event, arg) {
            pluginData(arg);
          });

          win.webContents.once('did-finish-load', function() {
            win.webContents.send("plugin-" + win.id, {
              method: "init",
              params: parsed
            });
          });

          parsed._process = win;

          model.plugins[parsed.id] = parsed;
          queryProgress.count += 1;
          break;
        }

        console.log(parsed.name + ": took " + elapsedTime(startTime));
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
      queryProgress.current = 0;
      process.nextTick(function () {
        pluginRequest({
          "method": "query",
          "params": query
        });
      });
    }
  };

  this.clearQuery = function() {
    lastQuery = null;
    queryResults = [];
    queryProgress.current = 0;
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
