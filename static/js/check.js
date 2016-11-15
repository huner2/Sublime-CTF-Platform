var clean = false;

setInterval(function(){
	if (clean){
		$("#submit").disabled = true;
	}
}, 1000);