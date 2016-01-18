const m = require('mithril');
var open = require('open');

var About = {
  controller: function controller() {
    var ctrl = this;

    ctrl.onLinkClick = function(e) {
      e.preventDefault();
      open(e.target.href);
      return false;
    };
  },

  view: function view(ctrl) {
    return <div>
      <h3 class="ui header">About</h3>
      <div class="ui center aligned container">
        <img src="../../icon.png" width="100" />
        <h1 class="ui header">GoLaunch</h1>
        <p>
        v1.0
        </p>
        <p>
        Made by Kevin Darlington
        </p>
        <p>
        <a class="ui black button" href="https://github.com/kdar" onclick={ctrl.onLinkClick}>Github</a>
        </p>
      </div>
    </div>;
  }
};

module.exports = About;
