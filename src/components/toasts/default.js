function getQueryVariable(variable) {
	var query = window.location.search.substring(1);
	var vars = query.split("&");
	for (var i=0;i<vars.length;i++) {
		var pair = vars[i].split("=");
		if (pair[0] === variable) {
			return decodeURIComponent(pair[1]);
		}
	}
}

function outerHeight(el) {
  var height = el.offsetHeight;
  var style = getComputedStyle(el);

  height += parseInt(style.marginTop) + parseInt(style.marginBottom);
  return height;
}

function outerWidth(el) {
	var width = el.offsetWidth;
	var style = getComputedStyle(el);

	width += parseInt(style.marginLeft) + parseInt(style.marginRight);
	return width;
}

var fixToContent = function() {
	var result = {
		height: outerHeight(document.body),
		width: outerWidth(document.body)
	}

	window.resizeTo(result.width, result.height);

	return result;
};

var onKeydown = function() {
	//window.close();
};

var onLoad = function load() {
	this.removeEventListener("load", load, false);

	fixToContent();

	this.setTimeout(function() {
		//this.close();
	}, parseInt(getQueryVariable("timeout")));

	document.addEventListener("keydown", onKeydown, false);
	//document.addEventListener("click", window.close);
};

document.querySelector(".title").innerHTML = getQueryVariable("title");
document.querySelector(".description").innerHTML = getQueryVariable("description");
window.addEventListener("load", onLoad, false);
