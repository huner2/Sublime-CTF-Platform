function countdown(){
  var now = new Date().getTime();
  var dist = countDownDate - now;
  var d = Math.floor(dist / (1000 * 60 * 60 * 24));
  var h = Math.floor((dist % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
  var m = Math.floor((dist % (1000 * 60 * 60)) / (1000 * 60));
  var s = Math.floor((dist % (1000 * 60)) / 1000);
  if (s < 0){
    return s;
  }

  $('#countdown').text(d + "D " + h + "H " + m + "M " + s + "S");

  return s;
}

var timestamp = $('#timestamp').data().timestamp;
var countDownDate = new Date(timestamp);
//countDownDate = new Date(countDownDate.valueOf() + new Date().getTimezoneOffset() * 60 * 1000); // Off by an hour for some reason
countdown();
var x = setInterval(function() {
  var dist = countdown();
  if (dist < 0){
    $('#countdown').text("0D 0H 0M 0S"); // Ensure no negatives
    clearInterval(x);
  }
}, 1000);
