var electron = require('electron');
var remote = electron.remote;
var m = require('mithril');
var mousetrap = require('mousetrap');
// var ripple = require('./js/ripple');

var Menu = remote.require('menu');
var MenuItem = remote.require('menu-item');
var baseSize = remote.getCurrentWindow().getSize();
var pluginManager = remote.getGlobal("pluginManager");
var themeManager = remote.getGlobal("themeManager");

// var ipc = electron.ipcRenderer;
// var msg = {
//   title : "Awesome!",
//   description : "Check this out!<br>Check this out!<br>Check this out!<br>Check this out!<br>Check this out!<br>Check this out!<br>",
//   width : 440,
//   timeout : 6000,
//   focus: true // set focus back to main window
// };
// ipc.send('golaunch-toast', msg);

var contextmenu = new Menu();

var LauncherVm = (function() {
  var vm = {
    //state: {}
  };
  vm.init = function() {
    vm.queryResults = m.prop([]);
    vm.progress = m.prop({
      count: 0,
      current: 0
    });
  }
  return vm;
}())

LauncherVm.init();

document.querySelector("#theme").innerHTML = themeManager.get();

var resultsCount = 0;
pluginManager.on('query-results', function(results) {
  LauncherVm.queryResults(results);
  m.redraw();
});

pluginManager.on('progress-update', function(progress) {
  LauncherVm.progress(progress);
  m.redraw();
});


function isScrolledIntoView(element, parent) {
  var elementTop    = element.getBoundingClientRect().top,
      elementBottom = element.getBoundingClientRect().bottom;

  var parentTop = parent.getBoundingClientRect().top,
      parentBottom = parent.getBoundingClientRect().bottom;

  return elementTop >= 0 &&  (Math.abs(elementBottom - parentTop) > 10) && (Math.abs(elementTop - parentBottom) > 10);
}

function debounce(func, wait, immediate) {
	var timeout;
	return function() {
		var context = this, args = arguments;
		var later = function() {
			timeout = null;
			if (!immediate) func.apply(context, args);
		};
		var callNow = immediate && !timeout;
		clearTimeout(timeout);
		timeout = setTimeout(later, wait);
		if (callNow) func.apply(context, args);
	};
};

var Launcher = {
  config: function(el, isInitialized) {
    var ctrl = this;
    ctrl.element = el;

    ctrl.searchText = document.getElementsByClassName("query-text")[0];
    ctrl.searchText.focus();
    ctrl.queryResults = document.getElementsByClassName("query-results")[0];

    if (!isInitialized) {
      function updown(dir) {
        if (dir == -1 && ctrl.selected() != 0) {
          ctrl.selected(ctrl.selected()+dir);
        } else if (dir == 1 && ctrl.selected() < LauncherVm.queryResults().length-1) {
          ctrl.selected((ctrl.selected()+dir));
        }

        if (LauncherVm.queryResults().length > 0) {
          var el = ctrl.queryResults.querySelectorAll("ul li")[ctrl.selected()];
          if (!isScrolledIntoView(el, ctrl.queryResults)) {
            ctrl.queryResults.scrollTop = dir == 1 ?
              ctrl.queryResults.scrollTop+el.offsetHeight :
              ctrl.queryResults.scrollTop-el.offsetHeight;
          }
        }
      }

      Mousetrap.bind('up', function(e) {
        updown(-1);
        return false;
      });

      Mousetrap.bind('down', function(e) {
        updown(1);
        return false;
      });

      remote.getCurrentWindow().on('focus', function() {
        ctrl.searchText.setSelectionRange(0, ctrl.searchText.value.length);
      });
    }
  },

  controller: function() {
    var ctrl = this;

    ctrl.selected = m.prop(0);

    ctrl.doAction = function() {
      pluginManager.pluginAction({
        type: "execute",
        queryResult: LauncherVm.queryResults()[ctrl.selected()]
      });

      // FIXME: should let plugins determine if the window hides.
      remote.getCurrentWindow().hide();
    };

    ctrl.onQueryKeyUp = function(event) {
      if (event.which == 13 || event.keyCode == 13) {
        ctrl.doAction();
        return false;
      }

      return true;
    };

    var debouncedQuery = debounce(function(query) {
      pluginManager.pluginQuery(query);
    }, 150);

    ctrl.onQueryInput = function() {
      if (ctrl.searchText.value != "") {
        debouncedQuery(ctrl.searchText.value);
        ctrl.selected(0);
      } else {
        LauncherVm.queryResults([]);
        LauncherVm.progress({
          count: 0,
          current: 0
        });
        ctrl.selected(0);
        pluginManager.clearQuery();
      }
    };

    ctrl.onQueryBlur = function() {
      ctrl.searchText.focus();
    };

    ctrl.queryResultsConfig = function(el, isInitialized) {
      // This is a workaround since electron does not support click-through with
      // transparent windows. Once that is implemented, we can get rid of resizing
      // the browser window.
      if (LauncherVm.queryResults().length != 0) {
        remote.getCurrentWindow().setSize(baseSize[0], baseSize[1] + ctrl.queryResults.offsetHeight);
      } else {
        remote.getCurrentWindow().setSize(baseSize[0], baseSize[1]);
      }

      if (!isInitialized) {
        // sub title scrolling
        var slideTimer, slide = function(el) {
          el.scrollLeft += 1;
          if (el.scrollLeft < el.scrollWidth) {
            slideTimer = setTimeout(function() {slide(el);}, 8);
          }
        };
        el.onmouseover = el.onmouseout = function(e) {
          if (e.target.tagName != "H2") {
            return;
          }
          var over = e.type === 'mouseover';
          clearTimeout(slideTimer);
          if (over) {
            e.target.classList.remove("hiding");
            slide(e.target);
          } else {
            e.target.classList.add("hiding");
            e.target.scrollLeft = 0;
          }
        }

        el.addEventListener('contextmenu', function(e) {
          e.preventDefault();

          var child = e.target;
          while (child != null && child.tagName != "LI") {
            child = child.parentNode;
          }

          if (!child) {
            return;
          }

          // ripple(child, e.pageX, e.pageY);

          var i = 0;
          while ((child = child.previousSibling) != null)
            i++;

          // this is not async. this will block the renderer:
          // https://github.com/atom/electron/issues/1854
          var cm = LauncherVm.queryResults()[i].contextmenu;
          if (cm) {
            contextmenu.clear();
            for (var x = 0; x < cm.length; x++) {
              contextmenu.append(new MenuItem(Object.assign({}, cm[x], {
                click: function(menuItem, browserWindow) {
                  pluginManager.pluginAction({
                    type: "contextmenu",
                    name: menuItem.label,
                    queryResult: LauncherVm.queryResults()[ctrl.selected()]
                  });
                }
              })));
            }
            contextmenu.popup(remote.getCurrentWindow());
          }
        }, false);
      }
    };

    ctrl.onQueryItemClick = function(e, index) {
      ctrl.selected(index);
    };

    ctrl.onQueryItemDblClick = function(e, index) {
      ctrl.doAction();
    };

    ctrl.onQueryResultsScroll = function(e) {
      // LauncherVm.state.scrollTop = e.target.scrollTop;
      // LauncherVm.state.height = e.target.offsetHeight;
      // m.redraw();
    };
  },

  view: function(ctrl) {
    // var scrollTop = LauncherVm.state.scrollTop;
    // var begin = scrollTop / 57 | 0
  	// var end = begin + (LauncherVm.state.height / 57 | 0 + 3)
  	// var offset = scrollTop % 57
    var progressValue = 0;
    var progressClass = "";
    var progress = LauncherVm.progress();
    if (progress.count != 0) {
      progressValue = Math.round(progress.current / progress.count * 100.0);

      if (progressValue == 100) {
        progressClass = "done";
      } else if (progressValue < 100) {
        progressClass = "working";
      }
    }

    return <div id="launcher" config={Launcher.config.bind(ctrl)}>
      <div class="query-text-wrapper">
        <input type="text" class="query-text mousetrap" onkeyup={ctrl.onQueryKeyUp} onblur={ctrl.onQueryBlur} oninput={ctrl.onQueryInput}/>
        <progress max="100" value={progressValue} class={progressClass}></progress>
      </div>

      <div class="query-results" config={ctrl.queryResultsConfig.bind(ctrl)} onscroll={ctrl.onQueryResultsScroll}>
        <ul>
        {LauncherVm.queryResults().map(function(item, index) {
          var cls = "";
          if (index == ctrl.selected()) {
            cls = "selected";
          }
          return <li class={cls} onclick={function(e) { ctrl.onQueryItemClick(e, index); }} ondblclick={function(e) { ctrl.onQueryItemDblClick(e, index); }}>
            <img src={item.icon}></img>
            <h1>{item.title}</h1>
            <h2>{item.subtitle}</h2>
          </li>;
        })}
        </ul>
      </div>
    </div>;
  }
};

module.exports = Launcher;
