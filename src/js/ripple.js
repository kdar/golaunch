module.exports = function(el, pageX, pageY) {
  var ink = el.querySelector('.ink');
  if (!ink) {
    el.insertAdjacentHTML('afterbegin', "<span class='ink'></span>");
    ink = el.querySelector('.ink');
  }

  ink.classList.remove('animate');

  var d = Math.max(el.offsetHeight, el.offsetWidth);
  ink.style.height = d+"px";
  ink.style.width = d+"px";

  var rect = el.getBoundingClientRect();
  var elOffsetLeft = rect.left + document.body.scrollLeft;
  var elOffsetTop = rect.top + document.body.scrollTop;

  var x = pageX - elOffsetLeft  - ink.offsetWidth/2;
	var y = pageY - elOffsetTop - ink.offsetHeight/2;

  ink.style.top = y+'px';
  ink.style.left = x+'px';
  ink.classList.add('animate');
};

// var parent, ink, d, x, y;
// $("ul li a").click(function(e){
// 	parent = $(this).parent();
// 	//create .ink element if it doesn't exist
// 	if(parent.find(".ink").length == 0)
// 		parent.prepend("<span class='ink'></span>");
//
// 	ink = parent.find(".ink");
// 	//incase of quick double clicks stop the previous animation
// 	ink.removeClass("animate");
//
// 	//set size of .ink
// 	if(!ink.height() && !ink.width())
// 	{
// 		//use parent's width or height whichever is larger for the diameter to make a circle which can cover the entire element.
// 		d = Math.max(parent.outerWidth(), parent.outerHeight());
// 		ink.css({height: d, width: d});
// 	}
//
// 	//get click coordinates
// 	//logic = click coordinates relative to page - parent's position relative to page - half of self height/width to make it controllable from the center;
// 	x = e.pageX - parent.offset().left - ink.width()/2;
// 	y = e.pageY - parent.offset().top - ink.height()/2;
//
// 	//set the position and add class .animate
// 	ink.css({top: y+'px', left: x+'px'}).addClass("animate");
// })
