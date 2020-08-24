/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
