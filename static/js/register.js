var clean = false;

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

window.setInterval(function(){
	if (clean){
		$("#submit").removeClass("disabled");
	}
	else if (!($("#submit").hasClass("disabled"))){
		$("#submit").addClass("disabled");
	}
}, 1000);
