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

function checkCTeam(team) {
	$.get(
    "/check/team/" + team,
    {},
    function(data) {
       if (data["taken"] == 1){
		   $("#cteam-name").css('border-color', "red");
		   $("#cteam-taken").slideDown(400);
	   }
	   else{
		   $("#cteam-taken").slideUp(400);
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
	items = [$("#cteam-name"),$("#jteam-name"),$("#jteam-code")];
	for (var i = 0;i < items.length;i++){
		if (items[i].css('border-color') == "rgb(255, 0, 0)"){
			return false;
		}
	}
	if ($("#cteam-name").val() == "" && $("#jteam-name").val() == "" && $("#jteam-code").val() == ""){
		return false;
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
	$("#jteam-name").val("");
	$("#jteam-code").val("");
	$("#team-code-warn").show();
	$(".team-creation").slideToggle(400);
});
$("#team-join").click(function(){
	if ($(".team-creation").is(":visible")){
		$(".team-creation").slideUp(400);
	}
	$("#cteam-name").val("")
	$(".team-join").slideToggle(400);
});
// Start Checking
$("#firstname").keyup(function(){
    if(!(/^[a-zA-Z]+$/.test($(this).val()))){
		$(this).css('border-color', "red");
		clean = false
	}
	else{
		$(this).css('border-color', "#ccc");
	}
});
$("#lastname").keyup(function(){
    if(!(/^[a-zA-Z]+$/.test($(this).val()))){
		$(this).css('border-color', "red");
	}
	else{
		$(this).css('border-color', "#ccc");
	}
});
$("#email").keyup(function(){
    if(!(validateEmail($(this).val()))){
		$(this).css('border-color', "red");
	}
	else{
		$(this).css('border-color', "#ccc");
	}
});
$("#username").keyup(function(){
	if(!(/^[a-zA-Z0-9]+$/.test($(this).val())) || $(this).val().length < parseInt($("#minulength").text()) || $(this).val().length > parseInt($("#maxulength").text())){
		$(this).css('border-color', "red");
	}
	else{
		$(this).css('border-color', "#ccc");
		checkUser($(this).val())
	}
});
$("#password").keyup(function(){
	if($(this).val().length < parseInt($("#minplength").text()) || $(this).val().length > parseInt($("#maxplength").text())){
		$(this).css('border-color', "red");
	}
	else{
		$(this).css('border-color', "#ccc");
	}
});
$("#cteam-name").keyup(function(){
	if(!(/^[a-zA-Z0-9]+$/.test($(this).val())) || $(this).val().length < parseInt($("#mintlength").text()) || $(this).val().length > parseInt($("#maxtlength").text())){
		$(this).css('border-color', "red");
	}
	else{
		$(this).css('border-color', "#ccc");
		checkCTeam($(this).val())
	}
});
window.setInterval(function(){
	if (isClean()){
		$("#submit").prop("disabled", false);
	}
	else{
		$("#submit").prop("disabled", true);
	}
}, 100);
