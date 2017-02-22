function validateEmail(email) {
  var re = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return re.test(email);
}
function checkUser(user) {
	$.get(
    "/check/username/" + user,
    {},
    function(data) {
       if (data["taken"] == 1){
		   $("#username").css('border-color', "red");
		   $("#username-taken").slideDown(400);
	   }
	   else{
		   $("#username-taken").slideUp(400);
	   }
    }
);
}

function isClean() {
	var items = [$("#username"),$("#password"),$("#email"),$("#firstname"),$("#lastname")];
	for (var i = 0;i < items.length;i++){
		if (items[i].css('border-color') == "rgb(255, 0, 0)" || items[i].val() == ""){
			return false;
		}
	}
	return true;
}

$("#username").focus(function(){
	$("#username-alert").slideDown(400);
});
$("#username").focusout(function(){
	$("#username-alert").slideUp(400);
});
$("#password").focus(function(){
	$("#password-alert").slideDown(400);
});
$("#password").focusout(function(){
	$("#password-alert").slideUp(400);
});
$("#cteam-name").focus(function(){
	$("#cteam-name-alert").slideDown(400);
});
$("#cteam-name").focusout(function(){
	$("#cteam-name-alert").slideUp(400);
});
$(".alert button.close").click(function (e){
	$(this).parent().slideUp('slow');
});
$("#team-creation").click(function(){
	if ($(".team-join").is(":visible")){
		$(".team-join").slideUp(400);
	}
	$("#team-code-warn").show();
	$(".team-creation").slideToggle(400);
});
$("#team-join").click(function(){
	if ($(".team-creation").is(":visible")){
		$(".team-creation").slideUp(400);
	}
	$(".team-join").slideToggle(400);
});
// Start Checking
$("#firstname").focusout(function(){
    if(!(/^[a-zA-Z]+$/.test($("#firstname").val()))){
		$("#firstname").css('border-color', "red");
		clean = false
	}
	else{
		$("#firstname").css('border-color', "#ccc");
	}
});
$("#lastname").focusout(function(){
    if(!(/^[a-zA-Z]+$/.test($("#lastname").val()))){
		$("#lastname").css('border-color', "red");
	}
	else{
		$("#lastname").css('border-color', "#ccc");
	}
});
$("#email").focusout(function(){
    if(!(validateEmail($("#email").val()))){
		$("#email").css('border-color', "red");
	}
	else{
		$("#email").css('border-color', "#ccc");
	}
});
$("#username").focusout(function(){
	if(!(/^[a-zA-Z0-9]+$/.test($("#username").val())) || $("#username").val().length < parseInt($("#minulength").text()) || $("#username").val().length > parseInt($("#maxulength").text())){
		$("#username").css('border-color', "red");
	}
	else{
		$("#username").css('border-color', "#ccc");
		checkUser($("#username").val())
	}
});
$("#password").focusout(function(){
	if($("#password").val().length < parseInt($("#minplength").text()) || $("#password").val().length > parseInt($("#maxplength").text())){
		$("#password").css('border-color', "red");
	}
	else{
		$("#password").css('border-color', "#ccc");
	}
});
window.setInterval(function(){
	if (isClean()){
		$("#submit").removeClass("disabled");
	}
	else if (!($("#submit").hasClass("disabled"))){
		$("#submit").addClass("disabled");
	}
}, 100);
