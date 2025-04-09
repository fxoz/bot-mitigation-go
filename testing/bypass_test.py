import requests


def main():
    res = requests.post(
        "http://localhost:9977/.__-INTERNAL-/api/__judge",
        headers={
            "Accept": "*/*",
            "Accept-Language": "en-US,en;q=0.9",
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "Content-Type": "application/json",
            "Origin": "http://localhost:9977",
            "Pragma": "no-cache",
            "Referer": "http://localhost:9977/",
            "Sec-Fetch-Dest": "empty",
            "Sec-Fetch-Mode": "cors",
            "Sec-Fetch-Site": "same-origin",
            "Sec-GPC": "1",
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36",
        },
        json={
            "userAgentFails": False,
            "susProperties": False,
            "usesWebDriver": False,
            "usesHeadlessChrome": False,
            "chromeDiscrepancy": False,
            "lackingCodecSupport": False,
            "playwrightStealthPixelRatio": False,
            "reportedUserAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36",
        },
    )
    print(res.text)
    print(res.status_code)


if __name__ == "__main__":
    main()
