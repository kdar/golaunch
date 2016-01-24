const m = require('mithril');
const General = require('./general')
const Plugins = require('./plugins')
const About = require('./about')
const electron = require('electron');
const remote = electron.remote;

var Settings = {
  config: function config(el, isInitialized) {
    if (!isInitialized) {
      remote.getCurrentWindow().setSize(800, 500);

      $(function() {
        $('.menu .item').tab();
      });
    }
  },

  controller: function controller() {
    var ctrl = this;

    ctrl.onDoneClick = function() {
      window.close();
    };
  },

  view: function view(ctrl) {
    // <div id="settings" config={Settings.config.bind(ctrl)} class="ui attached segment pushable">
    return <div id="settings" config={Settings.config.bind(ctrl)}>
      <div class="ui visible inverted left vertical sidebar menu">
        <div class="item">
          <h4 class="ui header blue inverted">User Settings</h4>
        </div>
        <a class="active item" data-tab="general">
          <i class="home icon"></i>
          General
        </a>
        <a class="item" data-tab="plugins">
          <i class="plug icon"></i>
          Plugins
        </a>
        <a class="item" data-tab="about">
          <i class="info circle icon"></i>
          About
        </a>
      </div>
      <div class="pusher">
        <div class="ui basic segment active tab" data-tab="general">
          {m.component(General)}
        </div>
        <div class="ui basic segment tab" data-tab="plugins">
          {m.component(Plugins)}
        </div>
        <div class="ui basic segment tab" data-tab="about">
          {m.component(About)}
        </div>
      </div>
      <div class="ui secondary vertical footer segment">
        <button type="submit" class="ui teal right floated button" onclick={ctrl.onDoneClick}>Done</button>
      </div>
    </div>;
  }
};

m.mount(document.getElementById('app'), Settings)
//module.exports = Settings;
