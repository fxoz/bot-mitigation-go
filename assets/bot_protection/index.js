function main() {
    userAgentFails = false;
    reportedUserAgent = null;
    try {
        navigator.userAgent
    } catch {
        userAgentFails = true;
    }

    const res = {
        userAgentFails: userAgentFails,
        usesWebDriver: !!navigator.webdriver,
        susProperties: window.__driver_unwrapped || window.__webdriver_script_fn || window.__driver_evaluate,
        usesHeadlessChrome: navigator.userAgent.includes("Headless"),
        chromeDiscrepancy: navigator.userAgent.includes("Chrome") && !window.chrome,
        lackingCodecSupport: document.createElement("video").canPlayType('video/mp4; codecs="avc1.42E01E, mp4a.40.2"') === "",
        playwrightStealthPixelRatio: window.devicePixelRatio === 1.0000000149011612,
        reportedUserAgent: reportedUserAgent,
    }

    fetch("/.__/api/__judge", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(res)
    }).then(response => response.json()).then(data => {
        if (data.verified) {
            window.location.reload();
        } else {
            window.location.href = "/.__captcha/";
        }
    }).catch(() => {
        alert("A critical error occurred while checking your browser!");
    });
}

document.addEventListener('DOMContentLoaded', () => {
    main();
});