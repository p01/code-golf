'use strict';

onload = function() {
    let selects = document.querySelectorAll('select');

    for (let select of selects)
        select.onchange = function() {
            location = '/scores/' + selects[0].value + '/' + selects[1].value;
        };
};
