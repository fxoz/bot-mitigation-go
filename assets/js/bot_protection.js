let isBot = false;

function botDetected() {
    isBot = true;
    mainLoading.style.display = 'none';
    mainBot.style.display = 'block';
}

function checkVideoCodecSupport(mimeTypeWithCodecs) {
    const video = document.createElement("video");
    return video.canPlayType(mimeTypeWithCodecs);
}

function main() {
    statusText.innerText = 'Starting';

    if (navigator.webdriver) {
        statusText.innerText = 'Failed navigator.webdriver';
        return botDetected();
    }
    statusText.innerText = 'Passed navigator.webdriver';

    if (!!(window.__driver_unwrapped || window.__webdriver_script_fn || window.__driver_evaluate)) {
        statusText.innerText = 'Failed window globals';
        return botDetected();
    }
    statusText.innerText = 'Passed window globals';

    if (navigator.userAgent.includes("Headless")) {
        statusText.innerText = 'Failed user agent (headless)';
        return botDetected();
    }
    statusText.innerText = 'Passed user agent (headless)';

    //! Experimental
    try {
        navigator.userAgent
    } catch {
        statusText.innerText = 'Failed user agent (exist)';
        return botDetected();
    }
    statusText.innerText = 'Passed user agent (exist)';

    if (navigator.userAgent.includes("Chrome") && !window.chrome) {
        statusText.innerText = 'Failed chrome spoofing';
        return botDetected();
    }
    statusText.innerText = 'Passed chrome spoofing';

    if (checkVideoCodecSupport('video/mp4; codecs="avc1.42E01E, mp4a.40.2"') === "") {
        statusText.innerText = 'Failed codec';
        return botDetected();
    }

    if (window.devicePixelRatio === 1.0000000149011612) {
        statusText.innerText = 'Failed device pixel ratio';
        return botDetected();
    }
    statusText.innerText = 'Passed';
}

document.addEventListener('DOMContentLoaded', function () {
    main();
});