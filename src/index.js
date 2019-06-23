require('./mystyles.scss');

["os", "arch", "version"].map(function(val){
    eventListener(val);
});

function eventListener(kind) {
    document.getElementById(kind).addEventListener("click", function(evt){
        let hideClass = kind + "-hide";
        let k = evt.target.dataset[kind];
        let notSelector = "tbody tr:not(."+k+")";
        let selector = "tbody tr."+k;

        var rows;

        // if this is the actively selected one then clear stuff
        if (evt.srcElement.classList.contains("is-success")) {
            evt.srcElement.classList.remove("is-success")
            rows = document.querySelectorAll(notSelector);
            for (let i=0; i < rows.length; i++) {
                rows.item(i).classList.remove(hideClass)
            }
            return
        }

        for (let i=0; i < this.children.length; i++) {
            if (this.children.item(i).dataset[kind] === evt.srcElement.dataset[kind]) {
                evt.srcElement.classList.add("is-success");
                continue
            }
            this.children.item(i).classList.remove("is-success");
        }
        rows = document.querySelectorAll(notSelector);
        for (let i=0; i < rows.length; i++) {
            rows.item(i).classList.add(hideClass)
        }
        rows = document.querySelectorAll(selector);
        for (let i=0; i < rows.length; i++) {
            rows.item(i).classList.remove(hideClass)
        }
    });
}

let links = document.getElementsByClassName("copy");
for (let i=0; i < links.length; i++) {
    links.item(i).addEventListener('click', function(evt) {
        evt.preventDefault();

        let el = document.createElement("textarea");
        el.value = this.href;
        el.setAttribute("readonly", "");
        el.style.position = 'absolute';
        el.style.left = '-9999px';
        document.body.appendChild(el);
        el.select();
        document.execCommand("copy");
        document.body.removeChild(el);
    });
}

