function onVerdict(data) {
    if (data.verified) {
        window.location.reload();
        return
    }

    localStorage.setItem("targetUrlAfterCaptcha", window.location.href);
    localStorage.setItem("preCaptchaTimestamp", Date.now());
    window.location.href = "/.__core_/captcha";
}

function safeCheck(fn, fallback) {
    try {
      return fn();
    } catch (error) {
      console.error("safeCheck error:", error);
      return fallback;
    }
}

function main() {
    const reportedUserAgent = safeCheck(() => navigator.userAgent, null);
    const userAgentFails = reportedUserAgent === null || reportedUserAgent === undefined;
    
    const res = {
    userAgentFails: userAgentFails,
    usesWebDriver: safeCheck(() => !!navigator.webdriver, false),
    susProperties: safeCheck(
        () => window.__driver_unwrapped || window.__webdriver_script_fn || window.__driver_evaluate,
        false
    ),
    usesHeadlessChrome: safeCheck(() => {
        return reportedUserAgent?.includes("Headless");
    }, false),
    chromeDiscrepancy: safeCheck(() => {
        return reportedUserAgent?.includes("Chrome") && !window.chrome;
    }, false),
    lackingCodecSupport: safeCheck(() => {
        return document.createElement("video")
        .canPlayType('video/mp4; codecs="avc1.42E01E, mp4a.40.2"') === "";
    }, false),
    playwrightStealthPixelRatio: safeCheck(() => window.devicePixelRatio === 1.0000000149011612, false),
    reportedUserAgent: reportedUserAgent
    };
    
    fetch("/.__core_/api/judge", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(res)
    }).then(response => response.json()).then(data => onVerdict(data)).catch(() => {
        alert("A critical error occurred while checking your browser!");
    });
}

document.addEventListener("DOMContentLoaded", () => {
    main();
});
