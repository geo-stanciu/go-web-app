function getHttpRequest() {
    var xmlhttp;
    if (window.XMLHttpRequest) { // code for IE7+, Firefox, Chrome, Opera, Safari
        xmlhttp = new XMLHttpRequest();
    } else { // code for IE6, IE5
        xmlhttp = new ActiveXObject("Microsoft.XMLHTTP");
    }

    return xmlhttp;
}

function getAJAX(path, params, callback) {
    var method = "GET";
    var url = path + "?lrt=" + (new Date().getTime());

    var str = [];

    for (var key in params) {
        if (params.hasOwnProperty(key)) {
            str.push(encodeURIComponent(key) + "=" + encodeURIComponent(params[key]));
        }
    }

    url += "&" + str.join("&");

    var xhr = getHttpRequest();
    xhr.open(method, url, true);

    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4 && xhr.status == 200) {
            if (callback != null && callback != undefined) {
                callback(xhr.responseText);
            }
        }
    }

    var token = document.getElementsByTagName("meta")["csrf.Token"];
    if (token != undefined) {
        xhr.setRequestHeader("X-CSRF-Token", token.getAttribute("content"));
    }

    xhr.send();
}

function postAJAX(path, params, callback) {
    var method = "POST";
    var url = path + "?lrt=" + (new Date().getTime());

    var str = [];

    for (var key in params) {
        if (params.hasOwnProperty(key)) {
            str.push(encodeURIComponent(key) + "=" + encodeURIComponent(params[key]));
        }
    }

    var xhr = getHttpRequest();
    xhr.open(method, url, true);

    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4 && xhr.status == 200) {
            if (callback != null && callback != undefined) {
                callback(xhr.responseText);
            }
        }
    }

    var token = document.getElementsByTagName("meta")["csrf.Token"];
    if (token != undefined) {
        xhr.setRequestHeader("X-CSRF-Token", token.getAttribute("content"));
    }

    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xhr.send(str.join("&"));
}

function sendPOST(path, params) {
    var method = "post";
    var url = path + "?lrt=" + (new Date().getTime());

    // The rest of this code assumes you are not using a library.
    // It can be made less wordy if you use one.
    var form = document.createElement("form");
    form.target = '_blank';
    form.setAttribute("method", method);
    form.setAttribute("action", url);

    for (var key in params) {
        if (params.hasOwnProperty(key)) {
            var hiddenField = document.createElement("input");
            hiddenField.setAttribute("type", "hidden");
            hiddenField.setAttribute("name", encodeURIComponent(key));
            hiddenField.setAttribute("value", encodeURIComponent(params[key]));

            form.appendChild(hiddenField);
        }
    }

    var token = document.getElementsByTagName("meta")["csrf.Token"];
    if (token != undefined) {
        var hiddenField = document.createElement("input");
        hiddenField.setAttribute("type", "hidden");
        hiddenField.setAttribute("name", encodeURIComponent(key));
        hiddenField.setAttribute("value", encodeURIComponent(token.getAttribute("content")));

        form.appendChild(hiddenField);
    }

    document.body.appendChild(form);
    form.submit();
    document.body.removeChild(form);
}