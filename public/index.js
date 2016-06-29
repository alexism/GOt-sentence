$(document).ready(function() {

    var loaded = [4];
    var intervalVar;

    function hmmDotting() {
        if (loaded[0] > 0) {
            $('#suggestion')[0].textContent = 'hmm' + ['', '.', '..', '...'][4 - loaded[0]];
            loaded[0] = loaded[0] - 1;
        } else {
            clearInterval(intervalVar);
        }
    }

    intervalVar = window.setInterval(function() {
        hmmDotting();
    }, 300);

    function changeWord() {
       
        $.get('http://'+document.location.hostname+':8080/ppf', function(data) {
            loaded[0] = 0;
            $('#suggestion')[0].textContent = data['sentence'];
            clearInterval(intervalVar);
        });
    }

    document.body.addEventListener("mousedown", tapOrClick, false);
    document.body.addEventListener("touchstart", tapOrClick, false);

    function tapOrClick(event) {
        changeWord();

        event.preventDefault();
        return false;
    }

    changeWord();

    var down = {};

    $(document).keydown(function(event){
        if (down[event.which] == null) {
            changeWord();
            down[event.which] = true;
        }
    });

    $(document).keyup(function(event) {
        down[event.which] = null;
    });

});
