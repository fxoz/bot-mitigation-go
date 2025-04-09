function main() {
    captchaImage = document.getElementById("captchaImage");

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
            .then(data => {
                if (data.verified) {
                    alert("Verified!");
                    window.location.href = "/";
                } else {
                    alert("Failed.");
                    window.location.reload();
                }
            }).catch(error => {
                alert(`Error verifying captcha: ${error}`);
            });
    });
}

document.addEventListener("DOMContentLoaded", () => {
    main()
});