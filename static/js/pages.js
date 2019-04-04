var editor;
var page;

$(document).ready(function() {
    var edit = document.getElementById('page-content');
    editor = CodeMirror.fromTextArea(edit, {
        mode: "htmlmixed",
        lineNumbers: true,
        theme: "gruvbox-dark"
    });
});

$("#pages").on("click", ".page-selector", function(e){
    e.preventDefault();
    page = $(this).attr('id').substring(5);
    $.get("/admin/pages/" + page).done(function(data){
        editor.setValue(data);
    }).fail(function(){
        console.log("failed");
    });
    $(".page-selector").each(function() {
        $(this).removeClass("active");
    });
    $(this).addClass("active");
    $("#save-page").prop("disabled", false);
    if (page != "index") {
        $("#delete-page").prop("disabled", false);
    } else {
        $("#delete-page").prop("disabled", true);
    }
});

$("#createPage").on("click", function(e){
    e.preventDefault();
    var name = $("#newFileName");
    if (/([^a-zA-Z\d_-])/.test(name.val()) || name.val() == "") {
        name.addClass("is-invalid");
        return;
    }
    name.removeClass("is-invalid");
    $.post("/admin/pages/" + name.val(), "{\"operation\": \"create\"}", function(){
        $("#pages").append("<a class=\"list-group-item list-group-item-action page-selector\" href=\"#\" id=\"page-" + name.val() + "\">" + name.val() + "</a>")
        $("#newModal").modal('hide');
    }).fail(function(res) {
        $("#incorrectText").text(res.responseText);
        $("#incorrect").show();
    });
})

$("#newModal").on('hidden.bs.modal', function(){
    $("#incorrect").hide();
})

$("#deleteModal").on('hidden.bs.modal', function(){
    $("#error").hide();
})

$("#delete-page").on("click", function(e) {
    e.preventDefault();
    if (page == "index" || page == null || page == undefined) {
        return;
    }
    $("#pageToDelete").text(page);
    $("#deleteModal").modal('show');
})

$("#deletePage").on("click", function(e) {
    e.preventDefault();
    if (page == "index" || page == null || page == undefined) {
        return;
    }
    $.post("/admin/pages/" + page, "{\"operation\": \"delete\"}", function(){
        $("#page-" + page).remove();
        $("#deleteModal").modal('hide');
    }).fail(function() {
        $("#error").show();
    });
})

$("#save-page").on("click", function(e) {
    e.preventDefault();
    if (page == null || page == undefined) {
        return;
    }
    data = JSON.stringify({operation: "update", contents: editor.getValue()})
    $.post("/admin/pages/" + page, data).fail(function() {
        $("#saveModal").modal('show');
    });
})