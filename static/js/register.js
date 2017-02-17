function validateEmail(email) {
  var re = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return re.test(email);
}

function isClean() {
	return false;
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
	$(".alert-dismissible").show();
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
	if(/^[a-zA-Z0-9]+$/.test($("#username")) || $("#username").val().length < parseInt($("#minulength").text()) || $("#username").val().length > parseInt($("#maxulength").text())){
		$("#username").css('border-color', "red");
	}
	else{
		$("#username").css('border-color', "#ccc");
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
}, 1000);
