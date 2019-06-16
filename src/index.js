require('./mystyles.scss');

const activeClass = "is-info";
const linux = "linux";
const osx = "osx";
const windows = "windows";
const arm = "arm64";
const amd = "amd64";
const client = "client";
const server = "server";
const node = "node";

const downloadInput = document.getElementById("dlurl");
const copyIcon = document.getElementById("clipboard");
const downloadButton = document.getElementById('downloadbutton');

function ButtonGroup(name, property) {
    this.name = name;
    this.property = property;

    // Computed
    this.activeData = "",
    this.buttons = [];

    this.update = function() {
        this.buttons = document.getElementById(this.name).children;
        this.activeData = document.getElementById(this.name).getElementsByClassName(activeClass)[0].dataset[this.property];
    },

    // populate activeData
    this.update();

    // set up active clicks for this button group
    enableActiveClicks(this);
}

let binary = new ButtonGroup("binaries", "binaries");
let os = new ButtonGroup("oses", "os");
let arch = new ButtonGroup("architectures", "architecture");
let versions = new ButtonGroup("versions", "version");

function init() {
    // set up the dependency between binaries and OS
    for (let i=0; i < binary.buttons.length; i++) {
        let child = binary.buttons[i];
        child.addEventListener("click", function(_) {
            // since binary buttons have requirements run those here
            updateOperatingSystemButton(os, this.dataset.binaries);

            // regenerate the URL after updating the active os.
            // TODO would prefer to use events to make this happen.
            os.update();
            generateURL();
        });
    }

    // set up dependency between arch and OS
    for (let i=0; i < arch.buttons.length; i++) {
        let child = arch.buttons[i];
        child.addEventListener("click", function(_){
            updateOperatingSystemButtonForArch(os, this.dataset.architecture);

            // regenerate the URL after updating the active os.
            // TODO would prefer to use events to make this happen.
            os.update();
            generateURL();
        });
    }

    // select all the text for easy copy/pasta
    var selectAll = true;
    downloadInput.addEventListener('click', function(_) {
        if (selectAll) {
            this.select();
        }
        selectAll = false;
    });
    downloadInput.addEventListener('blur', function(_){
        selectAll = true;
    })

    // make the checkmark work when you copy the url
    copyIcon.addEventListener('click', function(_) {
        downloadInput.select();
        document.execCommand("copy");
        let copyspan = copyIcon.querySelector('span');
        copyspan.innerHTML = `<i class="fas fa-check"></i>`;
        setTimeout(function() {
            copyspan.innerHTML = `<i class="far fa-clipboard"></i>`;
            unfade(copyspan);
        }, 500);
    });
}

init();
generateURL();

// Library functions

function unfade(element) {
    var op = 0.1;  // initial opacity
    element.style.display = 'block';
    var timer = setInterval(function () {
        if (op >= 1){
            clearInterval(timer);
        }
        element.style.opacity = op;
        element.style.filter = 'alpha(opacity=' + op * 100 + ")";
        op += op * 0.1;
    }, 10);
}

function generateURL() {
    // https://dl.k8s.io/v1.14.3/kubernetes-client-darwin-386.tar.gz
    let urlOS = osTransform(os.activeData);
    let url = `https://dl.k8s.io/v${versions.activeData}/kubernetes-${binary.activeData}-${urlOS}-${arch.activeData}.tar.gz`;
    downloadInput.value = url;
    downloadButton.href = url;
}

function osTransform(o) {
    if (o === osx) {
        return "darwin"
    }
    return o
}

function enableActiveClicks(bg) {
    let buttons = bg.buttons;
    for (let i=0; i < buttons.length; i++) {
        let child = buttons[i];
        child.addEventListener("click", function(_){
            this.classList.add(activeClass)
            for (let j=0; j < buttons.length; j++) {
                if (j == i) {
                    continue
                }
                buttons[j].classList.remove(activeClass);
            }
            bg.update();
            generateURL();
        });
    }
}

function updateOperatingSystemButton(buttonGroup, binary) {
    switch (binary) {
        case client:
            if (arch.activeData === arm) {
                break;
            }
            for (let i=0; i < buttonGroup.buttons.length; i++) {
                buttonGroup.buttons[i].disabled = false;
            }
            break;
        case server:
            // select linux, disable other two oses
            for (let i=0; i < buttonGroup.buttons.length; i++) {
                let child = buttonGroup.buttons[i];
                let os = child.dataset.os;
                if (os === linux) {
                    child.classList.add(activeClass);
                }
                if (os === windows || os === osx) {
                    child.classList.remove(activeClass);
                    child.disabled = true;
                }
            }
            break;
        case node:
            // enable windows but not if arm is selected
            for (let i=0; i < buttonGroup.buttons.length; i++) {
                if (arch.activeData === arm) {
                    continue
                }
                if (buttonGroup.buttons[i].dataset.os === windows) {
                    buttonGroup.buttons[i].disabled = false;
                }
            }
            // select linux, disable os x
            for (let i=0; i < buttonGroup.buttons.length; i++) {
                let child = buttonGroup.buttons[i];
                let os = child.dataset.os;
                if (os === linux) {
                    child.classList.add(activeClass);
                }
                if (os === osx) {
                    child.classList.remove(activeClass);
                    child.disabled = true;
                }
            }
            break;
        default:
            console.log("unknown binary");
    }
}

function updateOperatingSystemButtonForArch(buttonGroup, arch) {
    switch (arch) {
        case arm:
            // select linux, disable other two oses
            for (let i=0; i < buttonGroup.buttons.length; i++) {
                let child = buttonGroup.buttons[i];
                let os = child.dataset.os;
                if (os === linux) {
                    child.classList.add(activeClass);
                }
                if (os === windows || os === osx) {
                    child.classList.remove(activeClass);
                    child.disabled = true;
                }
            }
            break;
        case amd:
            // if binaires are node, enable windows
            if (binary.activeData === node) {
                for (let i=0; i < buttonGroup.buttons.length; i++) {
                    if (buttonGroup.buttons[i].dataset.os === windows) {
                        buttonGroup.buttons[i].disabled = false;
                    }
                }
            }
            break;
        default:
            break;
    }
}