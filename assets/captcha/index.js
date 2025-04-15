function onCaptchaVerdict(data) {
    if (data.verified) {
        const preCaptchaTimestamp = localStorage.getItem("preCaptchaTimestamp")
        if (preCaptchaTimestamp) {
            if ((Date.now() - preCaptchaTimestamp ) < 60000) {
                window.location.href = localStorage.getItem("targetUrlAfterCaptcha");
                return
            }
        }
        window.location.href = "/";
        return
    }

    if (data.exceeded) {
        exceededAttempts.style.display = "block";
        document.querySelector('main').style.display = "none";
        return;
    }

    const baseUrl = window.location.origin + window.location.pathname;
    window.location.href = `${baseUrl}?failed=1`;
}

function main() {
    captchaImage = document.getElementById("captchaImage");
    failedAttempt = document.getElementById("failedAttempt");
    exceededAttempts = document.getElementById("exceededAttempts");

    if (window.location.search.includes("failed=1")) {
        failedAttempt.style.display = "block";
    }

    fetch("/.__core_/api/captcha/generate")
        .then(res => {
            if (!res.ok) {
                throw new Error("Network response was not ok");
            }
            return res.json();
        })
        .then(data => {
            captchaImage.src = data.image;
        })
        .catch(error => {
            alert(`Error fetching captcha image: ${error}`);
        });

    captchaImage.addEventListener("click", (e) => {
        const x = e.offsetX / captchaImage.width
        const y = e.offsetY / captchaImage.height
        console.log(`Captcha clicked at coordinates: (${x}, ${y})`);

        fetch("/.__core_/api/captcha/verify", {
            method: "POST",
            body: JSON.stringify({
                x, y
            }),
            headers: {
                "Content-Type": "application/json"
            }
        })
            .then(res => {
                if (!res.ok) {
                    throw new Error("Network response was not ok")
                }

                return res.json();
            })
            .then(data => onCaptchaVerdict(data)).catch(error => {
                alert(`Error verifying captcha: ${error}`);
            });
    });
}

document.addEventListener("DOMContentLoaded", () => {
    main()
});