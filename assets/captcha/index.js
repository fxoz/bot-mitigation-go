function main() {
    captchaImage = document.getElementById("captchaImage");

    captchaImage.addEventListener("click", (e) => {
        const x = e.offsetX;
        const y = e.offsetY;

        alert(`Clicked at ${x}, ${y}`);
    });
}

document.addEventListener("DOMContentLoaded", () => {
    main()
});