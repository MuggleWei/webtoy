
const headers = new Headers({
    "Accept": "application/json",
    "Content-Type": "application/json",
});

function get(url) {
    return fetch(url, {
        method: "GET",
        headers: fillupSession(headers),
    }).then(response => {
        return handleResponse(url, response);
    }).catch(err => {
        let errMsg = `requests.get failed. url=${url}, msg=${err}`;
        console.error(errMsg);
        return { code: -1, msg: errMsg };
    })
}

function post(url, data) {
    return fetch(url, {
        method: "POST",
        headers: fillupSession(headers),
        body: JSON.stringify(data)
    }).then(response => {
        return handleResponse(url, response);
    }).catch(err => {
        let errMsg = `requests.post failed. url=${url}, msg=${err}`;
        console.error(errMsg);
        return { code: -1, msg: errMsg };
    })
}

function put(url, data) {
    return fetch(url, {
        method: "PUT",
        headers: fillupSession(headers),
        body: JSON.stringify(data)
    }).then(response => {
        return handleResponse(url, response);
    }).catch(err => {
        let errMsg = `requests.put failed. url=${url}, msg=${err}`;
        console.error(errMsg);
        return { code: -1, msg: errMsg };
    })
}

function fillupSession(headers) {
    let session = sessionStorage.getItem("session");
    let token = sessionStorage.getItem("token");
    if (!session || !token) {
        // session = Cookies.get("session");
        // token = Cookies.get("token");
        session = localStorage.getItem("session");
        token = localStorage.getItem("token");
        if (session && token) {
            sessionStorage.setItem("session", session);
            sessionStorage.setItem("token", token);
        }
    }

    if (session && token) {
        headers.set("session", session);
        headers.set("token", token);
    }
    return headers;
}

function handleResponse(url, response) {
    if (response.status === 200) {
        return response.json();
    } else {
        let errMsg = `requests response error. url=${url}, status=${response.status}, msg=${response.statusText}`;
        console.error(errMsg);
        return { code: -1, msg: errMsg };
    }
}

const requests = { get, post, put };
export default requests;