var electron = require('electron');
var app = electron.app;
var BrowserWindow = electron.BrowserWindow;
//electron.crashReporter.start();
var globalShortcut = electron.globalShortcut;
var plugins = require('./plugins');
var themes = require('./themes');

global["pluginManager"] = new plugins.PluginManager();
global["themeManager"] = new themes.ThemeManager();

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
var mainWindow = null;

// process.on('uncaughtException', function uncaughtException(err) {
//   console.error(err.stack);
//   app.quit();
// })

// Quit when all windows are closed.
app.on('window-all-closed', function() {
  // On OS X it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform != 'darwin') {
    app.quit();
  }
});

var shouldQuit = app.makeSingleInstance(function(commandLine, workingDirectory) {
  if (mainWindow) {
    mainWindow.show();
    mainWindow.focus();
  }
  return true;
});

if (shouldQuit) {
  app.quit();
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
app.on('ready', function() {
  // Create the browser window.
  mainWindow = new BrowserWindow({
    width: 800,
    height: 80,
    'min-width': 500,
    'min-height': 80,
    'accept-first-mouse': true,
    'title-bar-style': 'hidden',
    'always-on-top': true,
    frame: false,
    resizable: false,
    show: false,
    'web-preferences': {
      'directWrite': process.platform === 'win32',
      'experimental-features': true,
      'overlayScrollbars': process.platform === 'win32'
    }
  });

  pluginManager.init();
  themeManager.init();

  //mainWindow.setMenu(null);

  // and load the index.html of the app.
  mainWindow.loadURL('file://' + __dirname + '/index.html');

  // Open the DevTools.
  //mainWindow.openDevTools();

  globalShortcut.register('ctrl+enter', function() {
    mainWindow.show();
  });

  mainWindow.on('blur', function() {
    mainWindow.hide();
  });

  // Emitted when the window is closed.
  mainWindow.on('closed', function() {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    mainWindow = null;

    pluginManager.shutdown();
  });
});

app.on('will-quit', function() {
  // Unregister all shortcuts.
  globalShortcut.unregisterAll();
});
