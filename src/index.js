require('./mystyles.scss');

["os", "arch", "version"].map(function(val){
    eventListener(val);
});

function eventListener(kind) {
    let buttonGroupQuery = '#' + kind + ' > .button';
    let buttonGroup = document.querySelectorAll(buttonGroupQuery);

    let rowsQuery = 'tbody tr';
    let rows = document.querySelectorAll(rowsQuery);

    let hideClass = kind + "-hide";

    buttonGroup.forEach(button => {
        button.addEventListener('click', (evt) => {
            let buttonData = button.dataset[kind];
            buttonGroup.forEach(b => {
                b.classList.remove('is-success');
            })
            button.classList.add('is-success');
            rows.forEach(row => {
                if (row.classList.contains(buttonData)) {
                    row.classList.remove(hideClass);
                } else {
                    row.classList.add(hideClass);
                }
            });
        });
    });
}

// make the click link work
document.querySelectorAll(".copy").forEach(link => {
    link.addEventListener('click', (evt) => {
        evt.preventDefault();

        let el = document.createElement("textarea");
        el.value = link.href;
        el.setAttribute("readonly", "");
        el.style.position = 'absolute';
        el.style.left = '-9999px';
        document.body.appendChild(el);
        el.select();
        document.execCommand("copy");
        document.body.removeChild(el);
        return false;
    });
});
