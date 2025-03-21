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
    if (navigator.webdriver) {
        return botDetected();
    }

    if (!!(window.__driver_unwrapped || window.__webdriver_script_fn)) {
        return botDetected();
    }

    if (navigator.userAgent.includes("Chrome") && !window.chrome) {
        return botDetected();
    }

    if (navigator.userAgent.includes("Headless")) {
        return botDetected();
    }

    if (checkVideoCodecSupport('video/mp4; codecs="avc1.42E01E, mp4a.40.2"') === "") {
        return botDetected();
    }

    window.location.href = '{-URL-}';
}

document.addEventListener('DOMContentLoaded', function () {
    main();
});