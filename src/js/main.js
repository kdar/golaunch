var m = require('mithril');

var Launcher = require('./pages/launcher/launcher');
// var Settings = require('./pages/settings/settings');

function persist(component) {
  var persistor = {
    controller: function() {
      var output = component.controller.apply (this, arguments) || this;

      persistor.controller = function() {
        return output;
      };

      return output;
    },
    view: component.view
  };

  return persistor;
}

m.route(document.getElementById('app'), "/", {
  "/": persist(Launcher)
  // "/settings": persist(Settings)
});
