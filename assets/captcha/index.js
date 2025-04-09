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
        const x = e.offsetX | 0;
        const y = e.offsetY | 0;

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
                if (!data.verified) {
                    window.location.reload();
                } else {
                    window.location.href = "/";
                }
            })
    });
}

document.addEventListener("DOMContentLoaded", () => {
    main()
});