$(document).ready(function() {
    var edit = document.getElementById('page-content');
    var editor = CodeMirror.fromTextArea(edit, {
        mode: "htmlmixed",
        lineNumbers: true,
        theme: "gruvbox-dark"
    });
});